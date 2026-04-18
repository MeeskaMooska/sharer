package main

import (
	"fmt"
	"log"
	"net/http"
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
		log.Fatal(err)
	}
	defer database.Close()

	h := handlers.New(database)

	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("../frontend")))

	mux.HandleFunc("/api/users", h.Users)
	mux.HandleFunc("/api/items", h.Items)
	mux.HandleFunc("/api/transactions", h.Transactions)

	log.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
