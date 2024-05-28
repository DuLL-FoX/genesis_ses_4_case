package email

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"awesomeProject/internal/db"
)

func SendEmails(rate float64) {
	mu := db.GetMutex()
	mu.Lock()
	defer mu.Unlock()

	rows, err := db.GetDB().QueryContext(db.GetContext(), "SELECT email FROM subscribers")
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
		if err := SendEmail(email, rate); err != nil {
			log.Printf("failed to send email to %s: %v", email, err)
		}
	}
}

func SendEmail(to string, rate float64) error {
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
