package web

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// apiTracksBounds returns combined bounds for multiple tracks.
// GET /api/tracks/bounds?trackIDs=track1,track2,track3
func (s *Server) apiTracksBounds(w http.ResponseWriter, r *http.Request) {
	trackIDsParam := r.URL.Query().Get("trackIDs")
	if trackIDsParam == "" {
		http.Error(w, "trackIDs parameter required", http.StatusBadRequest)
		return
	}
	trackIDs := strings.Split(trackIDsParam, ",")
	if len(trackIDs) == 0 {
		http.Error(w, "No track IDs provided", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	var minLat, minLon, maxLat, maxLon float64
	first := true
	for _, trackID := range trackIDs {
		trackID = strings.TrimSpace(trackID)
		if trackID == "" {
			continue
		}
		var tMinLat, tMinLon, tMaxLat, tMaxLon sql.NullFloat64
		var query string
		if s.Config.DBType == "pgx" {
			query = `SELECT MIN(lat) as min_lat, MIN(lon) as min_lon, MAX(lat) as max_lat, MAX(lon) as max_lon 
			         FROM markers WHERE trackID = $1`
		} else {
			query = `SELECT MIN(lat) as min_lat, MIN(lon) as min_lon, MAX(lat) as max_lat, MAX(lon) as max_lon 
			         FROM markers WHERE trackID = ?`
		}
		err := s.DB.DB.QueryRowContext(ctx, query, trackID).Scan(&tMinLat, &tMinLon, &tMaxLat, &tMaxLon)
		if err != nil {
			log.Printf("Error querying bounds for track %s: %v", trackID, err)
			continue
		}
		if !tMinLat.Valid || !tMinLon.Valid || !tMaxLat.Valid || !tMaxLon.Valid {
			continue
		}
		if first {
			minLat, minLon = tMinLat.Float64, tMinLon.Float64
			maxLat, maxLon = tMaxLat.Float64, tMaxLon.Float64
			first = false
		} else {
			if tMinLat.Float64 < minLat {
				minLat = tMinLat.Float64
			}
			if tMinLon.Float64 < minLon {
				minLon = tMinLon.Float64
			}
			if tMaxLat.Float64 > maxLat {
				maxLat = tMaxLat.Float64
			}
			if tMaxLon.Float64 > maxLon {
				maxLon = tMaxLon.Float64
			}
		}
	}
	if first {
		http.Error(w, "No valid tracks found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"bounds": map[string]float64{
			"minLat": minLat,
			"minLon": minLon,
			"maxLat": maxLat,
			"maxLon": maxLon,
		},
		"trackIDs": trackIDs,
	})
}
