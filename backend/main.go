package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"nova_hackathon_2026/handlers"
	"os"

	_ "github.com/go-sql-driver/mysql" // Blank import required for the driver
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load the .env file from one directory up
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// 2. Build the connection string (DSN)
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Print DB connection info
	fmt.Printf("Connecting to DB: user=%s, host=%s, port=%s, name=%s\n", dbUser, dbHost, dbPort, dbName)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)

	// 3. Initialize the database connection pool
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}
	defer db.Close()

	// Best practice: configure connection limits
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// 4. Initialize Handlers
	// Pass the active DB connection pool to your handlers struct
	h := handlers.New(db)

	// 5. Define API routes

	// Keep the test route for debugging
	http.HandleFunc("/test-db", func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Database connection failed: %v\n", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully connected to the MySQL database at %s:%s!\n", dbHost, dbPort)
	})

	// Register the new Sign-Up endpoint
	http.HandleFunc("/api/sign-up", h.SignUp)
	http.HandleFunc("/api/sign-in", h.SignIn)
	http.HandleFunc("/api/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetItems(w, r) // Your pagination function
		case http.MethodPost:
			h.AddItem(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/transactions/request", h.RequestItem)
	http.HandleFunc("/api/transactions/respond", h.RespondToRequest)

	// 6. Start the server
	port := ":8080"
	fmt.Printf("API running...\n")
	fmt.Printf("-> Test DB at: http://localhost%s/test-db\n", port)
	fmt.Printf("-> Sign-up ready at: POST http://localhost%s/api/sign-up\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
