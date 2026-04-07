// handlers_trackinfo.go — GET /api/track-info/{id} returns upload metadata (username, detector, recordingDate).
package httpapi

import (
	"net/http"
	"strings"
)

// trackInfo returns metadata for a track.
//
// @Summary     Get track metadata
// @Description Returns username, detector, and recording date metadata for a track ID.
// @Tags        web
// @Produce     json
// @Param       id path string true "Track ID"
// @Success     200 {object} map[string]interface{} "Track info payload"
// @Failure     400 {object} map[string]string "Missing track ID"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /api/track-info/{id} [get]
func (s *Server) trackInfo(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	if s.DB == nil || s.DB.DB == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "database not available")
		return
	}
	trackID := strings.TrimPrefix(r.URL.Path, "/api/track-info/")
	if trackID == "" {
		writeJSONError(w, http.StatusBadRequest, "missing track ID")
		return
	}
	ctx := r.Context()
	info, err := s.Services.Track.GetTrackInfo(ctx, trackID)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	if err != nil || !info.Found {
		writeJSON(w, http.StatusOK, map[string]interface{}{"trackID": trackID})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"trackID":       trackID,
		"username":      info.Username,
		"detector":      info.Detector,
		"recordingDate": info.RecordingDate,
		"sourceURL":     info.SourceURL,
	})
}
