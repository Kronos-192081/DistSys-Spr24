package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"math/rand"
	"time"
	"errors"
	"prakhar/conhash"
	// "github.com/fsouza/go-dockerclient"
)

const Mod = 1e4 + 7
var c = conhash.NewConHash(512, 9)
var servNamePort = make(map[string]string)	// To be removed

func main() {
	fmt.Println("Starting load balancer")

	rand.NewSource(time.Now().UnixNano())

	// To be removed
	ids := []int{rand.Intn(Mod), rand.Intn(Mod), rand.Intn(Mod)}
	servNames := []string{"Server 1", "Server 2", "Server 3"}
	c.Add(ids, servNames)
	servNamePort["Server 1"] = "5001"
	servNamePort["Server 2"] = "5002"
	servNamePort["Server 3"] = "5003"

	http.HandleFunc("/rep", rep)
	repSrv := &http.Server{Addr: "127.0.0.1:5000"}

	http.HandleFunc("/add", add)
	addSrv := &http.Server{Addr: "127.0.0.1:5000"}

	http.HandleFunc("/rm", rm)
	rmSrv := &http.Server{Addr: "127.0.0.1:5000"}

	http.HandleFunc("/", path)
	pathSrv := &http.Server{Addr: "127.0.0.1:5000"}

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

func GenerateRandomString(length int) string {
	rand.NewSource(time.Now().UnixNano())

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := 0; i < length; i++ {
		result[i] = charset[rand.Intn(len(charset))]
	}

	return string(result)
}

// Fisher-Yates algorithm for random permutation
func permuteSlice(slice []string) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// Handler functions for incoming requests

func rep(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		servNames := []string{}
		for servName := range c.AllServers {
			servNames = append(servNames, servName)
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
		ids := []int{}
		for i := 0; i < len(payloadData.Hostnames); i++ {
			ids = append(ids, rand.Intn(Mod))
		}
		res := c.Add(ids, payloadData.Hostnames)
		if res == 0 {
			fmt.Println("Error:", "Server creation failed")
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if payloadData.N >= len(payloadData.Hostnames) {
			extraServ := payloadData.N - len(payloadData.Hostnames)
			for i := 0; i < extraServ; i++ {
				res := c.AddServer(rand.Intn(Mod), GenerateRandomString(10))
				if res == 0 {
					fmt.Println("Error:", "Server creation failed")
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			servNames := []string{}
			for servName := range c.AllServers {
				servNames = append(servNames, servName)
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
				Message: "<Error> Length of hostname list is more than newly added instances",
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
			res := c.RemoveServer(servName)
			if res == 0 {
				fmt.Println("Error:", "Server deletion failed")
				rw.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		if payloadData.N >= len(payloadData.Hostnames) {
			curServNames := []string{}
			for servName := range c.AllServers {
				curServNames = append(curServNames, servName)
			}
			permuteSlice(curServNames)
			extraServ := payloadData.N - len(payloadData.Hostnames)
			for i := 0; i < extraServ; i++ {
				res := c.RemoveServer(curServNames[i])
				if res == 0 {
					fmt.Println("Error:", "Server deletion failed")
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
			servNames := []string{}
			for servName := range c.AllServers {
				servNames = append(servNames, servName)
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
				Message: "<Error> Length of hostname list is more than newly added instances",
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
	id := rand.Intn(Mod)
	servName := c.GetServer(id)
	return servName
}

func serverHeartbeat() (string, error) {
	max_tries := 10000
	for max_tries != 0 {
		servName := GetServerName()
		servPort := servNamePort[servName]	// needs modification
		url := "http://127.0.0.1:" + servPort	// needs modification
		servResp, err := http.Get(url + "/heartbeat")
		if err == nil && servResp.StatusCode == http.StatusOK {
			return url, nil
		}
		res := c.RemoveServer(servName)
		if res == 0 {
			return "", errors.New("Inactive server deletion failed")
		}
		res = c.AddServer(rand.Intn(Mod), GenerateRandomString(10))
		if res == 0 {
			return "", errors.New("New server creation failed")
		}
		max_tries--
	}
	return "", errors.New("Server unavailable")
}

func path(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		if req.RequestURI != "/home" && req.RequestURI != "/heartbeat" {
			resp := Response{
				Message: "<Error> '/other' endpoint does not exist in server replicas",
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
