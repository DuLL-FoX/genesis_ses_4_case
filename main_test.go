package main

import (
	_ "context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS subscribers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT UNIQUE NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}

	return db
}

func TestInitDB(t *testing.T) {
	dbFile = ":memory:"
	err := initDB()
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
}

func TestGetExchangeRate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rates := []ExchangeRate{{Rate: 27.5}}
		json.NewEncoder(w).Encode(rates)
	}))
	defer server.Close()

	nbuAPI = server.URL
	rate, err := getExchangeRate()
	if err != nil {
		t.Fatalf("failed to get exchange rate: %v", err)
	}

	if rate != 27.5 {
		t.Errorf("expected rate 27.5, got %.2f", rate)
	}
}

func TestSubscribeHandler(t *testing.T) {
	db = setupTestDB(t)

	req, err := http.NewRequest("POST", "/subscribe", strings.NewReader("email=test@example.com"))
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(subscribeHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Subscribed email: test@example.com"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestSendEmail(t *testing.T) {
	err := sendEmail("recipient@example.com", 27.5)
	if err != nil {
		t.Errorf("failed to send email: %v", err)
	}
}

func TestMain(m *testing.M) {
	os.Setenv("DB_FILE", ":memory:")
	err := initDB()
	if err != nil {
		fmt.Printf("failed to initialize database: %v\n", err)
		os.Exit(1)
	}

	exitCode := m.Run()

	// Clean up code after running tests
	os.Unsetenv("DB_FILE")
	os.Unsetenv("ETHEREAL_EMAIL")
	os.Unsetenv("ETHEREAL_PASSWORD")

	os.Exit(exitCode)
}
