package main

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Connect to database
	db, err := sql.Open("pgx", "postgres://postgres@127.0.0.1:5432/safecast")
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	ctx := context.Background()

	log.Println("Adding internal_user_id column to uploads table...")

	// Add column if it doesn't exist
	_, err = db.ExecContext(ctx, `
		ALTER TABLE uploads ADD COLUMN IF NOT EXISTS internal_user_id TEXT
	`)
	if err != nil {
		log.Fatalf("Error adding column: %v", err)
	}

	log.Println("Creating index on internal_user_id...")
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_uploads_internal_user_id ON uploads(internal_user_id)
	`)
	if err != nil {
		log.Fatalf("Error creating index: %v", err)
	}

	log.Println("Column and index added successfully!")
}
