package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// MeasurementFilters holds version-neutral filter parameters for listing measurements.
type MeasurementFilters struct {
	ID                  *int64
	Latitude            *float64
	Longitude           *float64
	Distance            *int
	CapturedAfter       *int64
	CapturedBefore      *int64
	UserID              *int64
	DeviceID            *int64
	MeasurementImportID *int64
	OriginalID          *int64
	Since               *int64
	Until               *int64
	Unit                string
	Order               string
	Page                int
	PerPage             int
}

// MeasurementRow represents a marker mapped to measurement-like shape for Safecast API.
type MeasurementRow struct {
	ID                  int64
	Value               float64
	Height              *float64
	UserID              *int64
	Unit                string
	DeviceID            *int64
	LocationName        string
	OriginalID          *int64
	CapturedAt          int64
	Latitude            float64
	Longitude           float64
	MeasurementImportID *int64
	TrackID             string
}

// QueryMarkersAsMeasurements returns markers mapped to measurement-like rows with filters.
// Uses markers + uploads join for user_id and measurement_import_id when available.
func (db *Database) QueryMarkersAsMeasurements(ctx context.Context, filters MeasurementFilters, dbType string) ([]MeasurementRow, error) {
	if db == nil || db.DB == nil {
		return nil, fmt.Errorf("database unavailable")
	}

	perPage := filters.PerPage
	if perPage <= 0 {
		perPage = 100
	}
	if perPage > 500 {
		perPage = 500
	}
	page := filters.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage

	var conditions []string
	var args []interface{}
	paramCount := 0

	next := func() string {
		paramCount++
		if dbType == "pgx" || dbType == "duckdb" {
			return fmt.Sprintf("$%d", paramCount)
		}
		return "?"
	}

	if filters.Latitude != nil && filters.Longitude != nil && filters.Distance != nil {
		radius := float64(*filters.Distance)
		if radius <= 0 {
			radius = 1000
		}
		halfDeltaLat := radius / 111320.0
		lat := *filters.Latitude
		lon := *filters.Longitude
		minLat := lat - halfDeltaLat
		maxLat := lat + halfDeltaLat
		cosLat := 1.0
		if lat != 0 {
			cosLat = 0.9998 // approximate
		}
		halfDeltaLon := radius / (111320.0 * cosLat)
		minLon := lon - halfDeltaLon
		maxLon := lon + halfDeltaLon
		conditions = append(conditions, fmt.Sprintf("m.lat BETWEEN %s AND %s", next(), next()))
		args = append(args, minLat, maxLat)
		conditions = append(conditions, fmt.Sprintf("m.lon BETWEEN %s AND %s", next(), next()))
		args = append(args, minLon, maxLon)
	}

	if filters.CapturedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("m.date >= %s", next()))
		args = append(args, *filters.CapturedAfter)
	}
	if filters.CapturedBefore != nil {
		conditions = append(conditions, fmt.Sprintf("m.date <= %s", next()))
		args = append(args, *filters.CapturedBefore)
	}
	if filters.MeasurementImportID != nil {
		conditions = append(conditions, fmt.Sprintf("u.source_id = %s", next()))
		args = append(args, strconv.FormatInt(*filters.MeasurementImportID, 10))
	}
	if filters.ID != nil {
		conditions = append(conditions, fmt.Sprintf("m.id = %s", next()))
		args = append(args, *filters.ID)
	}
	if filters.OriginalID != nil {
		conditions = append(conditions, fmt.Sprintf("m.id = %s", next()))
		args = append(args, *filters.OriginalID)
	}
	if filters.Since != nil {
		conditions = append(conditions, fmt.Sprintf("m.date >= %s", next()))
		args = append(args, *filters.Since)
	}
	if filters.Until != nil {
		conditions = append(conditions, fmt.Sprintf("m.date <= %s", next()))
		args = append(args, *filters.Until)
	}
	if filters.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("u.internal_user_id = %s", next()))
		args = append(args, strconv.FormatInt(*filters.UserID, 10))
	}

	orderBy := "m.date DESC"
	if filters.Order != "" {
		lower := strings.ToLower(filters.Order)
		if strings.Contains(lower, "asc") {
			orderBy = "m.date ASC"
		}
	}

	// Column names: PostgreSQL uses lowercase
	trackCol := "trackid"
	if dbType != "pgx" && dbType != "duckdb" {
		trackCol = "trackID"
	}

	baseQuery := `
		SELECT m.id, m.doserate, m.lat, m.lon, m.date, m.` + trackCol + `,
		       COALESCE(u.internal_user_id, '') as internal_user_id,
		       COALESCE(u.source_id, '') as source_id
		FROM markers m
		LEFT JOIN uploads u ON u.track_id = m.` + trackCol + `
	`
	if len(conditions) > 0 {
		baseQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	// orderBy is "m.date DESC" or "m.date ASC"
	baseQuery += " ORDER BY m.date " + strings.TrimSuffix(strings.TrimPrefix(orderBy, "m.date "), "")
	baseQuery += fmt.Sprintf(" LIMIT %s OFFSET %s", next(), next())
	args = append(args, perPage, offset)

	rows, err := db.DB.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("query markers as measurements: %w", err)
	}
	defer rows.Close()

	var result []MeasurementRow
	for rows.Next() {
		var m MeasurementRow
		var internalUserID, sourceID string
		var doseRate, lat, lon float64
		var date int64
		var trackID string
		if err := rows.Scan(&m.ID, &doseRate, &lat, &lon, &date, &trackID, &internalUserID, &sourceID); err != nil {
			return nil, fmt.Errorf("scan marker: %w", err)
		}
		m.Value = doseRate
		m.Unit = "cpm"
		if doseRate > 0 && doseRate < 1000 {
			m.Unit = "usv"
		}
		m.Latitude = lat
		m.Longitude = lon
		m.CapturedAt = date
		m.TrackID = trackID
		if internalUserID != "" {
			if uid, err := strconv.ParseInt(internalUserID, 10, 64); err == nil {
				m.UserID = &uid
			}
		}
		if sourceID != "" {
			if sid, err := strconv.ParseInt(sourceID, 10, 64); err == nil {
				m.MeasurementImportID = &sid
			}
		}
		m.OriginalID = &m.ID
		result = append(result, m)
	}
	return result, rows.Err()
}

// CountMarkers returns total marker count, optionally with filters.
func (db *Database) CountMarkers(ctx context.Context, filters MeasurementFilters, dbType string) (int64, error) {
	if db == nil || db.DB == nil {
		return 0, fmt.Errorf("database unavailable")
	}

	var conditions []string
	var args []interface{}
	paramCount := 0

	next := func() string {
		paramCount++
		if dbType == "pgx" || dbType == "duckdb" {
			return fmt.Sprintf("$%d", paramCount)
		}
		return "?"
	}

	trackCol := "trackid"
	if dbType != "pgx" && dbType != "duckdb" {
		trackCol = "trackID"
	}

	if filters.Latitude != nil && filters.Longitude != nil && filters.Distance != nil {
		radius := float64(*filters.Distance)
		if radius <= 0 {
			radius = 1000
		}
		halfDeltaLat := radius / 111320.0
		lat := *filters.Latitude
		lon := *filters.Longitude
		cosLat := 1.0
		if lat != 0 {
			cosLat = 0.9998
		}
		halfDeltaLon := radius / (111320.0 * cosLat)
		minLat, maxLat := lat-halfDeltaLat, lat+halfDeltaLat
		minLon, maxLon := lon-halfDeltaLon, lon+halfDeltaLon
		conditions = append(conditions, fmt.Sprintf("m.lat BETWEEN %s AND %s", next(), next()))
		args = append(args, minLat, maxLat)
		conditions = append(conditions, fmt.Sprintf("m.lon BETWEEN %s AND %s", next(), next()))
		args = append(args, minLon, maxLon)
	}
	if filters.CapturedAfter != nil {
		conditions = append(conditions, fmt.Sprintf("m.date >= %s", next()))
		args = append(args, *filters.CapturedAfter)
	}
	if filters.CapturedBefore != nil {
		conditions = append(conditions, fmt.Sprintf("m.date <= %s", next()))
		args = append(args, *filters.CapturedBefore)
	}
	if filters.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("u.internal_user_id = %s", next()))
		args = append(args, strconv.FormatInt(*filters.UserID, 10))
	}
	if filters.MeasurementImportID != nil {
		conditions = append(conditions, fmt.Sprintf("u.source_id = %s", next()))
		args = append(args, strconv.FormatInt(*filters.MeasurementImportID, 10))
	}

	query := `SELECT COUNT(*) FROM markers m LEFT JOIN uploads u ON u.track_id = m.` + trackCol
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := db.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count markers: %w", err)
	}
	return count, nil
}

// BgeigieImportFilters holds filters for listing bgeigie imports.
type BgeigieImportFilters struct {
	ID             *int64
	Status         string
	ByUserID       *int64
	UploadedAfter  *int64
	UploadedBefore *int64
	Q              string
	Page           int
	PerPage        int
}

// BgeigieImportRow represents an upload mapped to bgeigie_import shape.
type BgeigieImportRow struct {
	ID                int64
	UserID            *int64
	Approved          bool
	CreatedAt         int64
	UpdatedAt         int64
	MeasurementsCount int
	MD5Sum            string
	Name              string
	Status            string
	SourceURL         string
	TrackID           string
}

// QueryUploadsAsBgeigieImports returns uploads mapped to bgeigie_imports schema.
func (db *Database) QueryUploadsAsBgeigieImports(ctx context.Context, filters BgeigieImportFilters, dbType string) ([]BgeigieImportRow, error) {
	if db == nil || db.DB == nil {
		return nil, fmt.Errorf("database unavailable")
	}

	perPage := filters.PerPage
	if perPage <= 0 {
		perPage = 25
	}
	if perPage > 100 {
		perPage = 100
	}
	page := filters.Page
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * perPage

	var conditions []string
	var args []interface{}
	paramCount := 0

	next := func() string {
		paramCount++
		if dbType == "pgx" || dbType == "duckdb" {
			return fmt.Sprintf("$%d", paramCount)
		}
		return "?"
	}

	if filters.ID != nil {
		conditions = append(conditions, fmt.Sprintf("u.id = %s", next()))
		args = append(args, *filters.ID)
	}
	if filters.ByUserID != nil {
		conditions = append(conditions, fmt.Sprintf("(u.internal_user_id = %s OR u.user_id = %s)", next(), next()))
		args = append(args, strconv.FormatInt(*filters.ByUserID, 10), strconv.FormatInt(*filters.ByUserID, 10))
	}
	if filters.UploadedAfter != nil {
		if dbType == "pgx" || dbType == "duckdb" {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(EPOCH FROM u.created_at) >= %s", next()))
		} else {
			conditions = append(conditions, fmt.Sprintf("u.created_at >= %s", next()))
		}
		args = append(args, *filters.UploadedAfter)
	}
	if filters.UploadedBefore != nil {
		if dbType == "pgx" || dbType == "duckdb" {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(EPOCH FROM u.created_at) <= %s", next()))
		} else {
			conditions = append(conditions, fmt.Sprintf("u.created_at <= %s", next()))
		}
		args = append(args, *filters.UploadedBefore)
	}
	if filters.Q != "" {
		pattern := "%" + filters.Q + "%"
		if dbType == "pgx" || dbType == "duckdb" {
			conditions = append(conditions, fmt.Sprintf("(u.filename ILIKE %s OR u.track_id ILIKE %s OR COALESCE(u.username,'') ILIKE %s)", next(), next(), next()))
		} else {
			conditions = append(conditions, "(u.filename LIKE ? OR u.track_id LIKE ? OR COALESCE(u.username,'') LIKE ?)")
		}
		for i := 0; i < 3; i++ {
			args = append(args, pattern)
		}
	}

	markerTrackCol := "trackID"
	if dbType == "pgx" || dbType == "duckdb" {
		markerTrackCol = "trackid"
	}
	createdSel := "u.created_at"
	if dbType == "pgx" || dbType == "duckdb" {
		createdSel = "EXTRACT(EPOCH FROM u.created_at)::BIGINT"
	}
	query := fmt.Sprintf(`
		SELECT u.id, u.filename, u.track_id, u.source_url, u.user_id, u.username, u.internal_user_id,
		       %s,
		       (SELECT COUNT(*) FROM markers m WHERE m.%s = u.track_id) as marker_count
		FROM uploads u
	`, createdSel, markerTrackCol)
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY u.created_at DESC"
	query += fmt.Sprintf(" LIMIT %s OFFSET %s", next(), next())
	args = append(args, perPage, offset)

	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query uploads as bgeigie_imports: %w", err)
	}
	defer rows.Close()

	var result []BgeigieImportRow
	for rows.Next() {
		var r BgeigieImportRow
		var filename, trackID, sourceURL, userID, username, internalUserID sql.NullString
		var createdAt interface{}
		var markerCount int

		if err := rows.Scan(&r.ID, &filename, &trackID, &sourceURL, &userID, &username, &internalUserID, &createdAt, &markerCount); err != nil {
			return nil, fmt.Errorf("scan upload: %w", err)
		}

		r.Name = filename.String
		r.TrackID = trackID.String
		r.SourceURL = sourceURL.String
		r.MeasurementsCount = markerCount
		r.Status = "done"
		r.Approved = true
		r.MD5Sum = ""
		switch v := createdAt.(type) {
		case int64:
			r.CreatedAt, r.UpdatedAt = v, v
		case float64:
			r.CreatedAt, r.UpdatedAt = int64(v), int64(v)
		case nil:
			r.CreatedAt, r.UpdatedAt = time.Now().Unix(), time.Now().Unix()
		default:
			r.CreatedAt, r.UpdatedAt = time.Now().Unix(), time.Now().Unix()
		}

		if internalUserID.Valid && internalUserID.String != "" {
			if uid, err := strconv.ParseInt(internalUserID.String, 10, 64); err == nil {
				r.UserID = &uid
			}
		} else if userID.Valid && userID.String != "" {
			if uid, err := strconv.ParseInt(userID.String, 10, 64); err == nil {
				r.UserID = &uid
			}
		}

		result = append(result, r)
	}
	return result, rows.Err()
}

// CreateSingleMeasurement inserts a single measurement as a one-marker track.
// Used by POST /measurements.json. Returns the created marker ID.
func (db *Database) CreateSingleMeasurement(ctx context.Context, value float64, unit string, lat, lon float64, capturedAt int64, userID string, dbType string) (int64, error) {
	if db == nil || db.DB == nil {
		return 0, fmt.Errorf("database unavailable")
	}
	trackID := fmt.Sprintf("safecast-%d-%d", time.Now().UnixNano(), time.Now().Unix()%100000)
	if err := db.EnsureTrackPresence(ctx, trackID, dbType); err != nil {
		return 0, fmt.Errorf("ensure track: %w", err)
	}
	m := Marker{
		DoseRate:  value,
		CountRate: value,
		Date:      capturedAt,
		Lat:       lat,
		Lon:       lon,
		Zoom:      0,
		Speed:     0,
		TrackID:   trackID,
		Detector:  "",
	}
	if unit == "cpm" {
		m.DoseRate = value / 333.0 // approximate cpm to µSv/h
	}
	err := db.withSerializedConnectionFor(ctx, WorkloadUserUpload, func(ctx context.Context, conn *sql.DB) error {
		return db.SaveMarkerAtomic(ctx, conn, m, dbType)
	})
	if err != nil {
		return 0, fmt.Errorf("save marker: %w", err)
	}
	// Fetch the created marker ID by track+date+lat+lon
	var id int64
	trackCol := "trackID"
	if dbType == "pgx" || dbType == "duckdb" {
		trackCol = "trackid"
	}
	var q string
	if dbType == "pgx" || dbType == "duckdb" {
		q = fmt.Sprintf("SELECT id FROM markers WHERE %s = $1 AND date = $2 AND lat = $3 AND lon = $4 LIMIT 1", trackCol)
	} else {
		q = fmt.Sprintf("SELECT id FROM markers WHERE %s = ? AND date = ? AND lat = ? AND lon = ? LIMIT 1", trackCol)
	}
	err = db.DB.QueryRowContext(ctx, q, trackID, capturedAt, lat, lon).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("fetch marker id: %w", err)
	}
	return id, nil
}
