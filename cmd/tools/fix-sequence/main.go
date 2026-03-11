// fix_sequence.go - Utility to sync PostgreSQL sequence with actual data
// Usage: go run tools/fix_sequence.go
//
// This fixes: ERROR: duplicate key value violates unique constraint "markers_pkey"

package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Get database URL from environment or use default
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Try common default
		dbURL = "postgres://safecast:safecast@localhost/safecast?sslmode=disable"
		log.Printf("Using default database URL: %s", dbURL)
		log.Println("Set DATABASE_URL environment variable to use a different connection")
	}

	ctx := context.Background()

	// Connect to database
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Check current sequence value
	var lastValue int64
	var isCalled bool
	err = db.QueryRowContext(ctx, "SELECT last_value, is_called FROM markers_id_seq").Scan(&lastValue, &isCalled)
	if err != nil {
		log.Fatalf("Failed to get sequence value: %v", err)
	}
	log.Printf("Current sequence: last_value=%d, is_called=%v", lastValue, isCalled)

	// Check actual MAX(id) in the table
	var maxID sql.NullInt64
	err = db.QueryRowContext(ctx, "SELECT MAX(id) FROM markers").Scan(&maxID)
	if err != nil {
		log.Fatalf("Failed to get MAX(id): %v", err)
	}

	if maxID.Valid {
		log.Printf("Actual MAX(id) in markers table: %d", maxID.Int64)

		if maxID.Int64 >= lastValue {
			log.Printf("WARNING: MAX(id) (%d) >= sequence last_value (%d)", maxID.Int64, lastValue)
			log.Println("This will cause duplicate key errors on insert!")

			// Fix: Reset the sequence to be after the maximum ID
			newSeqValue := maxID.Int64 + 1
			log.Printf("Fixing: Setting sequence to %d...", newSeqValue)

			_, err = db.ExecContext(ctx, "SELECT setval('markers_id_seq', $1, true)", newSeqValue)
			if err != nil {
				log.Fatalf("Failed to reset sequence: %v", err)
			}

			log.Printf("Sequence reset successfully to %d", newSeqValue)

			// Verify the fix
			err = db.QueryRowContext(ctx, "SELECT last_value, is_called FROM markers_id_seq").Scan(&lastValue, &isCalled)
			if err != nil {
				log.Fatalf("Failed to verify sequence: %v", err)
			}
			log.Printf("Verified: new sequence value is %d (is_called=%v)", lastValue, isCalled)
		} else {
			log.Println("Sequence is ahead of MAX(id) - no fix needed")
		}
	} else {
		log.Println("markers table is empty - no fix needed")
	}

	log.Println("Done!")
}
