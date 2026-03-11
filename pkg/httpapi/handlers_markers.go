// handlers_markers.go — GET /api/markers/spectra and POST /api/update-coordinates.
package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
)

// markersWithSpectra returns markers that have spectral data, optionally filtered by bounding box (minLat, maxLat, minLon, maxLon). Omitted params default to world bounds.
func (s *Server) markersWithSpectra(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	if s.DB == nil || s.DB.DB == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "database not available")
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
	markers, err := s.Services.Marker.GetMarkersWithSpectra(ctx, bounds)
	if err != nil {
		s.Logf("error fetching markers with spectra: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "error fetching markers")
		return
	}
	writeJSON(w, http.StatusOK, markers)
}

// updateCoordinates sets lat/lon for all markers in a track (admin-only). Body: JSON with trackID, lat, lon.
func (s *Server) updateCoordinates(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
		return
	}
	if !auth.IsStaticAdminAuthorized(r, s.Config.AdminPassword) {
		auth.ChallengeStaticAdminBasic(w)
		return
	}
	if s.DB == nil || s.DB.DB == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "database not available")
		return
	}
	var req struct {
		TrackID string  `json:"trackID"`
		Lat     float64 `json:"lat"`
		Lon     float64 `json:"lon"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.TrackID == "" {
		writeJSONError(w, http.StatusBadRequest, "trackID is required")
		return
	}
	if req.Lat < -90 || req.Lat > 90 {
		writeJSONError(w, http.StatusBadRequest, "invalid latitude")
		return
	}
	if req.Lon < -180 || req.Lon > 180 {
		writeJSONError(w, http.StatusBadRequest, "invalid longitude")
		return
	}
	ctx := r.Context()
	rowsAffected, err := s.Services.Marker.UpdateTrackCoordinates(ctx, req.TrackID, req.Lat, req.Lon)
	if err != nil {
		s.Logf("error updating coordinates for track %s: %v", req.TrackID, err)
		writeJSONError(w, http.StatusInternalServerError, "failed to update coordinates")
		return
	}
	s.Logf("Updated coordinates for track %s: %d markers updated to (%.6f, %.6f)",
		req.TrackID, rowsAffected, req.Lat, req.Lon)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":         "success",
		"markersUpdated": rowsAffected,
		"lat":            req.Lat,
		"lon":            req.Lon,
	})
}
