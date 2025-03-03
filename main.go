package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize a new Mux router
	r := mux.NewRouter()

	// Define the route for the webhook endpoint and bind it to the WebhookHandler function
	r.HandleFunc("/webhook", WebhookHandpackage main

	import (
		"database/sql"
		"encoding/json"
		"fmt"
		"io/ioutil"
		"log"
		"net/http"
	
		_ "github.com/mattn/go-sqlite3"
		"github.com/gorilla/mux"
	)
	
	var db *sql.DB
	
	func init() {
		// Open a SQLite database file (or create one if it doesn't exist)
		var err error
		db, err = sql.Open("sqlite3", "./webhook_logs.db")
		if err != nil {
			log.Fatal(err)
		}
	
		// Create a table for storing the webhook logs if it doesn't exist
		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			event TEXT,
			amount REAL,
			body TEXT
		)`)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	func main() {
		r := mux.NewRouter()
	
		// Handle the incoming webhook POST request
		r.HandleFunc("/webhook", WebhookHandler).Methods("POST")
	
		// Endpoint to fetch logs from the database
		r.HandleFunc("/logs", GetLogsHandler).Methods("GET")
	
		// Start the server
		fmt.Println("Magic is Happening on port 8080...")
		log.Fatal(http.ListenAndServe(":8080", r))
	}
	
	// WebhookHandler handles POST requests to /webhook
	func WebhookHandler(w http.ResponseWriter, r *http.Request) {
		// Read the incoming request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
	
		// For simplicity, we'll use hardcoded values for the event and amount
		// You can modify this to extract values dynamically from the body
		event := "payment_received"
		amount := 100.00
	
		// Insert the log into the database
		_, err = db.Exec("INSERT INTO logs (event, amount, body) VALUES (?, ?, ?)", event, amount, string(body))
		if err != nil {
			http.Error(w, "Failed to store log", http.StatusInternalServerError)
			return
		}
	
		// Respond to acknowledge receipt of the webhook
		fmt.Fprintf(w, "Webhook received!")
	}
	
	// GetLogsHandler fetches all the logs from the database and sends them as JSON
	func GetLogsHandler(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, event, amount, body FROM logs")
		if err != nil {
			http.Error(w, "Failed to retrieve logs", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
	
		var logs []map[string]interface{}
		for rows.Next() {
			var id int
			var event string
			var amount float64
			var body string
			if err := rows.Scan(&id, &event, &amount, &body); err != nil {
				http.Error(w, "Failed to read log", http.StatusInternalServerError)
				return
			}
			logData := map[string]interface{}{
				"id":     id,
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
	er).Methods("POST")

	// Start the server
	fmt.Println("Magic is Happening on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// WebhookHandler handles POST requests to /webhook
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	// Log the HTTP method, headers, and body content
	fmt.Println("HTTP Method:", r.Method)
	fmt.Println("Headers:", r.Header)

	// Read the body of the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Log the request body
	fmt.Println("Body:", string(body))

	// Respond to acknowledge the request was received
	fmt.Fprintf(w, "Webhook received!")
}
