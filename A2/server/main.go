package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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

	serverNo := os.Getenv("SERVER_NUMBER")
	msgStr := ""
	for i := 0; i < len(config.Shards)-1; i++ {
		msgStr += "Server" + serverNo + ":" + config.Shards[i] + ", "
	}
	msgStr += "Server" + serverNo + ":" + config.Shards[len(config.Shards)-1] + " configured"

	response := map[string]string{
		"message": msgStr,
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

func copyHandler(w http.ResponseWriter, r *http.Request) {
	var copyReq copyRequest
	err := json.NewDecoder(r.Body).Decode(&copyReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received copy request:", copyReq) // debug

	var err_resp map[string]string

	var final_result = make(map[string][]Row)

	for i := 0; i < len(copyReq.Shards); i++ {
		q := fmt.Sprintf("SELECT * FROM %s;", copyReq.Shards[i])
		rows, err := db.Query(q)
		if err != nil {
			fmt.Println("Error querying the database:", err)
			err_resp = map[string]string{
				"message": "Error querying the database",
				"status":  "error",
			}
			jsonResponse, err := json.Marshal(err_resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonResponse)
			return
		}

		result := make([]Row, 0)
		for rows.Next() {
			var row Row
			err = rows.Scan(&row.Stud_id, &row.Stud_name, &row.Stud_marks)
			if err != nil {
				fmt.Println("Error scanning the rows:", err)
				err_resp = map[string]string{
					"message": "Error scanning the rows",
					"status":  "error",
				}
				jsonResponse, err := json.Marshal(err_resp)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonResponse)
				return
			}
			result = append(result, row)
		}

		fmt.Println("Result:", result) // debug
		final_result[copyReq.Shards[i]] = result
	}

	response := make(map[string][]Row)

	for i := 0; i < len(copyReq.Shards); i++ {
		response[copyReq.Shards[i]] = final_result[copyReq.Shards[i]]
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

func readHandler(w http.ResponseWriter, r *http.Request) {
	var readReq readRequest
	err := json.NewDecoder(r.Body).Decode(&readReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received read request:", readReq) // debug

	var err_resp map[string]string

	q := fmt.Sprintf("SELECT * FROM %s WHERE %s BETWEEN %d AND %d;", readReq.Shard, "stud_id", readReq.Stud_id.Low, readReq.Stud_id.High)

	rows, err := db.Query(q)
	if err != nil {
		fmt.Println("Error querying the database:", err)
		err_resp = map[string]string{
			"message": "Error querying the database",
			"status":  "error",
		}
		jsonResponse, err := json.Marshal(err_resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	result := make([]Row, 0)
	for rows.Next() {
		var row Row
		err = rows.Scan(&row.Stud_id, &row.Stud_name, &row.Stud_marks)
		if err != nil {
			fmt.Println("Error scanning the rows:", err)
			err_resp = map[string]string{
				"message": "Error scanning the rows",
				"status":  "error",
			}
			jsonResponse, err := json.Marshal(err_resp)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(jsonResponse)
			return
		}
		result = append(result, row)
	}

	response := map[string]interface{}{
		"data":   result,
		"status": "success",
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

func writeHandler(w http.ResponseWriter, r *http.Request) {
	var writeReq writeRequest
	err := json.NewDecoder(r.Body).Decode(&writeReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received write request:", writeReq) // debug

	var err_resp map[string]string

	var values []string
	for _, data := range writeReq.Data {
		value := fmt.Sprintf("(%d, '%s', %d)", data.Stud_id, data.Stud_name, data.Stud_marks)
		values = append(values, value)
	}

	q := fmt.Sprintf("INSERT INTO %s VALUES %s;", writeReq.Shard, strings.Join(values, ","))
	_, err = db.Exec(q)
	if err != nil {
		fmt.Println("Error inserting into the database:", err)
		err_resp = map[string]string{
			"message": "Error inserting into the database",
			"status":  "error",
		}
		jsonResponse, err := json.Marshal(err_resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	response := writeServResponse{
		Message:  "Data entries added",
		Curr_idx: writeReq.Curr_idx + len(writeReq.Data),
		Status:   "success",
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

func updateHandler(w http.ResponseWriter, r *http.Request) {
	var updateReq updateRequest
	err := json.NewDecoder(r.Body).Decode(&updateReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received update request:", updateReq) // debug

	var err_resp map[string]string

	q := fmt.Sprintf("UPDATE %s SET %s = %d, %s = '%s', %s = %d WHERE %s = %d;", updateReq.Shard, "stud_id", updateReq.Data.Stud_id, "stud_name", updateReq.Data.Stud_name, "stud_marks", updateReq.Data.Stud_marks, "stud_id", updateReq.Stud_id)

	res, err := db.Exec(q)
	if err != nil {
		fmt.Println("Error updating the database:", err)
		err_resp = map[string]string{
			"message": "Error updating the database",
			"status":  "error",
		}
		jsonResponse, err := json.Marshal(err_resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		return
	}

	if count == 0 {
		err_resp = map[string]string{
			"message": fmt.Sprintf("No data entry found for Stud_id:%d", updateReq.Stud_id),
			"status":  "error",
		}
		jsonResponse, err := json.Marshal(err_resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResponse)
		return
	}

	response := map[string]string{
		"message": fmt.Sprintf("Data entry for Stud_id:%d updated", updateReq.Stud_id),
		"status":  "success",
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

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var deleteReq deleteRequest
	err := json.NewDecoder(r.Body).Decode(&deleteReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Received delete request:", deleteReq) // debug

	var err_resp map[string]string

	q := fmt.Sprintf("DELETE FROM %s WHERE %s = %d;", deleteReq.Shard, "stud_id", deleteReq.Stud_id)

	res, err := db.Exec(q)
	if err != nil {
		fmt.Println("Error deleting the database:", err)
		err_resp = map[string]string{
			"message": "Error deleting the database",
			"status":  "error",
		}
		jsonResponse, err := json.Marshal(err_resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	count, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error getting rows affected:", err)
		return
	}

	if count == 0 {
		err_resp = map[string]string{
			"message": fmt.Sprintf("No data entry found for Stud_id:%d", deleteReq.Stud_id),
			"status":  "error",
		}
		jsonResponse, err := json.Marshal(err_resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResponse)
		return
	}

	response := map[string]string{
		"message": fmt.Sprintf("Data entry for Stud_id:%d deleted", deleteReq.Stud_id),
		"status":  "success",
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
	http.HandleFunc("/read", readHandler)
	http.HandleFunc("/copy", copyHandler)
	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/del", deleteHandler)

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

type Row struct {
	Stud_id    int    `json:"Stud_id"`
	Stud_name  string `json:"Stud_name"`
	Stud_marks int    `json:"Stud_marks"`
}

type dbSchema struct {
	Columns []string `json:"columns"`
	Dtypes  []string `json:"dtypes"`
}

type dbConfig struct {
	Schema dbSchema `json:"schema"`
	Shards []string `json:"shards"`
}

type rangeID struct {
	Low  int
	High int
}

type copyRequest struct {
	Shards []string `json:"shards"`
}

type readRequest struct {
	Shard   string  `json:"shard"`
	Stud_id rangeID `json:"stud_id"`
}

type writeRequest struct {
	Shard    string `json:"shard"`
	Curr_idx int    `json:"curr_idx"`
	Data     []Row  `json:"data"`
}

type updateRequest struct {
	Shard   string `json:"shard"`
	Stud_id int    `json:"stud_id"`
	Data    Row    `json:"data"`
}

type deleteRequest struct {
	Shard   string `json:"shard"`
	Stud_id int    `json:"stud_id"`
}

type writeServResponse struct {
	Message  string `json:"message"`
	Curr_idx int    `json:"curr_idx"`
	Status   string `json:"status"`
}

func dbSetup(config dbConfig) bool {
	for i := 0; i < len(config.Shards); i++ {
		var q string = fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT PRIMARY KEY, %s VARCHAR(255), %s INT);", config.Shards[i], config.Schema.Columns[0], config.Schema.Columns[1], config.Schema.Columns[2])
		_, err := db.Exec(q)
		if err != nil {
			fmt.Println("Error creating table:", err)
			return false
		}
	}
	return true
}
