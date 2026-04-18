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
	rows, err := h.db.Query(`SELECT id, email, school, goodwill_points, profile_picture, created_at FROM users`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		ID             int64          `json:"id"`
		Email          string         `json:"email"`
		School         sql.NullString `json:"school"`
		GoodwillPoints int            `json:"goodwill_points"`
		ProfilePicture sql.NullString `json:"profile_picture"`
		CreatedAt      string         `json:"created_at"`
	}
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.School, &u.GoodwillPoints, &u.ProfilePicture, &u.CreatedAt); err != nil {
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
	rows, err := h.db.Query(`SELECT id, name, description, value, category, picture, created_at FROM items`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Item struct {
		ID          int64          `json:"id"`
		Name        string         `json:"name"`
		Description sql.NullString `json:"description"`
		Value       float64        `json:"value"`
		Category    sql.NullString `json:"category"`
		Picture     sql.NullString `json:"picture"`
		CreatedAt   string         `json:"created_at"`
	}
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Name, &i.Description, &i.Value, &i.Category, &i.Picture, &i.CreatedAt); err != nil {
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
	rows, err := h.db.Query(`SELECT id, user_giving, user_receiving, item_id, reviewed, review, created_at FROM transactions`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Transaction struct {
		ID            int64        `json:"id"`
		UserGiving    int64        `json:"user_giving"`
		UserReceiving int64        `json:"user_receiving"`
		ItemID        int64        `json:"item_id"`
		Reviewed      bool         `json:"reviewed"`
		Review        sql.NullBool `json:"review"`
		CreatedAt     string       `json:"created_at"`
	}
	var txs []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.UserGiving, &t.UserReceiving, &t.ItemID, &t.Reviewed, &t.Review, &t.CreatedAt); err != nil {
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
