package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	serverNumber := os.Getenv("SERVER_NUMBER")
	response := map[string]string{
		"message": "Hello from Server: " + serverNumber,
		"status":  "successful",
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	var config dbConfig
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received config:", config) // debug

	response := map[string]string{
		"message": "Config received and processed successfully",
		"status":  "success",
	}

	err_resp := map[string]string{
		"message": "Could not process the config",
		"status":  "error",
	}

	var dbStatus bool = dbSetup(config)
	if dbStatus {
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	} else {
		jsonResponse, err := json.Marshal(err_resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
	}
}

var db *sql.DB

func main() {

	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "20CS30061"
		dbname   = "testdb"
	)

	// Construct connection string
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// Connect to database
	d, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Connected to the database")
	}

	db = d

	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/heartbeat", heartbeatHandler)
	http.HandleFunc("/config", configHandler)

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "5000"
	}

	fmt.Println("Server listening on port:", httpPort)
	er := http.ListenAndServe(":"+httpPort, nil)
	if er != nil {
		fmt.Println("Error starting server:", err)
	}
}

type dbSchema struct {
	Columns []string `json:"columns"`
	Dtypes  []string `json:"dtypes"`
}

type dbConfig struct {
	Schema dbSchema `json:"schema"`
	Shards []string `json:"shards"`
}

func dbSetup(config dbConfig) bool {
	var q1 string = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT PRIMARY KEY, %s VARCHAR(255), %s INT);", config.Shards[0], config.Schema.Columns[0], config.Schema.Columns[1], config.Schema.Columns[2])
	var q2 string = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT PRIMARY KEY, %s VARCHAR(255), %s INT);", config.Shards[1], config.Schema.Columns[0], config.Schema.Columns[1], config.Schema.Columns[2])
	_, err := db.Exec(q1)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return false
	}

	_, _err := db.Exec(q2)
	if _err != nil {
		fmt.Println("Error creating table:", _err)
		return false
	}

	// var q3 string = "SELECT * FROM sh1;"
	// rows, err := db.Query(q3)
	// if err != nil {
	// 	fmt.Println("Error querying table:", err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var Stud_id int
	// 	var Stud_name string
	// 	var Stud_marks int
	// 	err = rows.Scan(&Stud_id, &Stud_name, &Stud_marks)
	// 	if err != nil {
	// 		fmt.Println("Error scanning rows:", err)
	// 	}
	// 	fmt.Println(Stud_id, Stud_name, Stud_marks)
	// }
	return true
}
