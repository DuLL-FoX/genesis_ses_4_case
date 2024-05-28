package handlers

import (
	"fmt"
	"log"
	"net/http"

	"awesomeProject/internal/db"
)

func SubscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	mu := db.GetMutex()
	mu.Lock()
	defer mu.Unlock()

	var exists bool
	err := db.GetDB().QueryRowContext(db.GetContext(), "SELECT EXISTS(SELECT 1 FROM subscribers WHERE email=?)", email).Scan(&exists)
	if err != nil {
		http.Error(w, "Failed to check subscription status", http.StatusInternalServerError)
		log.Printf("failed to check subscription status for email %s: %v", email, err)
		return
	}

	if exists {
		http.Error(w, "Email already subscribed", http.StatusConflict)
		return
	}

	_, err = db.GetDB().ExecContext(db.GetContext(), "INSERT INTO subscribers (email) VALUES (?)", email)
	if err != nil {
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		log.Printf("failed to subscribe email %s: %v", email, err)
		return
	}

	fmt.Fprintf(w, "Subscribed email: %s", email)
}
