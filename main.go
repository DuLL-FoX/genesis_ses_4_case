package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"sync"
)

// NBU API URL
const nbuAPI = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=USD&json"

// ExchangeRate represents the response from the NBU API
type ExchangeRate struct {
	Rate float64 `json:"rate"`
}

// Global variables to store subscribed emails and a mutex for synchronization
var (
	subscribedEmails = make(map[string]struct{})
	mu               sync.Mutex
)

// Function to get the current exchange rate from the NBU API
func getExchangeRate() (float64, error) {
	resp, err := http.Get(nbuAPI)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rates []ExchangeRate
	if err := json.NewDecoder(resp.Body).Decode(&rates); err != nil {
		return 0, err
	}

	if len(rates) == 0 {
		return 0, fmt.Errorf("no exchange rate found")
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
	subscribedEmails[email] = struct{}{}
	mu.Unlock()

	fmt.Fprintf(w, "Subscribed email: %s", email)
}

// Function to send emails to all subscribed users
func sendEmails(rate float64) {
	mu.Lock()
	defer mu.Unlock()

	for email := range subscribedEmails {
		go sendEmail(email, rate)
	}
}

// Function to send an individual email
func sendEmail(email string, rate float64) {
	from := os.Getenv("ETHEREAL_EMAIL")
	password := os.Getenv("ETHEREAL_PASSWORD")

	if from == "" || password == "" {
		log.Fatal("ETHEREAL_EMAIL or ETHEREAL_PASSWORD environment variable is not set")
	}

	to := email
	subject := "Daily Exchange Rate"
	body := fmt.Sprintf("Current USD to UAH exchange rate: %.2f", rate)

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	auth := smtp.PlainAuth("", from, password, "smtp.ethereal.email")
	err := smtp.SendMail("smtp.ethereal.email:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("failed to send email to %s: %v", to, err)
	}

	if err == nil {
		log.Printf("email sent to %s", to)
	}
}

// Main function to set up HTTP handlers and start the server and scheduler
func main() {
	http.HandleFunc("/rate", func(w http.ResponseWriter, r *http.Request) {
		rate, err := getExchangeRate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Current USD to UAH exchange rate: %.2f", rate)
		sendEmails(rate)
	})

	http.HandleFunc("/subscribe", subscribeHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
