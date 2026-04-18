package main

import (
	"log"
	"net/http"
	"os"

	"nova_hackathon_2026/db"
	"nova_hackathon_2026/handlers"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/nova?parseTime=true"
	}
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
