// Admin handlers for the /admin/qa-embeddings tab.
//
// Routes (wired in main.go):
//   GET  /api/admin/qa-embeddings                  → searchable/sortable list
//   GET  /api/admin/qa-embeddings/{id}             → full row including answer
//   POST /api/admin/qa-embeddings/{id}/promote     → feedback_score += 1, status='active'
//   POST /api/admin/qa-embeddings/{id}/demote      → feedback_score -= 1, status='demoted'
//   POST /api/admin/qa-embeddings/{id}/archive     → status='archived' (audit-only)
//   POST /api/admin/qa-embeddings/{id}/restore     → status='active'
//
// The cache lives in DuckLake (in-process DuckDB attached to the catalog), so
// we query the package-global duckDB handle, not the PostgreSQL `db` handle
// used by other admin pages.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// qaEmbeddingsColumns whitelists columns that may appear in ?sort=.
// `embedding` is intentionally excluded — it's neither sortable nor sent to
// the client (large, opaque to humans).
var qaEmbeddingsColumns = []string{
	"id", "chat_id", "question", "feedback_score",
	"used_count", "last_used_at", "status", "lang", "created_at",
}

// adminQAEmbeddingsListHandler serves the searchable/sortable JSON listing.
// GET /api/admin/qa-embeddings?search=&status=&lang=&score=&sort=&order=&limit=&offset=
func adminQAEmbeddingsListHandler(w http.ResponseWriter, r *http.Request) {
	if !duckDBAvailable() {
		http.Error(w, "Analytics database not available", http.StatusServiceUnavailable)
		return
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}

	sortCol := r.URL.Query().Get("sort")
	if !isValidColumn(sortCol, qaEmbeddingsColumns) {
		sortCol = "created_at"
	}
	order := strings.ToUpper(r.URL.Query().Get("order"))
	if order != "ASC" {
		order = "DESC"
	}

	search := r.URL.Query().Get("search")
	statusFilter := r.URL.Query().Get("status")
	langFilter := r.URL.Query().Get("lang")
	// score filter: "pos" (>0), "neg" (<0), "zero" (=0); empty = all
	scoreFilter := r.URL.Query().Get("score")

	// DuckDB's Go driver uses `?` placeholders here (matches the pattern used
	// elsewhere in semantic_cache.go).
	var conditions []string
	var args []interface{}

	if search != "" {
		conditions = append(conditions, "(question ILIKE ? OR answer ILIKE ?)")
		args = append(args, "%"+search+"%", "%"+search+"%")
	}
	if statusFilter != "" {
		conditions = append(conditions, "COALESCE(status, 'active') = ?")
		args = append(args, statusFilter)
	}
	if langFilter != "" {
		conditions = append(conditions, "COALESCE(lang, '') = ?")
		args = append(args, langFilter)
	}
	switch scoreFilter {
	case "pos":
		conditions = append(conditions, "COALESCE(feedback_score, 0) > 0")
	case "neg":
		conditions = append(conditions, "COALESCE(feedback_score, 0) < 0")
	case "zero":
		conditions = append(conditions, "COALESCE(feedback_score, 0) = 0")
	}

	whereSQL := ""
	if len(conditions) > 0 {
		whereSQL = "WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := duckDB.QueryRow(
		fmt.Sprintf("SELECT COUNT(*) FROM qa_embeddings %s", whereSQL),
		args...,
	).Scan(&total); err != nil {
		log.Printf("admin qa-embeddings count error: %v", err)
		writeQAJSONError(w, err)
		return
	}

	dataSQL := fmt.Sprintf(`
		SELECT id, chat_id, question, feedback_score,
		       COALESCE(used_count, 0)    AS used_count,
		       last_used_at,
		       COALESCE(status, 'active') AS status,
		       COALESCE(lang, '')         AS lang,
		       created_at,
		       LENGTH(answer)             AS answer_len
		FROM qa_embeddings %s
		ORDER BY %s %s NULLS LAST
		LIMIT %d OFFSET %d`,
		whereSQL, sortCol, order, limit, offset,
	)

	rows, err := duckDB.Query(dataSQL, args...)
	if err != nil {
		log.Printf("admin qa-embeddings query error: %v", err)
		writeQAJSONError(w, err)
		return
	}
	defer rows.Close()

	results := []map[string]interface{}{}
	for rows.Next() {
		var (
			id, chatID               int64
			question, status, lang   string
			feedbackScore, usedCount int
			lastUsedAt, createdAt    interface{}
			answerLen                int
		)
		if err := rows.Scan(&id, &chatID, &question, &feedbackScore, &usedCount, &lastUsedAt, &status, &lang, &createdAt, &answerLen); err != nil {
			log.Printf("admin qa-embeddings scan error: %v", err)
			continue
		}
		results = append(results, map[string]interface{}{
			"id":             id,
			"chat_id":        chatID,
			"question":       question,
			"feedback_score": feedbackScore,
			"used_count":     usedCount,
			"last_used_at":   lastUsedAt,
			"status":         status,
			"lang":           lang,
			"created_at":     createdAt,
			"answer_len":     answerLen,
		})
	}

	languages := distinctQAStrings("SELECT DISTINCT COALESCE(lang, '') FROM qa_embeddings ORDER BY 1")
	statuses := distinctQAStrings("SELECT DISTINCT COALESCE(status, 'active') FROM qa_embeddings ORDER BY 1")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":      results,
		"total":     total,
		"languages": languages,
		"statuses":  statuses,
	})
}

// distinctQAStrings returns one column of distinct strings; empty slice on error.
func distinctQAStrings(sqlStr string) []string {
	rows, err := duckDB.Query(sqlStr)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if rows.Scan(&s) == nil {
			out = append(out, s)
		}
	}
	return out
}

// adminQAEmbeddingsDetailHandler returns the full row (including the answer
// text) for the modal in the admin tab.
// GET /api/admin/qa-embeddings/{id}
func adminQAEmbeddingsDetailHandler(w http.ResponseWriter, r *http.Request, id int64) {
	if !duckDBAvailable() {
		http.Error(w, "Analytics database not available", http.StatusServiceUnavailable)
		return
	}

	row := duckDB.QueryRow(`
		SELECT id, chat_id, question, answer, feedback_score,
		       COALESCE(used_count, 0), last_used_at,
		       COALESCE(status, 'active'), COALESCE(lang, ''), created_at
		FROM qa_embeddings WHERE id = ?`, id)

	var (
		dbID, chatID             int64
		question, answer         string
		status, lang             string
		feedbackScore, usedCount int
		lastUsedAt, createdAt    interface{}
	)
	if err := row.Scan(&dbID, &chatID, &question, &answer, &feedbackScore, &usedCount, &lastUsedAt, &status, &lang, &createdAt); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":             dbID,
		"chat_id":        chatID,
		"question":       question,
		"answer":         answer,
		"feedback_score": feedbackScore,
		"used_count":     usedCount,
		"last_used_at":   lastUsedAt,
		"status":         status,
		"lang":           lang,
		"created_at":     createdAt,
	})
}

// adminQAEmbeddingsActionHandler applies promote/demote/archive/restore.
// POST /api/admin/qa-embeddings/{id}/{action}
func adminQAEmbeddingsActionHandler(w http.ResponseWriter, r *http.Request, id int64, action string) {
	if !duckDBAvailable() {
		http.Error(w, "Analytics database not available", http.StatusServiceUnavailable)
		return
	}

	var query string
	switch action {
	case "promote":
		query = `UPDATE qa_embeddings
		         SET feedback_score = COALESCE(feedback_score, 0) + 1,
		             status = 'active'
		         WHERE id = ?`
	case "demote":
		query = `UPDATE qa_embeddings
		         SET feedback_score = COALESCE(feedback_score, 0) - 1,
		             status = 'demoted'
		         WHERE id = ?`
	case "archive":
		query = `UPDATE qa_embeddings SET status = 'archived' WHERE id = ?`
	case "restore":
		query = `UPDATE qa_embeddings SET status = 'active' WHERE id = ?`
	default:
		http.Error(w, "Unknown action", http.StatusBadRequest)
		return
	}

	res, err := duckDB.Exec(query, id)
	if err != nil {
		log.Printf("admin qa-embeddings %s error: %v", action, err)
		writeQAJSONError(w, err)
		return
	}
	rowsAffected, _ := res.RowsAffected()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":            true,
		"action":        action,
		"id":            id,
		"rows_affected": rowsAffected,
	})
}

func writeQAJSONError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": fmt.Sprintf("Query failed: %v", err),
	})
}
