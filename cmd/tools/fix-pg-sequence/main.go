// Quick fix for PostgreSQL sequence sync
// Run: go run fix-pg-sequence.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set. Please set it, e.g.:")
		log.Fatal("  export DATABASE_URL=\"postgres://user:pass@host:port/db?sslmode=disable\"")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot ping database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL")

	// Get current sequence value
	var seqVal, maxID int64
	err = db.QueryRow("SELECT last_value FROM markers_id_seq").Scan(&seqVal)
	if err != nil {
		log.Fatalf("Failed to get sequence: %v", err)
	}

	err = db.QueryRow("SELECT COALESCE(MAX(id), 0) FROM markers").Scan(&maxID)
	if err != nil {
		log.Fatalf("Failed to get MAX(id): %v", err)
	}

	fmt.Printf("Current sequence value: %d\n", seqVal)
	fmt.Printf("Actual MAX(id) in table: %d\n", maxID)

	if maxID >= seqVal {
		fmt.Printf("\n⚠️  PROBLEM: MAX(id) >= sequence value\n")
		fmt.Printf("Fixing: Setting sequence to %d...\n", maxID+1)

		_, err = db.Exec("SELECT setval('markers_id_seq', $1, true)", maxID+1)
		if err != nil {
			log.Fatalf("Failed to fix sequence: %v", err)
		}

		fmt.Println("✓ Sequence fixed successfully!")

		// Verify
		var newVal int64
		db.QueryRow("SELECT last_value FROM markers_id_seq").Scan(&newVal)
		fmt.Printf("New sequence value: %d\n", newVal)
	} else {
		fmt.Println("\n✓ Sequence is OK (ahead of MAX(id))")
	}
}
