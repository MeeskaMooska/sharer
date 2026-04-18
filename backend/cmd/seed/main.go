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

var emojis = []string{
	"🐶", "🐱", "🐭", "🐹", "🐰", "🦊", "🐻", "🐼",
	"🐨", "🐯", "🦁", "🐮", "🐸", "🐙", "🦋", "🐺",
	"🦄", "🐲", "🌵", "🍄",
}

var eduDomains = []string{
	"mit.edu", "stanford.edu", "harvard.edu", "unc.edu",
	"gatech.edu", "umich.edu", "berkeley.edu", "nyu.edu",
}

var schools = []string{
	"MIT", "Stanford", "Harvard", "UNC Chapel Hill",
	"Georgia Tech", "University of Michigan", "UC Berkeley", "NYU",
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
	for _, p := range []string{".env", "../.env", "../../.env"} {
		if godotenv.Load(p) == nil {
			break
		}
	}

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
		INSERT INTO users (email, school, hashed_password, goodwill_points, profile_picture)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			school          = VALUES(school),
			hashed_password = VALUES(hashed_password),
			goodwill_points = VALUES(goodwill_points),
			profile_picture = VALUES(profile_picture)`

	hashedPw := hash("password123")
	var ids []int64

	for i := 0; i < numUsers; i++ {
		domain := eduDomains[gofakeit.Number(0, len(eduDomains)-1)]
		email := fmt.Sprintf("%s.%s@%s", gofakeit.FirstName(), gofakeit.LastName(), domain)
		school := schools[gofakeit.Number(0, len(schools)-1)]
		goodwill := gofakeit.Number(0, 50)
		emoji := emojis[gofakeit.Number(0, len(emojis)-1)]

		res, err := db.Exec(q, email, school, hashedPw, goodwill, emoji)
		if err != nil {
			log.Fatalf("insert user: %v", err)
		}
		id, _ := res.LastInsertId()
		ids = append(ids, id)
		log.Printf("inserted user %s %s (%s)", emoji, email, school)
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
		cat := itemCatalog[gofakeit.Number(0, len(itemCatalog)-1)]
		name := cat.names[gofakeit.Number(0, len(cat.names)-1)]
		description := gofakeit.Sentence(8)
		value := gofakeit.Price(1, 150)

		res, err := db.Exec(q, name, description, value, cat.category)
		if err != nil {
			log.Fatalf("insert item: %v", err)
		}
		id, _ := res.LastInsertId()
		ids = append(ids, id)
		log.Printf("inserted item [%s] %q ($%.2f)", cat.category, name, value)
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
