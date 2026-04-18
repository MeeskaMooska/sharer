package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

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
	seedUsers(db)
	seedItems(db)
	seedTransactions(db)

	log.Println("seed complete")
}

func truncate(db *sql.DB) {
	// FK order: transactions -> items -> users
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

func seedUsers(db *sql.DB) {
	users := []struct {
		email    string
		school   string
		password string
		goodwill int
	}{
		{"alice@university.edu", "State University", "password123", 10},
		{"bob@college.edu", "City College", "password123", 5},
	}

	const q = `
		INSERT INTO users (email, school, hashed_password, goodwill_points)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			school          = VALUES(school),
			hashed_password = VALUES(hashed_password),
			goodwill_points = VALUES(goodwill_points)`

	for _, u := range users {
		if _, err := db.Exec(q, u.email, u.school, hash(u.password), u.goodwill); err != nil {
			log.Fatalf("insert user %s: %v", u.email, err)
		}
		log.Printf("upserted user %s", u.email)
	}
}

func seedItems(db *sql.DB) {
	const q = `
		INSERT INTO items (name, description, value, category)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			description = VALUES(description),
			value       = VALUES(value),
			category    = VALUES(category)`

	if _, err := db.Exec(q, "Calculus Textbook", "8th edition, minor highlights", 35.00, "textbooks"); err != nil {
		log.Fatal("insert item:", err)
	}
	log.Println("upserted item: Calculus Textbook")
}

func seedTransactions(db *sql.DB) {
	var giverID, receiverID, itemID int64

	if err := db.QueryRow(`SELECT id FROM users WHERE email = ?`, "alice@university.edu").Scan(&giverID); err != nil {
		log.Fatal("lookup alice:", err)
	}
	if err := db.QueryRow(`SELECT id FROM users WHERE email = ?`, "bob@college.edu").Scan(&receiverID); err != nil {
		log.Fatal("lookup bob:", err)
	}
	if err := db.QueryRow(`SELECT id FROM items WHERE name = ?`, "Calculus Textbook").Scan(&itemID); err != nil {
		log.Fatal("lookup item:", err)
	}

	// reviewed=1, review=1 (as advertised)
	const q = `
		INSERT INTO transactions (user_giving, user_receiving, item_id, reviewed, review)
		VALUES (?, ?, ?, 1, 1)`

	if _, err := db.Exec(q, giverID, receiverID, itemID); err != nil {
		log.Fatal("insert transaction:", err)
	}
	log.Printf("inserted transaction: user %d -> user %d, item %d", giverID, receiverID, itemID)
}
