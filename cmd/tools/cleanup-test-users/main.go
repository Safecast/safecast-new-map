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

	// Delete related records first (foreign keys)
	log.Println("Deleting email verification tokens...")
	result, err := db.ExecContext(ctx, "DELETE FROM email_verification_tokens WHERE user_id >= 20000")
	if err != nil {
		log.Printf("Error deleting email verification tokens: %v", err)
	} else {
		rows, _ := result.RowsAffected()
		log.Printf("Deleted %d email verification tokens", rows)
	}

	log.Println("Deleting password reset tokens...")
	result, err = db.ExecContext(ctx, "DELETE FROM password_reset_tokens WHERE user_id >= 20000")
	if err != nil {
		log.Printf("Error deleting password reset tokens: %v", err)
	} else {
		rows, _ := result.RowsAffected()
		log.Printf("Deleted %d password reset tokens", rows)
	}

	log.Println("Deleting sessions...")
	result, err = db.ExecContext(ctx, "DELETE FROM sessions WHERE user_id >= 20000")
	if err != nil {
		log.Printf("Error deleting sessions: %v", err)
	} else {
		rows, _ := result.RowsAffected()
		log.Printf("Deleted %d sessions", rows)
	}

	log.Println("Deleting users...")
	result, err = db.ExecContext(ctx, "DELETE FROM users WHERE id >= 20000")
	if err != nil {
		log.Fatalf("Error deleting users: %v", err)
	}
	rows, _ := result.RowsAffected()
	log.Printf("Deleted %d users", rows)

	// Verify
	var count int
	err = db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE id >= 20000").Scan(&count)
	if err != nil {
		log.Fatalf("Error verifying deletion: %v", err)
	}
	log.Printf("Remaining users with ID >= 20000: %d", count)

	log.Println("Cleanup complete!")
}
