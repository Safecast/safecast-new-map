// DuckDB Analytics Initialization for Unified Server
// DuckDB is used for analytics queries, attaching to PostgreSQL for read-only access

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/marcboeker/go-duckdb"
)

var duckDB *sql.DB

// initDuckDBAnalytics initializes DuckDB for analytics queries
// It attaches to the PostgreSQL database for read-only analytics
func initDuckDBAnalytics() error {
	// Get DuckDB file path from environment
	duckPath := os.Getenv("DUCKDB_PATH")
	if duckPath == "" {
		duckPath = "./analytics.duckdb"
	}

	var err error
	duckDB, err = sql.Open("duckdb", duckPath+"?access_mode=READ_WRITE")
	if err != nil {
		return fmt.Errorf("failed to open duckdb: %w", err)
	}

	// Configure connection pool for DuckDB
	duckDB.SetMaxOpenConns(1)  // Single writer
	duckDB.SetMaxIdleConns(1)
	duckDB.SetConnMaxLifetime(0)

	if err := duckDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping duckdb: %w", err)
	}

	log.Printf("DuckDB initialized at %s", duckPath)

	// Enable WAL for durability
	duckDB.Exec("PRAGMA wal_autocheckpoint=1000;")

	// Try to load postgres extension and attach to PostgreSQL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		if err := attachPostgres(databaseURL); err != nil {
			log.Printf("Warning: PostgreSQL attach failed: %v (analytics will use local DuckDB only)", err)
		} else {
			log.Println("DuckDB attached to PostgreSQL for analytics queries")
		}
	}

	// Create analytics schema
	if err := createDuckDBSchema(); err != nil {
		log.Printf("Warning: failed to create DuckDB schema: %v", err)
	}

	return nil
}

// attachPostgres attaches DuckDB to PostgreSQL for read-only analytics queries
func attachPostgres(databaseURL string) error {
	// Install and load postgres extension
	if _, err := duckDB.Exec("INSTALL postgres;"); err != nil {
		return fmt.Errorf("install postgres extension: %w", err)
	}
	if _, err := duckDB.Exec("LOAD postgres;"); err != nil {
		return fmt.Errorf("load postgres extension: %w", err)
	}

	// Parse PostgreSQL URL and create attachment string
	// Format: postgres://user:pass@host:port/dbname?sslmode=...
	// Need to extract query params for proper attachment
	attachStr := databaseURL
	if !strings.Contains(databaseURL, "?") {
		attachStr = databaseURL + "?sslmode=prefer"
	}

	createSchemaQuery := fmt.Sprintf("ATTACH '%s' AS postgres (TYPE POSTGRES, READ_ONLY);", attachStr)

	if _, err := duckDB.Exec(createSchemaQuery); err != nil {
		return fmt.Errorf("attach postgres: %w", err)
	}

	return nil
}

// createDuckDBSchema creates the local analytics tables
func createDuckDBSchema() error {
	createSchemaQuery := `
		CREATE TABLE IF NOT EXISTS mcp_query_log (
			tool_name VARCHAR,
			timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			duration_ms BIGINT,
			result_count INTEGER,
			client VARCHAR,
			user_id VARCHAR,
			user_email VARCHAR
		);
		
		CREATE TABLE IF NOT EXISTS mcp_ai_query_log (
			user_id VARCHAR,
			user_email VARCHAR,
			session_id VARCHAR,
			timestamp TIMESTAMP,
			tool_name VARCHAR,
			generated_query VARCHAR,
			duration_ms BIGINT,
			commit_hash VARCHAR,
			error VARCHAR
		);

		CREATE SEQUENCE IF NOT EXISTS seq_chat_questions START 1;
		CREATE TABLE IF NOT EXISTS chat_questions (
			id BIGINT DEFAULT nextval('seq_chat_questions'),
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
			cloudfront BOOLEAN
		);
	`

	if _, err := duckDB.Exec(createSchemaQuery); err != nil {
		return fmt.Errorf("create schema: %w", err)
	}

	// Create indexes for common queries
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_query_log_tool ON mcp_query_log(tool_name);",
		"CREATE INDEX IF NOT EXISTS idx_query_log_timestamp ON mcp_query_log(timestamp);",
		"CREATE INDEX IF NOT EXISTS idx_ai_log_tool ON mcp_ai_query_log(tool_name);",
		"CREATE INDEX IF NOT EXISTS idx_ai_log_timestamp ON mcp_ai_query_log(timestamp);",
		"CREATE INDEX IF NOT EXISTS idx_chat_q_timestamp ON chat_questions(timestamp);",
		"CREATE INDEX IF NOT EXISTS idx_chat_q_source ON chat_questions(source);",
	}

	for _, idx := range indexes {
		duckDB.Exec(idx)
	}

	log.Println("DuckDB schema ready")
	return nil
}

// duckDBAvailable returns true if DuckDB is initialized
func duckDBAvailable() bool {
	return duckDB != nil
}
