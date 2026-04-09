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
//
// @Summary     Admin translations list
// @Description Returns paginated translation rows for admin editing.
// @Tags        admin
// @Produce     json
// @Param       limit query int false "Page size"
// @Param       offset query int false "Offset"
// @Param       sort query string false "Sort column"
// @Param       order query string false "Sort order"
// @Param       search query string false "Search term"
// @Param       lang query string false "Language code filter"
// @Success     200 {object} map[string]interface{} "Translation rows"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/translations [get]
func adminTranslationsDataHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
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
	svc := newTranslationService()
	listResult, err := svc.List(limit, offset, sortCol, order, search, langFilter)
	if err != nil {
		log.Printf("admin translations list error: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]interface{}{
			"error": fmt.Sprintf("Query failed: %v", err),
			"data":  []interface{}{},
			"total": 0,
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":      listResult.Rows,
		"total":     listResult.Total,
		"languages": listResult.Languages,
	})
}

// adminTranslationUpdateHandler updates a single translation.
// PUT /api/admin/translations/{id}
//
// @Summary     Admin translation update
// @Description Updates translation value by translation row ID.
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       id path int true "Translation row ID"
// @Success     200 {object} map[string]string "Update result"
// @Failure     400 {string} string "Invalid request"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/translations/{id} [put]
func adminTranslationUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}

	// Extract ID from path: /api/admin/translations/123
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "Missing translation ID")
		return
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid translation ID")
		return
	}

	var payload struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if err := newTranslationService().Update(id, payload.Value); err != nil {
		log.Printf("admin translation update error: %v", err)
		writeError(w, http.StatusInternalServerError, "Update failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// adminTranslationCreateHandler creates a new translation.
// POST /api/admin/translations
//
// @Summary     Admin translation create
// @Description Creates a translation row or updates existing key/language pair.
// @Tags        admin
// @Accept      json
// @Produce     json
// @Success     201 {object} map[string]string "Create result"
// @Failure     400 {string} string "Invalid request"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/translations [post]
func adminTranslationCreateHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}

	var payload struct {
		LanguageCode string `json:"language_code"`
		Key          string `json:"key"`
		Value        string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if payload.LanguageCode == "" || payload.Key == "" {
		writeError(w, http.StatusBadRequest, "language_code and key are required")
		return
	}

	if err := newTranslationService().Create(payload.LanguageCode, payload.Key, payload.Value); err != nil {
		log.Printf("admin translation create error: %v", err)
		writeError(w, http.StatusInternalServerError, "Create failed")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"status": "ok"})
}

// adminTranslationDeleteHandler deletes a translation.
// DELETE /api/admin/translations/{id}
//
// @Summary     Admin translation delete
// @Description Deletes translation row by ID.
// @Tags        admin
// @Produce     json
// @Param       id path int true "Translation row ID"
// @Success     200 {object} map[string]string "Delete result"
// @Failure     400 {string} string "Invalid request"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/translations/{id} [delete]
func adminTranslationDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}

	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		writeError(w, http.StatusBadRequest, "Missing translation ID")
		return
	}
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid translation ID")
		return
	}

	if err := newTranslationService().Delete(id); err != nil {
		log.Printf("admin translation delete error: %v", err)
		writeError(w, http.StatusInternalServerError, "Delete failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// adminTranslationsReloadHandler reloads translations from DB into memory.
// POST /api/admin/translations/reload
//
// @Summary     Admin translations reload
// @Description Reloads in-memory translations from database.
// @Tags        admin
// @Produce     json
// @Success     200 {object} map[string]string "Reload result"
// @Failure     500 {string} string "Reload failed"
// @Router      /api/admin/translations/reload [post]
func adminTranslationsReloadHandler(w http.ResponseWriter, r *http.Request) {
	if loadTranslationsFromDB() {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "message": "Translations reloaded from database"})
	} else {
		writeError(w, http.StatusInternalServerError, "Failed to reload translations")
	}
}
