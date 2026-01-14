package main

import (
	"context"
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/email"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

func main() {
	// Command-line flags
	csvFile := flag.String("csv", "", "Path to CSV file with user data (required)")
	dbType := flag.String("db-type", "pgx", "Database type: pgx, sqlite, chai, duckdb")
	dbConn := flag.String("db-conn", "", "Database connection string (required)")
	smtpHost := flag.String("smtp-host", "", "SMTP server hostname (required for sending emails)")
	smtpPort := flag.Int("smtp-port", 587, "SMTP server port")
	smtpUsername := flag.String("smtp-username", "", "SMTP username")
	smtpPassword := flag.String("smtp-password", "", "SMTP password")
	smtpFrom := flag.String("smtp-from", "", "Email from address")
	smtpFromName := flag.String("smtp-from-name", "Safecast", "Email from name")
	baseURL := flag.String("base-url", "https://map.safecast.org", "Base URL for password setup links")
	dryRun := flag.Bool("dry-run", false, "Dry run - don't actually insert users or send emails")
	sendEmails := flag.Bool("send-emails", true, "Send password setup emails to users")

	flag.Parse()

	// Validate required flags
	if *csvFile == "" {
		log.Fatal("Error: -csv flag is required")
	}
	if *dbConn == "" {
		log.Fatal("Error: -db-conn flag is required")
	}
	if *sendEmails && *smtpHost == "" {
		log.Fatal("Error: -smtp-host is required when sending emails")
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

	log.Printf("Connected to database (type: %s)", *dbType)

	// Setup email sender if needed
	var emailSender *email.Sender
	if *sendEmails && !*dryRun {
		emailConfig := email.SMTPConfig{
			Host:     *smtpHost,
			Port:     *smtpPort,
			Username: *smtpUsername,
			Password: *smtpPassword,
			From:     *smtpFrom,
			FromName: *smtpFromName,
			UseTLS:   true,
		}
		emailSender = email.NewSender(emailConfig)
		log.Printf("Email sender configured (SMTP: %s:%d)", *smtpHost, *smtpPort)
	}

	// Parse header row to find column indices
	header := records[0]
	userIDIdx := findColumnIndex(header, "user_id", "id", "external_id")
	emailIdx := findColumnIndex(header, "email", "user_email")
	usernameIdx := findColumnIndex(header, "username", "user_name", "name")

	if userIDIdx == -1 {
		log.Fatal("Error: Could not find user_id column in CSV")
	}
	if emailIdx == -1 {
		log.Fatal("Error: Could not find email column in CSV")
	}

	log.Printf("CSV columns found - user_id: %d, email: %d, username: %d", userIDIdx, emailIdx, usernameIdx)
	log.Printf("Found %d users to migrate (excluding header)", len(records)-1)

	if *dryRun {
		log.Println("DRY RUN MODE - No changes will be made")
	}

	// Process each user
	ctx := context.Background()
	successCount := 0
	skipCount := 0
	errorCount := 0

	for i, record := range records[1:] {
		lineNum := i + 2 // +2 because we skip header and arrays are 0-indexed

		// Extract user data
		if len(record) <= userIDIdx || len(record) <= emailIdx {
			log.Printf("Line %d: Skipping - insufficient columns", lineNum)
			skipCount++
			continue
		}

		externalUserID := strings.TrimSpace(record[userIDIdx])
		userEmail := strings.TrimSpace(record[emailIdx])
		var username string
		if usernameIdx != -1 && len(record) > usernameIdx {
			username = strings.TrimSpace(record[usernameIdx])
		}

		// Validate email
		if userEmail == "" {
			log.Printf("Line %d: Skipping - empty email", lineNum)
			skipCount++
			continue
		}

		if err := auth.ValidateEmail(userEmail); err != nil {
			log.Printf("Line %d: Skipping - invalid email '%s': %v", lineNum, userEmail, err)
			skipCount++
			continue
		}

		// Check if user already exists
		existingUser, err := auth.GetUserByEmail(ctx, db, *dbType, userEmail)
		if err != nil && err != sql.ErrNoRows {
			log.Printf("Line %d: Error checking existing user '%s': %v", lineNum, userEmail, err)
			errorCount++
			continue
		}

		if existingUser != nil {
			log.Printf("Line %d: User '%s' already exists (ID: %d), skipping", lineNum, userEmail, existingUser.ID)
			skipCount++
			continue
		}

		if *dryRun {
			log.Printf("Line %d: [DRY RUN] Would create user: email=%s, username=%s, external_id=%s",
				lineNum, userEmail, username, externalUserID)
			successCount++
			continue
		}

		// Create user with temporary password (they'll need to set one via email)
		tempPassword, err := auth.GenerateRandomPassword(16)
		if err != nil {
			log.Printf("Line %d: Error generating password for '%s': %v", lineNum, userEmail, err)
			errorCount++
			continue
		}

		user := &auth.User{
			Email:                 userEmail,
			Username:              username,
			EmailVerified:         false, // They'll verify when they set their password
			IsActive:              true,
			ExternalID:            externalUserID,
			ExternalSource:        "csv-migration",
			RequiresPasswordSetup: true,
		}

		userID, err := auth.CreateUser(ctx, db, *dbType, user, tempPassword)
		if err != nil {
			log.Printf("Line %d: Error creating user '%s': %v", lineNum, userEmail, err)
			errorCount++
			continue
		}

		log.Printf("Line %d: Created user '%s' (ID: %d)", lineNum, userEmail, userID)

		// Link existing uploads to this user
		if externalUserID != "" {
			err = auth.LinkUserIDToUploads(ctx, db, *dbType, externalUserID, fmt.Sprintf("%d", userID))
			if err != nil {
				log.Printf("Line %d: Warning: Could not link uploads for external ID '%s': %v",
					lineNum, externalUserID, err)
			} else {
				log.Printf("Line %d: Linked existing uploads for external ID '%s'", lineNum, externalUserID)
			}
		}

		// Send password setup email
		if emailSender != nil {
			// Create password reset token for setup
			token, err := auth.GenerateToken()
			if err != nil {
				log.Printf("Line %d: Error generating token for '%s': %v", lineNum, userEmail, err)
				errorCount++
				continue
			}

			// Token expires in 7 days for initial setup
			expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix()
			resetToken := &auth.PasswordResetToken{
				UserID:    userID,
				Token:     token,
				ExpiresAt: expiresAt,
				Used:      false,
			}

			if err := auth.CreatePasswordResetToken(ctx, db, *dbType, resetToken); err != nil {
				log.Printf("Line %d: Error creating password reset token for '%s': %v", lineNum, userEmail, err)
				errorCount++
				continue
			}

			setupURL := fmt.Sprintf("%s/reset-password?token=%s", *baseURL, token)

			// Send email
			if err := emailSender.SendPasswordSetupEmail(userEmail, username, setupURL); err != nil {
				log.Printf("Line %d: Error sending email to '%s': %v", lineNum, userEmail, err)
				errorCount++
				continue
			}

			log.Printf("Line %d: Sent password setup email to '%s'", lineNum, userEmail)
		}

		successCount++
	}

	// Print summary
	log.Println("\n========== Migration Summary ==========")
	log.Printf("Total users processed: %d", len(records)-1)
	log.Printf("Successfully migrated: %d", successCount)
	log.Printf("Skipped (already exists or invalid): %d", skipCount)
	log.Printf("Errors: %d", errorCount)
	log.Println("=======================================")

	if *dryRun {
		log.Println("\nThis was a DRY RUN - no actual changes were made")
	}
}

// findColumnIndex finds the index of a column in the header row.
// It tries multiple possible column names (case-insensitive).
func findColumnIndex(header []string, names ...string) int {
	for i, col := range header {
		// Remove BOM if present and trim whitespace
		colClean := strings.TrimPrefix(col, "\uFEFF")
		colLower := strings.ToLower(strings.TrimSpace(colClean))
		for _, name := range names {
			if colLower == strings.ToLower(name) {
				return i
			}
		}
	}
	return -1
}
