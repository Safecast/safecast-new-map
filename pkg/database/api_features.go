package database

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
)

// TrackInfo stores lightweight upload metadata for API responses.
type TrackInfo struct {
	Username      string
	Detector      string
	RecordingDate int64
	SourceURL     string
	Filename      string
	FileType      string
	Name          string
	Comment       string
	UploadedAt    int64
	Source        string
	SourceID      string
	ImportID      string
	SimplemapUser string
	UploadIP      string
	Found         bool
}

// IsSafecastDevice reports whether a track came from a Safecast bGeigie (Nano/Zen).
// CzechRad firmware reports itself as a bGeigie, and Radiacode is not a Safecast device.
func IsSafecastDevice(detector, username string) bool {
	combined := strings.ToLower(detector + " " + username)
	for _, bad := range []string{"czechrad", "radiacode"} {
		if strings.Contains(combined, bad) {
			return false
		}
	}
	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(detector)), "rc-") {
		return false
	}
	d := strings.ToLower(detector)
	return strings.Contains(d, "bgeigie") || strings.Contains(d, "nano") || strings.Contains(d, "zen")
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
		         COALESCE(source_url, ''), COALESCE(filename, ''), COALESCE(file_type, ''),
		         COALESCE(name, ''), COALESCE(comment, ''),
		         COALESCE(EXTRACT(EPOCH FROM created_at)::BIGINT, 0),
		         COALESCE(source, ''), COALESCE(source_id, ''), COALESCE(safecast_import_id, ''),
		         COALESCE((SELECT usr.username FROM users usr WHERE usr.id::text = uploads.internal_user_id), ''),
		         COALESCE(upload_ip, '')
		         FROM uploads WHERE track_id = $1 LIMIT 1`
	default:
		query = `SELECT COALESCE(username, ''), COALESCE(detector, ''),
		         COALESCE(recording_date, 0),
		         COALESCE(source_url, ''), COALESCE(filename, ''), COALESCE(file_type, ''),
		         COALESCE(name, ''), COALESCE(comment, ''), COALESCE(created_at, 0),
		         COALESCE(source, ''), COALESCE(source_id, ''), COALESCE(safecast_import_id, ''),
		         COALESCE((SELECT usr.username FROM users usr WHERE CAST(usr.id AS TEXT) = uploads.internal_user_id), ''),
		         COALESCE(upload_ip, '')
		         FROM uploads WHERE track_id = ? LIMIT 1`
	}

	err := db.DB.QueryRowContext(ctx, query, trackID).Scan(&out.Username, &out.Detector, &out.RecordingDate, &out.SourceURL, &out.Filename, &out.FileType, &out.Name, &out.Comment, &out.UploadedAt, &out.Source, &out.SourceID, &out.ImportID, &out.SimplemapUser, &out.UploadIP)
	if err != nil {
		if err == sql.ErrNoRows {
			return out, nil
		}
		return out, fmt.Errorf("query track info: %w", err)
	}
	out.Found = true
	// Radiacode files carry no detector name; infer it from the file extension.
	if strings.TrimSpace(out.Detector) == "" {
		switch strings.ToLower(filepath.Ext(out.Filename)) {
		case ".rctrk", ".rcspg", ".rcxml":
			out.Detector = "Radiacode (" + strings.ToLower(filepath.Ext(out.Filename)) + " file)"
		}
	}
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

