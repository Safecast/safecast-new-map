package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Upload represents a file upload record for tracking purposes.
type Upload struct {
	ID             int64  `json:"id"`
	Filename       string `json:"filename"`
	FileType       string `json:"fileType"`
	TrackID        string `json:"trackID"`
	FileSize       int64  `json:"fileSize"`
	UploadIP       string `json:"uploadIP"`
	CreatedAt      int64  `json:"createdAt"`
	RecordingDate  int64  `json:"recordingDate,omitempty"`  // Earliest marker date for this track
	Source         string `json:"source,omitempty"`         // Import source (e.g., "safecast-api", "user-upload")
	SourceID       string `json:"sourceID,omitempty"`       // External reference ID (e.g., Safecast import ID)
	SourceURL      string `json:"sourceURL,omitempty"`      // Source file URL (e.g., S3 URL)
	UserID         string `json:"userID,omitempty"`         // User ID from source (e.g., Safecast user ID)
	Username       string `json:"username,omitempty"`       // Username from source (fetched from API)
	InternalUserID string `json:"internalUserID,omitempty"` // Internal user ID from users table (for authenticated uploads)
	Detector       string `json:"detector,omitempty"`       // Detector info from track_statistics
	Name           string `json:"name,omitempty"`           // Display name (admin-editable)
	Notes          string `json:"notes,omitempty"`          // Admin notes
	Comment        string `json:"comment,omitempty"`        // User-supplied description from Safecast API

	SafecastImportID      string `json:"safecastImportID,omitempty"`      // bgeigie_imports id, once submitted to api.safecast.org
	SubmittedToSafecastAt int64  `json:"submittedToSafecastAt,omitempty"` // Unix timestamp of successful submit
	SafecastSubmitError   string `json:"safecastSubmitError,omitempty"`   // Last submit failure reason, if any
}

// UpdateUploadSafecastStatus records the outcome of submitting an upload to
// api.safecast.org, keyed by the upload row's own id (not track_id, which can be
// shared across files in one request or reused when duplicate detection finds an
// existing track). A successful submit sets importID and clears any previous
// error; a failed submit sets submitErr and leaves importID empty.
func (db *Database) UpdateUploadSafecastStatus(ctx context.Context, uploadID int64, importID, submitErr string) error {
	var submittedAt sql.NullInt64
	if importID != "" {
		submittedAt = sql.NullInt64{Int64: time.Now().Unix(), Valid: true}
	}

	var query string
	var args []interface{}

	switch db.Driver {
	case "pgx", "duckdb":
		query = `UPDATE uploads SET safecast_import_id = $1, submitted_to_safecast_at = to_timestamp($2), safecast_submit_error = $3 WHERE id = $4`
		args = []interface{}{importID, submittedAt, submitErr, uploadID}
	case "sqlite", "chai":
		query = `UPDATE uploads SET safecast_import_id = ?, submitted_to_safecast_at = ?, safecast_submit_error = ? WHERE id = ?`
		args = []interface{}{importID, submittedAt, submitErr, uploadID}
	default:
		return fmt.Errorf("unsupported database driver: %s", db.Driver)
	}

	if _, err := db.DB.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("update upload safecast status: %w", err)
	}
	return nil
}

// InsertUpload records a file upload in the uploads table.
func (db *Database) InsertUpload(ctx context.Context, upload Upload) (int64, error) {
	createdAt := upload.CreatedAt
	if createdAt == 0 {
		createdAt = time.Now().Unix()
	}

	var query string
	var args []interface{}
	var id int64

	switch db.Driver {
	case "pgx":
		query = `
			INSERT INTO uploads (filename, file_type, track_id, file_size, upload_ip, created_at, recording_date, source, source_id, source_url, user_id, username, internal_user_id, detector, comment)
			VALUES ($1, $2, $3, $4, $5, to_timestamp($6), to_timestamp($7), $8, $9, $10, $11, $12, $13, $14, $15)
			RETURNING id
		`
		recordingDateArg := upload.RecordingDate
		if recordingDateArg == 0 {
			recordingDateArg = createdAt
		}
		args = []interface{}{
			upload.Filename, upload.FileType, upload.TrackID,
			upload.FileSize, upload.UploadIP, createdAt, recordingDateArg,
			upload.Source, upload.SourceID, upload.SourceURL, upload.UserID, upload.Username, upload.InternalUserID, upload.Detector, upload.Comment,
		}
		err := db.DB.QueryRowContext(ctx, query, args...).Scan(&id)
		return id, err

	case "duckdb":
		query = `
			INSERT INTO uploads (filename, file_type, track_id, file_size, upload_ip, created_at, recording_date, source, source_id, source_url, user_id, username, internal_user_id, detector, comment)
			VALUES ($1, $2, $3, $4, $5, to_timestamp($6), to_timestamp($7), $8, $9, $10, $11, $12, $13, $14, $15)
			RETURNING id
		`
		recordingDateArg := upload.RecordingDate
		if recordingDateArg == 0 {
			recordingDateArg = createdAt
		}
		args = []interface{}{
			upload.Filename, upload.FileType, upload.TrackID,
			upload.FileSize, upload.UploadIP, createdAt, recordingDateArg,
			upload.Source, upload.SourceID, upload.SourceURL, upload.UserID, upload.Username, upload.InternalUserID, upload.Detector, upload.Comment,
		}
		err := db.DB.QueryRowContext(ctx, query, args...).Scan(&id)
		return id, err

	case "sqlite", "chai":
		query = `
			INSERT INTO uploads (filename, file_type, track_id, file_size, upload_ip, created_at, recording_date, source, source_id, source_url, user_id, username, internal_user_id, detector, comment)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`
		recordingDateArg := upload.RecordingDate
		if recordingDateArg == 0 {
			recordingDateArg = createdAt
		}
		args = []interface{}{
			upload.Filename, upload.FileType, upload.TrackID,
			upload.FileSize, upload.UploadIP, createdAt, recordingDateArg,
			upload.Source, upload.SourceID, upload.SourceURL, upload.UserID, upload.Username, upload.InternalUserID, upload.Detector, upload.Comment,
		}
		result, err := db.DB.ExecContext(ctx, query, args...)
		if err != nil {
			return 0, err
		}
		return result.LastInsertId()

	default:
		return 0, fmt.Errorf("unsupported database driver: %s", db.Driver)
	}
}

// GetUploads retrieves upload records, ordered by most recent first.
// If userID is not empty, only returns uploads from that user.
func (db *Database) GetUploads(ctx context.Context, limit int, userID string) ([]Upload, error) {
	return db.GetUploadsPaginated(ctx, limit, 0, userID, "")
}

func (db *Database) GetUploadsPaginated(ctx context.Context, limit int, offset int, userID string, search string) ([]Upload, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var query string
	var args []interface{}
	var whereConditions []string
	paramCount := 0

	// Build WHERE conditions
	if userID != "" {
		paramCount++
		if db.Driver == "pgx" || db.Driver == "duckdb" {
			whereConditions = append(whereConditions, fmt.Sprintf("u.internal_user_id = $%d", paramCount))
		} else {
			whereConditions = append(whereConditions, "u.internal_user_id = ?")
		}
	}

	if search != "" {
		paramCount++
		if db.Driver == "pgx" || db.Driver == "duckdb" {
			// PostgreSQL: use ILIKE for case-insensitive search, also search numeric fields by converting to text
			whereConditions = append(whereConditions, fmt.Sprintf(
				"(u.track_id ILIKE $%d OR u.filename ILIKE $%d OR u.file_type ILIKE $%d OR COALESCE(u.user_id, '') ILIKE $%d OR COALESCE(u.username, '') ILIKE $%d OR COALESCE(u.source, '') ILIKE $%d OR COALESCE(u.source_id, '') ILIKE $%d OR COALESCE(u.detector, '') ILIKE $%d OR CAST(u.id AS TEXT) ILIKE $%d OR TO_CHAR(u.recording_date, 'YYYY-MM-DD HH24:MI:SS') ILIKE $%d OR COALESCE(u.comment, '') ILIKE $%d)",
				paramCount, paramCount, paramCount, paramCount, paramCount, paramCount, paramCount, paramCount, paramCount, paramCount, paramCount))
		} else {
			// SQLite: use LIKE (case-insensitive by default), also search numeric fields by converting to text
			whereConditions = append(whereConditions,
				"(u.track_id LIKE ? OR u.filename LIKE ? OR u.file_type LIKE ? OR COALESCE(u.user_id, '') LIKE ? OR COALESCE(u.username, '') LIKE ? OR COALESCE(u.source, '') LIKE ? OR COALESCE(u.source_id, '') LIKE ? OR COALESCE(u.detector, '') LIKE ? OR CAST(u.id AS TEXT) LIKE ? OR strftime('%Y-%m-%d %H:%M:%S', u.recording_date, 'unixepoch') LIKE ? OR COALESCE(u.comment, '') LIKE ?)")
		}
	}

	// Build the base query - read detector directly from uploads table
	// LEFT JOIN with users to get internal user info
	baseSelect := `
		SELECT u.id, u.filename, u.file_type, u.track_id, u.file_size, u.upload_ip,`

	if db.Driver == "pgx" || db.Driver == "duckdb" {
		baseSelect += `
		       EXTRACT(EPOCH FROM u.created_at)::BIGINT,
		       COALESCE(EXTRACT(EPOCH FROM u.recording_date)::BIGINT, 0) as recording_date,
		       u.source, u.source_id, u.source_url, u.user_id, u.username,
		       COALESCE(u.detector, '') as detector, u.internal_user_id,
		       usr.username as internal_username, usr.email as internal_email,
		       COALESCE(u.name, '') as name, COALESCE(u.notes, '') as notes,
		       COALESCE(u.comment, '') as comment
		FROM uploads u
		LEFT JOIN users usr ON u.internal_user_id = usr.id::text`
	} else {
		baseSelect += `
		       u.created_at,
		       COALESCE(u.recording_date, 0) as recording_date,
		       u.source, u.source_id, u.source_url, u.user_id, u.username,
		       COALESCE(u.detector, '') as detector, u.internal_user_id,
		       usr.username as internal_username, usr.email as internal_email,
		       COALESCE(u.name, '') as name, COALESCE(u.notes, '') as notes,
		       COALESCE(u.comment, '') as comment
		FROM uploads u
		LEFT JOIN users usr ON u.internal_user_id = CAST(usr.id AS TEXT)`
	}

	// Add WHERE clause if we have conditions
	if len(whereConditions) > 0 {
		query = baseSelect + "\nWHERE " + whereConditions[0]
		for i := 1; i < len(whereConditions); i++ {
			query += " AND " + whereConditions[i]
		}
	} else {
		query = baseSelect
	}

	query += "\nORDER BY u.created_at DESC\n"

	// Add LIMIT and OFFSET with correct parameter numbering
	if db.Driver == "pgx" || db.Driver == "duckdb" {
		paramCount++
		limitParam := paramCount
		paramCount++
		offsetParam := paramCount
		query += fmt.Sprintf("LIMIT $%d OFFSET $%d", limitParam, offsetParam)
	} else {
		query += "LIMIT ? OFFSET ?"
	}

	// Build args array
	if userID != "" {
		args = append(args, userID)
	}
	if search != "" {
		searchPattern := "%" + search + "%"
		if db.Driver == "pgx" || db.Driver == "duckdb" {
			args = append(args, searchPattern)
		} else {
			// SQLite needs the pattern repeated once per ? placeholder (track_id, filename, file_type, user_id, username, source, source_id, detector, id, recording_date, comment)
			for i := 0; i < 11; i++ {
				args = append(args, searchPattern)
			}
		}
	}
	args = append(args, limit, offset)

	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query uploads: %w", err)
	}
	defer rows.Close()

	var uploads []Upload
	for rows.Next() {
		var u Upload
		var source, sourceID, sourceURL, userID, username, detector, internalUserID, internalUsername, internalEmail *string
		var name, notes, comment *string
		err := rows.Scan(
			&u.ID, &u.Filename, &u.FileType, &u.TrackID,
			&u.FileSize, &u.UploadIP, &u.CreatedAt, &u.RecordingDate,
			&source, &sourceID, &sourceURL, &userID, &username, &detector, &internalUserID,
			&internalUsername, &internalEmail, &name, &notes, &comment,
		)
		if err != nil {
			continue
		}
		if source != nil {
			u.Source = *source
		}
		if sourceID != nil {
			u.SourceID = *sourceID
		}
		if sourceURL != nil {
			u.SourceURL = *sourceURL
		}
		if detector != nil {
			u.Detector = *detector
		}
		if userID != nil {
			u.UserID = *userID
		}
		if username != nil {
			u.Username = *username
		}
		if internalUserID != nil {
			u.InternalUserID = *internalUserID
		}
		// If we have an internal user, prefer their username over the external one
		if internalUsername != nil && *internalUsername != "" {
			u.Username = *internalUsername
		}
		if name != nil {
			u.Name = *name
		}
		if notes != nil {
			u.Notes = *notes
		}
		if comment != nil {
			u.Comment = *comment
		}
		uploads = append(uploads, u)
	}

	return uploads, nil
}

// DeleteTrack removes all data associated with a track ID.
// This includes markers, spectra, and upload records.
func (db *Database) DeleteTrack(ctx context.Context, trackID string) error {
	// Delete from spectra first (via CASCADE from markers)
	// Then delete markers, tracks, and upload records

	queries := []string{
		"DELETE FROM markers WHERE trackID = ?",
		"DELETE FROM tracks WHERE trackID = ?",
		"DELETE FROM uploads WHERE track_id = ?",
	}

	if db.Driver == "pgx" || db.Driver == "duckdb" {
		queries = []string{
			"DELETE FROM markers WHERE trackID = $1",
			"DELETE FROM tracks WHERE trackID = $1",
			"DELETE FROM uploads WHERE track_id = $1",
		}
	}

	for _, query := range queries {
		_, err := db.DB.ExecContext(ctx, query, trackID)
		if err != nil {
			return fmt.Errorf("delete track data: %w", err)
		}
	}

	return nil
}

// TrackBelongsToUser reports whether the given track was uploaded by internalUserID.
func (db *Database) TrackBelongsToUser(ctx context.Context, trackID, internalUserID string) (bool, error) {
	if trackID == "" || internalUserID == "" {
		return false, nil
	}

	query := "SELECT 1 FROM uploads WHERE track_id = ? AND internal_user_id = ? LIMIT 1"
	if db.Driver == "pgx" || db.Driver == "duckdb" {
		query = "SELECT 1 FROM uploads WHERE track_id = $1 AND internal_user_id = $2 LIMIT 1"
	}

	var exists int
	err := db.DB.QueryRowContext(ctx, query, trackID, internalUserID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckImportExists returns true if a Safecast import ID has already been imported AND has data points.
func (db *Database) CheckImportExists(ctx context.Context, sourceType string, importID int64) (bool, error) {
	var query string
	var args []interface{}

	sourceIDStr := fmt.Sprintf("%d", importID)

	switch db.Driver {
	case "pgx", "duckdb":
		// Check if upload exists AND has markers (actual data points)
		query = `SELECT COUNT(*) FROM uploads u
		         WHERE u.source = $1 AND u.source_id = $2
		         AND EXISTS (SELECT 1 FROM markers m WHERE m.trackid = u.track_id LIMIT 1)`
		args = []interface{}{sourceType, sourceIDStr}
	case "sqlite", "chai":
		// Check if upload exists AND has markers (actual data points)
		query = `SELECT COUNT(*) FROM uploads u
		         WHERE u.source = ? AND u.source_id = ?
		         AND EXISTS (SELECT 1 FROM markers m WHERE m.trackid = u.track_id LIMIT 1)`
		args = []interface{}{sourceType, sourceIDStr}
	default:
		return false, fmt.Errorf("unsupported database driver: %s", db.Driver)
	}

	var count int
	err := db.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("check import exists: %w", err)
	}

	return count > 0, nil
}

// GetLastImportedSafecastID returns the highest Safecast import ID that has been processed.
func (db *Database) GetLastImportedSafecastID(ctx context.Context, sourceType string) (int64, error) {
	var query string
	var args []interface{}

	switch db.Driver {
	case "pgx", "duckdb":
		// PostgreSQL and DuckDB: CAST to INTEGER
		query = `
			SELECT COALESCE(MAX(CAST(source_id AS INTEGER)), 0)
			FROM uploads
			WHERE source = $1 AND source_id IS NOT NULL AND source_id != ''
		`
		args = []interface{}{sourceType}
	case "sqlite", "chai":
		// SQLite: CAST to INTEGER
		query = `
			SELECT COALESCE(MAX(CAST(source_id AS INTEGER)), 0)
			FROM uploads
			WHERE source = ? AND source_id IS NOT NULL AND source_id != ''
		`
		args = []interface{}{sourceType}
	default:
		return 0, fmt.Errorf("unsupported database driver: %s", db.Driver)
	}

	var lastID int64
	err := db.DB.QueryRowContext(ctx, query, args...).Scan(&lastID)
	if err != nil {
		return 0, fmt.Errorf("get last imported safecast ID: %w", err)
	}

	return lastID, nil
}

// GetMinImportedSafecastID returns the lowest Safecast import ID that has been processed.
// This is used for backfill to know where to continue from.
func (db *Database) GetMinImportedSafecastID(ctx context.Context, sourceType string) (int64, error) {
	var query string
	var args []interface{}

	switch db.Driver {
	case "pgx", "duckdb":
		query = `
			SELECT COALESCE(MIN(CAST(source_id AS INTEGER)), 0)
			FROM uploads
			WHERE source = $1 AND source_id IS NOT NULL AND source_id != ''
		`
		args = []interface{}{sourceType}
	case "sqlite", "chai":
		query = `
			SELECT COALESCE(MIN(CAST(source_id AS INTEGER)), 0)
			FROM uploads
			WHERE source = ? AND source_id IS NOT NULL AND source_id != ''
		`
		args = []interface{}{sourceType}
	default:
		return 0, fmt.Errorf("unsupported database driver: %s", db.Driver)
	}

	var minID int64
	err := db.DB.QueryRowContext(ctx, query, args...).Scan(&minID)
	if err != nil {
		return 0, fmt.Errorf("get min imported safecast ID: %w", err)
	}

	return minID, nil
}

// CountImportsBySource returns the number of uploads from a specific source.
func (db *Database) CountImportsBySource(ctx context.Context, sourceType string) (int, error) {
	var query string
	var args []interface{}

	switch db.Driver {
	case "pgx", "duckdb":
		query = `SELECT COUNT(*) FROM uploads WHERE source = $1`
		args = []interface{}{sourceType}
	case "sqlite", "chai":
		query = `SELECT COUNT(*) FROM uploads WHERE source = ?`
		args = []interface{}{sourceType}
	default:
		return 0, fmt.Errorf("unsupported database driver: %s", db.Driver)
	}

	var count int
	err := db.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count imports by source: %w", err)
	}

	return count, nil
}

// GetLatestImportDate returns the latest created_at date for imports from a specific source.
// Returns empty string if no imports exist.
func (db *Database) GetLatestImportDate(ctx context.Context, sourceType string) (string, error) {
	var query string
	var args []interface{}

	switch db.Driver {
	case "pgx":
		query = `SELECT TO_CHAR(MAX(created_at), 'YYYY-MM-DD') FROM uploads WHERE source = $1`
		args = []interface{}{sourceType}
	case "duckdb":
		query = `SELECT strftime(MAX(created_at), '%Y-%m-%d') FROM uploads WHERE source = $1`
		args = []interface{}{sourceType}
	case "sqlite", "chai":
		query = `SELECT strftime('%Y-%m-%d', MAX(created_at), 'unixepoch') FROM uploads WHERE source = ?`
		args = []interface{}{sourceType}
	default:
		return "", fmt.Errorf("unsupported database driver: %s", db.Driver)
	}

	var dateStr sql.NullString
	err := db.DB.QueryRowContext(ctx, query, args...).Scan(&dateStr)
	if err != nil {
		return "", fmt.Errorf("get latest import date: %w", err)
	}

	if !dateStr.Valid {
		return "", nil
	}

	return dateStr.String, nil
}

// CountUploads returns the total number of uploads, optionally filtered by user_id and search term
func (db *Database) CountUploads(ctx context.Context, userID string, search string) (int, error) {
	var query string
	var args []interface{}
	var whereConditions []string
	paramCount := 0

	// Build WHERE conditions
	if userID != "" {
		paramCount++
		if db.Driver == "pgx" || db.Driver == "duckdb" {
			whereConditions = append(whereConditions, fmt.Sprintf("internal_user_id = $%d", paramCount))
		} else {
			whereConditions = append(whereConditions, "internal_user_id = ?")
		}
		args = append(args, userID)
	}

	if search != "" {
		paramCount++
		searchPattern := "%" + search + "%"
		if db.Driver == "pgx" || db.Driver == "duckdb" {
			whereConditions = append(whereConditions, fmt.Sprintf(
				"(track_id ILIKE $%d OR filename ILIKE $%d OR file_type ILIKE $%d OR COALESCE(user_id, '') ILIKE $%d OR COALESCE(username, '') ILIKE $%d OR COALESCE(source, '') ILIKE $%d OR COALESCE(source_id, '') ILIKE $%d OR COALESCE(detector, '') ILIKE $%d OR CAST(id AS TEXT) ILIKE $%d OR TO_CHAR(recording_date, 'YYYY-MM-DD HH24:MI:SS') ILIKE $%d)",
				paramCount, paramCount, paramCount, paramCount, paramCount, paramCount, paramCount, paramCount, paramCount, paramCount))
			args = append(args, searchPattern)
		} else {
			whereConditions = append(whereConditions,
				"(track_id LIKE ? OR filename LIKE ? OR file_type LIKE ? OR COALESCE(user_id, '') LIKE ? OR COALESCE(username, '') LIKE ? OR COALESCE(source, '') LIKE ? OR COALESCE(source_id, '') LIKE ? OR COALESCE(detector, '') LIKE ? OR CAST(id AS TEXT) LIKE ? OR strftime('%Y-%m-%d %H:%M:%S', recording_date, 'unixepoch') LIKE ?)")
			for i := 0; i < 10; i++ {
				args = append(args, searchPattern)
			}
		}
	}

	// Build query
	query = "SELECT COUNT(*) FROM uploads"
	if len(whereConditions) > 0 {
		query += " WHERE " + whereConditions[0]
		for i := 1; i < len(whereConditions); i++ {
			query += " AND " + whereConditions[i]
		}
	}

	var count int
	err := db.DB.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count uploads: %w", err)
	}

	return count, nil
}
