package main

import (
	"fmt"
	"log"
	"net/http"
	"nova_hackathon_2026/handlers"
	"os"

	"github.com/joho/godotenv"
	"nova_hackathon_2026/db"
	"nova_hackathon_2026/handlers"
)

func main() {
	_ = godotenv.Load("../.env")

	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	if user == "" {
		user = "root"
	}
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "3306"
	}
	if name == "" {
		name = "nova"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)
	database, err := db.Init(dsn)
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
			h.AddItem(w, r) // The new function we just wrote
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 6. Start the server
	port := ":8080"
	fmt.Printf("API running...\n")
	fmt.Printf("-> Test DB at: http://localhost%s/test-db\n", port)
	fmt.Printf("-> Sign-up ready at: POST http://localhost%s/api/sign-up\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
