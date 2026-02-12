package safecast

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"safecast-new-map/pkg/database"
)

func parseBgeigieFilters(r *http.Request) database.BgeigieImportFilters {
	f := database.BgeigieImportFilters{PerPage: 25, Page: 1}
	f.Q = r.URL.Query().Get("q")
	f.Status = r.URL.Query().Get("status")
	if v := r.URL.Query().Get("by_user_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			f.ByUserID = &id
		}
	}
	if v := r.URL.Query().Get("uploaded_after"); v != "" {
		if t, err := parseTime(v); err == nil {
			f.UploadedAfter = &t
		}
	}
	if v := r.URL.Query().Get("uploaded_before"); v != "" {
		if t, err := parseTime(v); err == nil {
			f.UploadedBefore = &t
		}
	}
	if v := r.URL.Query().Get("page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p >= 1 {
			f.Page = p
		}
	}
	if v := r.URL.Query().Get("per_page"); v != "" {
		if p, err := strconv.Atoi(v); err == nil && p > 0 {
			f.PerPage = p
		}
	}
	return f
}

func bgeigieRowToRails(r database.BgeigieImportRow) BgeigieImportRails {
	out := BgeigieImportRails{
		ID:                r.ID,
		UserID:            r.UserID,
		Approved:          r.Approved,
		MeasurementsCount: r.MeasurementsCount,
		MD5Sum:            r.MD5Sum,
		Name:              r.Name,
		Status:            r.Status,
	}
	out.CreatedAt = time.Unix(r.CreatedAt, 0).Format("2006-01-02T15:04:05Z07:00")
	out.UpdatedAt = time.Unix(r.UpdatedAt, 0).Format("2006-01-02T15:04:05Z07:00")
	if r.SourceURL != "" {
		out.Source = &SourceRails{URL: r.SourceURL}
	}
	return out
}

func (h *Handler) adapterBgeigieImports(w http.ResponseWriter, r *http.Request) {
	if !wantsJSON(r) {
		http.Error(w, "Accept application/json", http.StatusNotAcceptable)
		return
	}
	addCORS(w, r)

	switch r.Method {
	case http.MethodGet:
		h.adapterBgeigieImportsList(w, r)
	case http.MethodPost:
		h.adapterBgeigieImportsCreate(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) adapterBgeigieImportsList(w http.ResponseWriter, r *http.Request) {
	filters := parseBgeigieFilters(r)
	rows, err := h.DB.QueryUploadsAsBgeigieImports(r.Context(), filters, h.DBType)
	if err != nil {
		h.Logf("safecast bgeigie_imports list: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	out := make([]BgeigieImportRails, len(rows))
	for i := range rows {
		out[i] = bgeigieRowToRails(rows[i])
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

func (h *Handler) adapterBgeigieImportsCreate(w http.ResponseWriter, r *http.Request) {
	apiKey := r.URL.Query().Get("api_key")
	if apiKey == "" {
		apiKey = r.Header.Get("X-API-Key")
	}
	if apiKey == "" && h.Auth != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (h *Handler) adapterBgeigieImportByID(w http.ResponseWriter, r *http.Request) {
	if !wantsJSON(r) {
		http.Error(w, "Accept application/json", http.StatusNotAcceptable)
		return
	}
	addCORS(w, r)
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.Trim(r.URL.Path, "/")
	path = strings.TrimSuffix(path, ".json")
	parts := strings.Split(path, "/")
	idStr := ""
	for i, p := range parts {
		if p == "bgeigie_imports" && i+1 < len(parts) {
			idStr = parts[i+1]
			if idStr == "submit" {
				return
			}
			break
		}
	}
	if idStr == "" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	filters := database.BgeigieImportFilters{ID: &id, PerPage: 1, Page: 1}
	rows, err := h.DB.QueryUploadsAsBgeigieImports(r.Context(), filters, h.DBType)
	if err != nil {
		h.Logf("safecast bgeigie_import get: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if len(rows) == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	out := bgeigieRowToRails(rows[0])
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}
