package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Init(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	if err := Migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id              INTEGER PRIMARY KEY AUTOINCREMENT,
			username        TEXT    NOT NULL UNIQUE,
			email           TEXT    NOT NULL UNIQUE,
			profile_picture TEXT,
			goodwill_points INTEGER NOT NULL DEFAULT 0,
			created_at      DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS items (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			name        TEXT    NOT NULL,
			user_id     INTEGER NOT NULL REFERENCES users(id),
			description TEXT,
			price       REAL    NOT NULL DEFAULT 0,
			category    TEXT    NOT NULL DEFAULT 'misc',
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id     INTEGER NOT NULL REFERENCES users(id),
			item_id     INTEGER NOT NULL REFERENCES items(id),
			quantity    INTEGER NOT NULL DEFAULT 1,
			total_price REAL    NOT NULL,
			completed   INTEGER NOT NULL DEFAULT 0,
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
