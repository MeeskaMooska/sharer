package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

const (
	numUsers        = 20
	numItems        = 30
	numTransactions = 50
)

var schools = []string{
	"State University", "City College", "Tech Institute",
	"Riverside University", "Northern College", "Eastside Community College",
}

var categories = []string{
	"textbooks", "electronics", "clothing", "furniture", "sports", "music", "other",
}

func main() {
	_ = godotenv.Load("../../.env")

	if os.Getenv("APP_ENV") == "production" {
		log.Fatal("refusing to seed in production")
	}

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

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("could not connect to database:", err)
	}

	truncate(db)
	userIDs := seedUsers(db)
	itemIDs := seedItems(db)
	seedTransactions(db, userIDs, itemIDs)

	log.Println("seed complete")
}

func truncate(db *sql.DB) {
	for _, table := range []string{"transactions", "items", "users"} {
		if _, err := db.Exec("DELETE FROM " + table); err != nil {
			log.Fatalf("truncate %s: %v", table, err)
		}
		log.Printf("cleared %s", table)
	}
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
		INSERT INTO users (email, school, hashed_password, goodwill_points)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			school          = VALUES(school),
			hashed_password = VALUES(hashed_password),
			goodwill_points = VALUES(goodwill_points)`

	hashedPw := hash("password123")
	var ids []int64

	for i := 0; i < numUsers; i++ {
		email := gofakeit.Email()
		school := schools[gofakeit.Number(0, len(schools)-1)]
		goodwill := gofakeit.Number(0, 50)

		res, err := db.Exec(q, email, school, hashedPw, goodwill)
		if err != nil {
			log.Fatalf("insert user: %v", err)
		}
		id, _ := res.LastInsertId()
		ids = append(ids, id)
		log.Printf("inserted user %s", email)
	}
	return ids
}

func seedItems(db *sql.DB) []int64 {
	const q = `
		INSERT INTO items (name, description, value, category)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			description = VALUES(description),
			value       = VALUES(value),
			category    = VALUES(category)`

	var ids []int64

	for i := 0; i < numItems; i++ {
		name := fmt.Sprintf("%s %s", gofakeit.AdjectiveDescriptive(), gofakeit.NounAbstract())
		description := gofakeit.Sentence(8)
		value := gofakeit.Price(1, 200)
		category := categories[gofakeit.Number(0, len(categories)-1)]

		res, err := db.Exec(q, name, description, value, category)
		if err != nil {
			log.Fatalf("insert item: %v", err)
		}
		id, _ := res.LastInsertId()
		ids = append(ids, id)
		log.Printf("inserted item %q ($%.2f)", name, value)
	}
	return ids
}

func seedTransactions(db *sql.DB, userIDs, itemIDs []int64) {
	const q = `
		INSERT INTO transactions (user_giving, user_receiving, item_id, reviewed, review)
		VALUES (?, ?, ?, ?, ?)`

	for i := 0; i < numTransactions; i++ {
		giver := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		receiver := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		for receiver == giver {
			receiver = userIDs[gofakeit.Number(0, len(userIDs)-1)]
		}
		item := itemIDs[gofakeit.Number(0, len(itemIDs)-1)]

		reviewed := gofakeit.Number(0, 1)
		var review interface{}
		if reviewed == 1 {
			review = gofakeit.Number(0, 1)
		}

		if _, err := db.Exec(q, giver, receiver, item, reviewed, review); err != nil {
			log.Fatalf("insert transaction: %v", err)
		}
		log.Printf("inserted transaction: user %d -> user %d, item %d", giver, receiver, item)
	}
}
