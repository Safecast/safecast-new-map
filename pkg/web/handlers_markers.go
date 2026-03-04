// handlers_markers.go — markers and spectra
//
// Handles marker-related API: listing markers that have gamma spectra, and
// updating coordinates for markers that were uploaded without GPS.
package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"safecast-new-map/pkg/database"
)

// markersWithSpectra returns markers that have associated spectral data,
// optionally filtered by a geographic bounding box.
//
// Route: GET /api/markers/spectra?minLat=...&maxLat=...&minLon=...&maxLon=...
//
// Query params (all optional): minLat, maxLat, minLon, maxLon. If omitted,
// the whole world (-90..90, -180..180) is used.
func (s *Server) markersWithSpectra(w http.ResponseWriter, r *http.Request) {
	if s.DB == nil || s.DB.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}
	ctx, cancel := WithMinimumDeadline(r.Context(), 30*time.Second)
	defer cancel()
	q := r.URL.Query()
	minLat, _ := strconv.ParseFloat(q.Get("minLat"), 64)
	minLon, _ := strconv.ParseFloat(q.Get("minLon"), 64)
	maxLat, _ := strconv.ParseFloat(q.Get("maxLat"), 64)
	maxLon, _ := strconv.ParseFloat(q.Get("maxLon"), 64)
	_, hasMinLat := q["minLat"]
	_, hasMaxLat := q["maxLat"]
	_, hasMinLon := q["minLon"]
	_, hasMaxLon := q["maxLon"]
	if !hasMinLat && !hasMaxLat && !hasMinLon && !hasMaxLon {
		minLat, maxLat = -90, 90
		minLon, maxLon = -180, 180
	}
	bounds := database.Bounds{
		MinLat: minLat,
		MaxLat: maxLat,
		MinLon: minLon,
		MaxLon: maxLon,
	}
	markers, err := s.DB.GetMarkersWithSpectra(ctx, bounds)
	if err != nil {
		log.Printf("Error fetching markers with spectra: %v", err)
		http.Error(w, "Error fetching markers", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(markers)
}

// updateCoordinates sets lat/lon for all markers in a track. Used when spectra
// were uploaded without GPS; an admin can later set the correct location.
//
// Route: POST /api/update-coordinates
//
// Body: JSON with trackID (required), lat (-90..90), lon (-180..180).
func (s *Server) updateCoordinates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.Config.AdminPassword != "" {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != s.Config.AdminPassword {
			w.Header().Set("WWW-Authenticate", `Basic realm="safecast"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	if s.DB == nil || s.DB.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}
	var req struct {
		TrackID string  `json:"trackID"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.TrackID == "" {
		http.Error(w, "trackID is required", http.StatusBadRequest)
		return
	}
	if req.Lat < -90 || req.Lat > 90 {
		http.Error(w, "Invalid latitude", http.StatusBadRequest)
		return
	}
	if req.Lon < -180 || req.Lon > 180 {
		http.Error(w, "Invalid longitude", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	query := "UPDATE markers SET lat = ?, lon = ? WHERE trackID = ?"
	args := []interface{}{req.Lat, req.Lon, req.TrackID}
	if s.Config.DBType == "pgx" || s.Config.DBType == "duckdb" {
		query = "UPDATE markers SET lat = $1, lon = $2 WHERE trackID = $3"
	}
	result, err := s.DB.DB.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("Error updating coordinates for track %s: %v", req.TrackID, err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "error",
			"error":  "Failed to update coordinates",
		})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	s.Logf("Updated coordinates for track %s: %d markers updated to (%.6f, %.6f)",
		req.TrackID, rowsAffected, req.Lat, req.Lon)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":         "success",
		"markersUpdated": rowsAffected,
		"lat":            req.Lat,
		"lon":            req.Lon,
	})
}