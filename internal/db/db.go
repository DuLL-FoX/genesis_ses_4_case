package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB  *sql.DB
	mu  sync.Mutex
	ctx = context.Background()
)

func InitDB() error {
	var err error
	dbFile := os.Getenv("DB_FILE")
	DB, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	return nil
}

func GetDB() *sql.DB {
	return DB
}

func GetContext() context.Context {
	return ctx
}

func GetMutex() *sync.Mutex {
	return &mu
}
