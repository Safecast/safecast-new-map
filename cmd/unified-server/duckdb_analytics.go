// DuckDB Analytics Initialization for Unified Server
// Uses DuckLake with PostgreSQL catalog for shared analytics across all services

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/duckdb/duckdb-go/v2"
)

var duckDB *sql.DB

// initDuckDBAnalytics initializes DuckDB with DuckLake catalog backed by PostgreSQL.
// This allows multiple services to share the same analytics tables concurrently.
func initDuckDBAnalytics() error {
	// Open in-memory DuckDB — all persistent data lives in DuckLake (PostgreSQL + Parquet)
	var err error
	duckDB, err = sql.Open("duckdb", "")
	if err != nil {
		return fmt.Errorf("failed to open duckdb: %w", err)
	}

	duckDB.SetMaxOpenConns(1)
	duckDB.SetMaxIdleConns(1)
	duckDB.SetConnMaxLifetime(0)

	if err := duckDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping duckdb: %w", err)
	}

	log.Println("DuckDB initialized (in-memory)")

	// Install and load required extensions
	for _, ext := range []string{"ducklake", "postgres"} {
		if _, err := duckDB.Exec(fmt.Sprintf("INSTALL %s;", ext)); err != nil {
			log.Printf("Warning: INSTALL %s failed: %v", ext, err)
		}
		if _, err := duckDB.Exec(fmt.Sprintf("LOAD %s;", ext)); err != nil {
			return fmt.Errorf("LOAD %s: %w", ext, err)
		}
	}

	// Attach DuckLake catalog via PostgreSQL
	ducklakePGURL := os.Getenv("DUCKLAKE_PG_URL")
	if ducklakePGURL == "" {
		ducklakePGURL = "dbname=ducklake_catalog host=localhost user=ducklake_rw"
	}
	dataPath := os.Getenv("DUCKLAKE_DATA_PATH")
	if dataPath == "" {
		dataPath = "/var/lib/safecast/ducklake/"
	}

	attachQuery := fmt.Sprintf(
		"ATTACH 'ducklake:postgres:%s' AS analytics (DATA_PATH '%s');",
		ducklakePGURL, dataPath,
	)
	if _, err := duckDB.Exec(attachQuery); err != nil {
		return fmt.Errorf("attach DuckLake: %w", err)
	}
	log.Printf("DuckLake attached (catalog=PostgreSQL, data=%s)", dataPath)

	// Use analytics as default database
	if _, err := duckDB.Exec("USE analytics;"); err != nil {
		return fmt.Errorf("USE analytics: %w", err)
	}

	// Also attach main Safecast PostgreSQL for cross-database queries (read-only)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Unified server uses -db-conn flag; construct URL from that
		// Fall back to environment variable
	}
	if databaseURL != "" {
		if err := attachPostgres(databaseURL); err != nil {
			log.Printf("Warning: PostgreSQL attach failed: %v (cross-db analytics disabled)", err)
		} else {
			log.Println("Safecast PostgreSQL attached for cross-database queries")
		}
	}

	// Create analytics schema in DuckLake
	if err := createDuckDBSchema(); err != nil {
		log.Printf("Warning: failed to create DuckLake schema: %v", err)
	}

	return nil
}

// attachPostgres attaches the main Safecast PostgreSQL for read-only cross-database queries
func attachPostgres(databaseURL string) error {
	attachStr := databaseURL
	if !strings.Contains(databaseURL, "?") {
		attachStr = databaseURL + "?sslmode=prefer"
	}

	query := fmt.Sprintf("ATTACH '%s' AS postgres_db (TYPE POSTGRES, READ_ONLY);", attachStr)
	if _, err := duckDB.Exec(query); err != nil {
		return fmt.Errorf("attach postgres: %w", err)
	}
	return nil
}

// createDuckDBSchema creates the shared analytics tables in DuckLake
func createDuckDBSchema() error {
	// DuckLake doesn't support multi-statement exec, so run each separately
	tables := []string{
		`CREATE TABLE IF NOT EXISTS mcp_query_log (
			tool_name VARCHAR,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			duration_ms BIGINT,
			result_count INTEGER,
			client VARCHAR,
			user_id VARCHAR,
			user_email VARCHAR
		)`,
		`CREATE TABLE IF NOT EXISTS mcp_ai_query_log (
			user_id VARCHAR,
			user_email VARCHAR,
			session_id VARCHAR,
			timestamp TIMESTAMP,
			tool_name VARCHAR,
			generated_query VARCHAR,
			duration_ms BIGINT,
			commit_hash VARCHAR,
			error VARCHAR
		)`,
		`CREATE TABLE IF NOT EXISTS chat_questions (
			id BIGINT,
			timestamp TIMESTAMPTZ DEFAULT now(),
			question VARCHAR,
			source VARCHAR,
			ip_address VARCHAR,
			user_agent VARCHAR,
			is_mobile BOOLEAN,
			os VARCHAR,
			browser VARCHAR,
			country VARCHAR,
			accept_language VARCHAR,
			referer VARCHAR,
			session_id VARCHAR,
			history_length INTEGER,
			model VARCHAR,
			cloudfront BOOLEAN,
			client_timestamp TIMESTAMPTZ,
			answer VARCHAR
		)`,
	}

	for _, ddl := range tables {
		if _, err := duckDB.Exec(ddl); err != nil {
			return fmt.Errorf("create schema: %w", err)
		}
	}

	log.Println("DuckLake analytics schema ready")
	return nil
}

// duckDBAvailable returns true if DuckDB is initialized
func duckDBAvailable() bool {
	return duckDB != nil
}
