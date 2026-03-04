// handlers_trackinfo.go — track metadata
//
// Returns lightweight metadata about a track (username, detector, recording
// date). Used by the map UI to show who recorded the data and when.
package web

import (
	"encoding/json"
	"net/http"
	"strings"
)

// trackInfo returns upload metadata for a track: username, detector, recordingDate.
//
// Route: GET /api/track-info/{trackID}
//
// If the track is not found, returns minimal JSON with just the trackID.
// Responses are cached for 1 hour (Cache-Control).
func (s *Server) trackInfo(w http.ResponseWriter, r *http.Request) {
	if s.DB == nil || s.DB.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}
	trackID := strings.TrimPrefix(r.URL.Path, "/api/track-info/")
	if trackID == "" {
		http.Error(w, "Missing track ID", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	var username, detector string
	var recordingDate int64
	var query string
	switch s.Config.DBType {
	case "pgx", "duckdb":
		query = `SELECT COALESCE(username, ''), COALESCE(detector, ''),
		         COALESCE(EXTRACT(EPOCH FROM recording_date)::BIGINT, 0)
		         FROM uploads WHERE track_id = $1 LIMIT 1`
	default:
		query = `SELECT COALESCE(username, ''), COALESCE(detector, ''),
		         COALESCE(recording_date, 0)
		         FROM uploads WHERE track_id = ? LIMIT 1`
	}
	err := s.DB.DB.QueryRowContext(ctx, query, trackID).Scan(&username, &detector, &recordingDate)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		json.NewEncoder(w).Encode(map[string]interface{}{"trackID": trackID})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"trackID":       trackID,
		"username":      username,
		"detector":      detector,
		"recordingDate": recordingDate,
	})
}
