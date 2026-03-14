package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// adminMCPDataHandler returns JSON data for MCP analytics tables.
// GET /api/admin/mcp/data?table=chat_questions&limit=50&offset=0&sort=timestamp&order=desc&search=...
func adminMCPDataHandler(w http.ResponseWriter, r *http.Request) {
	if !duckDBAvailable() {
		http.Error(w, "Analytics not available", http.StatusServiceUnavailable)
		return
	}

	tableName := r.URL.Query().Get("table")
	columns, ok := mcpTableColumns[tableName]
	if !ok {
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 1000 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	sortCol := r.URL.Query().Get("sort")
	if !isValidColumn(sortCol, columns) {
		sortCol = "timestamp"
	}
	order := strings.ToUpper(r.URL.Query().Get("order"))
	if order != "ASC" {
		order = "DESC"
	}

	search := r.URL.Query().Get("search")

	// Build WHERE clause for search
	var whereClauses []string
	if search != "" {
		for _, col := range columns {
			whereClauses = append(whereClauses, fmt.Sprintf("CAST(%s AS VARCHAR) ILIKE '%%%s%%'", col, escapeLike(search)))
		}
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " OR ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s %s", tableName, whereSQL)
	var total int
	if err := duckDB.QueryRow(countQuery).Scan(&total); err != nil {
		log.Printf("admin mcp count error (table=%s): %v", tableName, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("DuckDB query failed: %v", err),
			"data":  []interface{}{},
			"total": 0,
		})
		return
	}

	// Fetch data — use LEFT() to work around DuckLake Go driver bug where large
	// inlined VARCHAR values are truncated to a single byte on read.
	// LEFT(col, 100000) forces DuckDB to materialize a new string that the driver reads correctly.
	longTextCols := map[string]bool{"question": true, "answer": true, "generated_query": true, "user_agent": true}
	castCols := make([]string, len(columns))
	for i, col := range columns {
		if longTextCols[col] {
			castCols[i] = fmt.Sprintf("LEFT(%s, 100000) AS %s", col, col)
		} else {
			castCols[i] = col
		}
	}
	colList := strings.Join(castCols, ", ")
	dataQuery := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY %s %s LIMIT %d OFFSET %d",
		colList, tableName, whereSQL, sortCol, order, limit, offset)

	rows, err := duckDB.Query(dataQuery)
	if err != nil {
		log.Printf("admin mcp query error: %v", err)
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			log.Printf("admin mcp scan error: %v", err)
			continue
		}
		row := make(map[string]interface{})
		for i, col := range columns {
			// DuckDB Go driver may return []byte for VARCHAR columns from DuckLake;
			// convert to string so JSON encoding works correctly.
			if b, ok := values[i].([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = values[i]
			}
		}
		results = append(results, row)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":    results,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"columns": columns,
	})
}

// adminMCPExportHandler exports MCP analytics data as CSV.
// GET /api/admin/mcp/export?table=chat_questions&search=...
func adminMCPExportHandler(w http.ResponseWriter, r *http.Request) {
	if !duckDBAvailable() {
		http.Error(w, "Analytics not available", http.StatusServiceUnavailable)
		return
	}

	tableName := r.URL.Query().Get("table")
	columns, ok := mcpTableColumns[tableName]
	if !ok {
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	search := r.URL.Query().Get("search")

	var whereClauses []string
	if search != "" {
		for _, col := range columns {
			whereClauses = append(whereClauses, fmt.Sprintf("CAST(%s AS VARCHAR) ILIKE '%%%s%%'", col, escapeLike(search)))
		}
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " OR ")
	}

	// Use LEFT() workaround for long text columns (same DuckLake driver bug as data handler)
	exportLongTextCols := map[string]bool{"question": true, "answer": true, "generated_query": true, "user_agent": true}
	exportCastCols := make([]string, len(columns))
	for i, col := range columns {
		if exportLongTextCols[col] {
			exportCastCols[i] = fmt.Sprintf("LEFT(%s, 100000) AS %s", col, col)
		} else {
			exportCastCols[i] = col
		}
	}
	colList := strings.Join(exportCastCols, ", ")
	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY timestamp DESC", colList, tableName, whereSQL)

	rows, err := duckDB.Query(query)
	if err != nil {
		log.Printf("admin mcp export error: %v", err)
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.csv", tableName))

	writer := csv.NewWriter(w)
	writer.Write(columns) // header row

	for rows.Next() {
		values := make([]interface{}, len(columns))
		ptrs := make([]interface{}, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			continue
		}
		record := make([]string, len(columns))
		for i, v := range values {
			if v == nil {
				record[i] = ""
			} else if b, ok := v.([]byte); ok {
				record[i] = string(b)
			} else {
				record[i] = fmt.Sprintf("%v", v)
			}
		}
		writer.Write(record)
	}
	writer.Flush()
}

// mcpTableColumns defines the valid tables and their columns for the admin MCP page.
var mcpTableColumns = map[string][]string{
	"chat_questions": {
		"id", "timestamp", "question", "answer", "source", "model",
		"ip_address", "country", "is_mobile", "os", "browser",
		"user_agent", "accept_language", "referer",
		"session_id", "history_length", "cloudfront", "client_timestamp",
	},
	"mcp_query_log": {
		"tool_name", "timestamp", "duration_ms", "result_count",
		"client", "user_id", "user_email",
	},
	"mcp_ai_query_log": {
		"user_id", "user_email", "session_id", "timestamp",
		"tool_name", "generated_query", "duration_ms",
		"commit_hash", "error",
	},
}

func isValidColumn(col string, validCols []string) bool {
	for _, c := range validCols {
		if c == col {
			return true
		}
	}
	return false
}

// adminMCPDeleteHandler deletes rows from MCP analytics tables.
// DELETE /api/admin/mcp/delete?table=chat_questions&ids=123,456,789
// DELETE /api/admin/mcp/delete?table=chat_questions&all=true&search=...
func adminMCPDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !duckDBAvailable() {
		http.Error(w, "Analytics not available", http.StatusServiceUnavailable)
		return
	}

	tableName := r.URL.Query().Get("table")
	columns, ok := mcpTableColumns[tableName]
	if !ok {
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}

	// Determine the ID/key column for each table
	keyCol := mcpTableKeyColumn[tableName]

	if r.URL.Query().Get("all") == "true" {
		// Delete all (optionally filtered by search)
		search := r.URL.Query().Get("search")
		var whereClauses []string
		if search != "" {
			for _, col := range columns {
				whereClauses = append(whereClauses, fmt.Sprintf("CAST(%s AS VARCHAR) ILIKE '%%%s%%'", col, escapeLike(search)))
			}
		}
		whereSQL := ""
		if len(whereClauses) > 0 {
			whereSQL = "WHERE " + strings.Join(whereClauses, " OR ")
		}

		query := fmt.Sprintf("DELETE FROM %s %s", tableName, whereSQL)
		result, err := duckDB.Exec(query)
		if err != nil {
			log.Printf("admin mcp delete all error: %v", err)
			http.Error(w, "Delete failed", http.StatusInternalServerError)
			return
		}
		affected, _ := result.RowsAffected()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{"deleted": affected})
		return
	}

	// Delete specific rows by key values
	ids := r.URL.Query().Get("ids")
	if ids == "" {
		http.Error(w, "Missing ids parameter", http.StatusBadRequest)
		return
	}

	idList := strings.Split(ids, ",")
	placeholders := make([]string, len(idList))
	for i := range idList {
		placeholders[i] = "'" + escapeLike(strings.TrimSpace(idList[i])) + "'"
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE CAST(%s AS VARCHAR) IN (%s)",
		tableName, keyCol, strings.Join(placeholders, ","))
	result, err := duckDB.Exec(query)
	if err != nil {
		log.Printf("admin mcp delete error: %v", err)
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}
	affected, _ := result.RowsAffected()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"deleted": affected})
}

// mcpTableKeyColumn maps each table to its primary key / unique identifier column
var mcpTableKeyColumn = map[string]string{
	"chat_questions":   "id",
	"mcp_query_log":    "timestamp",
	"mcp_ai_query_log": "timestamp",
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "%", "")
	s = strings.ReplaceAll(s, "_", "")
	return s
}
