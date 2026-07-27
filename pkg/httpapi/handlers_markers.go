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

// markersWithSpectra returns markers that have spectral data.
//
// @Summary     List markers with spectra
// @Description Returns markers that contain spectroscopy data, optionally filtered by bounding box.
// @Tags        web
// @Produce     json
// @Param       minLat query number false "Minimum latitude"
// @Param       maxLat query number false "Maximum latitude"
// @Param       minLon query number false "Minimum longitude"
// @Param       maxLon query number false "Maximum longitude"
// @Success     200 {array} map[string]interface{} "Markers"
// @Failure     500 {object} map[string]string "Server error"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /api/markers/spectra [get]
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

// updateCoordinates sets lat/lon for all markers in a track. Allowed for the
// static admin (Basic Auth) or the logged-in user who owns the track's upload.
//
// @Summary     Update coordinates for a track
// @Description Updates all marker coordinates for a track. Requires admin Basic Auth or session ownership of the upload.
// @Tags        admin
// @Accept      json
// @Produce     json
// @Param       body body object true "trackID, lat, lon"
// @Success     200 {object} map[string]interface{} "Update result"
// @Failure     400 {object} map[string]string "Invalid request"
// @Failure     403 {object} map[string]string "Forbidden"
// @Failure     500 {object} map[string]string "Server error"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /api/update-coordinates [post]
func (s *Server) updateCoordinates(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodPost) {
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

	authorized := auth.IsStaticAdminAuthorized(r, s.Config.AdminPassword)
	sessionUser, hasSessionUser := auth.GetUserFromContext(ctx)
	hasSessionUser = hasSessionUser && sessionUser != nil
	if !authorized && hasSessionUser {
		owns, err := s.DB.TrackBelongsToUser(ctx, req.TrackID, strconv.FormatInt(sessionUser.ID, 10))
		if err != nil {
			s.Logf("error checking track ownership for %s: %v", req.TrackID, err)
		}
		authorized = owns
	}
	if !authorized {
		if hasSessionUser {
			// Logged in but doesn't own this track: plain forbidden, no Basic Auth popup.
			writeJSONError(w, http.StatusForbidden, "not authorized to update this track")
			return
		}
		// No session at all: preserve the static-admin Basic Auth challenge for anonymous/tooling callers.
		auth.ChallengeStaticAdminBasic(w)
		return
	}
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
