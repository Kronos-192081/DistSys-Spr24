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

const Mod = 1e4 + 7
var c = conhash.NewConHash(512, 9) // Needs change
var mtx sync.Mutex

func main() {
	fmt.Println("Starting load balancer")
	// pulling server image
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil{
		fmt.Println("Client creation failed", err)
	}

	err = pullImage(client, "alutnopk/go-http-server:latest")
	if err != nil {
		fmt.Println("Could not Pull Image", err)
	}

	rand.NewSource(time.Now().UnixNano())

	ids := []int{rand.Intn(Mod), rand.Intn(Mod), rand.Intn(Mod)}
	servNames := []string{"Server_1", "Server_2", "Server_3"}
	addServerContainer(servNames[0], ids[0])
	addServerContainer(servNames[1], ids[1])
	addServerContainer(servNames[2], ids[2])

	listServerContainers()
	c.GetConfig()

	listServerContainers()

	http.HandleFunc("/rep", rep)
	repSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/add", add)
	addSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/rm", rm)
	rmSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/", path)
	pathSrv := &http.Server{Addr: "0.0.0.0:5000"}

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

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

	defer func() {
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

	sig := <-sigs
	fmt.Println(sig)

	cancel()

	fmt.Println("Shutting down load balancer")
}

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

// Utility functions

func GenerateRandomString(num int) string {

	// rand.NewSource(time.Now().UnixNano())

	// const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// result := make([]byte, length)

	// for i := 0; i < length; i++ {
	// 	result[i] = charset[rand.Intn(len(charset))]
	// }

	name := "spawned_server_"+strconv.Itoa(num)

	return name
}

// Fisher-Yates algorithm for random permutation
func permuteSlice(slice []string) {
	rand.NewSource(time.Now().UnixNano())
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Handler functions for incoming requests

func rep(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		servNames, err := listServerContainers()
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
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
		rw.WriteHeader(http.StatusNotFound)
	}
}

func add(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
        var payloadData Payload
        err := json.NewDecoder(req.Body).Decode(&payloadData)
        if err != nil {
            fmt.Println("Error:", err)
            rw.WriteHeader(http.StatusInternalServerError)
            return
        }
		
		rand.NewSource(time.Now().UnixNano())

		for i := 0; i < len(payloadData.Hostnames); i++ {
			err := addServerContainer(payloadData.Hostnames[i], rand.Intn(Mod))
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		if payloadData.N >= len(payloadData.Hostnames) {
			extraServ := payloadData.N - len(payloadData.Hostnames)
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

			servNames, err := listServerContainers()
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
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
		rw.WriteHeader(http.StatusNotFound)
	}
}

func rm(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodDelete:
        var payloadData Payload
        err := json.NewDecoder(req.Body).Decode(&payloadData)
        if err != nil {
            fmt.Println("Error:", err)
            rw.WriteHeader(http.StatusInternalServerError)
            return
        }
		for _, servName := range payloadData.Hostnames {
			err := killServerContainer(servName)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if payloadData.N >= len(payloadData.Hostnames) {
			curServNames, err := listServerContainers()
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			permuteSlice(curServNames)
			extraServ := payloadData.N - len(payloadData.Hostnames)
			for i := 0; i < extraServ; i++ {
				err := killServerContainer(curServNames[i])
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			servNames, err := listServerContainers()
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
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
		rw.WriteHeader(http.StatusNotFound)
	}
}

func GetServerName() string {
	rand.NewSource(time.Now().UnixNano())
	id := rand.Intn(Mod)
	servName := c.GetServer(id)
	return servName
}

func serverHeartbeat() (string, error) {
	rand.NewSource(time.Now().UnixNano())
	max_tries := 10000
	for max_tries != 0 {
		mtx.Lock()
		servName := GetServerName()
		url := "http://" + servName + ":5000"
		servResp, err := http.Get(url + "/heartbeat")
		if err == nil && servResp.StatusCode == http.StatusOK {
			mtx.Unlock()
			return url, nil
		}
		res := c.RemoveServer(servName)
		if res == 0 {
			mtx.Unlock()
			return "", errors.New("Inactive server deletion failed")
		}
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

func path(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if req.RequestURI != "/home" && req.RequestURI != "/heartbeat" {
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
			url, err := serverHeartbeat()
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			rw.WriteHeader(http.StatusOK)
            if req.RequestURI == "/home" {
				servResp, err := http.Get(url + req.RequestURI)
				if err != nil || servResp.StatusCode != http.StatusOK {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
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
		rw.WriteHeader(http.StatusNotFound)
	}
}

func listServerContainers() ([]string, error) {

    endpoint := "unix:///var/run/docker.sock"
    client, err := docker.NewClient(endpoint)
    if err != nil {
        return []string{}, err
    }

    containers, err := client.ListContainers(docker.ListContainersOptions{All: false})
    if err != nil {
        return []string{}, err
    }
    // currentHostname, err := os.Hostname()
    // if err != nil {
	//     return []string{}, err
	// }
	// fmt.Println(currentHostname)
		
	hostnames := []string{}
	for _, container := range containers {
		for network := range container.NetworkSettings.Networks {
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
		
func addServerContainer(serverName string, serverNumber int) error {
			
	res := c.AddServer(serverNumber, serverName)
	if res == 0 {
		return errors.New("Server already exists")
	}
			
	endpoint := "unix:///var/run/docker.sock"
    client, err := docker.NewClient(endpoint)
    if err != nil {
		c.RemoveServer(serverName)
		return err
    }

	createContainerOptions := docker.CreateContainerOptions{
        Name: serverName,
        Config: &docker.Config{
            Image: "alutnopk/go-http-server",
            Env: []string{"SERVER_NUMBER=" + strconv.Itoa(serverNumber)},
        },
        HostConfig: &docker.HostConfig{
            AutoRemove: true,
            // Tty:        true,
            // OpenStdin:  true,
            NetworkMode: "net1",
        },
    }
    container, err := client.CreateContainer(createContainerOptions)
    if err != nil {
		fmt.Println("Container could not be created\n", err)
		c.RemoveServer(serverName)
        return err
    }

    err = client.StartContainer(container.ID, nil)
    if err != nil {
		fmt.Println("Container could not be started\n", err)
		c.RemoveServer(serverName)
        return err
    }
	// TODO: figure out callback
	time.Sleep(1*time.Second)
	return nil
}

func killServerContainer(serverName string) error {

	endpoint := "unix:///var/run/docker.sock"
    client, err := docker.NewClient(endpoint)
    if err != nil {
        return err
    }

	killOptions := docker.KillContainerOptions{ID: serverName}
    err = client.KillContainer(killOptions)
    if err != nil {
        return err
    }

	res := c.RemoveServer(serverName)
	if res == 0 {
		return errors.New("Server not found")
	}
	time.Sleep(1*time.Second)
    return nil
}

func pullImage(client *docker.Client, imageName string) error {
	fmt.Println("Pulling image", imageName)
    pullOptions := docker.PullImageOptions{
        Repository: imageName,
    }
    authConfiguration := docker.AuthConfiguration{}
    err := client.PullImage(pullOptions, authConfiguration)
    if err != nil {
		fmt.Println("Could not pull image from Docker Hub\n")
        return err
    }
	fmt.Println("Image pulled\n")
    return nil
}