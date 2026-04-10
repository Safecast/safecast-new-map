package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"safecast-new-map/pkg/database"
	"safecast-new-map/pkg/httpapi"
)

// =====================
// Парсинг файлов
// =====================

// countingReader forwards Read calls while emitting byte deltas over a channel.
// We prefer this tiny helper over mutex-protected counters so the progress
// logger can stay decoupled and responsive even when the network stream stalls.
type countingReader struct {
	r       io.Reader
	updates chan<- int64
}

func (c *countingReader) Read(p []byte) (int, error) {
	n, err := c.r.Read(p)
	if n > 0 && c.updates != nil {
		select {
		case c.updates <- int64(n):
		default:
		}
	}
	return n, err
}

// logArchiveImportProgress aggregates download and parse progress for remote
// tgz imports. A ticker throttles updates so huge archives do not overwhelm the
// logs while still giving operators confidence that the stream is moving.
func logArchiveImportProgress(
	ctx context.Context,
	logf func(string, ...any),
	source string,
	contentLength int64,
	byteUpdates <-chan int64,
	entryUpdates <-chan archiveProgress,
	done chan<- struct{},
) {
	defer close(done)

	// Emit an immediate status line so operators see that the goroutine is alive
	// even before bytes start flowing. This keeps "no logs" confusion at bay when
	// network buffers stall at the start of a long transfer.
	logf("tgz import [%s]: starting (size=%d bytes)", source, contentLength)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	var downloaded int64
	var entries int
	lastFile := ""
	lastLoggedBytes := int64(-1)
	lastLoggedEntries := -1
	lastLoggedFile := ""
	lastLoggedPercent := -1.0
	var lastLogTime time.Time
	lastLogLine := ""

	byteCh := byteUpdates
	entryCh := entryUpdates

	logSnapshot := func(force bool) {
		percent := float64(0)
		if contentLength > 0 && downloaded > 0 {
			percent = (float64(downloaded) / float64(contentLength)) * 100
		}

		progressed := percent != lastLoggedPercent || entries != lastLoggedEntries || lastFile != lastLoggedFile
		enoughTime := lastLogTime.IsZero() || time.Since(lastLogTime) >= 30*time.Second
		percentJump := percent >= 0 && lastLoggedPercent >= 0 && (percent-lastLoggedPercent) >= 0.5
		byteJump := contentLength == 0 && lastLoggedBytes >= 0 && (downloaded-lastLoggedBytes) >= 8*1024*1024

		if !force && (!progressed || !(enoughTime || percentJump || byteJump)) {
			return
		}

		line := fmt.Sprintf("tgz import [%s]: %.1f%% (%d/%d bytes) entries=%d last=%s", source, percent, downloaded, contentLength, entries, lastFile)

		// Avoid emitting identical consecutive snapshots so operators do not see doubled lines
		// when the final forced update matches the last timed tick.
		if line == lastLogLine {
			return
		}

		logf("%s", line)
		lastLogLine = line
		lastLoggedBytes = downloaded
		lastLoggedEntries = entries
		lastLoggedFile = lastFile
		lastLoggedPercent = percent
		lastLogTime = time.Now()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case delta, ok := <-byteCh:
			if !ok {
				byteCh = nil
				continue
			}
			downloaded += delta
		case progress, ok := <-entryCh:
			if !ok {
				entryCh = nil
				continue
			}
			entries = progress.entries
			lastFile = progress.filename
		case <-ticker.C:
			logSnapshot(false)
		}

		// When both channels close we break out so the caller sees the done
		// signal instead of blocking forever on a goroutine that cannot progress.
		if byteCh == nil && entryCh == nil {
			break
		}
	}

	logSnapshot(true)
}

// queueDuckDBMaintenanceAfterImport schedules the maintenance pass that keeps DuckDB files
// compact after large TGZ loads. We keep the coordination asynchronous so uploads and HTTP
// handlers stay responsive, while the database package serialises the actual PRAGMA calls.
func queueDuckDBMaintenanceAfterImport(driver string, db *database.Database, logf func(string, ...any), label string) {
	if !strings.EqualFold(strings.TrimSpace(driver), "duckdb") || db == nil {
		return
	}
	if logf == nil {
		logf = func(string, ...any) {}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := db.ScheduleDuckDBMaintenance(ctx, logf)
	go func(ch <-chan error) {
		if ch == nil {
			return
		}
		if err, ok := <-ch; ok {
			if err != nil {
				logf("duckdb maintenance after %s failed: %v", label, err)
				return
			}
			logf("duckdb maintenance after %s finished", label)
		}
	}(done)
}

// importArchiveFromFile streams a local tgz through the shared parser so offline
// operators can preload bundles without relying on HTTP. The function mirrors
// the remote helper by emitting byte and entry progress over channels, keeping
// the UI responsive on slow disks without extra mutexes.
func importArchiveFromFile(
	ctx context.Context,
	path string,
	trackID string,
	db *database.Database,
	dbType string,
	logf func(string, ...any),
) error {
	if logf == nil {
		logf = func(string, ...any) {}
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open tgz file: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("stat tgz file: %w", err)
	}

	bytesCh := make(chan int64, 256)
	entriesCh := make(chan archiveProgress, 64)
	done := make(chan struct{})
	go logArchiveImportProgress(ctx, logf, path, info.Size(), bytesCh, entriesCh, done)

	reader := &countingReader{r: file, updates: bytesCh}
	bounds, finalTrack, imported, err := processTrackExportArchiveReader(ctx, reader, trackID, db, dbType, entriesCh)
	close(bytesCh)
	close(entriesCh)
	<-done
	if err != nil {
		return fmt.Errorf("local tgz import: %w", err)
	}

	logf("local tgz import complete: imported=%v track=%s bounds=%v", imported, finalTrack, bounds)
	if imported {
		queueDuckDBMaintenanceAfterImport(dbType, db, logf, path)
	}
	return nil
}

// importArchiveFromURL streams a remote tgz into the existing archive parser so
// operators can refresh a deployment directly from an external weekly bundle.
// The function runs synchronously to match the explicit CLI flag and exits the
// program after finishing, keeping behaviour predictable for automation.
func importArchiveFromURL(
	ctx context.Context,
	sourceURL string,
	trackID string,
	db *database.Database,
	dbType string,
	logf func(string, ...any),
) error {
	if logf == nil {
		logf = func(string, ...any) {}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("download tgz: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("download tgz: unexpected status %s", resp.Status)
	}

	bytesCh := make(chan int64, 256)
	entriesCh := make(chan archiveProgress, 64)
	done := make(chan struct{})
	go logArchiveImportProgress(ctx, logf, sourceURL, resp.ContentLength, bytesCh, entriesCh, done)

	reader := &countingReader{r: resp.Body, updates: bytesCh}
	bounds, finalTrack, imported, err := processTrackExportArchiveReader(ctx, reader, trackID, db, dbType, entriesCh)
	close(bytesCh)
	close(entriesCh)
	<-done
	if err != nil {
		return fmt.Errorf("remote tgz import: %w", err)
	}

	logf("remote tgz import complete: imported=%v track=%s bounds=%v", imported, finalTrack, bounds)
	if imported {
		queueDuckDBMaintenanceAfterImport(dbType, db, logf, sourceURL)
	}
	return nil
}

// startBackgroundArchiveImport kicks off a non-blocking import pipeline so startup
// flags cannot prevent the HTTP listener from coming up. Progress and completion
// are reported through the provided logger, keeping the goroutine simple and free
// of shared state.
func startBackgroundArchiveImport(
	ctx context.Context,
	label string,
	importer func(context.Context) error,
	logf func(string, ...any),
) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		if logf == nil {
			logf = func(string, ...any) {}
		}
		logf("background tgz import queued: %s", label)
		if err := importer(ctx); err != nil {
			logf("background tgz import failed (%s): %v", label, err)
			return
		}
		logf("background tgz import finished (%s)", label)
	}()
	return done
}

// isSingleUserDriver reports whether the selected database driver relies on a
// single process-local connection. We gate import shielding on this so that
// multi-tenant engines like PostgreSQL keep serving traffic while the archive
// loader runs without any extra branching.
func isSingleUserDriver(dbType string) bool {
	switch strings.ToLower(strings.TrimSpace(dbType)) {
	case "sqlite", "chai", "duckdb":
		return true
	default:
		return false
	}
}

// importStillRunning checks the done channel without blocking, letting HTTP
// handlers quickly decide whether to short-circuit DB-heavy paths. Using a
// select keeps the check goroutine-free and stays faithful to "Don't
// communicate by sharing memory; share memory by communicating."
func importStillRunning(done <-chan struct{}) bool {
	if done == nil {
		return false
	}
	select {
	case <-done:
		return false
	default:
		return true
	}
}

// importShield builds a middleware that gently nudges DB-backed HTTP requests
// while an import is inflight against single-user file databases. The earlier
// version deprioritised uploads and URL shortener calls, but now that all
// engines flow through the serialized channel pipeline we can let requests
// proceed with a bounded deadline instead of outright rejecting them. This
// keeps user interactions responsive while still protecting the importer from
// unbounded queues.
// withMinimumDeadline ensures requests get a chance to wait in the serialized
// pipeline instead of failing instantly when an import is running. We only add
// a timeout when the caller has none or when it is too short to let the round-
// robin scheduler pick the job. This keeps the single-writer engines feeling
// multi-user without hiding cancellations from the client. The helper follows
// the Go Proverb "Don't communicate by sharing memory; share memory by
// communicating" by letting the channel-owned pipeline enforce ordering while
// we merely bound total wait time.
func importShield(done <-chan struct{}, dbType string, logf func(string, ...any)) func(http.Handler) http.Handler {
	if !isSingleUserDriver(dbType) || done == nil {
		return nil
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !importStillRunning(done) {
				next.ServeHTTP(w, r)
				return
			}

			// Keep the middleware non-blocking by always passing work through while
			// allowing enough time for the queue to drain. This mirrors the round-
			// robin lane scheduler in the database package so imports, uploads, and
			// reads all get a turn without starving each other.
			notice := "Идет импорт новых данных; ответы могут задерживаться."
			w.Header().Set("Retry-After", "1")
			w.Header().Set("X-Import-Notice", notice)

			// Give uploads and map queries a healthy window to reach the worker
			// goroutine instead of cancelling after one second when a TGZ import is
			// active. We still cap the wait to avoid hiding client disconnects.
			ctx, cancel := httpapi.WithMinimumDeadline(r.Context(), 2*time.Minute)
			defer cancel()
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// mergeBounds combines multiple bounding boxes while tracking whether we already
// have a baseline. This keeps archive imports from misreporting coordinates when
// the first few entries happen to be duplicates.
func mergeBounds(current database.Bounds, incoming database.Bounds, have bool) (database.Bounds, bool) {
	if incoming == (database.Bounds{}) {
		return current, have
	}
	if !have {
		return incoming, true
	}
	if incoming.MinLat < current.MinLat {
		current.MinLat = incoming.MinLat
	}
	if incoming.MinLon < current.MinLon {
		current.MinLon = incoming.MinLon
	}
	if incoming.MaxLat > current.MaxLat {
		current.MaxLat = incoming.MaxLat
	}
	if incoming.MaxLon > current.MaxLon {
		current.MaxLon = incoming.MaxLon
	}
	return current, true
}

func processSafecastTrackJSON(
	data []byte,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {
	var payload struct {
		Format  string `json:"format"`
		Version int    `json:"version"`
		Track   struct {
			TrackID        string   `json:"trackID"`
			DetectorName   string   `json:"detectorName"`
			DetectorType   string   `json:"detectorType"`
			RadiationTypes []string `json:"radiationTypes"`
		} `json:"track"`
		Markers []struct {
			ID                 int64    `json:"id"`
			TrackID            string   `json:"trackID"`
			TimeUnix           int64    `json:"timeUnix"`
			TimeUTC            string   `json:"timeUTC"`
			Lat                float64  `json:"lat"`
			Lon                float64  `json:"lon"`
			AltitudeM          *float64 `json:"altitudeM"`
			DoseMicroSvH       float64  `json:"doseRateMicroSvH"`
			DoseMicroRoentgenH float64  `json:"doseRateMicroRh"`
			DoseMilliSvH       float64  `json:"doseRateMilliSvH"`
			DoseMilliRH        float64  `json:"doseRateMilliRH"`
			CountRateCPS       float64  `json:"countRateCPS"`
			SpeedMS            float64  `json:"speedMS"`
			SpeedKMH           float64  `json:"speedKMH"`
			TemperatureC       *float64 `json:"temperatureC"`
			HumidityPercent    *float64 `json:"humidityPercent"`
			DetectorName       string   `json:"detectorName"`
			DetectorType       string   `json:"detectorType"`
			RadiationTypes     []string `json:"radiationTypes"`
		} `json:"markers"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return database.Bounds{}, trackID, errNotSafecastTrackJSON
	}
	if !strings.EqualFold(payload.Format, "safecast-track-json") {
		return database.Bounds{}, trackID, errNotSafecastTrackJSON
	}
	if len(payload.Markers) == 0 {
		return database.Bounds{}, trackID, fmt.Errorf("safecast track json: no markers")
	}

	candidateTrackID := strings.TrimSpace(payload.Track.TrackID)
	defaultDetectorType := strings.TrimSpace(payload.Track.DetectorType)
	defaultDetectorName := strings.TrimSpace(payload.Track.DetectorName)
	defaultRadiation := normalizeRadiationList(payload.Track.RadiationTypes)

	markers := make([]database.Marker, 0, len(payload.Markers))
	for _, item := range payload.Markers {
		ts := extractUnixSeconds(item.TimeUnix, item.TimeUTC)
		dose := item.DoseMicroSvH
		if dose == 0 && item.DoseMicroRoentgenH != 0 {
			dose = item.DoseMicroRoentgenH / microRoentgenPerMicroSievert
		}
		if dose == 0 && item.DoseMilliSvH != 0 {
			dose = item.DoseMilliSvH * 1000.0
		}
		if dose == 0 && item.DoseMilliRH != 0 {
			dose = item.DoseMilliRH * 10.0
		}

		speed := item.SpeedMS
		if speed == 0 && item.SpeedKMH != 0 {
			speed = item.SpeedKMH / 3.6
		}

		detectorName := strings.TrimSpace(item.DetectorName)
		if detectorName == "" {
			detectorName = defaultDetectorName
		}

		detector := strings.TrimSpace(item.DetectorType)
		if detector == "" {
			detector = defaultDetectorType
		}
		if detector == "" {
			detector = detectorTypeFromName(detectorName)
		}

		radiationList := normalizeRadiationList(item.RadiationTypes)
		if len(radiationList) == 0 {
			radiationList = defaultRadiation
		}

		var altitude float64
		var altitudeValid bool
		if item.AltitudeM != nil {
			altitude = *item.AltitudeM
			altitudeValid = true
		}
		var temperature float64
		var temperatureValid bool
		if item.TemperatureC != nil {
			temperature = *item.TemperatureC
			temperatureValid = true
		}
		var humidity float64
		var humidityValid bool
		if item.HumidityPercent != nil {
			humidity = *item.HumidityPercent
			humidityValid = true
		}

		markers = append(markers, database.Marker{
			ID:               item.ID,
			DoseRate:         dose,
			Date:             ts,
			Lon:              item.Lon,
			Lat:              item.Lat,
			CountRate:        item.CountRateCPS,
			Speed:            speed,
			Altitude:         altitude,
			Temperature:      temperature,
			Humidity:         humidity,
			Detector:         detector,
			Radiation:        strings.Join(radiationList, ","),
			AltitudeValid:    altitudeValid,
			TemperatureValid: temperatureValid,
			HumidityValid:    humidityValid,
		})

		if candidateTrackID == "" {
			candidateTrackID = strings.TrimSpace(item.TrackID)
		}
	}

	if candidateTrackID != "" {
		trackID = candidateTrackID
	}

	logT(trackID, "SafecastJSON", "parsed %d markers", len(markers))
	return processAndStoreMarkers(markers, trackID, db, dbType)
}

func extractUnixSeconds(timeUnix int64, timeUTC string) int64 {
	if timeUnix > 1_000_000_000_000 {
		return timeUnix / 1000
	}
	if timeUnix > 0 {
		return timeUnix
	}
	if strings.TrimSpace(timeUTC) != "" {
		if ts, err := time.Parse(time.RFC3339Nano, strings.TrimSpace(timeUTC)); err == nil {
			return ts.Unix()
		}
	}
	return 0
}

func normalizeRadiationList(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	out := make([]string, 0, len(values))
	for _, raw := range values {
		channel := strings.ToLower(strings.TrimSpace(raw))
		if channel == "" {
			continue
		}
		if _, ok := seen[channel]; ok {
			continue
		}
		seen[channel] = struct{}{}
		out = append(out, channel)
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// detectorTypeFromName extracts the type hint that we encode into detectorName
// during export. We slice on the first ':' because stableDetectorName prefixes
// the trackID before the reported detector model.
func detectorTypeFromName(detectorName string) string {
	detectorName = strings.TrimSpace(detectorName)
	if detectorName == "" {
		return ""
	}
	if idx := strings.Index(detectorName, ":"); idx >= 0 && idx+1 < len(detectorName) {
		candidate := strings.TrimSpace(detectorName[idx+1:])
		if candidate != "" {
			return candidate
		}
	}
	return ""
}

// parseBGeigieCoord parses coordinates that may have hemisphere suffix.
func parseBGeigieCoord(s string) float64 {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	r := s[len(s)-1]
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		base := s[:len(s)-1]
		v, _ := strconv.ParseFloat(base, 64)
		switch strings.ToUpper(string(r)) {
		case "S", "W":
			return -v
		default:
			return v
		}
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// parseDMM parses degrees+minutes (DDMM.MMMM or DDDMM.MMMM) with hemisphere.
func parseDMM(val, hemi string, degDigits int) float64 {
	val = strings.TrimSpace(val)
	hemi = strings.TrimSpace(hemi)
	if val == "" {
		return 0
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0
	}
	deg := int(f / 100.0)
	minutes := f - float64(deg*100)
	d := float64(deg) + minutes/60.0
	switch strings.ToUpper(hemi) {
	case "S", "W":
		d = -d
	}
	if degDigits == 2 {
		if d < -90 || d > 90 {
			return 0
		}
	} else {
		if d < -180 || d > 180 {
			return 0
		}
	}
	return d
}

// processAndStoreMarkers is the common pipeline:
// 0. bbox calculation             • O(N)
// 1. fast duplicate-track probe   • O(K·q), K ≪ N (early-exit)
// 2. assign final TrackID         • O(N)
// 3. basic filters                • O(N)
// 4. speed calculation            • O(N)
// 5. pre-compute 20 zoom levels   • O(N) parallel
// 6. batch-insert into DB         • one transaction, multi-row VALUES
//
// Concurrency-friendly: no mutexes; DB/sql pool handles connection safety.
func processAndStoreMarkers(
	markers []database.Marker,
	initTrackID string, // initially generated ID
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {
	// Keep legacy callers working by delegating to the context-aware variant.
	return processAndStoreMarkersWithContext(context.Background(), markers, initTrackID, db, dbType)
}

func processAndStoreMarkersWithContext(
	ctx context.Context,
	markers []database.Marker,
	initTrackID string, // initially generated ID
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {

	// ── step 0: bounding box (cheap) ────────────────────────────────
	bbox := database.Bounds{MinLat: 90, MinLon: 180, MaxLat: -90, MaxLon: -180}
	for _, m := range markers {
		if err := observeContext(ctx); err != nil {
			return bbox, initTrackID, err
		}
		if m.Lat < bbox.MinLat {
			bbox.MinLat = m.Lat
		}
		if m.Lat > bbox.MaxLat {
			bbox.MaxLat = m.Lat
		}
		if m.Lon < bbox.MinLon {
			bbox.MinLon = m.Lon
		}
		if m.Lon > bbox.MaxLon {
			bbox.MaxLon = m.Lon
		}
	}

	trackID := initTrackID

	// ── step 1: fast probe instead of full-scan ─────────────────────
	// Limit DB random lookups to a tiny sample (e.g. 128 points).
	if err := observeContext(ctx); err != nil {
		return bbox, trackID, err
	}

	probe := pickIdentityProbe(markers, 128)
	if existing, err := db.DetectExistingTrackID(probe, 10, dbType); err != nil {
		return bbox, trackID, err
	} else if existing != "" {
		logT(trackID, "Store", "⚠ detected existing trackID %s — reusing", existing)
		trackID = existing
	} else {
		logT(trackID, "Store", "unique track, proceed with new trackID")
	}

	// ── step 2: attach FINAL TrackID ────────────────────────────────
	for i := range markers {
		if err := observeContext(ctx); err != nil {
			return bbox, trackID, err
		}
		markers[i].TrackID = trackID
	}

	// ── step 3: light filters ───────────────────────────────────────
	markers = filterZeroMarkers(markers)
	if err := observeContext(ctx); err != nil {
		return bbox, trackID, err
	}
	markers = filterInvalidDateMarkers(markers)
	if len(markers) == 0 {
		return bbox, trackID, fmt.Errorf("all markers filtered out")
	}

	// Keep the track registry in sync so pagination queries avoid expensive DISTINCT scans.
	if err := db.EnsureTrackPresence(ctx, trackID, dbType); err != nil {
		return bbox, trackID, err
	}

	// ── step 4: speed calculation (pure Go) ─────────────────────────
	markers = calculateSpeedForMarkers(markers)
	if err := observeContext(ctx); err != nil {
		return bbox, trackID, err
	}

	// ── step 5: store only raw markers (zoom=0) ─────────────────────
	// OPTIMIZATION: Instead of pre-computing 20 zoom levels (10x storage),
	// we store only raw markers and cluster on-the-fly at query time.
	// This reduces database size by ~90% and speeds up uploads significantly.
	if err := observeContext(ctx); err != nil {
		return bbox, trackID, err
	}
	// Ensure all markers have zoom=0 (raw/original resolution)
	for i := range markers {
		markers[i].Zoom = 0
	}
	allZoom := markers
	logT(trackID, "Store", "storing %d raw markers (on-the-fly clustering enabled)", len(allZoom))

	// Initialize progress tracking using initTrackID (not trackID which may have changed)
	// This ensures the client monitoring the original trackID sees the progress updates
	uploadProgress.Lock()
	if prog, exists := uploadProgress.tracks[initTrackID]; exists {
		// Update existing progress tracker (created by upload handler)
		prog.mu.Lock()
		prog.Total = len(allZoom)
		prog.mu.Unlock()
	} else {
		// Create new progress tracker
		uploadProgress.tracks[initTrackID] = &UploadProgress{
			Total:   len(allZoom),
			Current: 0,
		}
	}
	uploadProgress.Unlock()

	progressCh := make(chan database.MarkerBatchProgress, 16)
	progressDone := make(chan struct{})

	go func(total int, progressTrackID string) {
		defer close(progressDone)
		started := time.Now()
		lastLog := time.Time{}
		for p := range progressCh {
			// Update global progress tracker using the original trackID
			uploadProgress.RLock()
			prog, exists := uploadProgress.tracks[progressTrackID]
			uploadProgress.RUnlock()
			if exists {
				prog.mu.Lock()
				prog.Current = p.Done
				prog.mu.Unlock()
			}

			if lastLog.IsZero() || time.Since(lastLog) >= 5*time.Second || p.Done >= total {
				logT(trackID, "Store", "storing markers %d/%d (+%d via %s) in %s elapsed %s",
					p.Done, p.Total, p.Batch, p.Mode, p.Duration.Truncate(time.Millisecond),
					time.Since(started).Truncate(time.Second))
				lastLog = time.Now()
			}
		}
	}(len(allZoom), initTrackID)

	// ── step 6: single transaction + multi-row VALUES ───────────────
	if err := observeContext(ctx); err != nil {
		close(progressCh)
		<-progressDone
		return bbox, trackID, err
	}
	if strings.EqualFold(dbType, "clickhouse") {
		if err := db.InsertMarkersBulk(ctx, nil, allZoom, dbType, 1000, progressCh, database.WorkloadUserUpload); err != nil {
			close(progressCh)
			<-progressDone
			return bbox, trackID, fmt.Errorf("bulk insert: %w", err)
		}
	} else {
		useTx := !(strings.EqualFold(dbType, "duckdb") || strings.EqualFold(dbType, "sqlite") || strings.EqualFold(dbType, "chai"))
		var tx *sql.Tx
		var err error
		if useTx {
			tx, err = db.DB.Begin()
			if err != nil {
				close(progressCh)
				<-progressDone
				return bbox, trackID, err
			}
			// Batch size 500–1000 usually gives a good balance on large B-Trees and keeps DuckDB in
			// a single transaction so inserts do not pay per-chunk commit costs.
			if err := observeContext(ctx); err != nil {
				_ = tx.Rollback()
				close(progressCh)
				<-progressDone
				return bbox, trackID, err
			}
		}
		if err := db.InsertMarkersBulk(ctx, tx, allZoom, dbType, 1000, progressCh, database.WorkloadUserUpload); err != nil {
			if tx != nil {
				_ = tx.Rollback()
			}
			close(progressCh)
			<-progressDone
			return bbox, trackID, fmt.Errorf("bulk insert: %w", err)
		}
		if tx != nil {
			if err := tx.Commit(); err != nil {
				close(progressCh)
				<-progressDone
				return bbox, trackID, err
			}
		}
	}

	close(progressCh)
	<-progressDone

	// Clean up progress tracker after a delay to allow final SSE read
	go func() {
		time.Sleep(2 * time.Second)
		uploadProgress.Lock()
		delete(uploadProgress.tracks, trackID)
		uploadProgress.Unlock()
	}()

	// Invalidate cache entries that might be affected by the new markers
	// Since we don't know exactly which tiles are affected, we clear the entire cache
	// In a production system, we might want to be more selective about cache invalidation
	tileCacheMu.Lock()
	tileCache.Purge() // Clear all cached entries
	tileCacheMu.Unlock()

	logT(trackID, "Store", "✔ stored (new %d markers)", len(allZoom))
	return bbox, trackID, nil
}

