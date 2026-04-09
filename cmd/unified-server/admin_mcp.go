package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// adminMCPDataHandler returns JSON data for MCP analytics tables.
// GET /api/admin/mcp/data?table=chat_questions&limit=50&offset=0&sort=timestamp&order=desc&search=...
//
// @Summary     Admin MCP analytics data
// @Description Returns paginated MCP analytics rows for a selected table. For table=mcp_query_log, the columns are tool_name, created_at, duration_ms, result_count, client_info, and params.
// @Tags        admin
// @Produce     json
// @Param       table query string true "Analytics table name" Enums(chat_questions,mcp_query_log,mcp_ai_query_log)
// @Param       limit query int false "Page size"
// @Param       offset query int false "Offset"
// @Param       sort query string false "Sort column (for mcp_query_log, default is created_at)"
// @Param       order query string false "Sort order: asc or desc"
// @Param       search query string false "Search term"
// @Success     200 {object} map[string]interface{} "Analytics rows"
// @Failure     400 {string} string "Invalid table"
// @Failure     503 {string} string "Analytics unavailable"
// @Router      /api/admin/mcp/data [get]
func adminMCPDataHandler(w http.ResponseWriter, r *http.Request) {
	if !duckDBAvailable() {
		writeError(w, http.StatusServiceUnavailable, "Analytics not available")
		return
	}

	tableName := r.URL.Query().Get("table")
	columns, ok := mcpTableColumns[tableName]
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid table name")
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
		switch tableName {
		case "mcp_query_log":
			sortCol = "created_at"
		case "qa_embeddings":
			sortCol = "feedback_score"
		case "location_knowledge":
			sortCol = "created_at"
		default:
			sortCol = "timestamp"
		}
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
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": fmt.Sprintf("DuckDB query failed: %v", err),
			"data":  []interface{}{},
			"total": 0,
		})
		return
	}

	// Fetch data — use LEFT() to work around DuckLake Go driver bug where large
	// inlined VARCHAR values are truncated to a single byte on read.
	// LEFT(col, 100000) forces DuckDB to materialize a new string that the driver reads correctly.
	longTextCols := map[string]bool{"question": true, "answer": true, "generated_query": true, "user_agent": true, "params": true, "note": true}
	virtualCols := map[string]bool{"thumbs_up": true, "thumbs_down": true}

	var dataQuery string
	if tableName == "chat_questions" {
		// chat_questions uses table alias q for the LEFT JOIN
		castCols := make([]string, 0, len(columns))
		for _, col := range columns {
			if virtualCols[col] {
				continue
			} else if longTextCols[col] {
				castCols = append(castCols, fmt.Sprintf("LEFT(q.%s, 100000) AS %s", col, col))
			} else {
				castCols = append(castCols, "q."+col)
			}
		}
		castCols = append(castCols,
			"COALESCE(f.thumbs_up, 0) AS thumbs_up",
			"COALESCE(f.thumbs_down, 0) AS thumbs_down",
		)
		colList := strings.Join(castCols, ", ")
		whereForJoin := whereSQL
		if whereForJoin != "" {
			whereForJoin = strings.ReplaceAll(whereForJoin, "CAST(", "CAST(q.")
		}
		dataQuery = fmt.Sprintf(`
			SELECT %s
			FROM chat_questions q
			LEFT JOIN (
				SELECT chat_id,
				       SUM(CASE WHEN score > 0 THEN 1 ELSE 0 END) AS thumbs_up,
				       SUM(CASE WHEN score < 0 THEN 1 ELSE 0 END) AS thumbs_down
				FROM chat_feedback GROUP BY chat_id
			) f ON f.chat_id = q.id
			%s ORDER BY q.%s %s LIMIT %d OFFSET %d`,
			colList, whereForJoin, sortCol, order, limit, offset)
	} else {
		// All other tables: no alias, use bare column names
		castCols := make([]string, 0, len(columns))
		for _, col := range columns {
			if virtualCols[col] {
				continue
			} else if longTextCols[col] {
				castCols = append(castCols, fmt.Sprintf("LEFT(%s, 100000) AS %s", col, col))
			} else {
				castCols = append(castCols, col)
			}
		}
		colList := strings.Join(castCols, ", ")
		dataQuery = fmt.Sprintf("SELECT %s FROM %s %s ORDER BY %s %s LIMIT %d OFFSET %d",
			colList, tableName, whereSQL, sortCol, order, limit, offset)
	}

	rows, err := duckDB.Query(dataQuery)
	if err != nil {
		log.Printf("admin mcp query error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
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

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    results,
		"total":   total,
		"limit":   limit,
		"offset":  offset,
		"columns": columns,
	})
}

// adminMCPExportHandler exports MCP analytics data as CSV.
// GET /api/admin/mcp/export?table=chat_questions&search=...
//
// @Summary     Admin MCP analytics export
// @Description Exports MCP analytics rows as CSV for a selected table. For table=mcp_query_log, rows are ordered by created_at descending.
// @Tags        admin
// @Produce     text/csv
// @Param       table query string true "Analytics table name" Enums(chat_questions,mcp_query_log,mcp_ai_query_log)
// @Param       search query string false "Search term"
// @Success     200 {file} file "CSV export"
// @Failure     400 {string} string "Invalid table"
// @Failure     503 {string} string "Analytics unavailable"
// @Router      /api/admin/mcp/export [get]
func adminMCPExportHandler(w http.ResponseWriter, r *http.Request) {
	if !duckDBAvailable() {
		writeError(w, http.StatusServiceUnavailable, "Analytics not available")
		return
	}

	tableName := r.URL.Query().Get("table")
	columns, ok := mcpTableColumns[tableName]
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid table name")
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
	exportLongTextCols := map[string]bool{"question": true, "answer": true, "generated_query": true, "user_agent": true, "params": true}
	exportCastCols := make([]string, len(columns))
	for i, col := range columns {
		if exportLongTextCols[col] {
			exportCastCols[i] = fmt.Sprintf("LEFT(%s, 100000) AS %s", col, col)
		} else {
			exportCastCols[i] = col
		}
	}
	colList := strings.Join(exportCastCols, ", ")
	orderCol := "timestamp"
	switch tableName {
	case "mcp_query_log", "location_knowledge":
		orderCol = "created_at"
	case "qa_embeddings":
		orderCol = "feedback_score"
	}
	query := fmt.Sprintf("SELECT %s FROM %s %s ORDER BY %s DESC", colList, tableName, whereSQL, orderCol)

	rows, err := duckDB.Query(query)
	if err != nil {
		log.Printf("admin mcp export error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
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
		"thumbs_up", "thumbs_down",
	},
	"mcp_query_log": {
		"tool_name", "created_at", "duration_ms", "result_count",
		"client_info", "params",
	},
	"mcp_ai_query_log": {
		"user_id", "user_email", "session_id", "timestamp",
		"tool_name", "generated_query", "duration_ms",
		"commit_hash", "error",
	},
	"qa_embeddings": {
		"id", "chat_id", "question", "answer", "feedback_score",
	},
	"location_knowledge": {
		"id", "lat", "lon", "radius_m", "note", "source_chat_id", "created_at",
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
//
// @Summary     Admin MCP analytics delete
// @Description Deletes selected or filtered MCP analytics rows. For table=mcp_query_log, ids map to created_at values.
// @Tags        admin
// @Produce     json
// @Param       table query string true "Analytics table name" Enums(chat_questions,mcp_query_log,mcp_ai_query_log)
// @Param       ids query string false "Comma-separated row IDs"
// @Param       all query boolean false "Delete all filtered rows"
// @Param       search query string false "Search term when all=true"
// @Success     200 {object} map[string]interface{} "Delete result"
// @Failure     400 {string} string "Invalid request"
// @Failure     503 {string} string "Analytics unavailable"
// @Router      /api/admin/mcp/delete [delete]
// @Router      /api/admin/mcp/delete [post]
func adminMCPDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}
	if !duckDBAvailable() {
		writeError(w, http.StatusServiceUnavailable, "Analytics not available")
		return
	}

	tableName := r.URL.Query().Get("table")
	columns, ok := mcpTableColumns[tableName]
	if !ok {
		writeError(w, http.StatusBadRequest, "Invalid table name")
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
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()
		result, err := duckDB.ExecContext(ctx, query)
		if err != nil {
			log.Printf("admin mcp delete all error: %v", err)
			writeError(w, http.StatusInternalServerError, "Delete failed")
			return
		}
		affected, _ := result.RowsAffected()
		writeJSON(w, http.StatusOK, map[string]interface{}{"deleted": affected})
		return
	}

	// Delete specific rows by key values
	ids := r.URL.Query().Get("ids")
	if ids == "" {
		writeError(w, http.StatusBadRequest, "Missing ids parameter")
		return
	}

	// Integer key columns: use numeric IN list (avoids CAST issues in DuckLake
	// and keeps IDs within JS MAX_SAFE_INTEGER range).
	// Timestamp key columns: fall back to string CAST comparison.
	integerKeyCols := map[string]bool{"id": true}

	idList := strings.Split(ids, ",")
	var query string
	if integerKeyCols[keyCol] {
		nums := make([]string, 0, len(idList))
		for _, raw := range idList {
			raw = strings.TrimSpace(raw)
			if _, err := strconv.ParseInt(raw, 10, 64); err == nil {
				nums = append(nums, raw)
			}
		}
		if len(nums) == 0 {
			http.Error(w, "No valid integer IDs provided", http.StatusBadRequest)
			return
		}
		query = fmt.Sprintf("DELETE FROM %s WHERE %s IN (%s)",
			tableName, keyCol, strings.Join(nums, ","))
	} else {
		placeholders := make([]string, len(idList))
		for i, raw := range idList {
			placeholders[i] = "'" + escapeLike(strings.TrimSpace(raw)) + "'"
		}
		query = fmt.Sprintf("DELETE FROM %s WHERE CAST(%s AS VARCHAR) IN (%s)",
			tableName, keyCol, strings.Join(placeholders, ","))
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	result, err := duckDB.ExecContext(ctx, query)
	if err != nil {
		log.Printf("admin mcp delete error: %v", err)
		writeError(w, http.StatusInternalServerError, "Delete failed")
		return
	}
	affected, _ := result.RowsAffected()
	if affected < 0 {
		affected = int64(len(idList)) // DuckDB doesn't always report rows affected
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"deleted": affected})
}

// mcpTableKeyColumn maps each table to its primary key / unique identifier column
var mcpTableKeyColumn = map[string]string{
	"chat_questions":     "id",
	"mcp_query_log":      "created_at",
	"mcp_ai_query_log":   "timestamp",
	"qa_embeddings":      "id",
	"location_knowledge": "id",
}

// mcpEditableColumns defines which columns can be edited per table.
// Only text/content fields are included — IDs, timestamps, and computed
// columns (thumbs_up/thumbs_down) are intentionally excluded.
var mcpEditableColumns = map[string][]string{
	"chat_questions":     {"question", "answer", "source", "model", "country", "browser", "os"},
	"mcp_query_log":      {"params", "client_info"},
	"mcp_ai_query_log":   {"generated_query", "error"},
	"qa_embeddings":      {"question", "answer"},
	"location_knowledge": {"note", "lat", "lon", "radius_m"},
}

// adminMCPUpdateHandler updates editable fields of a single row.
// PUT /api/admin/mcp/update?table=chat_questions&id=<key_value>
// Body: JSON object with field names → new string values.
func adminMCPUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !duckDBAvailable() {
		http.Error(w, "Analytics not available", http.StatusServiceUnavailable)
		return
	}

	tableName := r.URL.Query().Get("table")
	if _, ok := mcpTableColumns[tableName]; !ok {
		http.Error(w, "Invalid table name", http.StatusBadRequest)
		return
	}
	keyVal := r.URL.Query().Get("id")
	if keyVal == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Build SET clause using only whitelisted editable columns.
	var setClauses []string
	var args []interface{}
	for _, col := range mcpEditableColumns[tableName] {
		if val, ok := payload[col]; ok {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", col))
			args = append(args, val)
		}
	}
	if len(setClauses) == 0 {
		http.Error(w, "No editable fields provided", http.StatusBadRequest)
		return
	}

	keyCol := mcpTableKeyColumn[tableName]
	args = append(args, keyVal)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE CAST(%s AS VARCHAR) = ?",
		tableName, strings.Join(setClauses, ", "), keyCol)

	if _, err := duckDB.Exec(query, args...); err != nil {
		log.Printf("admin mcp update error: %v", err)
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"}) //nolint:errcheck
}

func escapeLike(s string) string {
	s = strings.ReplaceAll(s, "'", "''")
	s = strings.ReplaceAll(s, "%", "")
	s = strings.ReplaceAll(s, "_", "")
	return s
}
