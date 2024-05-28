package scheduler

import (
	"log"
	"time"

	"awesomeProject/internal/email"
	"awesomeProject/internal/handlers"
)

func StartScheduler() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		rate, err := handlers.GetExchangeRate()
		if err != nil {
			log.Printf("failed to get exchange rate: %v", err)
			continue
		}
		email.SendEmails(rate)
	}
}
