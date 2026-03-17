// MCP database helper functions

package main

import (
	"context"
	"database/sql"
	"time"
)

// queryRows executes a query on PostgreSQL and returns results as a slice of maps
func queryRows(ctx context.Context, query string, args ...any) ([]map[string]any, error) {
	if db == nil || db.DB == nil {
		return nil, sql.ErrNoRows
	}

	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		valuePtrs := make([]any, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(columns))
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return results, rows.Err()
}

// queryRow executes a query on PostgreSQL and returns a single row as a map
func queryRow(ctx context.Context, query string, args ...any) (map[string]any, error) {
	rows, err := queryRows(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, sql.ErrNoRows
	}
	return rows[0], nil
}

// dbAvailable returns true if the database connection is available
func dbAvailable() bool {
	return db != nil && db.DB != nil
}

// LogQueryAsync logs MCP tool usage to DuckDB analytics
func LogQueryAsync(name string, args map[string]any, resultCount int, duration time.Duration, client string) {
	if duckDB == nil {
		return
	}

	// Extract user info if available
	userID := ""
	userEmail := ""
	if args != nil {
		if v, ok := args["user_id"].(string); ok {
			userID = v
		}
		if v, ok := args["user_email"].(string); ok {
			userEmail = v
		}
	}

	_, err := duckDB.Exec(`
		INSERT INTO mcp_query_log (tool_name, duration_ms, result_count, client, user_id, user_email)
		VALUES (?, ?, ?, ?, ?, ?)
	`, name, duration.Milliseconds(), resultCount, client, userID, userEmail)

	if err != nil {
		// Silently ignore logging errors
		return
	}
}
