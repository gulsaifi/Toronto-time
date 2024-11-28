package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Database connection string
const dbDSN = "root:1234@tcp(127.0.0.1:3306)/time_api"

// Struct for JSON response
type TimeResponse struct {
	CurrentTime string `json:"current_time"`
}

func getCurrentTime(w http.ResponseWriter, r *http.Request) {
	// Connect to database
	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		log.Fatal(err)
	}
	defer db.Close()

	// Get Toronto time
	loc, err := time.LoadLocation("America/Toronto")
	if err != nil {
		http.Error(w, "Failed to load timezone", http.StatusInternalServerError)
		log.Fatal(err)
	}
	currentTime := time.Now().In(loc)

	// Insert time into database
	_, err = db.Exec("INSERT INTO time_log (timestamp) VALUES (?)", currentTime)
	if err != nil {
		http.Error(w, "Failed to log time", http.StatusInternalServerError)
		log.Fatal(err)
	}

	// Respond with JSON
	response := TimeResponse{
		CurrentTime: currentTime.Format("2006-01-02 15:04:05"),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Set up router
	r := mux.NewRouter()
	r.HandleFunc("/current-time", getCurrentTime).Methods("GET")

	// Start server
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
