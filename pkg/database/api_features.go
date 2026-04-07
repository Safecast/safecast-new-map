package database

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// TrackInfo stores lightweight upload metadata for API responses.
type TrackInfo struct {
	Username      string
	Detector      string
	RecordingDate int64
	SourceURL     string
	Found         bool
}

// GetTrackInfo returns username/detector/date/sourceURL metadata for a track upload.
func (db *Database) GetTrackInfo(ctx context.Context, trackID string) (TrackInfo, error) {
	out := TrackInfo{}
	if strings.TrimSpace(trackID) == "" {
		return out, fmt.Errorf("trackID is required")
	}

	var query string
	switch db.Driver {
	case "pgx", "duckdb":
		query = `SELECT COALESCE(username, ''), COALESCE(detector, ''),
		         COALESCE(EXTRACT(EPOCH FROM recording_date)::BIGINT, 0),
		         COALESCE(source_url, '')
		         FROM uploads WHERE track_id = $1 LIMIT 1`
	default:
		query = `SELECT COALESCE(username, ''), COALESCE(detector, ''),
		         COALESCE(recording_date, 0),
		         COALESCE(source_url, '')
		         FROM uploads WHERE track_id = ? LIMIT 1`
	}

	err := db.DB.QueryRowContext(ctx, query, trackID).Scan(&out.Username, &out.Detector, &out.RecordingDate, &out.SourceURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return out, nil
		}
		return out, fmt.Errorf("query track info: %w", err)
	}
	out.Found = true
	return out, nil
}

// UpdateTrackCoordinates updates all marker coordinates for a track.
func (db *Database) UpdateTrackCoordinates(ctx context.Context, trackID string, lat, lon float64) (int64, error) {
	if strings.TrimSpace(trackID) == "" {
		return 0, fmt.Errorf("trackID is required")
	}

	query := "UPDATE markers SET lat = ?, lon = ? WHERE trackID = ?"
	args := []any{lat, lon, trackID}
	if db.Driver == "pgx" || db.Driver == "duckdb" {
		query = "UPDATE markers SET lat = $1, lon = $2 WHERE trackID = $3"
	}

	result, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("update track coordinates: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	return rowsAffected, nil
}

// GetTrackBounds returns min/max lat/lon for one track.
func (db *Database) GetTrackBounds(ctx context.Context, trackID string) (Bounds, bool, error) {
	out := Bounds{}
	if strings.TrimSpace(trackID) == "" {
		return out, false, nil
	}

	var query string
	if db.Driver == "pgx" || db.Driver == "duckdb" {
		query = `SELECT MIN(lat), MIN(lon), MAX(lat), MAX(lon)
		         FROM markers WHERE trackID = $1`
	} else {
		query = `SELECT MIN(lat), MIN(lon), MAX(lat), MAX(lon)
		         FROM markers WHERE trackID = ?`
	}

	var minLat, minLon, maxLat, maxLon sql.NullFloat64
	err := db.DB.QueryRowContext(ctx, query, trackID).Scan(&minLat, &minLon, &maxLat, &maxLon)
	if err != nil {
		return out, false, fmt.Errorf("query track bounds: %w", err)
	}
	if !minLat.Valid || !minLon.Valid || !maxLat.Valid || !maxLon.Valid {
		return out, false, nil
	}

	out.MinLat = minLat.Float64
	out.MinLon = minLon.Float64
	out.MaxLat = maxLat.Float64
	out.MaxLon = maxLon.Float64
	return out, true, nil
}
