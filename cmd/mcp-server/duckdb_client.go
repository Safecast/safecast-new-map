package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
)

var duckDB *sql.DB

func initDuckDB() error {
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
	pgURL := os.Getenv("DATABASE_URL")
	if pgURL != "" {
		query := fmt.Sprintf(
			"ATTACH '%s' AS postgres_db (TYPE POSTGRES, READ_ONLY)",
			pgURL,
		)
		if _, err := duckDB.Exec(query); err != nil {
			log.Printf("Warning: failed to attach postgres: %v", err)
		} else {
			log.Println("Safecast PostgreSQL attached as postgres_db")
		}
	}

	// Create analytics schema in DuckLake
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
			user_id TEXT,
			user_email TEXT,
			session_id TEXT,
			timestamp TIMESTAMP,
			tool_name TEXT,
			generated_query TEXT,
			duration_ms BIGINT,
			commit_hash TEXT,
			error TEXT
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
			log.Printf("Warning: DuckLake DDL failed: %v (table may already exist)", err)
		}
	}

	log.Println("DuckLake analytics schema ready")
	return nil
}

// LogQueryAsync logs a tool execution to DuckLake asynchronously.
func LogQueryAsync(toolName string, params map[string]any, resultCount int, duration time.Duration, clientInfo string) {
	if duckDB == nil {
		return
	}

	go func() {
		paramsJSON, err := json.Marshal(params)
		if err != nil {
			log.Printf("Error marshaling params to JSON: %v", err)
			return
		}
		paramsStr := string(paramsJSON)

		_, execErr := duckDB.Exec(`
			INSERT INTO mcp_query_log (tool_name, params, result_count, duration_ms, client_info)
			VALUES (?, ?, ?, ?, ?)
		`, toolName, paramsStr, resultCount, float64(duration.Milliseconds()), clientInfo)

		if execErr != nil {
			log.Printf("Error logging query to DuckLake: %v", execErr)
		}
	}()
}

// Analytics Functions

// GetToolUsageStats returns usage statistics for tools.
func GetToolUsageStats() ([]map[string]any, error) {
	if duckDB == nil {
		return nil, fmt.Errorf("duckdb not initialized")
	}

	rows, err := duckDB.Query(`
		SELECT tool_name, COUNT(*) as calls, AVG(duration_ms) as avg_duration, MAX(duration_ms) as max_duration
		FROM mcp_query_log
		GROUP BY tool_name
		ORDER BY calls DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []map[string]any
	for rows.Next() {
		var toolName string
		var calls int64
		var avgDur, maxDur float64
		if err := rows.Scan(&toolName, &calls, &avgDur, &maxDur); err != nil {
			return nil, err
		}
		stats = append(stats, map[string]any{
			"tool":   toolName,
			"calls":  calls,
			"avg_ms": avgDur,
			"max_ms": maxDur,
		})
	}
	return stats, nil
}

// QueryPostgresAnalytics executes an analytical query on the attached Postgres DB.
func QueryPostgresAnalytics(query string, args ...any) ([]map[string]any, error) {
	if duckDB == nil {
		return nil, fmt.Errorf("duckdb not initialized")
	}

	rows, err := duckDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]any
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		row := make(map[string]any)
		for i, colName := range cols {
			val := columns[i]
			if b, ok := val.([]byte); ok {
				row[colName] = string(b)
			} else {
				row[colName] = val
			}
		}
		results = append(results, row)
	}
	return results, nil
}
