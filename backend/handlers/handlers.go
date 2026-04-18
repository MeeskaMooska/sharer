package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// Handler holds the database connection pool
type Handler struct {
	DB *sql.DB
}

func New(db *sql.DB) *Handler {
	return &Handler{DB: db}
}

// --- SIGN UP ---

type SignUpRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	School         string `json:"school"`
	ProfilePicture string `json:"profile_picture"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Insert into DB
	query := `INSERT INTO users (email, school, hashed_password, profile_picture) VALUES (?, ?, ?, ?)`
	result, err := h.DB.Exec(query, req.Email, req.School, string(hashedPassword), req.ProfilePicture)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		http.Error(w, "Failed to create user (email might already exist)", http.StatusConflict)
		return
	}

	userID, _ := result.LastInsertId()

	// THE HACKATHON COOKIE DROP
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprintf("%d", userID),
		Path:     "/",
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created and signed in successfully",
		"user_id": userID,
	})
}

// --- SIGN IN ---

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 1. Fetch the user's ID and hashed password from the database
	var userID int64
	var dbHash string

	query := `SELECT id, hashed_password FROM users WHERE email = ?`
	err := h.DB.QueryRow(query, req.Email).Scan(&userID, &dbHash)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't tell them if it's the email or password that was wrong
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		log.Printf("Database error during sign in: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 2. Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(req.Password))
	if err != nil {
		// Passwords didn't match
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3. Password matches! Drop the cookie.
	http.SetCookie(w, &http.Cookie{
		Name:     "user_id",
		Value:    fmt.Sprintf("%d", userID),
		Path:     "/",
		HttpOnly: true,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Signed in successfully",
		"user_id": userID,
	})
}

// Item maps directly to your MySQL items table
type Item struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"` // Pointer allows this to safely be NULL
	Value       float64 `json:"value"`       // Maps nicely from decimal(10,2)
	Category    *string `json:"category"`
	Picture     *string `json:"picture"`
	CreatedAt   string  `json:"created_at"`
}

// GetItems handles the GET /api/items route
func (h *Handler) GetItems(w http.ResponseWriter, r *http.Request) {
	// 1. Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Parse the 'page' parameter from the URL (e.g., /api/items?page=2)
	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if they don't provide it or send garbage
	}

	// 3. The Pagination Math
	limit := 20
	offset := (page - 1) * limit // Page 1 skips 0, Page 2 skips 20, etc.

	// 4. Query the database
	// We order by created_at DESC so the newest marketplace items show up first
	query := `
		SELECT id, name, description, value, category, picture, created_at 
		FROM items 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := h.DB.Query(query, limit, offset)
	if err != nil {
		log.Printf("Database query failed: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// 5. Parse the rows into our Go slice
	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Value,
			&item.Category,
			&item.Picture,
			&item.CreatedAt,
		); err != nil {
			log.Printf("Failed to scan item row: %v", err)
			continue // Skip broken rows so the whole API doesn't crash
		}
		items = append(items, item)
	}

	// Handle the edge case where the database is empty
	if items == nil {
		items = []Item{}
	}

	// 6. Return the paginated response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"current_page": page,
		"items_count":  len(items),
		"items":        items,
	})
}

// AddItemRequest defines the incoming JSON payload
type AddItemRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Value       float64 `json:"value"`
	Category    *string `json:"category"`
	Picture     *string `json:"picture"`
}

// AddItem handles POST /api/items
func (h *Handler) AddItem(w http.ResponseWriter, r *http.Request) {
	// 1. Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. THE HACKATHON AUTH CHECK
	// FIX: Actually save the cookie variable instead of using '_'
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized - Please sign in to post items", http.StatusUnauthorized)
		return
	}

	// 3. Parse the incoming JSON
	var req AddItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.Name == "" {
		http.Error(w, "Item name is required", http.StatusBadRequest)
		return
	}

	// 4. Insert into the database
	userID := cookie.Value

	query := `
		INSERT INTO items (user_id, name, description, value, category, picture) 
		VALUES (?, ?, ?, ?, ?, ?)
	`

	// FIX: Add userID as the very first argument to match the ? placeholders
	result, err := h.DB.Exec(query, userID, req.Name, req.Description, req.Value, req.Category, req.Picture)
	if err != nil {
		log.Printf("Failed to insert item: %v", err)
		http.Error(w, "Failed to create item (name might already exist)", http.StatusConflict)
		return
	}

	// 5. Get the new Item ID
	itemID, _ := result.LastInsertId()

	// 6. Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Item added to marketplace",
		"item_id": itemID,
	})
}

// --- 1. REQUEST AN ITEM ---

type RequestItemPayload struct {
	ItemID int64 `json:"item_id"`
}

func (h *Handler) RequestItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Auth Check: Who is asking for the item?
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	requesterID, _ := strconv.ParseInt(cookie.Value, 10, 64)

	var req RequestItemPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// Find out who owns the item
	var ownerID int64
	err = h.DB.QueryRow("SELECT user_id FROM items WHERE id = ?", req.ItemID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if ownerID == requesterID {
		http.Error(w, "You cannot request your own item", http.StatusBadRequest)
		return
	}

	// Insert the pending transaction
	query := `INSERT INTO transactions (user_giving, user_receiving, item_id) VALUES (?, ?, ?)`
	result, err := h.DB.Exec(query, ownerID, requesterID, req.ItemID)
	if err != nil {
		log.Printf("Failed to insert transaction: %v", err)
		http.Error(w, "Could not create request", http.StatusInternalServerError)
		return
	}

	txID, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Item requested successfully! Waiting for owner approval.",
		"transaction_id": txID,
		"status":         "pending",
	})
}

// --- 2. RESPOND TO A REQUEST ---

type RespondPayload struct {
	TransactionID int64  `json:"transaction_id"`
	Action        string `json:"action"` // Expects "accept" or "reject"
}

func (h *Handler) RespondToRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Auth Check: Who is responding?
	cookie, err := r.Cookie("user_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	responderID, _ := strconv.ParseInt(cookie.Value, 10, 64)

	var req RespondPayload
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	if req.Action != "accept" && req.Action != "reject" {
		http.Error(w, "Action must be 'accept' or 'reject'", http.StatusBadRequest)
		return
	}

	// Verify the transaction exists AND the responder is actually the owner (user_giving)
	var ownerID int64
	var currentStatus string
	err = h.DB.QueryRow("SELECT user_giving, status FROM transactions WHERE id = ?", req.TransactionID).Scan(&ownerID, &currentStatus)

	if err == sql.ErrNoRows {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	if ownerID != responderID {
		http.Error(w, "Forbidden: You do not own this item", http.StatusForbidden)
		return
	}

	if currentStatus != "pending" {
		http.Error(w, "This request has already been processed", http.StatusBadRequest)
		return
	}

	// Update the transaction status
	newStatus := "accepted"
	if req.Action == "reject" {
		newStatus = "rejected"
	}

	_, err = h.DB.Exec("UPDATE transactions SET status = ? WHERE id = ?", newStatus, req.TransactionID)
	if err != nil {
		log.Printf("Failed to update transaction: %v", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":        "Transaction " + newStatus,
		"transaction_id": req.TransactionID,
	})
}
