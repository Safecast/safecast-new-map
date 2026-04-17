package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// tourStepsPublicHandler returns enabled tour steps with text resolved for
// the requested language. Used by map.html.
//
// @Summary     Public tour steps
// @Description Returns enabled, ordered tour steps with text resolved for the requested language.
// @Tags        tour
// @Produce     json
// @Param       lang query string false "Language code (defaults to 'en')"
// @Success     200 {array}  map[string]interface{} "Tour steps"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/tour/steps [get]
func tourStepsPublicHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = getPreferredLanguage(r)
	}
	steps, err := newTourService().ListPublic(lang)
	if err != nil {
		log.Printf("tour public list error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"steps": steps})
}

// adminTourStepsListHandler returns every row (enabled + disabled) for the
// admin UI.
//
// @Summary     Admin tour steps list
// @Description Returns all tour step rows for admin editing.
// @Tags        admin
// @Produce     json
// @Param       sort  query string false "Sort column"
// @Param       order query string false "Sort order"
// @Success     200 {object} map[string]interface{} "Tour step rows"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/tour/steps [get]
func adminTourStepsListHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	sortCol := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	rows, err := newTourService().List(sortCol, order)
	if err != nil {
		log.Printf("admin tour list error: %v", err)
		writeError(w, http.StatusInternalServerError, "Query failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  rows,
		"total": len(rows),
	})
}

// adminTourStepsCreateHandler inserts a new step at the end of the order.
//
// @Summary     Admin tour step create
// @Description Creates a new tour step row.
// @Tags        admin
// @Accept      json
// @Produce     json
// @Success     201 {object} map[string]interface{} "Create result"
// @Failure     400 {string} string "Invalid request"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/tour/steps [post]
func adminTourStepsCreateHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	var in tourStepInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	id, err := newTourService().Create(in)
	if err != nil {
		log.Printf("admin tour create error: %v", err)
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{"status": "ok", "id": id})
}

// adminTourStepsReorderHandler rewrites sort_order to match a given slice of
// row IDs in a single transaction.
//
// @Summary     Admin tour steps reorder
// @Description Rewrites sort_order to match a supplied ID order.
// @Tags        admin
// @Accept      json
// @Produce     json
// @Success     200 {object} map[string]string "Reorder result"
// @Failure     400 {string} string "Invalid request"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/tour/steps/reorder [post]
func adminTourStepsReorderHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	var payload struct {
		IDs []int64 `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	if err := newTourService().Reorder(payload.IDs); err != nil {
		log.Printf("admin tour reorder error: %v", err)
		writeError(w, http.StatusInternalServerError, "Reorder failed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// adminTourStepsByIDHandler routes PUT/DELETE on /api/admin/tour/steps/{id}.
//
// @Summary     Admin tour step update/delete
// @Description Update (PUT) or delete (DELETE) a tour step row by ID.
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       id path int true "Tour step ID"
// @Success     200 {object} map[string]string "Result"
// @Failure     400 {string} string "Invalid request"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/admin/tour/steps/{id} [put]
func adminTourStepsByIDHandler(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		writeError(w, http.StatusServiceUnavailable, "Database not available")
		return
	}
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	idStr := parts[len(parts)-1]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid step ID")
		return
	}

	switch r.Method {
	case http.MethodPut:
		var in tourStepInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}
		if err := newTourService().Update(id, in); err != nil {
			log.Printf("admin tour update error: %v", err)
			writeError(w, http.StatusInternalServerError, "Update failed")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	case http.MethodDelete:
		if err := newTourService().Delete(id); err != nil {
			log.Printf("admin tour delete error: %v", err)
			writeError(w, http.StatusInternalServerError, "Delete failed")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
