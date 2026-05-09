package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
	"safecast-new-map/pkg/mapimport"
	"safecast-new-map/pkg/spectrum"
)

func enqueueArchiveImport(fh *multipart.FileHeader, trackID string, db *database.Database, dbType string) error {
	src, err := fh.Open()
	if err != nil {
		return fmt.Errorf("open upload: %w", err)
	}
	defer src.Close()

	tmp, err := os.CreateTemp("", "upload-*.tgz")
	if err != nil {
		return fmt.Errorf("temp tgz: %w", err)
	}
	defer tmp.Close()

	if _, err := io.Copy(tmp, src); err != nil {
		return fmt.Errorf("copy tgz: %w", err)
	}

	path := tmp.Name()

	go func(localPath string) {
		defer os.Remove(localPath)
		ctx := context.Background()
		logf := func(format string, v ...any) {
			logT(trackID, "Upload", format, v...)
		}
		logf("queued tgz import in background: %s", fh.Filename)
		if err := importArchiveFromFile(ctx, localPath, trackID, db, dbType, logf); err != nil {
			logf("background tgz import failed: %v", err)
		}
	}(path)

	return nil
}

// progressHandler streams upload progress via Server-Sent Events.
//
// @Summary     Stream upload progress
// @Description Streams upload progress events for a previously submitted upload track ID.
// @Tags        map
// @Produce     text/event-stream
// @Param       trackid query string true "Upload tracking ID returned by /upload"
// @Success     200 {string} string "SSE stream"
// @Failure     400 {string} string "Missing trackid"
// @Router      /upload/progress [get]
func progressHandler(w http.ResponseWriter, r *http.Request) {
	trackID := r.URL.Query().Get("trackid")
	if trackID == "" {
		http.Error(w, "trackid required", http.StatusBadRequest)
		return
	}

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	// Send progress updates every 300ms for smoother animation
	ticker := time.NewTicker(300 * time.Millisecond)
	defer ticker.Stop()

	timeout := time.After(10 * time.Minute)
	lastProgress := -1

	for {
		select {
		case <-timeout:
			return
		case <-r.Context().Done():
			return
		case <-ticker.C:
			uploadProgress.RLock()
			prog, exists := uploadProgress.tracks[trackID]
			uploadProgress.RUnlock()

			if !exists {
				// Send 0% if not started yet
				fmt.Fprintf(w, "data: {\"current\":0,\"total\":100,\"percent\":0,\"complete\":false}\n\n")
				flusher.Flush()
				continue
			}

			prog.mu.RLock()
			current := prog.Current
			total := prog.Total
			complete := prog.Complete
			errMsg := prog.Error
			redirectURL := prog.RedirectURL
			needsCoordinates := prog.NeedsCoordinates
			progTrackID := prog.TrackID
			fileName := prog.FileName
			pointCount := prog.PointCount
			dateFrom := prog.DateFrom
			dateTo := prog.DateTo
			avgDoseRate := prog.AvgDoseRate
			minDoseRate := prog.MinDoseRate
			maxDoseRate := prog.MaxDoseRate
			filenames := prog.Filenames
			spectrumMarkerID := prog.SpectrumMarkerID
			prog.mu.RUnlock()

			percent := 0
			if total > 0 {
				percent = int((float64(current) / float64(total)) * 100)
			} else if complete {
				// If complete but total is 0, show 100%
				percent = 100
			}

			// Send update if progress changed OR if complete
			if percent != lastProgress || complete {
				if errMsg != "" {
					fmt.Fprintf(w, "data: {\"current\":%d,\"total\":%d,\"percent\":%d,\"complete\":true,\"error\":%q}\n\n",
						current, total, percent, errMsg)
				} else if complete && needsCoordinates {
					type sseNeedsCoords struct {
						Current          int      `json:"current"`
						Total            int      `json:"total"`
						Percent          int      `json:"percent"`
						Complete         bool     `json:"complete"`
						RedirectURL      string   `json:"redirectURL"`
						NeedsCoordinates bool     `json:"needsCoordinates"`
						TrackID          string   `json:"trackID"`
						FileName         string   `json:"fileName"`
						PointCount       int      `json:"pointCount"`
						DateFrom         int64    `json:"dateFrom"`
						DateTo           int64    `json:"dateTo"`
						AvgDoseRate      float64  `json:"avgDoseRate"`
						MinDoseRate      float64  `json:"minDoseRate"`
						MaxDoseRate      float64  `json:"maxDoseRate"`
						Filenames        []string `json:"filenames"`
						SpectrumMarkerID int64    `json:"spectrumMarkerID,omitempty"`
					}
					payload, _ := json.Marshal(sseNeedsCoords{
						Current: current, Total: total, Percent: 100, Complete: true,
						RedirectURL: redirectURL, NeedsCoordinates: true,
						TrackID: progTrackID, FileName: fileName,
						PointCount:       pointCount,
						DateFrom:         dateFrom,
						DateTo:           dateTo,
						AvgDoseRate:      avgDoseRate,
						MinDoseRate:      minDoseRate,
						MaxDoseRate:      maxDoseRate,
						Filenames:        filenames,
						SpectrumMarkerID: spectrumMarkerID,
					})
					fmt.Fprintf(w, "data: %s\n\n", payload)
				} else if complete {
					type sseComplete struct {
						Current          int      `json:"current"`
						Total            int      `json:"total"`
						Percent          int      `json:"percent"`
						Complete         bool     `json:"complete"`
						RedirectURL      string   `json:"redirectURL"`
						PointCount       int      `json:"pointCount"`
						DateFrom         int64    `json:"dateFrom"`
						DateTo           int64    `json:"dateTo"`
						AvgDoseRate      float64  `json:"avgDoseRate"`
						MinDoseRate      float64  `json:"minDoseRate"`
						MaxDoseRate      float64  `json:"maxDoseRate"`
						Filenames        []string `json:"filenames"`
						SpectrumMarkerID int64    `json:"spectrumMarkerID,omitempty"`
					}
					payload, _ := json.Marshal(sseComplete{
						Current: current, Total: total, Percent: 100, Complete: true,
						RedirectURL:      redirectURL,
						PointCount:       pointCount,
						DateFrom:         dateFrom,
						DateTo:           dateTo,
						AvgDoseRate:      avgDoseRate,
						MinDoseRate:      minDoseRate,
						MaxDoseRate:      maxDoseRate,
						Filenames:        filenames,
						SpectrumMarkerID: spectrumMarkerID,
					})
					fmt.Fprintf(w, "data: %s\n\n", payload)
				} else {
					fmt.Fprintf(w, "data: {\"current\":%d,\"total\":%d,\"percent\":%d,\"complete\":false}\n\n",
						current, total, percent)
				}
				flusher.Flush()
				lastProgress = percent
			}

			// Stop when complete
			if complete {
				return
			}
		}
	}
}

// bytesFile wraps bytes.Reader to implement multipart.File interface
type bytesFile struct {
	*bytes.Reader
}

func (b *bytesFile) Close() error {
	return nil
}

func newBytesFile(data []byte) *bytesFile {
	return &bytesFile{Reader: bytes.NewReader(data)}
}

// uploadHandler accepts multipart uploads and starts async processing.
//
// @Summary     Upload radiation files
// @Description Accepts one or more files in multipart form field `files[]` and returns a track ID for async progress polling.
// @Tags        map
// @Accept      mpfd
// @Produce     json
// @Success     200 {object} map[string]interface{} "Upload accepted"
// @Failure     400 {string} string "Invalid upload"
// @Failure     500 {string} string "Server error"
// @Router      /upload [post]
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent CloudFront from caching upload responses (user-specific, dynamic)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, "multipart parse error", http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["files[]"]
	if len(files) == 0 {
		http.Error(w, "no files selected", http.StatusBadRequest)
		return
	}

	trackID := GenerateSerialNumber()
	logT(trackID, "Upload", "▶ start, total=%d", len(files))

	// Read all files into memory so we can return quickly and process async
	type fileData struct {
		filename string
		ext      string
		size     int64
		content  []byte
	}
	var fileDataList []fileData
	for _, fh := range files {
		f, err := fh.Open()
		if err != nil {
			http.Error(w, "failed to open file", http.StatusInternalServerError)
			return
		}
		content, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			http.Error(w, "failed to read file", http.StatusInternalServerError)
			return
		}
		fileDataList = append(fileDataList, fileData{
			filename: fh.Filename,
			ext:      strings.ToLower(filepath.Ext(fh.Filename)),
			size:     fh.Size,
			content:  content,
		})
		logT(trackID, "Upload", "file received: %s (%d bytes)", fh.Filename, len(content))
	}

	// Get client IP for upload tracking
	clientIP := r.RemoteAddr
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		clientIP = forwarded
	}

	// Extract authenticated user if present
	var internalUserID string
	if user, ok := auth.GetUserFromContext(r.Context()); ok {
		internalUserID = fmt.Sprintf("%d", user.ID)
		displayName := user.Username
		if displayName == "" {
			displayName = user.Email
		}
		logT(trackID, "Upload", "authenticated user: %s (ID: %s, Email: %s, API Key: %s)",
			displayName, internalUserID, user.Email, user.APIKey)
	}

	// Initialize progress tracker
	uploadProgress.Lock()
	uploadProgress.tracks[trackID] = &UploadProgress{Total: 0, Current: 0, Complete: false}
	uploadProgress.Unlock()

	// Return immediately with trackID so client can start monitoring progress
	w.Header().Set("Content-Type", "application/json")
	response := map[string]any{
		"status":  "processing",
		"trackID": trackID,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("upload response write error: %v", err)
	}

	// Process files in background goroutine
	// Save original trackID for progress tracking - the trackID may change if duplicate detection finds existing track
	originalTrackID := trackID
	go func() {
		defer func() {
			// Cleanup progress tracker after 5 minutes
			time.AfterFunc(5*time.Minute, func() {
				uploadProgress.Lock()
				delete(uploadProgress.tracks, originalTrackID)
				uploadProgress.Unlock()
			})
		}()

		global := database.Bounds{MinLat: 90, MinLon: 180, MaxLat: -90, MaxLon: -180}
		hasBounds := false
		backgroundImport := false
		isSpectrumUpload := false
		var lastError error
		var successFilenames []string

		for _, fd := range fileDataList {
			reader := newBytesFile(fd.content)

			var (
				bbox database.Bounds
				err  error
			)
			switch fd.ext {
			case ".kml":
				bbox, trackID, err = processKMLFile(reader, trackID, db, *dbType)
			case ".kmz":
				bbox, trackID, err = processKMZFile(reader, trackID, db, *dbType)
			case ".gpx":
				bbox, trackID, err = processGPXFile(reader, trackID, db, *dbType)
			case ".csv":
				bbox, trackID, err = processCSVFile(reader, trackID, db, *dbType)
			case ".rctrk":
				bbox, trackID, err = processRCTRKFile(reader, trackID, db, *dbType)
			case ".n42":
				isSpectrumUpload = true
				bbox, trackID, err = processN42File(reader, fd.filename, trackID, db, *dbType)
			case ".spe":
				isSpectrumUpload = true
				bbox, trackID, err = processSPEFile(reader, fd.filename, trackID, db, *dbType)
			case ".xml":
				// Check if this is RadiaCode XML format
				if spectrum.IsRCXMLFormat(fd.content) {
					isSpectrumUpload = true
					bbox, trackID, err = processRCXMLFile(reader, fd.filename, trackID, db, *dbType)
				} else {
					logT(trackID, "Upload", "unsupported XML format: %s", fd.filename)
					lastError = fmt.Errorf("unsupported XML format (not RadiaCode)")
					continue
				}
			case ".cim":
				var imported bool
				bbox, trackID, imported, err = processTrackExportFile(context.Background(), reader, trackID, db, *dbType)
				if err == nil && !imported {
					logT(trackID, "Upload", "skipped duplicate legacy payload")
				}
			case ".json":
				var imported bool
				bbox, trackID, imported, err = processTrackExportPayload(context.Background(), fd.content, trackID, db, *dbType)
				if errors.Is(err, mapimport.ErrInvalidExportFormat) {
					bbox, trackID, err = processSafecastTrackJSON(fd.content, trackID, db, *dbType)
				}
				if errors.Is(err, errNotSafecastTrackJSON) {
					bbox, trackID, err = processAtomFastData(fd.content, trackID, db, *dbType)
				}
				if err == nil && !imported {
					logT(trackID, "Upload", "skipped duplicate track export JSON")
				}
			case ".tgz":
				backgroundImport = true
				// For tgz, we need to create a temp file since enqueueArchiveImport expects a file header
				logT(trackID, "Upload", "tgz archive - processing inline")
				continue
			case ".log", ".txt":
				bbox, trackID, err = processBGeigieZenFile(reader, trackID, db, *dbType)
			default:
				// Sniff for bGeigie $BNRDD content
				content := string(fd.content[:min(1024, len(fd.content))])
				if strings.Contains(content, "$BNRDD") {
					logT(trackID, "Upload", "sniffed bGeigie content for %s (ext=%q)", fd.filename, fd.ext)
					bbox, trackID, err = processBGeigieZenFile(newBytesFile(fd.content), trackID, db, *dbType)
				} else if strings.HasSuffix(strings.ToLower(fd.filename), ".tar.gz") {
					backgroundImport = true
					continue
				} else {
					logT(trackID, "Upload", "unsupported file type: %s (ext=%q)", fd.filename, fd.ext)
					lastError = fmt.Errorf("unsupported file type: %s", fd.ext)
					continue
				}
			}

			if err != nil {
				logT(trackID, "Upload", "error processing %s: %v", fd.filename, err)
				lastError = err
				continue
			}

			// Track upload in database
			var recordingDate int64
			var minDateQuery string
			if *dbType == "pgx" {
				minDateQuery = `SELECT MIN(date) FROM markers WHERE trackID = $1`
			} else {
				minDateQuery = `SELECT MIN(date) FROM markers WHERE trackID = ?`
			}
			_ = db.DB.QueryRowContext(context.Background(), minDateQuery, trackID).Scan(&recordingDate)

			// Get detector from markers (first non-empty value)
			var detector string
			var detectorQuery string
			if *dbType == "pgx" {
				detectorQuery = `SELECT COALESCE(detector, '') FROM markers WHERE trackID = $1 AND detector != '' LIMIT 1`
			} else {
				detectorQuery = `SELECT COALESCE(detector, '') FROM markers WHERE trackID = ? AND detector != '' LIMIT 1`
			}
			_ = db.DB.QueryRowContext(context.Background(), detectorQuery, trackID).Scan(&detector)

			upload := database.Upload{
				Filename:       fd.filename,
				FileType:       fd.ext,
				TrackID:        trackID,
				FileSize:       fd.size,
				UploadIP:       clientIP,
				RecordingDate:  recordingDate,
				CreatedAt:      0,
				Detector:       detector,
				InternalUserID: internalUserID, // Internal user ID from authenticated session
				Source:         "user-upload",
			}
			logT(trackID, "Upload", "inserting upload record: user_id=%s, filename=%s", internalUserID, fd.filename)
			if _, uploadErr := db.InsertUpload(context.Background(), upload); uploadErr != nil {
				logT(trackID, "Upload", "warning: failed to track upload: %v", uploadErr)
			} else {
				logT(trackID, "Upload", "✓ upload record inserted successfully")
			}

			successFilenames = append(successFilenames, fd.filename)

			// Update bounds
			if bbox.MinLat != 90 || bbox.MaxLat != -90 || bbox.MinLon != 180 || bbox.MaxLon != -180 {
				hasBounds = true
				if bbox.MinLat < global.MinLat {
					global.MinLat = bbox.MinLat
				}
				if bbox.MaxLat > global.MaxLat {
					global.MaxLat = bbox.MaxLat
				}
				if bbox.MinLon < global.MinLon {
					global.MinLon = bbox.MinLon
				}
				if bbox.MaxLon > global.MaxLon {
					global.MaxLon = bbox.MaxLon
				}
			}
		}

		// Query summary stats for the imported track
		var summaryPointCount int
		var summaryDateFrom, summaryDateTo int64
		var summaryAvg, summaryMin, summaryMax float64
		var summarySpectrumMarkerID int64
		if lastError == nil && len(successFilenames) > 0 {
			var q string
			if *dbType == "pgx" {
				q = `SELECT COUNT(*), COALESCE(MIN(date),0), COALESCE(MAX(date),0),
				     COALESCE(AVG(doseRate),0), COALESCE(MIN(doseRate),0), COALESCE(MAX(doseRate),0)
				     FROM markers WHERE trackID = $1`
			} else {
				q = `SELECT COUNT(*), COALESCE(MIN(date),0), COALESCE(MAX(date),0),
				     COALESCE(AVG(doseRate),0), COALESCE(MIN(doseRate),0), COALESCE(MAX(doseRate),0)
				     FROM markers WHERE trackID = ?`
			}
			_ = db.DB.QueryRowContext(context.Background(), q, trackID).Scan(
				&summaryPointCount, &summaryDateFrom, &summaryDateTo,
				&summaryAvg, &summaryMin, &summaryMax,
			)
			// Find the first marker with spectrum data for this track
			var specQ string
			if *dbType == "pgx" {
				specQ = `SELECT id FROM markers WHERE trackID = $1 AND has_spectrum = true ORDER BY id LIMIT 1`
			} else {
				specQ = `SELECT id FROM markers WHERE trackID = ? AND has_spectrum = 1 ORDER BY id LIMIT 1`
			}
			_ = db.DB.QueryRowContext(context.Background(), specQ, trackID).Scan(&summarySpectrumMarkerID)
		}

		// Build redirect URL
		needsCoordinates := false
		if hasBounds && global.MinLat == 0 && global.MaxLat == 0 && global.MinLon == 0 && global.MaxLon == 0 {
			needsCoordinates = true
			logT(trackID, "Upload", "spectrum file uploaded without GPS coordinates - needs manual location")
		}

		trackURL := "/"
		if hasBounds && !needsCoordinates {
			trackURL = fmt.Sprintf(
				"/trackid/%s?minLat=%f&minLon=%f&maxLat=%f&maxLon=%f&zoom=18&layer=%s",
				trackID, global.MinLat, global.MinLon, global.MaxLat, global.MaxLon,
				"OpenStreetMap")
		} else if isSpectrumUpload && needsCoordinates {
			trackURL = fmt.Sprintf("/trackid/%s", trackID)
		} else if backgroundImport {
			trackURL = "/"
		}

		// Mark complete with redirect URL or error
		// Use originalTrackID because that's what the client is monitoring
		uploadProgress.Lock()
		if prog, exists := uploadProgress.tracks[originalTrackID]; exists {
			prog.mu.Lock()
			prog.Complete = true
			if lastError != nil {
				prog.Error = lastError.Error()
			} else {
				prog.RedirectURL = trackURL
				prog.PointCount = summaryPointCount
				prog.DateFrom = summaryDateFrom
				prog.DateTo = summaryDateTo
				prog.AvgDoseRate = summaryAvg
				prog.MinDoseRate = summaryMin
				prog.MaxDoseRate = summaryMax
				prog.Filenames = successFilenames
				prog.SpectrumMarkerID = summarySpectrumMarkerID
				if isSpectrumUpload && needsCoordinates {
					prog.NeedsCoordinates = true
					prog.TrackID = trackID
					// Get first filename for display
					if len(fileDataList) > 0 {
						prog.FileName = fileDataList[0].filename
					}
				}
			}
			prog.mu.Unlock()
		}
		uploadProgress.Unlock()

		if lastError != nil {
			logT(trackID, "Upload", "✘ completed with error: %v", lastError)
		} else {
			logT(trackID, "Upload", "✔ completed, redirecting to: %s", trackURL)
		}
	}()
}

