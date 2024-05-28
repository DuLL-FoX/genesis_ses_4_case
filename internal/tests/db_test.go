package tests

import (
	"awesomeProject/internal/db"
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	os.Setenv("DB_FILE", "test.db")
	err := db.InitDB()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestRunMigrations(t *testing.T) {
	os.Setenv("DB_FILE", "test.db")
	err := db.InitDB()
	if err != nil {
		t.Fatalf("Expected no error initializing DB, got %v", err)
	}

	err = db.RunMigrations()
	if err != nil {
		t.Fatalf("Expected no error running migrations, got %v", err)
	}
}
