package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"safecast-new-map/pkg/database"
)

// markersWithSpectra returns markers that have associated spectral data.
// GET /api/markers/spectra?minLat=...&maxLat=...&minLon=...&maxLon=...
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
	if minLat == 0 && maxLat == 0 && minLon == 0 && maxLon == 0 {
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

// updateCoordinates updates marker coordinates for spectrum files uploaded without GPS.
// POST /api/update-coordinates Body: {"trackID": "...", "lat": 34.488, "lon": 136.166}
func (s *Server) updateCoordinates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
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
