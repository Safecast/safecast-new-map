package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// translationColumns defines the valid columns for the translations table.
var translationColumns = []string{
	"id", "language_code", "key", "value", "updated_at",
}

// adminTranslationsDataHandler returns JSON data for the admin translations page.
// GET /api/admin/translations?limit=50&offset=0&sort=key&order=asc&search=...&lang=...
func adminTranslationsDataHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
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
	if !isValidColumn(sortCol, translationColumns) {
		sortCol = "key"
	}
	order := strings.ToUpper(r.URL.Query().Get("order"))
	if order != "DESC" {
		order = "ASC"
	}

	search := r.URL.Query().Get("search")
	langFilter := r.URL.Query().Get("lang")

	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIdx := 1

	if langFilter != "" {
		conditions = append(conditions, fmt.Sprintf("language_code = $%d", argIdx))
		args = append(args, langFilter)
		argIdx++
	}
	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(key ILIKE $%d OR value ILIKE $%d)", argIdx, argIdx+1))
		args = append(args, "%"+search+"%", "%"+search+"%")
		argIdx += 2
	}

	whereSQL := ""
	if len(conditions) > 0 {
		whereSQL = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM translations %s", whereSQL)
	var total int
	if err := db.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		log.Printf("admin translations count error: %v", err)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("Query failed: %v", err),
			"data":  []interface{}{},
			"total": 0,
		})
		return
	}

	// Fetch data
	dataQuery := fmt.Sprintf("SELECT id, language_code, key, value, updated_at FROM translations %s ORDER BY %s %s LIMIT %d OFFSET %d",
		whereSQL, sortCol, order, limit, offset)

	rows, err := db.DB.Query(dataQuery, args...)
	if err != nil {
		log.Printf("admin translations query error: %v", err)
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id int64
		var langCode, key, value string
		var updatedAt interface{}
		if err := rows.Scan(&id, &langCode, &key, &value, &updatedAt); err != nil {
			log.Printf("admin translations scan error: %v", err)
			continue
		}
		row := map[string]interface{}{
			"id":            id,
			"language_code": langCode,
			"key":           key,
			"value":         value,
			"updated_at":    updatedAt,
		}
		results = append(results, row)
	}

	if results == nil {
		results = []map[string]interface{}{}
	}

	// Get available languages for filter dropdown
	langRows, err := db.DB.Query("SELECT DISTINCT language_code FROM translations ORDER BY language_code")
	var languages []string
	if err == nil {
		defer langRows.Close()
		for langRows.Next() {
			var lang string
			if langRows.Scan(&lang) == nil {
				languages = append(languages, lang)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":      results,
		"total":     total,
		"languages": languages,
	})
}

// adminTranslationUpdateHandler updates a single translation.
// PUT /api/admin/translations/{id}
func adminTranslationUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	// Extract ID from path: /api/admin/translations/123
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.Error(w, "Missing translation ID", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid translation ID", http.StatusBadRequest)
		return
	}

	var payload struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("UPDATE translations SET value = $1, updated_at = NOW() WHERE id = $2", payload.Value, id)
	if err != nil {
		log.Printf("admin translation update error: %v", err)
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// adminTranslationCreateHandler creates a new translation.
// POST /api/admin/translations
func adminTranslationCreateHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	var payload struct {
		LanguageCode string `json:"language_code"`
		Key          string `json:"key"`
		Value        string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if payload.LanguageCode == "" || payload.Key == "" {
		http.Error(w, "language_code and key are required", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec(
		"INSERT INTO translations (language_code, key, value) VALUES ($1, $2, $3) ON CONFLICT (language_code, key) DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()",
		payload.LanguageCode, payload.Key, payload.Value,
	)
	if err != nil {
		log.Printf("admin translation create error: %v", err)
		http.Error(w, "Create failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// adminTranslationDeleteHandler deletes a translation.
// DELETE /api/admin/translations/{id}
func adminTranslationDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}

	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		http.Error(w, "Missing translation ID", http.StatusBadRequest)
		return
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid translation ID", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("DELETE FROM translations WHERE id = $1", id)
	if err != nil {
		log.Printf("admin translation delete error: %v", err)
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// adminTranslationsReloadHandler reloads translations from DB into memory.
// POST /api/admin/translations/reload
func adminTranslationsReloadHandler(w http.ResponseWriter, r *http.Request) {
	if loadTranslationsFromDB() {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Translations reloaded from database"})
	} else {
		http.Error(w, "Failed to reload translations", http.StatusInternalServerError)
	}
}
