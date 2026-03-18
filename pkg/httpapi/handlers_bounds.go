// handlers_bounds.go — GET /api/tracks/bounds?trackIDs=id1,id2,... returns combined geographic bounds.
package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

// apiTracksBounds returns merged geographic bounds for track IDs.
//
// @Summary     Get combined bounds for tracks
// @Description Returns a combined min/max lat/lon bounding box for the provided track IDs.
// @Tags        web
// @Produce     json
// @Param       trackIDs query string true "Comma-separated track IDs"
// @Success     200 {object} map[string]interface{} "Bounds payload"
// @Failure     400 {object} map[string]string "Invalid request"
// @Failure     404 {object} map[string]string "No valid tracks found"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /api/tracks/bounds [get]
func (s *Server) apiTracksBounds(w http.ResponseWriter, r *http.Request) {
	if !requireMethod(w, r, http.MethodGet) {
		return
	}
	if s.DB == nil || s.DB.DB == nil {
		writeJSONError(w, http.StatusServiceUnavailable, "database not available")
		return
	}
	trackIDsParam := r.URL.Query().Get("trackIDs")
	if trackIDsParam == "" {
		writeJSONError(w, http.StatusBadRequest, "trackIDs parameter required")
		return
	}
	rawTrackIDs := strings.Split(trackIDsParam, ",")
	trackIDs := normalizeTrackIDs(rawTrackIDs)
	if len(trackIDs) == 0 {
		writeJSONError(w, http.StatusBadRequest, "no track IDs provided")
		return
	}
	ctx := r.Context()

	first := true
	var combined struct {
		minLat float64
		minLon float64
		maxLat float64
		maxLon float64
	}
	for _, trackID := range trackIDs {
		bounds, found, err := s.Services.Track.GetTrackBounds(ctx, trackID)
		if err != nil {
			s.Logf("error querying bounds for track %s: %v", trackID, err)
			continue
		}
		if !found {
			continue
		}
		if first {
			combined.minLat, combined.minLon = bounds.MinLat, bounds.MinLon
			combined.maxLat, combined.maxLon = bounds.MaxLat, bounds.MaxLon
			first = false
		} else {
			if bounds.MinLat < combined.minLat {
				combined.minLat = bounds.MinLat
			}
			if bounds.MinLon < combined.minLon {
				combined.minLon = bounds.MinLon
			}
			if bounds.MaxLat > combined.maxLat {
				combined.maxLat = bounds.MaxLat
			}
			if bounds.MaxLon > combined.maxLon {
				combined.maxLon = bounds.MaxLon
			}
		}
	}
	if first {
		writeJSONError(w, http.StatusNotFound, "no valid tracks found")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"bounds": map[string]float64{
			"minLat": combined.minLat,
			"minLon": combined.minLon,
			"maxLat": combined.maxLat,
			"maxLon": combined.maxLon,
		},
		"trackIDs": rawTrackIDs,
	})
}
