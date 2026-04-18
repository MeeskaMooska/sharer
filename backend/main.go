package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql" // Blank import required for the driver
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load the .env file from two directories up
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

	// 4. Define the single endpoint to test connectivity
	http.HandleFunc("/test-db", func(w http.ResponseWriter, r *http.Request) {
		// Ping actively verifies the connection is still alive
		err := db.Ping()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Database connection failed: %v\n", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully connected to the MySQL database at %s:%s!\n", dbHost, dbPort)
	})

	// 5. Start the server
	port := ":8080"
	fmt.Printf("API running...\nTest your DB connection at: http://localhost%s/test-db\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
