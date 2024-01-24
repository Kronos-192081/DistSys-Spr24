package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"math/rand"
	"time"
	"errors"
	"prakhar/conhash"
	"strings"
	"strconv"
	"github.com/fsouza/go-dockerclient"
)

// Environment variables for configuring the load balancer
var num_serv , _ = strconv.Atoi(os.Getenv("NUM_SERV"))
var num_slots, _ = strconv.Atoi(os.Getenv("NUM_SLOTS"))
var num_virt_serv, _ = strconv.Atoi(os.Getenv("NUM_VIRT_SERV"))

// Constants
const Mod = 1e4 + 7

// Global consistent hash instance
var c = conhash.NewConHash(num_slots, num_virt_serv)

// Mutex for thread-safe operations
var mtx sync.Mutex

// Main function
func main() {
	fmt.Println("Starting load balancer")

	// Seed for randomization
	rand.NewSource(time.Now().UnixNano())

	// Add server containers based on environment variables
	for i := 0; i < num_serv; i++ {
		addServerContainer("Server_" + strconv.Itoa(i + 1), rand.Intn(Mod))
	}

	// List all server containers
	listServerContainers()

	// Setup HTTP servers for different endpoints
	http.HandleFunc("/rep", rep)
	repSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/add", add)
	addSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/rm", rm)
	rmSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/", path)
	pathSrv := &http.Server{Addr: "0.0.0.0:5000"}

	// Setup context and signal handling
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	// Start HTTP servers in separate goroutines
	go func() {
		repSrv.ListenAndServe()
	}()

	go func() {
		addSrv.ListenAndServe()
	}()

	go func() {
		rmSrv.ListenAndServe()
	}()

	go func() {
		pathSrv.ListenAndServe()
	}()

	// Defer shutdown of servers
	defer func() {
		// Graceful shutdown of servers
		if err := repSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the rep server: ", err)
		}
		if err := addSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the add server: ", err)
		}
		if err := rmSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the rm server: ", err)
		}
		if err := pathSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the path server: ", err)
		}
	}()

	// Wait for SIGINT signal
	sig := <-sigs
	fmt.Println(sig)

	// Cancel the context to initiate shutdown
	cancel()

	fmt.Println("Shutting down load balancer")
}

// Structs for representing JSON responses
type ServDetails struct {
	N        int
	Replicas []string `json:"replicas"`
}

type Payload struct {
    N           int `json:"n"`
    Hostnames   []string `json:"hostnames"`
}

type ResponseSuccess struct {
	Message ServDetails `json:"message"`
	Status  string `json:"status"`
}

type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// Utility function to generate a random server name
func GenerateRandomString(num int) string {
	name := "spawned_server_" + strconv.Itoa(num)

	return name
}

// Fisher-Yates algorithm for random permutation of a slice
func permuteSlice(slice []string) {
	rand.NewSource(time.Now().UnixNano())
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Handler functions for incoming requests

// Handler for /rep endpoint (GET)
func rep(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		// Get the list of server containers
		servNames, err := listServerContainers()
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Prepare and send JSON response
		servData := ServDetails{
			N:        c.Nserv,
			Replicas: servNames,
		}
		resp := ResponseSuccess{
			Message: servData,
			Status:  "successful",
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.Header().Set("Content-Type", "application/json")
        rw.WriteHeader(http.StatusOK)
		rw.Write(jsonResp)
	default:
		// Handle unsupported methods
		rw.WriteHeader(http.StatusNotFound)
	}
}

// Handler for /add endpoint (POST)
func add(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// Decode the JSON payload
        var payloadData Payload
        err := json.NewDecoder(req.Body).Decode(&payloadData)
        if err != nil {
            fmt.Println("Error:", err)
            rw.WriteHeader(http.StatusInternalServerError)
            return
        }
		
		rand.NewSource(time.Now().UnixNano())

		// Add server containers based on the payload
		for i := 0; i < len(payloadData.Hostnames); i++ {
			err := addServerContainer(payloadData.Hostnames[i], rand.Intn(Mod))
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// Check if the number of servers requested is greater than the added servers
		if payloadData.N >= len(payloadData.Hostnames) {
			// Calculate the extra servers needed
			extraServ := payloadData.N - len(payloadData.Hostnames)

			// Add randomly generated servers
			for i := 0; i < extraServ; i++ {
				for {
					num := rand.Intn(Mod)
					name := GenerateRandomString(num)
					if _, ok := c.AllServers[name]; ok {
						continue
					}
					err = addServerContainer(name, num)
					if err != nil {
						fmt.Println("Error:", err)
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}
					break
				}
			}

			// Get the updated list of server containers
			servNames, err := listServerContainers()
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Prepare and send JSON response
			servData := ServDetails{
				N:        c.Nserv,
				Replicas: servNames,
			}
			resp := ResponseSuccess{
				Message: servData,
				Status:  "successful",
			}
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(jsonResp)
		} else {
			// If the number of servers requested is less than the added servers, return an error response
			resp := Response{
				Message: "ERROR: Length of hostname list is more than newly added instances",
				Status:  "failure",
			}
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(jsonResp)
		}
	default:
		// Handle unsupported methods
		rw.WriteHeader(http.StatusNotFound)
	}
}

// Handler for /rm endpoint (DELETE)
func rm(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodDelete:
		// Decode the JSON payload
        var payloadData Payload
        err := json.NewDecoder(req.Body).Decode(&payloadData)
        if err != nil {
            fmt.Println("Error:", err)
            rw.WriteHeader(http.StatusInternalServerError)
            return
        }
		// Remove specified server containers
		for _, servName := range payloadData.Hostnames {
			err := killServerContainer(servName)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		// Check if the number of servers requested is greater than the removed servers
		if payloadData.N >= len(payloadData.Hostnames) {
			// Get the current list of server containers
			curServNames, err := listServerContainers()
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Randomly permute the list of server containers
			permuteSlice(curServNames)

			// Calculate the extra servers needed
			extraServ := payloadData.N - len(payloadData.Hostnames)

			// Remove the extra servers
			for i := 0; i < extraServ; i++ {
				err := killServerContainer(curServNames[i])
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}

			// Get the updated list of server containers
			servNames, err := listServerContainers()
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			// Prepare and send JSON response
			servData := ServDetails{
				N:        c.Nserv,
				Replicas: servNames,
			}
			resp := ResponseSuccess{
				Message: servData,
				Status:  "successful",
			}
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rw.Write(jsonResp)
		} else {
			// If the number of servers requested is less than the removed servers, return an error response
			resp := Response{
				Message: "ERROR: Length of hostname list is more than newly added instances",
				Status:  "failure",
			}
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(jsonResp)
		}
	default:
		// Handle unsupported methods
		rw.WriteHeader(http.StatusNotFound)
	}
}

// Function to get a server name based on consistent hashing
func GetServerName() string {
	rand.NewSource(time.Now().UnixNano())
	id := rand.Intn(Mod)
	servName := c.GetServer(id)
	return servName
}

// Function to perform server heartbeat and return a reachable server URL
func serverHeartbeat() (string, error) {
	rand.NewSource(time.Now().UnixNano())
	max_tries := 10000

	// Attempt to find a reachable server within a limit
	for max_tries != 0 {
		mtx.Lock()
		servName := GetServerName()
		url := "http://" + servName + ":5000"
		servResp, err := http.Get(url + "/heartbeat")
		if err == nil && servResp.StatusCode == http.StatusOK {
			mtx.Unlock()
			return url, nil
		}
		// Remove the inactive server and add a new one
		res := c.RemoveServer(servName)
		if res == 0 {
			mtx.Unlock()
			return "", errors.New("Inactive server deletion failed")
		}
		// Retry until a new server is successfully added
		for {
			num := rand.Intn(Mod)
			name := GenerateRandomString(num)
			if _, ok := c.AllServers[name]; ok {
				continue
			}
			err = addServerContainer(name, num)
			if err != nil {
				mtx.Unlock()
				return "", errors.New("New server creation failed")
			}
			break
		}
		mtx.Unlock()
		max_tries--
	}
	return "", errors.New("Server unavailable")
}

// Handler for the default endpoint "/<path>" and other paths
func path(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if req.RequestURI != "/home" && req.RequestURI != "/heartbeat" {
			// Return an error for unsupported endpoints
			resp := Response{
				Message: "ERROR: '/other' endpoint does not exist in server replicas",
				Status:  "failure",
			}
			jsonResp, err := json.Marshal(resp)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(jsonResp)
		} else {
			// Perform server heartbeat and route the request to a reachable server
			url, err := serverHeartbeat()
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.WriteHeader(http.StatusOK)
            if req.RequestURI == "/home" {
				// Forward the request to the chosen server
				servResp, err := http.Get(url + req.RequestURI)
				if err != nil || servResp.StatusCode != http.StatusOK {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				// Decode and forward the response
				var resp Response
                err = json.NewDecoder(servResp.Body).Decode(&resp)
                if err != nil {
                    fmt.Println("Error:", err)
                    rw.WriteHeader(http.StatusInternalServerError)
                    return
                }
                jsonResp, err := json.Marshal(resp)
                if err != nil {
                    fmt.Println("Error:", err)
                    rw.WriteHeader(http.StatusInternalServerError)
                    return
                }
                rw.Header().Set("Content-Type", "application/json")
                rw.Write(jsonResp)
			}
		}
	default:
		// Handle unsupported methods
		rw.WriteHeader(http.StatusNotFound)
	}
}

// Function to list all server containers
func listServerContainers() ([]string, error) {
	// Docker API endpoint

    endpoint := "unix:///var/run/docker.sock"
    client, err := docker.NewClient(endpoint)
    if err != nil {
        return []string{}, err
    }

	// Get the list of containers
    containers, err := client.ListContainers(docker.ListContainersOptions{All: false})
    if err != nil {
        return []string{}, err
    }
		
	// Extract hostnames of server containers in the "net1" network
	hostnames := []string{}
	for _, container := range containers {
		containerInfo, _ := client.InspectContainer(container.ID)

		for network := range containerInfo.NetworkSettings.Networks {
			if network == "net1" {
				for _, name := range container.Names {
					cleanName := strings.TrimPrefix(name, "/")
					if cleanName != "lb" {
						hostnames = append(hostnames, cleanName)
					}
				} 
			}
		}
	}
	fmt.Println("Hostnames: ", hostnames)
	return hostnames, nil
}
 	
// Function to add a new server container
func addServerContainer(serverName string, serverNumber int) error {
	// Add the server to the consistent hash ring
			
	res := c.AddServer(serverNumber, serverName)
	if res == 0 {
		return errors.New("Server already exists")
	}

	// Docker API endpoint
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
    if err != nil {
		// If adding server to the hash ring failed, remove it and return an error
		c.RemoveServer(serverName)
		return err
	}

	// Create Docker container options
	createContainerOptions := docker.CreateContainerOptions{
	    Name: serverName,
        Config: &docker.Config{
            Image: "server",
			// Assuming "server" is the Docker image for your server
            Env: []string{"SERVER_NUMBER=" + strconv.Itoa(serverNumber)},
        },
        HostConfig: &docker.HostConfig{
            AutoRemove: true,
            // Tty:        true,
            // OpenStdin:  true,
            NetworkMode: "net1",
        },
    }
	// Create the Docker container
	container, err := client.CreateContainer(createContainerOptions)
	if err != nil {
		fmt.Println("Container could not be created\n", err)
		// If container creation fails, remove the server from the hash ring and return an error
		c.RemoveServer(serverName)
		return err
	}

	// Start the Docker container
	err = client.StartContainer(container.ID, nil)
	if err != nil {
		fmt.Println("Container could not be started\n", err)
		// If starting the container fails, remove the server from the hash ring and return an error
		c.RemoveServer(serverName)
		return err
	}

	// Allow some time for the container to start
	time.Sleep(1*time.Second)
	return nil
}

// Function to kill an existing server container
func killServerContainer(serverName string) error {

	// Docker API endpoint
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return err
	}

	// Kill the Docker container
	killOptions := docker.KillContainerOptions{ID: serverName}
	err = client.KillContainer(killOptions)
	if err != nil {
		return err
	}

	// Remove the server from the hash ring
	res := c.RemoveServer(serverName)
	if res == 0 {
		// If server not found in the hash ring, return an error
		return errors.New("Server not found")
	}

	// Allow some time for the container to stop
	time.Sleep(1*time.Second)
	return nil
}