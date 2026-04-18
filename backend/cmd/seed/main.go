package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	dbpkg "nova_hackathon_2026/db"
)

const (
	numUsers        = 20
	numItems        = 30
	numTransactions = 50
)

var emojis = []string{
	"🐶", "🐱", "🐭", "🐹", "🐰", "🦊", "🐻", "🐼",
	"🐨", "🐯", "🦁", "🐮", "🐸", "🐙", "🦋", "🐺",
	"🦄", "🐲", "🌵", "🍄",
}


func main() {
	_ = godotenv.Load("../../.env")

	if os.Getenv("APP_ENV") == "production" {
		log.Fatal("refusing to seed in production")
	}

	database, err := dbpkg.Init("app.db")
	if err != nil {
		log.Fatal("could not init database:", err)
	}
	defer database.Close()

	truncate(database)
	userIDs := seedUsers(database)
	itemIDs := seedItems(database, userIDs)
	seedTransactions(database, userIDs, itemIDs)

	log.Println("seed complete")
}

func truncate(db *sql.DB) {
	drops := []string{"transactions", "items", "users"}
	for _, table := range drops {
		if _, err := db.Exec("DROP TABLE IF EXISTS " + table); err != nil {
			log.Fatalf("drop %s: %v", table, err)
		}
		log.Printf("dropped %s", table)
	}
	if err := dbpkg.Migrate(db); err != nil {
		log.Fatal("re-migrate:", err)
	}
	log.Println("tables recreated")
}

func hash(password string) string {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}

func seedUsers(db *sql.DB) []int64 {
	const q = `
		INSERT OR REPLACE INTO users (username, email, profile_picture)
		VALUES (?, ?, ?)`

	var ids []int64
	for i := 0; i < numUsers; i++ {
		username := gofakeit.Username()
		email := gofakeit.Email()
		emoji := emojis[gofakeit.Number(0, len(emojis)-1)]

		res, err := db.Exec(q, username, email, emoji)
		if err != nil {
			log.Fatalf("insert user: %v", err)
		}
		id, _ := res.LastInsertId()
		ids = append(ids, id)
		log.Printf("inserted user %s %s", emoji, username)
	}
	return ids
}

func seedItems(db *sql.DB, userIDs []int64) []int64 {
	const q = `
		INSERT OR REPLACE INTO items (name, user_id, description, price)
		VALUES (?, ?, ?, ?)`

	var ids []int64
	for i := 0; i < numItems; i++ {
		name := fmt.Sprintf("%s %s", gofakeit.AdjectiveDescriptive(), gofakeit.NounAbstract())
		userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		description := gofakeit.Sentence(8)
		price := gofakeit.Price(1, 200)

		res, err := db.Exec(q, name, userID, description, price)
		if err != nil {
			log.Fatalf("insert item: %v", err)
		}
		id, _ := res.LastInsertId()
		ids = append(ids, id)
		log.Printf("inserted item %q ($%.2f)", name, price)
	}
	return ids
}

func seedTransactions(db *sql.DB, userIDs, itemIDs []int64) {
	const q = `
		INSERT INTO transactions (user_id, item_id, quantity, total_price)
		VALUES (?, ?, ?, ?)`

	for i := 0; i < numTransactions; i++ {
		userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		itemID := itemIDs[gofakeit.Number(0, len(itemIDs)-1)]
		quantity := gofakeit.Number(1, 5)

		var price float64
		db.QueryRow(`SELECT price FROM items WHERE id = ?`, itemID).Scan(&price)
		total := price * float64(quantity)

		if _, err := db.Exec(q, userID, itemID, quantity, total); err != nil {
			log.Fatalf("insert transaction: %v", err)
		}
		log.Printf("inserted transaction: user %d, item %d, qty %d", userID, itemID, quantity)
	}
}
