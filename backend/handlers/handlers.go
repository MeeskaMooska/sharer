package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{ db *sql.DB }

func New(db *sql.DB) *Handler { return &Handler{db: db} }

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) Users(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT id, username, email, profile_picture, goodwill_points, created_at FROM users`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type User struct {
		ID             int64          `json:"id"`
		Username       string         `json:"username"`
		Email          string         `json:"email"`
		ProfilePicture sql.NullString `json:"profile_picture"`
		GoodwillPoints int            `json:"goodwill_points"`
		CreatedAt      string         `json:"created_at"`
	}
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.ProfilePicture, &u.GoodwillPoints, &u.CreatedAt); err != nil {
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
	rows, err := h.db.Query(`SELECT id, name, user_id, description, price, created_at FROM items`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Item struct {
		ID          int64   `json:"id"`
		Name        string  `json:"name"`
		UserID      int64   `json:"user_id"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		CreatedAt   string  `json:"created_at"`
	}
	var items []Item
	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.ID, &i.Name, &i.UserID, &i.Description, &i.Price, &i.CreatedAt); err != nil {
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

// Transactions handles GET /api/transactions and PATCH /api/transactions/{id}
func (h *Handler) Transactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listTransactions(w, r)
	case http.MethodPatch:
		h.completeTransaction(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) listTransactions(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query(`SELECT id, user_id, item_id, quantity, total_price, completed, created_at FROM transactions`)
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
		Completed  bool    `json:"completed"`
		CreatedAt  string  `json:"created_at"`
	}
	var txs []Transaction
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.ItemID, &t.Quantity, &t.TotalPrice, &t.Completed, &t.CreatedAt); err != nil {
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

// completeTransaction marks a transaction complete and awards the giver 1 goodwill point.
func (h *Handler) completeTransaction(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/transactions/")
	txID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid transaction id", http.StatusBadRequest)
		return
	}

	var alreadyDone bool
	var giverID int64
	err = h.db.QueryRow(`SELECT completed, user_id FROM transactions WHERE id = ?`, txID).
		Scan(&alreadyDone, &giverID)
	if err == sql.ErrNoRows {
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if alreadyDone {
		http.Error(w, "transaction already completed", http.StatusConflict)
		return
	}

	if _, err := h.db.Exec(`UPDATE transactions SET completed = 1 WHERE id = ?`, txID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := awardGiverPoints(h.db, giverID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]any{"id": txID, "giver_id": giverID, "points_awarded": 1})
}

func awardGiverPoints(db *sql.DB, userID int64) error {
	_, err := db.Exec(`UPDATE users SET goodwill_points = goodwill_points + 1 WHERE id = ?`, userID)
	return err
}
