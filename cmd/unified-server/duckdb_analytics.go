//go:build duckdb

package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/duckdb/duckdb-go/v2"
)

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

	// Attach DuckLake catalog via PostgreSQL (optional — falls back to in-memory for local dev)
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
	ducklakeOK := false
	if _, err := duckDB.Exec(attachQuery); err != nil {
		log.Printf("Warning: DuckLake attach failed (%v) — falling back to in-memory tables (local dev mode)", err)
	} else {
		log.Printf("DuckLake attached (catalog=PostgreSQL, data=%s)", dataPath)
		if _, err := duckDB.Exec("USE analytics;"); err != nil {
			log.Printf("Warning: USE analytics failed: %v — using in-memory tables", err)
		} else {
			ducklakeOK = true
		}
	}
	if !ducklakeOK {
		log.Println("DuckDB running in-memory only (analytics will not persist across restarts)")
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
			params JSON,
			result_count INTEGER,
			duration_ms DOUBLE,
			client_info VARCHAR,
			created_at TIMESTAMPTZ DEFAULT now()
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
		`CREATE TABLE IF NOT EXISTS chat_feedback (
			chat_id    BIGINT,
			score      INTEGER,
			created_at TIMESTAMPTZ DEFAULT now()
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
		// Semantic cache: stores embeddings + user feedback for past Q&A pairs.
		// embedding is stored as a JSON array of float32 values (VARCHAR).
		`CREATE TABLE IF NOT EXISTS qa_embeddings (
			id BIGINT,
			chat_id BIGINT,
			question VARCHAR,
			answer VARCHAR,
			embedding VARCHAR,
			feedback_score INTEGER DEFAULT 0,
			created_at TIMESTAMPTZ DEFAULT now()
		)`,
		// Curated geographic knowledge: confirmed explanations for elevated/unusual
		// readings at specific locations, auto-populated from positively-rated answers.
		`CREATE TABLE IF NOT EXISTS location_knowledge (
			id BIGINT,
			lat DOUBLE,
			lon DOUBLE,
			radius_m DOUBLE,
			note VARCHAR,
			source_chat_id BIGINT,
			created_at TIMESTAMPTZ DEFAULT now()
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
