package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

func NewPostgresDB(dsn string) (*sql.DB, error) {
	if dsn == "" {
		log.Println("warning: DATABASE_DSN is not set. Server will start without db connection.")
		return nil, nil
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to the database.")

	return db, nil
}

func RunMigrations(db *sql.DB, migrationsDir string) error {
	if db == nil {
		return fmt.Errorf("cannot run migrations, db is nil")
	}

	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		return fmt.Errorf("migrations directory not found: %s", migrationsDir)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("goose migration failed: %v", err)
	}

	log.Println("Migrations applied successfully.")
	return nil
}
