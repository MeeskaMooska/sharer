package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Init(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id         BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
			username   VARCHAR(255) NOT NULL UNIQUE,
			email      VARCHAR(255) NOT NULL UNIQUE,
			created_at DATETIME     DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB`,
		`CREATE TABLE IF NOT EXISTS items (
			id          BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name        VARCHAR(255) NOT NULL,
			description TEXT,
			price       DECIMAL(10,2) NOT NULL DEFAULT 0,
			created_at  DATETIME      DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id          BIGINT        NOT NULL AUTO_INCREMENT PRIMARY KEY,
			user_id     BIGINT        NOT NULL,
			item_id     BIGINT        NOT NULL,
			quantity    INT           NOT NULL DEFAULT 1,
			total_price DECIMAL(10,2) NOT NULL,
			created_at  DATETIME      DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (item_id) REFERENCES items(id)
		) ENGINE=InnoDB`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
