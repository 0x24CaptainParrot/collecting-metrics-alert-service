package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgresDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Println("warning: DATABASE_DSN is not set. Server will start without db connection.")
		return nil, nil
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to the database.")

	return db, nil
}
