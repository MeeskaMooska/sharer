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

var eduDomains = []string{
	"mit.edu", "stanford.edu", "harvard.edu", "unc.edu",
	"gatech.edu", "umich.edu", "berkeley.edu", "nyu.edu",
}

type categoryItems struct {
	category string
	names    []string
}

var itemCatalog = []categoryItems{
	{
		category: "dorm supplies",
		names: []string{
			"Mini Fridge", "Desk Lamp", "Twin XL Bedding Set", "Storage Bins",
			"Box Fan", "Power Strip", "Shower Caddy", "Over-Door Organizer",
			"Laundry Hamper", "Whiteboard",
		},
	},
	{
		category: "school supplies",
		names: []string{
			"Scientific Calculator", "Backpack", "Notebook Set", "Binder Pack",
			"Highlighter Set", "Pencil Case", "Laptop Stand", "Planner",
			"Stapler", "Wireless Mouse",
		},
	},
	{
		category: "misc",
		names: []string{
			"Coffee Maker", "Bluetooth Speaker", "Bike Lock", "Umbrella",
			"Yoga Mat", "Reusable Water Bottle", "Headphones", "Desk Chair",
			"Floor Lamp", "Microwave",
		},
	},
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
	for _, table := range []string{"transactions", "items", "users"} {
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
	const q = `INSERT OR REPLACE INTO users (username, email, profile_picture) VALUES (?, ?, ?)`

	hashedPw := hash("password123")
	_ = hashedPw

	var ids []int64
	for i := 0; i < numUsers; i++ {
		name := fmt.Sprintf("%s %s", gofakeit.FirstName(), gofakeit.LastName())
		emoji := emojis[gofakeit.Number(0, len(emojis)-1)]

		var email string
		if gofakeit.Bool() {
			domain := eduDomains[gofakeit.Number(0, len(eduDomains)-1)]
			email = fmt.Sprintf("%s.%s@%s",
				gofakeit.FirstName(), gofakeit.LastName(), domain)
		} else {
			email = gofakeit.Email()
		}

		res, err := db.Exec(q, name, email, emoji)
		if err != nil {
			log.Fatalf("insert user: %v", err)
		}
		id, _ := res.LastInsertId()
		ids = append(ids, id)
		log.Printf("inserted user %s %s (%s)", emoji, name, email)
	}
	return ids
}

func seedItems(db *sql.DB, userIDs []int64) []int64 {
	const q = `INSERT OR REPLACE INTO items (name, user_id, description, price, category) VALUES (?, ?, ?, ?, ?)`

	var ids []int64
	for i := 0; i < numItems; i++ {
		cat := itemCatalog[gofakeit.Number(0, len(itemCatalog)-1)]
		name := cat.names[gofakeit.Number(0, len(cat.names)-1)]
		userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		description := gofakeit.Sentence(8)
		price := gofakeit.Price(1, 150)

		res, err := db.Exec(q, name, userID, description, price, cat.category)
		if err != nil {
			log.Fatalf("insert item: %v", err)
		}
		id, _ := res.LastInsertId()
		ids = append(ids, id)
		log.Printf("inserted item [%s] %q ($%.2f)", cat.category, name, price)
	}
	return ids
}

func seedTransactions(db *sql.DB, userIDs, itemIDs []int64) {
	const q = `INSERT INTO transactions (user_id, item_id, quantity, total_price) VALUES (?, ?, ?, ?)`

	for i := 0; i < numTransactions; i++ {
		userID := userIDs[gofakeit.Number(0, len(userIDs)-1)]
		itemID := itemIDs[gofakeit.Number(0, len(itemIDs)-1)]
		quantity := gofakeit.Number(1, 3)

		var price float64
		db.QueryRow(`SELECT price FROM items WHERE id = ?`, itemID).Scan(&price)
		total := price * float64(quantity)

		if _, err := db.Exec(q, userID, itemID, quantity, total); err != nil {
			log.Fatalf("insert transaction: %v", err)
		}
		log.Printf("inserted transaction: user %d, item %d, qty %d", userID, itemID, quantity)
	}
}
