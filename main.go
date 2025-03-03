package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

// Initialize database connection as a global variable
var db *sql.DB

// Initialization function
func init() {
	var err error
	// Open a SQLite database file (or create one if it doesn't exist)
	db, err = sql.Open("sqlite3", "./webhook_logs.db")
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	// Create a table for storing the webhook logs if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		source TEXT,
		event TEXT,
		amount REAL,
		body TEXT
	)`)
	if err != nil {
		log.Fatal("Error creating table: ", err)
	}
}

// Main entry point
func main() {
	// Log that the server is starting
	fmt.Println("Server is starting...")

	// Set up routes and handlers
	r := mux.NewRouter()

	// Handle the incoming webhook POST request
	r.HandleFunc("/webhook", WebhookHandler).Methods("POST")

	// Endpoint to fetch logs from the database
	r.HandleFunc("/logs", GetLogsHandler).Methods("GET")

	// Allow all origins for testing purposes
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (for testing)
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(r)

	// Start the server and log any errors that occur
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

// Webhook handler function
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received a webhook request") // Log that we received the request

	// Ensure the content type is application/json
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Invalid content type, expected application/json", http.StatusBadRequest)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Log the body for debugging
	fmt.Printf("Received Body: %s\n", body)

	// Try parsing the body as a JSON object
	var jsonBody map[string]interface{}
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Log the parsed JSON body
	fmt.Printf("Parsed JSON Body: %#v\n", jsonBody)

	// Hardcoded values for simplicity (adjust this to your needs)
	event := "payment_received"
	amount := 100.00
	source := "GitHub"

	// Insert the log into the database
	_, err = db.Exec("INSERT INTO logs (source, event, amount, body) VALUES (?, ?, ?, ?)", source, event, amount, string(body))
	if err != nil {
		fmt.Printf("Error inserting log into database: %v\n", err)  // Log the error message
		http.Error(w, "Failed to store log", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	fmt.Fprintf(w, "Webhook received successfully!")
}

// Handler function to retrieve all logs from the database
func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, source, event, amount, body FROM logs")
	if err != nil {
		http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var logs []map[string]interface{}
	for rows.Next() {
		var id int
		var source, event, body string
		var amount float64
		if err := rows.Scan(&id, &source, &event, &amount, &body); err != nil {
			http.Error(w, "Failed to read log", http.StatusInternalServerError)
			return
		}
		logData := map[string]interface{}{
			"id":     id,
			"source": source,
			"event":  event,
			"amount": amount,
			"body":   body,
		}
		logs = append(logs, logData)
	}

	// Send the logs as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
