package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func main() {
	// Command-line flags
	csvFile := flag.String("csv", "", "Path to CSV file with user API keys (required)")
	dbType := flag.String("db-type", "pgx", "Database type: pgx, sqlite, chai, duckdb")
	dbConn := flag.String("db-conn", "", "Database connection string (required)")
	dryRun := flag.Bool("dry-run", false, "Dry run - don't actually update users")

	flag.Parse()

	// Validate required flags
	if *csvFile == "" {
		log.Fatal("Error: -csv flag is required")
	}
	if *dbConn == "" {
		log.Fatal("Error: -db-conn flag is required")
	}

	// Open CSV file
	file, err := os.Open(*csvFile)
	if err != nil {
		log.Fatalf("Error opening CSV file: %v", err)
	}
	defer file.Close()

	// Parse CSV
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV file: %v", err)
	}

	if len(records) < 2 {
		log.Fatal("Error: CSV file must have at least a header row and one data row")
	}

	// Find column indices from header
	header := records[0]
	emailIdx := -1
	apiKeyIdx := -1
	for i, col := range header {
		col = strings.ToLower(strings.TrimSpace(col))
		if col == "email" {
			emailIdx = i
		}
		if col == "authentication_token" || col == "api_key" || col == "apikey" {
			apiKeyIdx = i
		}
	}

	if emailIdx == -1 {
		log.Fatal("Error: CSV file must have an 'email' column")
	}
	if apiKeyIdx == -1 {
		log.Fatal("Error: CSV file must have an 'authentication_token' or 'api_key' column")
	}

	// Connect to database
	db, err := sql.Open(*dbType, *dbConn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Printf("Connected to %s database", *dbType)
	log.Printf("Processing %d records from CSV (excluding header)", len(records)-1)

	if *dryRun {
		log.Println("DRY RUN MODE - no changes will be made")
	}

	// Process records
	ctx := context.Background()
	updated := 0
	notFound := 0
	skipped := 0
	errors := 0

	for i, record := range records[1:] {
		if len(record) <= emailIdx || len(record) <= apiKeyIdx {
			log.Printf("Row %d: skipping - insufficient columns", i+2)
			skipped++
			continue
		}

		email := strings.TrimSpace(record[emailIdx])
		apiKey := strings.TrimSpace(record[apiKeyIdx])

		if email == "" {
			log.Printf("Row %d: skipping - empty email", i+2)
			skipped++
			continue
		}

		if apiKey == "" {
			log.Printf("Row %d: skipping %s - empty API key", i+2, email)
			skipped++
			continue
		}

		// Check if user exists
		var userID int64
		var existingAPIKey sql.NullString
		var checkQuery string
		if *dbType == "pgx" || *dbType == "duckdb" {
			checkQuery = "SELECT id, api_key FROM users WHERE email = $1"
		} else {
			checkQuery = "SELECT id, api_key FROM users WHERE email = ?"
		}

		err := db.QueryRowContext(ctx, checkQuery, email).Scan(&userID, &existingAPIKey)
		if err == sql.ErrNoRows {
			log.Printf("Row %d: user not found: %s", i+2, email)
			notFound++
			continue
		} else if err != nil {
			log.Printf("Row %d: error checking user %s: %v", i+2, email, err)
			errors++
			continue
		}

		// Skip if user already has an API key
		if existingAPIKey.Valid && existingAPIKey.String != "" {
			log.Printf("Row %d: user %s already has API key, skipping", i+2, email)
			skipped++
			continue
		}

		// Update user's API key
		if !*dryRun {
			var updateQuery string
			if *dbType == "pgx" || *dbType == "duckdb" {
				updateQuery = "UPDATE users SET api_key = $1 WHERE id = $2"
			} else {
				updateQuery = "UPDATE users SET api_key = ? WHERE id = ?"
			}

			_, err := db.ExecContext(ctx, updateQuery, apiKey, userID)
			if err != nil {
				log.Printf("Row %d: error updating API key for %s: %v", i+2, email, err)
				errors++
				continue
			}
		}

		log.Printf("Row %d: %s API key for %s (ID: %d)", i+2, func() string {
			if *dryRun {
				return "would update"
			}
			return "updated"
		}(), email, userID)
		updated++
	}

	log.Println("---")
	log.Printf("Summary:")
	log.Printf("  Updated: %d", updated)
	log.Printf("  Not found: %d", notFound)
	log.Printf("  Skipped: %d", skipped)
	log.Printf("  Errors: %d", errors)

	if *dryRun {
		log.Println("\nDRY RUN completed - no changes were made")
		log.Println("Run without -dry-run to apply changes")
	}
}
