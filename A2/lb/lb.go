package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"strconv"
	"regexp"
	"strings"
	"bytes"
	"errors"
	"prakhar/conhash"
	"github.com/fsouza/go-dockerclient"
)

// Structs for representing JSON responses and payloads
type Response struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

type schema struct {
	Columns []string `json:"columns"`
	Dtypes  []string `json:"dtypes"`
}

type shard struct {
	Stud_id_low int
	Shard_id    string
	Shard_size  int
}

type configPayload struct {
	Schema schema  	`json:"schema"`
	Shards []string `json:"shards"`
}

type initPayload struct {
	N       int
	Schema  schema              `json:"schema"`
	Shards  []shard             `json:"shards"`
	Servers map[string][]string `json:"servers"`
}

type addPayload struct {
	N          int                 `json:"n"`
	New_shards []shard             `json:"new_shards"`
	Servers    map[string][]string `json:"servers"`
}

type rmPayload struct {
	N          int                 `json:"n"`
	Servers    []string 		   `json:"servers"`
}

type editResponse struct {
	N 	  	int
	Message string `json:"message"`
	Status  string `json:"status"`
}

type Range struct {
	Low 	int `json:"low"`
	High 	int `json:"high"`
}

type readPayload struct {
	Stud_id Range
}

type data struct {
	Stud_id  	int
	Stud_name 	string
	Stud_marks 	string
}

type readResponse struct {
	Shards_queried 	[]string 	`json:"shards_queried"`
	Data 		 	[]data 		`json:"data"`
	Status 		 	string 		`json:"status"`
}

type readServPayload struct {
	Shard 		string
	Stud_id 	Range
}

type readServResponse struct {
	Data 	[]data 	`json:"data"`
	Status 	string 	`json:"status"`
}

type writePayload struct {
	Data 	[]data 	`json:"data"`
}

type writeServPayload struct {
	Shard 		string	`json:"shard"`
	Curr_idx 	int		`json:"curr_idx"`
	Data 		[]data 	`json:"data"`
}

type writeServResponse struct {
	Message 	string 	`json:"message"`
	Curr_idx 	int		`json:"curr_idx"`
	Status 		string 	`json:"status"`
}

type updatePayload struct {
	Stud_id int
	Data 	data 	`json:"data"`
}

type updateServPayload struct {
	Shard 		string	`json:"shard"`
	Stud_id 	int		`json:"stud_id"`
	Data 		data 	`json:"data"`
}

type delPayload struct {
	Stud_id int
}

type delServPayload struct {
	Shard 		string	`json:"shard"`
	Stud_id 	int		`json:"stud_id"`
}

// Environment variables for configuring the load balancer
var num_slots, _ = strconv.Atoi(os.Getenv("NUM_SLOTS"))	// 512
var num_virt_serv, _ = strconv.Atoi(os.Getenv("NUM_VIRT_SERV"))	// log(512) = 9

// Constants
const Mod = 1e4 + 7

// Global consistent hash instance
var ConHashList = make(map[string](*(conhash.ConHash)))

// Server List for storing names of all active server containers
var ServerList = make(map[string]bool)

var serv_schema schema

// Mutex for thread-safe operations
var mtx sync.Mutex

// Main function
func main() {
	fmt.Println("Starting load balancer")

	// Seed for randomization
	rand.NewSource(time.Now().UnixNano())

	// Setup HTTP servers for different endpoints
	http.HandleFunc("/init", init_)
	initSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/status", status)
	statusSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/add", add)
	addSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/rm", rm)
	rmSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/read", read)
	readSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/write", write)
	writeSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/update", update)
	updateSrv := &http.Server{Addr: "0.0.0.0:5000"}

	http.HandleFunc("/del", del)
	delSrv := &http.Server{Addr: "0.0.0.0:5000"}

	// Setup context and signal handling
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	// Start HTTP servers in separate goroutines
	go func() {
		initSrv.ListenAndServe()
	}()

	go func() {
		statusSrv.ListenAndServe()
	}()

	go func() {
		addSrv.ListenAndServe()
	}()

	go func() {
		rmSrv.ListenAndServe()
	}()

	go func() {
		readSrv.ListenAndServe()
	}()

	go func() {
		writeSrv.ListenAndServe()
	}()

	go func() {
		updateSrv.ListenAndServe()
	}()

	go func() {
		delSrv.ListenAndServe()
	}()

	// Defer shutdown of servers
	defer func() {
		// Graceful shutdown of servers
		if err := initSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the init server: ", err)
		}
		if err := statusSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the status server: ", err)
		}
		if err := addSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the add server: ", err)
		}
		if err := rmSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the rm server: ", err)
		}
		if err := readSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the read server: ", err)
		}
		if err := writeSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the write server: ", err)
		}
		if err := updateSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the update server: ", err)
		}
		if err := delSrv.Shutdown(ctx); err != nil {
			fmt.Println("error when shutting down the del server: ", err)
		}
	}()

	// Wait for SIGINT signal
	sig := <-sigs
	fmt.Println(sig)

	// Cancel the context to initiate shutdown
	cancel()

	fmt.Println("Shutting down load balancer")
}

// Handler functions for incoming requests

// Handler for /init endpoint (POST)
func init_(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		var payloadData initPayload
		err := json.NewDecoder(req.Body).Decode(&payloadData)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		serv_schema = payloadData.Schema

		for _, shard := range payloadData.Shards {
			ConHashList[shard.Shard_id] = conhash.NewConHash(num_slots, num_virt_serv)
			// add to db
		}

		for servName, shards := range payloadData.Servers {
			configServData := configPayload{
				Schema: serv_schema,
				Shards: shards,
			}

			jsonBody, err := json.Marshal(configServData)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			url := "http://" + servName + ":5000"
			servResp, err := http.Post(url + "/config", "application/json", bytes.NewReader(jsonBody))
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if servResp.StatusCode != http.StatusOK {
				fmt.Println("Error: Server failed")
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			for _, shard := range shards {
				ConHashList[shard].AddServer(rand.Intn(Mod), servName)
			}
		}

		resp := Response{
			Message: "Configured Database",
			Status:  "success",
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

// Handler for /status endpoint (GET)
func status(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:

		// fetch data

		resp := initPayload{
			N:      3,
			Schema: schema{Columns: []string{"Stud_id", "Stud_name", "Stud_marks"}, Dtypes: []string{"Number", "String", "String"}},
			Shards: []shard{{Stud_id_low: 0, Shard_id: "sh1", Shard_size: 4096},
				{Stud_id_low: 4096, Shard_id: "sh2", Shard_size: 4096},
				{Stud_id_low: 8192, Shard_id: "sh3", Shard_size: 4096}},
			Servers: map[string][]string{"Server0": {"sh1", "sh2"}, "Server1": {"sh2", "sh3"}, "Server2": {"sh1", "sh3"}},
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

// Handler for /add endpoint (POST)
func add(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// Decode the JSON payload
        var payloadData addPayload
        err := json.NewDecoder(req.Body).Decode(&payloadData)
        if err != nil {
            fmt.Println("Error:", err)
            rw.WriteHeader(http.StatusInternalServerError)
            return
        }
		
		rand.NewSource(time.Now().UnixNano())

		// Check if the number of servers requested is greater than the added servers
		if payloadData.N == len(payloadData.Servers) {
			// Add server containers based on the payload
			server_names := []string{}

			for _, shard := range payloadData.New_shards {
				ConHashList[shard.Shard_id] = conhash.NewConHash(num_slots, num_virt_serv)
			}

			for k, v := range payloadData.Servers {
				var servName string
				if match, _ := regexp.MatchString("Server\\[[0-9]+\\]", k); match {
					for {
						num := rand.Intn(Mod)
						name := GenerateRandomString(num)
						if _, ok := ServerList[name]; ok {
							continue
						}
						err = addServerContainer(name, num)
						if err != nil {
							fmt.Println("Error:", err)
							rw.WriteHeader(http.StatusInternalServerError)
							return
						}
						ServerList[name] = true
						server_names = append(server_names, name)
						servName = name

						break
					}
				} else {
					err := addServerContainer(k, rand.Intn(Mod))
					if err != nil {
						fmt.Println("Error:", err)
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}
					ServerList[k] = true
					server_names = append(server_names, k)
					servName = k
				}

				configServData := configPayload{
					Schema: serv_schema,
					Shards: v,
				}
	
				jsonBody, err := json.Marshal(configServData)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				url := "http://" + servName + ":5000"
				servResp, err := http.Post(url + "/config", "application/json", bytes.NewReader(jsonBody))
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				if servResp.StatusCode != http.StatusOK {
					fmt.Println("Error: Server failed")
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				for _, shard := range v {
					// copy if shard exists
					ConHashList[shard].AddServer(rand.Intn(Mod), servName)
				}
			}

			// Prepare and send JSON response
			resp := editResponse{
				N:	   	len(ServerList),
				Message: "Added " + strings.Join(server_names, ", "),
				Status:  "success",
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
				Message: "ERROR: Number of new servers (n) is greater than newly added instances",
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
        var payloadData rmPayload
        err := json.NewDecoder(req.Body).Decode(&payloadData)
        if err != nil {
            fmt.Println("Error:", err)
            rw.WriteHeader(http.StatusInternalServerError)
            return
        }

		// Check if the number of servers requested is greater than the removed servers
		if payloadData.N >= len(payloadData.Servers) {
			server_names := []string{}
			// Remove specified server containers
			for _, servName := range payloadData.Servers {
				err := killServerContainer(servName)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				delete(ServerList, servName)
				server_names = append(server_names, servName)
			}

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
			extraServ := payloadData.N - len(payloadData.Servers)

			// Remove the extra servers
			for i := 0; i < extraServ; i++ {
				err := killServerContainer(curServNames[i])
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				delete(ServerList, curServNames[i])
				server_names = append(server_names, curServNames[i])
			}

			// Prepare and send JSON response
			resp := editResponse{
				N:	   	len(ServerList),
				Message: "Removed " + strings.Join(server_names, ", "),
				Status:  "success",
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
				Message: "ERROR: Length of server list is more than removable instances",
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

// Handler for /read endpoint (POST)
func read(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// Decode the JSON payload
		var payloadData readPayload
		err := json.NewDecoder(req.Body).Decode(&payloadData)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Range parsing and get shard id list
		shard_ids := []string{}
		data_entries := []data{}
		for _, shard := range shard_ids {
			//get serv name from shard
			//get partial range
			servName := ""
			readServData := readServPayload{
				Shard: shard,
				Stud_id: Range{Low: -1, High: -1},
			}

			jsonBody, err := json.Marshal(readServData)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			url := "http://" + servName + ":5000"
			servResp, err := http.Post(url + "/read", "application/json", bytes.NewReader(jsonBody))
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if servResp.StatusCode != http.StatusOK {
				fmt.Println("Error: Server failed")
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			var readServResp readServResponse
			err = json.NewDecoder(servResp.Body).Decode(&readServResp)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			
			data_entries = append(data_entries, readServResp.Data...)
		}

		// Prepare and send JSON response
		resp := readResponse{
			Shards_queried: shard_ids,
			Data: 			data_entries,
			Status:  		"success",
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

// Handler for /write endpoint (POST)
func write(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		// Decode the JSON payload
		var payloadData writePayload
		err := json.NewDecoder(req.Body).Decode(&payloadData)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Range parsing and get shard id list and bucket out the data entries
		shard_ids := []string{}
		data_entries := [][]data{}
		for i, shard := range shard_ids {
			//get serv name from shard
			//get curr_idx
			servName := ""
			curr_idx := -1
			writeServData := writeServPayload{
				Shard: shard,
				Curr_idx: curr_idx,
				Data: data_entries[i],
			}

			jsonBody, err := json.Marshal(writeServData)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			url := "http://" + servName + ":5000"
			servResp, err := http.Post(url + "/write", "application/json", bytes.NewReader(jsonBody))
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if servResp.StatusCode != http.StatusOK {
				fmt.Println("Error: Server failed")
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			var writeServResp writeServResponse
			err = json.NewDecoder(servResp.Body).Decode(&writeServResp)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			//update curr_idx
			//write fault-tolerance handling - some shards write committed, some not
		}

		// Prepare and send JSON response
		resp := Response{
			Message: 		strconv.Itoa(len(payloadData.Data)) + " Data entries added",
			Status:  		"success",
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

// Handler for /update endpoint (PUT)
func update(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		// Decode the JSON payload
		var payloadData updatePayload
		err := json.NewDecoder(req.Body).Decode(&payloadData)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		//get shard id
		shard_id := ""
		//get serv name from shard
		servName := ""
		updateServData := updateServPayload{
			Shard: 		shard_id,
			Stud_id: 	payloadData.Stud_id,
			Data: 		payloadData.Data,
		}

		jsonBody, err := json.Marshal(updateServData)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := "http://" + servName + ":5000"
		servResp, err := http.Post(url + "/update", "application/json", bytes.NewReader(jsonBody))
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if servResp.StatusCode != http.StatusOK {
			fmt.Println("Error: Server failed")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var updateServResp Response
		err = json.NewDecoder(servResp.Body).Decode(&updateServResp)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Prepare and send JSON response
		resp := Response{
			Message: 		updateServResp.Message,
			Status:  		"success",
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

// Handler for /del endpoint (DELETE)
func del(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodDelete:
		// Decode the JSON payload
		var payloadData delPayload
		err := json.NewDecoder(req.Body).Decode(&payloadData)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		//get shard id
		shard_id := ""
		//get serv name from shard
		servName := ""
		delServData := delServPayload{
			Shard: 		shard_id,
			Stud_id: 	payloadData.Stud_id,
		}

		jsonBody, err := json.Marshal(delServData)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := "http://" + servName + ":5000"
		servResp, err := http.Post(url + "/del", "application/json", bytes.NewReader(jsonBody))
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if servResp.StatusCode != http.StatusOK {
			fmt.Println("Error: Server failed")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var delServResp Response
		err = json.NewDecoder(servResp.Body).Decode(&delServResp)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Prepare and send JSON response
		resp := Response{
			Message: 		delServResp.Message + " from all replicas",
			Status:  		"success",
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
		// Retry until a new server is successly added
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