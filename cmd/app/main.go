package main

import (
	"fmt"
	"log"
	"net/http"

	"awesomeProject/internal/db"
	"awesomeProject/internal/handlers"
	"awesomeProject/internal/scheduler"
)

func main() {
	if err := db.InitDB(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	http.HandleFunc("/rate", handlers.ExchangeRateHandler)
	http.HandleFunc("/subscribe", handlers.SubscribeHandler)

	go scheduler.StartScheduler()

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
