package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

var (
	nbuAPI = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=USD&json"
	db     *sql.DB
	mu     sync.Mutex
	ctx    = context.Background()
	dbFile = os.Getenv("DB_FILE") // use an environment variable for the DB file path
)

// ExchangeRate represents the response from the NBU API
type ExchangeRate struct {
	Rate float64 `json:"rate"`
}

func initDB() error {
	var err error

	db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

// Function to get the current exchange rate from the NBU API
func getExchangeRate() (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, nbuAPI, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rates []ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(rates) == 0 {
		return 0, errors.New("no exchange rate found")
	}

	return rates[0].Rate, nil
}

// Handler to subscribe users with their email addresses
func subscribeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	var exists bool
	err := db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM subscribers WHERE email=?)", email).Scan(&exists)
	if err != nil {
		http.Error(w, "Failed to check subscription status", http.StatusInternalServerError)
		log.Printf("failed to check subscription status for email %s: %v", email, err)
		return
	}

	if exists {
		http.Error(w, "Email already subscribed", http.StatusConflict)
		return
	}

	_, err = db.ExecContext(ctx, "INSERT INTO subscribers (email) VALUES (?)", email)
	if err != nil {
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		log.Printf("failed to subscribe email %s: %v", email, err)
		return
	}

	fmt.Fprintf(w, "Subscribed email: %s", email)
}

// Function to send emails to all subscribed users
func sendEmails(rate float64) {
	mu.Lock()
	defer mu.Unlock()

	rows, err := db.QueryContext(ctx, "SELECT email FROM subscribers")
	if err != nil {
		log.Printf("failed to get subscribers: %v", err)
		return
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Printf("failed to scan email: %v", err)
			continue
		}
		emails = append(emails, email)
	}

	for _, email := range emails {
		if err := sendEmail(email, rate); err != nil {
			log.Printf("failed to send email to %s: %v", email, err)
		}
	}
}

// Function to send an individual email
func sendEmail(to string, rate float64) error {
	from := os.Getenv("ETHEREAL_EMAIL")
	password := os.Getenv("ETHEREAL_PASSWORD")

	if from == "" || password == "" {
		return errors.New("ETHEREAL_EMAIL or ETHEREAL_PASSWORD environment variable is not set")
	}

	subject := "Daily Exchange Rate"
	body := fmt.Sprintf("Current USD to UAH exchange rate: %.2f", rate)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, "smtp.ethereal.email")
	err := smtp.SendMail("smtp.ethereal.email:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("email sent to %s", to)
	return nil
}

// StartScheduler Function to start the scheduler that runs every 24 hours
func StartScheduler() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		rate, err := getExchangeRate()
		if err != nil {
			log.Printf("failed to get exchange rate: %v", err)
			continue
		}
		sendEmails(rate)
	}
}

// Main function to set up HTTP handlers and start the server and scheduler
func main() {
	if err := initDB(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	http.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		rate, err := getExchangeRate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Current USD to UAH exchange rate: %.2f", rate)
		sendEmails(rate) // For testing purposes only
	})

	http.HandleFunc("/subscribe", subscribeHandler)

	go StartScheduler()

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
