package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"math/rand"
	"math"
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
	"databases/sql"
	_ "github.com/mattn/go-sqlite3"
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

type copyPayload struct {
	Shards 		[]string 	`json:"shards"`
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

var db *sql.DB

// Mutex for thread-safe operations
var mtx sync.Mutex

// Main function
func main() {
	fmt.Println("Starting load balancer")

	db, err := sql.Open("sqlite3", "lb.db")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS ShardT (Stud_id_low INTEGER PRIMARY KEY, Shard_id TEXT, Shard_size INTEGER, valid_idx INTEGER)")
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS MapT (Shard_id TEXT, Server_id TEXT)")
	if err != nil {
		fmt.Println(err)
		return
	}

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
			_, err := db.Exec("INSERT INTO ShardT (Stud_id_low, Shard_id, Shard_size, valid_idx) VALUES (?, ?, ?, ?)", shard.Stud_id_low, shard.Shard_id, shard.Shard_size, 0)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			ConHashList[shard.Shard_id] = conhash.NewConHash(num_slots, num_virt_serv)
		}

		for servName, shards := range payloadData.Servers {
			err = addServerContainer(servName, rand.Intn(Mod), shards)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			ServerList[servName] = true

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

			// heartbeat required
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
				_, err = db.Exec("INSERT INTO MapT (Shard_id, Server_id) VALUES (?, ?)", shard, servName)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
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

		rows, err := db.Query("SELECT Stud_id_low, Shard_id, Shard_size FROM ShardT")
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		shard_list := []shard{}
		for rows.Next() {
			var shard shard
			err = rows.Scan(&shard.Stud_id_low, &shard.Shard_id, &shard.Shard_size)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			shard_list = append(shard_list, shard)
		}

		server_dict := map[string][]string{}
		for servName, _ := range ServerList {
			rows, err := db.Query("SELECT Shard_id FROM MapT WHERE Server_id = ?", servName)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			shard_list := []string{}
			for rows.Next() {
				var shard_id string
				err = rows.Scan(&shard_id)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				shard_list = append(shard_list, shard_id)
			}
			server_dict[servName] = shard_list
		}

		resp := initPayload{
			N:      len(ServerList),
			Schema: serv_schema,
			Shards: shard_list,
			Servers: server_dict,
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

		// Check if the number of servers requested is equal to the added servers
		if payloadData.N == len(payloadData.Servers) {
			// Add server containers based on the payload
			server_names := []string{}

			for _, shard := range payloadData.New_shards {
				_, err := db.Exec("INSERT INTO ShardT (Stud_id_low, Shard_id, Shard_size, valid_idx) VALUES (?, ?, ?, ?)", shard.Stud_id_low, shard.Shard_id, shard.Shard_size, 0)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
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
						err = addServerContainer(name, num, v)
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
					err := addServerContainer(k, rand.Intn(Mod), v)
					if err != nil {
						fmt.Println("Error:", err)
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}
					ServerList[k] = true
					server_names = append(server_names, k)
					servName = k
				}

				for _, shard_id := range v {
					_, err = db.Exec("INSERT INTO MapT (Shard_id, Server_id) VALUES (?, ?)", shard_id, servName)
					if err != nil {
						fmt.Println("Error:", err)
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}
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

				// heartbeat required
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
					exists := false
					for _, new_shard := range payloadData.New_shards {
						if shard == new_shard.Shard_id {
							exists = true
							break
						}
					}
					if exists {
						continue
					}

					// copy from old data
					copyServData := copyPayload{
						Shards: []string{shard},
					}

					jsonBody, err := json.Marshal(copyServData)
					if err != nil {
						fmt.Println("Error:", err)
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}

					// heartbeat required
					oldServName := GetServerName(shard)
					url := "http://" + oldServName + ":5000"
					servResp, err := http.Post(url + "/copy", "application/json", bytes.NewReader(jsonBody))
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

					var shard_data map[string]interface{}
					err = json.NewDecoder(servResp.Body).Decode(&shard_data)
					if err != nil {
						fmt.Println("Error:", err)
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}

					if _, ok := shard_data["status"]; !ok {
						fmt.Println("Error: Server failed")
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}

					if shard_data["status"] != "success" {
						fmt.Println("Error: Server failed")
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}

					for k, v := range shard_data {
						if k == "status" {
							continue
						}

						writeServData := writeServPayload{
							Shard: k,
							Curr_idx: 0,
							Data: v.([]data),
						}

						jsonBody, err := json.Marshal(writeServData)
						if err != nil {
							fmt.Println("Error:", err)
							rw.WriteHeader(http.StatusInternalServerError)
							return
						}

						// heartbeat required maybe
						url = "http://" + servName + ":5000"
						servResp, err = http.Post(url + "/write", "application/json", bytes.NewReader(jsonBody))
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
						// write failure handling
					}

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
				rows, err := db.Query("SELECT Shard_id FROM MapT WHERE Server_id = ?", servName)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				rmShardList := []string{}
				for rows.Next() {
					var shard_id string
					err = rows.Scan(&shard_id)
					if err != nil {
						fmt.Println("Error:", err)
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}
					rmShardList = append(rmShardList, shard_id)
				}

				err = killServerContainer(servName, rmShardList)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				_, err = db.Exec("DELETE FROM MapT WHERE Server_id = ?", servName)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				delete(ServerList, servName)
				server_names = append(server_names, servName)
			}

			// Get the current list of server containers
			curServNames := []string{}
			for k := range ServerList {
				curServNames = append(curServNames, k)
			}

			// Randomly permute the list of server containers
			permuteSlice(curServNames)

			// Calculate the extra servers needed
			extraServ := payloadData.N - len(payloadData.Servers)

			// Remove the extra servers
			for i := 0; i < extraServ; i++ {
				rows, err := db.Query("SELECT Shard_id FROM MapT WHERE Server_id = ?", curServNames[i])
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				rmShardList := []string{}
				for rows.Next() {
					var shard_id string
					err = rows.Scan(&shard_id)
					if err != nil {
						fmt.Println("Error:", err)
						rw.WriteHeader(http.StatusInternalServerError)
						return
					}
					rmShardList = append(rmShardList, shard_id)
				}

				err = killServerContainer(curServNames[i], rmShardList)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				_, err = db.Exec("DELETE FROM MapT WHERE Server_id = ?", curServNames[i])
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

		// Range parsing and obtaining shard id list
		rows, err := db.Query("SELECT Shard_id, Stud_id_low, Shard_size FROM ShardT WHERE Stud_id_low >= ? AND Stud_id_low <= ?", payloadData.Stud_id.Low, payloadData.Stud_id.High)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		shard_ids := []string{}
		shard_stud_id_low := []int{}
		shard_stud_id_size := []int{}
		for rows.Next() {
			var shard_id string
			err = rows.Scan(&shard_id, &shard_stud_id_low, &shard_stud_id_size)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			shard_ids = append(shard_ids, shard_id)
		}

		data_entries := []data{}
		for i, shard := range shard_ids {
			// heartbeat required
			servName := GetServerName(shard)
			
			readServData := readServPayload{
				Shard: shard,
				Stud_id: Range{Low: max(shard_stud_id_low[i], payloadData.Stud_id.Low), High: min(shard_stud_id_low[i] + shard_stud_id_size[i] - 1, payloadData.Stud_id.High)},
			}

			jsonBody, err := json.Marshal(readServData)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			// heartbeat required maybe
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

		// Range parsing and get shard id list and bucket the data entries
		data_entries := map[string][]data{}
		valid_idx_dict := map[string]int{}
		for _, data_entry := range payloadData.Data {
			rows, err := db.Query("SELECT Shard_id, valid_idx FROM ShardT WHERE Stud_id_low <= ? AND Stud_id_low + Shard_size > ?", data_entry.Stud_id, data_entry.Stud_id)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !rows.Next() {
				fmt.Println("Error: Invalid data entry found")
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			for rows.Next() {
				var shard_id string
				var valid_idx int
				err = rows.Scan(&shard_id, &valid_idx)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				data_entries[shard_id] = append(data_entries[shard_id], data_entry)
				valid_idx_dict[shard_id] = valid_idx
			}
		}

		for shard, data_list := range data_entries {

			servNameList := []string{}
			rows, err := db.Query("SELECT Server_id FROM MapT WHERE Shard_id = ?", shard)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}

			for rows.Next() {
				var servName string
				err = rows.Scan(&servName)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				servNameList = append(servNameList, servName)
			}

			failed := false
			curr_idx := -1
			for _, servName := range servNameList {

				writeServData := writeServPayload{
					Shard: shard,
					Curr_idx: valid_idx_dict[shard],
					Data: data_list,
				}

				jsonBody, err := json.Marshal(writeServData)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				
				// heartbeat required
				url := "http://" + servName + ":5000"
				servResp, err := http.Post(url + "/write", "application/json", bytes.NewReader(jsonBody))
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				if servResp.StatusCode != http.StatusOK {
					fmt.Println("Error: Server failed")
					failed = true
					break
				}

				var writeServResp writeServResponse
				err = json.NewDecoder(servResp.Body).Decode(&writeServResp)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}

				if curr_idx != -1 && curr_idx != writeServResp.Curr_idx {
					failed = true
					break
				}
				curr_idx = writeServResp.Curr_idx
			}
			
			if !failed && curr_idx != -1 {
				_, err = db.Exec("UPDATE ShardT SET valid_idx = ? WHERE Shard_id = ?", curr_idx, shard)
				if err != nil {
					fmt.Println("Error:", err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {

			}
			//fault-tolerance handling - some shards write committed, some not, kill and respawn done servers
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

		var shard_id string
		rows, err := db.Query("SELECT Shard_id FROM ShardT WHERE Stud_id_low <= ? AND Stud_id_low + Shard_size > ?", payloadData.Stud_id, payloadData.Stud_id)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !rows.Next() {
			fmt.Println("Error: Invalid data entry found")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			err = rows.Scan(&shard_id)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		servNameList := []string{}
		rows, err = db.Query("SELECT Server_id FROM MapT WHERE Shard_id = ?", shard_id)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var servName string
			err = rows.Scan(&servName)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			servNameList = append(servNameList, servName)
		}

		failed := false
		response := ""
		for _, servName := range servNameList {
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
	
			// heartbeat required
			url := "http://" + servName + ":5000"
			servResp, err := http.Post(url + "/update", "application/json", bytes.NewReader(jsonBody))
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if servResp.StatusCode != http.StatusOK {
				fmt.Println("Error: Server failed")
				failed = true
				return
			}
	
			var updateServResp Response
			err = json.NewDecoder(servResp.Body).Decode(&updateServResp)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			
			if response != "" && response != updateServResp.Message {
				failed = true
				break
			}
			response = updateServResp.Message
		}
		
		if failed {

		}
		//fault-tolerance handling - some shards update committed, some not, kill and respawn done servers

		// Prepare and send JSON response
		resp := Response{
			Message: 		response,
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

		var shard_id string
		rows, err := db.Query("SELECT Shard_id FROM ShardT WHERE Stud_id_low <= ? AND Stud_id_low + Shard_size > ?", payloadData.Stud_id, payloadData.Stud_id)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !rows.Next() {
			fmt.Println("Error: Invalid data entry found")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			err = rows.Scan(&shard_id)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		servNameList := []string{}
		rows, err = db.Query("SELECT Server_id FROM MapT WHERE Shard_id = ?", shard_id)
		if err != nil {
			fmt.Println("Error:", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		for rows.Next() {
			var servName string
			err = rows.Scan(&servName)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			servNameList = append(servNameList, servName)
		}

		failed := false
		response := ""
		for _, servName := range servNameList {
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
	
			// heartbeat required
			url := "http://" + servName + ":5000"
			servResp, err := http.Post(url + "/del", "application/json", bytes.NewReader(jsonBody))
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			if servResp.StatusCode != http.StatusOK {
				fmt.Println("Error: Server failed")
				failed = true
				return
			}
	
			var delServResp Response
			err = json.NewDecoder(servResp.Body).Decode(&delServResp)
			if err != nil {
				fmt.Println("Error:", err)
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
			
			if response != "" && response != delServResp.Message {
				failed = true
				break
			}
			response = delServResp.Message
		}
		
		if failed {

		}
		//fault-tolerance handling - some shards del committed, some not, kill and respawn done servers

		// Prepare and send JSON response
		resp := Response{
			Message: 		response + " from all replicas",
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
func GetServerName(shard string) string {
	rand.NewSource(time.Now().UnixNano())
	id := rand.Intn(Mod)
	servName := ConHashList[shard].GetServer(id)
	return servName
}

// Function to perform server heartbeat and return a reachable server URL
func serverHeartbeat(shard_id string) (string, error) {
	rand.NewSource(time.Now().UnixNano())
	max_tries := 10000

	// Attempt to find a reachable server within a limit
	for max_tries != 0 {
		mtx.Lock()
		servName := GetServerName(shard_id)
		url := "http://" + servName + ":5000"
		servResp, err := http.Get(url + "/heartbeat")
		if err == nil && servResp.StatusCode == http.StatusOK {
			mtx.Unlock()
			return url, nil
		}

		delete(ServerList, servName)

		shard_list := []string{}
		rows, err := db.Query("SELECT Shard_id FROM MapT WHERE Server_id = ?", servName)
		if err != nil {
			mtx.Unlock()
			return "", err
		}

		for rows.Next() {
			var shard string
			err = rows.Scan(&shard)
			if err != nil {
				mtx.Unlock()
				return "", err
			}
			shard_list = append(shard_list, shard)
		}

		for _, shard := range shard_list {
			res := ConHashList[shard].RemoveServer(servName)
			if res == 0 {
				mtx.Unlock()
				return "", errors.New("Inactive server deletion failed")
			}
		}
		
		// Retry until a new server is successly added
		for {
			num := rand.Intn(Mod)
			name := GenerateRandomString(num)
			if _, ok := ServerList[name]; ok {
				continue
			}
			err = addServerContainer(name, num, shard_list)
			if err != nil {
				mtx.Unlock()
				return "", err
			}

			// new server updated
			for _, shard := range shard_list {

				// copy from old data
				copyServData := copyPayload{
					Shards: []string{shard},
				}

				jsonBody, err := json.Marshal(copyServData)
				if err != nil {
					mtx.Unlock()
					return "", err
				}

				// heartbeat required
				oldServName := GetServerName(shard)
				url := "http://" + oldServName + ":5000"
				servResp, err := http.Post(url + "/copy", "application/json", bytes.NewReader(jsonBody))
				if err != nil {
					mtx.Unlock()
					return "", err
				}
				if servResp.StatusCode != http.StatusOK {
					mtx.Unlock()
					return "", errors.New("Server failed")
				}

				var shard_data map[string]interface{}
				err = json.NewDecoder(servResp.Body).Decode(&shard_data)
				if err != nil {
					mtx.Unlock()
					return "", err
				}

				if _, ok := shard_data["status"]; !ok {
					mtx.Unlock()
					return "", errors.New("Server failed")
				}

				if shard_data["status"] != "success" {
					mtx.Unlock()
					return "", errors.New("Server failed")
				}

				for k, v := range shard_data {
					if k == "status" {
						continue
					}

					writeServData := writeServPayload{
						Shard: k,
						Curr_idx: 0,
						Data: v.([]data),
					}

					jsonBody, err := json.Marshal(writeServData)
					if err != nil {
						mtx.Unlock()
						return "", err
					}

					// heartbeat required maybe
					url = "http://" + name + ":5000"
					servResp, err = http.Post(url + "/write", "application/json", bytes.NewReader(jsonBody))
					if err != nil {
						mtx.Unlock()
						return "", err
					}
					if servResp.StatusCode != http.StatusOK {
						mtx.Unlock()
						return "", errors.New("Server failed")
					}
					// write failure handling
				}
			}

			_, err = db.Exec("UPDATE MapT SET Server_id = ? WHERE Server_id = ?", name, servName)
			if err != nil {
				mtx.Unlock()
				return "", err
			}
			ServerList[name] = true

			break
		}

		mtx.Unlock()
		max_tries--
	}
	return "", errors.New("Server unavailable")
}
	
// Function to add a new server container
func addServerContainer(serverName string, serverNumber int, shards []string) error {
	// Add the server to the consistent hash ring
	for _, shard := range shards {
		ConHashList[shard].AddServer(serverNumber, serverName)
	}

	// Docker API endpoint
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
    if err != nil {
		// If adding server to the hash ring failed, remove it and return an error
		for _, shard := range shards {
			ConHashList[shard].RemoveServer(serverName)
		}
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
		for _, shard := range shards {
			ConHashList[shard].RemoveServer(serverName)
		}
		return err
	}

	// Start the Docker container
	err = client.StartContainer(container.ID, nil)
	if err != nil {
		fmt.Println("Container could not be started\n", err)
		// If starting the container fails, remove the server from the hash ring and return an error
		for _, shard := range shards {
			ConHashList[shard].RemoveServer(serverName)
		}
		return err
	}

	// Allow some time for the container to start
	time.Sleep(1*time.Second)
	return nil
}

// Function to kill an existing server container
func killServerContainer(serverName string, shards []string) error {

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
	for _, shard := range shards {
		ConHashList[shard].RemoveServer(serverName)
	}

	// Allow some time for the container to stop
	time.Sleep(1*time.Second)
	return nil
}

/***

TO-DO:

1. Logical verification of code
2. Error handling
3. valid_idx update in /del ?
4. /write, /update and /del fault-tolerance
5. /heartbeat before every server access
6. Testing

***/