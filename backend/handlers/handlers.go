package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type Handler struct{ db *sql.DB }

func New(db *sql.DB) *Handler { return &Handler{db: db} }

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT id, username, email, created_at FROM users`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		ID        int64  `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		CreatedAt string `json:"created_at"`
	}
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}
	if users == nil {
		users = []User{}
	}
	writeJSON(w, users)
}

func (h *Handler) Items(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT id, name, description, price, created_at FROM items`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Item struct {
		ID          int64   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		CreatedAt   string  `json:"created_at"`
	}
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Description, &i.Price, &i.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, i)
	}
	if items == nil {
		items = []Item{}
	}
	writeJSON(w, items)
}

func (h *Handler) Transactions(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT id, user_id, item_id, quantity, total_price, created_at FROM transactions`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Transaction struct {
		ID         int64   `json:"id"`
		UserID     int64   `json:"user_id"`
		ItemID     int64   `json:"item_id"`
		Quantity   int     `json:"quantity"`
		TotalPrice float64 `json:"total_price"`
		CreatedAt  string  `json:"created_at"`
	}
	var txs []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.ItemID, &t.Quantity, &t.TotalPrice, &t.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		txs = append(txs, t)
	}
	if txs == nil {
		txs = []Transaction{}
	}
	writeJSON(w, txs)
}
