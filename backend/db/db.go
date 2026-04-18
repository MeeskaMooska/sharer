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
			id              BIGINT       NOT NULL AUTO_INCREMENT PRIMARY KEY,
			email           VARCHAR(255) NOT NULL UNIQUE,
			school          VARCHAR(255),
			hashed_password VARCHAR(255) NOT NULL,
			goodwill_points INT          NOT NULL DEFAULT 0,
			profile_picture VARCHAR(512),
			created_at      DATETIME     DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB`,
		`CREATE TABLE IF NOT EXISTS items (
			id          BIGINT        NOT NULL AUTO_INCREMENT PRIMARY KEY,
			name        VARCHAR(255)  NOT NULL UNIQUE,
			description TEXT,
			value       DECIMAL(10,2) NOT NULL DEFAULT 0,
			category    VARCHAR(100),
			picture     VARCHAR(512),
			created_at  DATETIME      DEFAULT CURRENT_TIMESTAMP
		) ENGINE=InnoDB`,
		`CREATE TABLE IF NOT EXISTS transactions (
			id             BIGINT     NOT NULL AUTO_INCREMENT PRIMARY KEY,
			user_giving    BIGINT     NOT NULL,
			user_receiving BIGINT     NOT NULL,
			item_id        BIGINT     NOT NULL,
			reviewed       TINYINT(1) NOT NULL DEFAULT 0,
			review         TINYINT(1),
			created_at     DATETIME   DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_giving)    REFERENCES users(id),
			FOREIGN KEY (user_receiving) REFERENCES users(id),
			FOREIGN KEY (item_id)        REFERENCES items(id)
		) ENGINE=InnoDB`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
