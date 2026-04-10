package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/vmihailenco/msgpack/v5"
	"safecast-new-map/pkg/database"
	safecastrealtime "safecast-new-map/pkg/safecast-realtime"
)

func trackHandler(w http.ResponseWriter, r *http.Request) {
	lang := getPreferredLanguage(r)

	// /trackid/<ID>
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "TrackID not provided", http.StatusBadRequest)
		return
	}
	trackID := parts[2] // всё равно понадобится в JS

	// --- шаблон ----------------------------------------------------------------
	tmpl := template.Must(template.New("map.html").Funcs(template.FuncMap{
		"translate": func(key string) string {
			if v, ok := translations[lang][key]; ok {
				return v
			}
			return translations["en"][key]
		},
	}).ParseFS(content, "public_html/map.html"))

	translationsJSON, err := marshalTemplateJS(translationsForLang(translations, lang))
	if err != nil {
		log.Printf("track handler: marshal translations failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	emptyMarkers := []database.Marker{}
	markersJSON, err := marshalTemplateJS(emptyMarkers)
	if err != nil {
		log.Printf("track handler: marshal markers failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// отдаём пустой срез маркеров
	data := struct {
		Version           string
		Translations      map[string]map[string]string
		Lang              string
		DefaultLat        float64
		DefaultLon        float64
		DefaultZoom       int
		DefaultLayer      string
		AutoLocateDefault bool
		RealtimeAvailable bool
		SupportEmail      string
		TranslationsJSON  template.JS
		MarkersJSON       template.JS
		DebugEnabled      bool
	}{
		Version:           CompileVersion,
		Translations:      translations,
		Lang:              lang,
		DefaultLat:        *defaultLat,
		DefaultLon:        *defaultLon,
		DefaultZoom:       *defaultZoom,
		DefaultLayer:      *defaultLayer,
		AutoLocateDefault: *autoLocateDefault,
		RealtimeAvailable: *safecastRealtimeEnabled,
		SupportEmail:      strings.TrimSpace(*supportEmail),
		TranslationsJSON:  translationsJSON,
		MarkersJSON:       markersJSON,
		DebugEnabled:      debugEnabledForRequest(r),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := buf.WriteTo(w); err != nil {
		if isClientDisconnect(err) {
			log.Printf("client disconnected while writing response")
		} else {
			log.Printf("write resp: %v", err)
		}
	}

	// Ради отладки: показываем, что HTML отдали без тяжёлых данных
	log.Printf("Track page %s rendered.", trackID)
}

// tracksHandler serves the multi-track map page.
//
// @Summary     Render map page for multiple tracks
// @Description Serves the HTML map page with multiple track IDs from a comma-separated path segment.
// @Tags        web
// @Produce     html
// @Param       ids path string true "Comma-separated track IDs"
// @Success     200 {string} string "HTML page"
// @Failure     400 {string} string "Track IDs missing"
// @Router      /tracks/{ids} [get]
func tracksHandler(w http.ResponseWriter, r *http.Request) {
	lang := getPreferredLanguage(r)

	// /tracks/<IDs>
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Track IDs not provided", http.StatusBadRequest)
		return
	}
	trackIDsParam := parts[2]

	// Parse template
	tmpl := template.Must(template.New("map.html").Funcs(template.FuncMap{
		"translate": func(key string) string {
			if v, ok := translations[lang][key]; ok {
				return v
			}
			return translations["en"][key]
		},
	}).ParseFS(content, "public_html/map.html"))

	translationsJSON, err := marshalTemplateJS(translationsForLang(translations, lang))
	if err != nil {
		log.Printf("tracks handler: marshal translations failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	emptyMarkers := []database.Marker{}
	markersJSON, err := marshalTemplateJS(emptyMarkers)
	if err != nil {
		log.Printf("tracks handler: marshal markers failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Version           string
		Translations      map[string]map[string]string
		Lang              string
		DefaultLat        float64
		DefaultLon        float64
		DefaultZoom       int
		DefaultLayer      string
		AutoLocateDefault bool
		RealtimeAvailable bool
		SupportEmail      string
		TranslationsJSON  template.JS
		MarkersJSON       template.JS
		DebugEnabled      bool
	}{
		Version:           CompileVersion,
		Translations:      translations,
		Lang:              lang,
		DefaultLat:        *defaultLat,
		DefaultLon:        *defaultLon,
		DefaultZoom:       *defaultZoom,
		DefaultLayer:      *defaultLayer,
		AutoLocateDefault: *autoLocateDefault,
		RealtimeAvailable: *safecastRealtimeEnabled,
		SupportEmail:      strings.TrimSpace(*supportEmail),
		TranslationsJSON:  translationsJSON,
		MarkersJSON:       markersJSON,
		DebugEnabled:      debugEnabledForRequest(r),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := buf.WriteTo(w); err != nil {
		if isClientDisconnect(err) {
			log.Printf("client disconnected while writing response")
		} else {
			log.Printf("write resp: %v", err)
		}
	}

	log.Printf("Multi-track page rendered for: %s", trackIDsParam)
}

// apiTracksBoundsHandler, qrPngHandler in pkg/httpapi

// generateTileCacheKey creates a unique key for caching clustered markers based on request parameters
// The key includes zoom level, bounding box, track ID, speed filters, and date filters
func generateTileCacheKey(zoom int, minLat, minLon, maxLat, maxLon float64, trackID, trackIDsParam string, speedRanges []database.SpeedRange, dateFrom, dateTo int64) string {
	// Create a string representation of speed ranges
	speedStr := ""
	for _, sr := range speedRanges {
		speedStr += fmt.Sprintf("%.2f-%.2f,", sr.Min, sr.Max)
	}

	// Format the key with all relevant parameters
	return fmt.Sprintf("tile:%d:%.6f:%.6f:%.6f:%.6f:%s:%s:%s:%d:%d",
		zoom, minLat, minLon, maxLat, maxLon, trackID, trackIDsParam, speedStr, dateFrom, dateTo)
}

// getMarkersHandler returns markers for the requested map viewport.
//
// @Summary     Get markers for viewport
// @Description Returns filtered markers for map bounds and optional track/speed/time filters.
// @Tags        map
// @Produce     json
// @Param       zoom     query int    false "Requested map zoom"
// @Param       minLat   query number true  "Minimum latitude"
// @Param       minLon   query number true  "Minimum longitude"
// @Param       maxLat   query number true  "Maximum latitude"
// @Param       maxLon   query number true  "Maximum longitude"
// @Param       trackID  query string false "Single track ID filter"
// @Param       trackIDs query string false "Comma-separated track IDs filter"
// @Param       speeds   query string false "Comma-separated speed categories"
// @Param       dateFrom query int    false "Start unix timestamp (seconds)"
// @Param       dateTo   query int    false "End unix timestamp (seconds)"
// @Success     200 {array} map[string]interface{} "Markers"
// @Failure     500 {string} string "Server error"
// @Router      /get_markers [get]
func getMarkersHandler(w http.ResponseWriter, r *http.Request) {
	// Use the request context so map tiles cancel promptly when the browser closes,
	// freeing the serialized DuckDB lane for ongoing imports.
	ctx := r.Context()

	q := r.URL.Query()
	zoom, _ := strconv.Atoi(q.Get("zoom"))
	minLat, _ := strconv.ParseFloat(q.Get("minLat"), 64)
	minLon, _ := strconv.ParseFloat(q.Get("minLon"), 64)
	maxLat, _ := strconv.ParseFloat(q.Get("maxLat"), 64)
	maxLon, _ := strconv.ParseFloat(q.Get("maxLon"), 64)
	trackID := q.Get("trackID")
	trackIDsParam := q.Get("trackIDs") // Multiple tracks

	// ----- ✈️🚗🚶 фильтр скорости  ---------------------------------
	var sr []database.SpeedRange
	if s := q.Get("speeds"); s != "" {
		for _, tag := range strings.Split(s, ",") {
			if r, ok := speedCatalog[tag]; ok {
				sr = append(sr, database.SpeedRange(r))
			}
		}
	}
	if len(sr) == 0 && q.Get("speeds") != "" { // все выключены
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
		return
	}

	// ----- ⏱️  фильтр времени  ------------------------------------
	var (
		dateFrom int64
		dateTo   int64
	)
	if s := q.Get("dateFrom"); s != "" {
		dateFrom, _ = strconv.ParseInt(s, 10, 64)
	}
	if s := q.Get("dateTo"); s != "" {
		dateTo, _ = strconv.ParseInt(s, 10, 64)
	}

	// ----- запрос к БД  ------------------------------------------
	// ON-THE-FLY CLUSTERING: Always query zoom=0 (raw markers), then cluster client-side
	// This allows storing only raw markers in DB (~10x storage savings)
	const rawZoom = 0
	var (
		markers []database.Marker
		err     error
	)
	if trackID != "" {
		markers, err = db.GetMarkersByTrackIDZoomBoundsSpeed(
			ctx,
			trackID, rawZoom, minLat, minLon, maxLat, maxLon,
			dateFrom, dateTo, sr, *dbType)
	} else if trackIDsParam != "" {
		// Handle multiple tracks
		trackIDs := strings.Split(trackIDsParam, ",")
		for _, tid := range trackIDs {
			tid = strings.TrimSpace(tid)
			if tid == "" {
				continue
			}
			trackMarkers, trackErr := db.GetMarkersByTrackIDZoomBoundsSpeed(
				ctx,
				tid, rawZoom, minLat, minLon, maxLat, maxLon,
				dateFrom, dateTo, sr, *dbType)
			if trackErr != nil {
				log.Printf("Error fetching markers for track %s: %v", tid, trackErr)
				continue
			}
			markers = append(markers, trackMarkers...)
		}
	} else {
		markers, err = db.GetMarkersByZoomBoundsSpeed(
			ctx,
			rawZoom, minLat, minLon, maxLat, maxLon,
			dateFrom, dateTo, sr, *dbType)
	}
	if err != nil {
		log.Printf("Error fetching markers: %v (zoom=%d, bounds=[%f,%f,%f,%f], dbType=%s)", err, rawZoom, minLat, minLon, maxLat, maxLon, *dbType)
		http.Error(w, "Error fetching markers", http.StatusInternalServerError)
		return
	}

	// Generate cache key based on request parameters
	cacheKey := generateTileCacheKey(zoom, minLat, minLon, maxLat, maxLon, trackID, trackIDsParam, sr, dateFrom, dateTo)

	// Check if clustered markers are already in cache
	tileCacheMu.RLock()
	cachedMarkers, found := tileCache.Get(cacheKey)
	tileCacheMu.RUnlock()

	if found {
		// Use cached markers if available
		markers = cachedMarkers

		// Add realtime markers if enabled (these shouldn't be cached as they change frequently)
		if *safecastRealtimeEnabled {
			// We only touch realtime tables when the operator explicitly enables the feature.
			if rt, err := db.GetLatestRealtimeByBounds(ctx, minLat, minLon, maxLat, maxLon, *dbType); err == nil {
				for i := range rt {
					// Sanitise detector names on the fly so legacy rows without
					// the new resolver still produce friendly popups.
					rt[i].Tube = safecastrealtime.DetectorLabel(rt[i].Tube, rt[i].Transport, rt[i].DeviceName)
				}
				markers = append(markers, rt...)
			} else {
				log.Printf("realtime query: %v", err)
			}
		}
	} else {
		// Apply on-the-fly clustering based on requested zoom level (this is the expensive operation)
		markers = clusterMarkersForZoom(markers, zoom)

		// Add to cache for future requests (but don't cache if there are too many markers to avoid memory issues)
		if len(markers) < 10000 { // arbitrary limit to prevent cache bloat
			tileCacheMu.Lock()
			tileCache.Add(cacheKey, markers)
			tileCacheMu.Unlock()
		}

		// Add realtime markers if enabled
		if *safecastRealtimeEnabled {
			// We only touch realtime tables when the operator explicitly enables the feature.
			if rt, err := db.GetLatestRealtimeByBounds(ctx, minLat, minLon, maxLat, maxLon, *dbType); err == nil {
				for i := range rt {
					// Sanitise detector names on the fly so legacy rows without
					// the new resolver still produce friendly popups.
					rt[i].Tube = safecastrealtime.DetectorLabel(rt[i].Tube, rt[i].Transport, rt[i].DeviceName)
				}
				markers = append(markers, rt...)
			} else {
				log.Printf("realtime query: %v", err)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(markers)
}

// ========
// Streaming markers via SSE
// ========

// aggregateMarkers chooses the most radioactive marker per grid cell while merging the
// static query feed with an optional live stream. Keeping the grid map inside the goroutine
// lets us reuse previous emissions and drop later duplicates without mutexes.
// Markers with spectral data are prioritized over those without.
func aggregateMarkers(ctx context.Context, base <-chan database.Marker, updates <-chan database.Marker, zoom int) <-chan database.Marker {
	out := make(chan database.Marker)
	go func() {
		defer close(out)
		cells := make(map[string]database.Marker)
		scale := math.Pow(2, float64(zoom))
		baseCh := base
		updateCh := updates

		emit := func(m database.Marker) {
			key := fmt.Sprintf("%d:%d", int(m.Lat*scale), int(m.Lon*scale))
			prev, exists := cells[key]

			// Decide if this marker should replace the previous one:
			// 1. No previous marker exists, OR
			// 2. New marker has spectrum and previous doesn't, OR
			// 3. Both have same spectrum status and new has higher dose rate
			shouldReplace := !exists ||
				(m.HasSpectrum && !prev.HasSpectrum) ||
				(m.HasSpectrum == prev.HasSpectrum && m.DoseRate > prev.DoseRate)

			if shouldReplace {
				cells[key] = m
				select {
				case out <- m:
				case <-ctx.Done():
				}
			}
		}

		for {
			select {
			case <-ctx.Done():
				return
			case m, ok := <-baseCh:
				if !ok {
					baseCh = nil
					if baseCh == nil && updateCh == nil {
						return
					}
					continue
				}
				emit(m)
			case m, ok := <-updateCh:
				if !ok {
					updateCh = nil
					if baseCh == nil && updateCh == nil {
						return
					}
					continue
				}
				emit(m)
			}
		}
	}()
	return out
}

// calculateSampleRate determines what percentage of markers to send based on zoom level.
// Lower zoom = higher density = lower sample rate to reduce network traffic and client load.
func calculateSampleRate(zoom int, minLat, minLon, maxLat, maxLon float64) float64 {
	// Always send all markers for track view or high zoom levels
	if zoom >= 8 {
		return 1.0
	}

	// Calculate viewport area (approximate)
	latSpan := maxLat - minLat
	lonSpan := maxLon - minLon
	area := latSpan * lonSpan

	// For very large areas at low zoom, more aggressive sampling is needed
	// At zoom 7, Europe view is ~9.5 lat × 19.5 lon ≈ 185 sq degrees
	_ = area // Used for future density calculations

	// Sample rate based on zoom level to achieve target performance
	// More aggressive sampling for smoother performance
	// Zoom 7: ~15% (1 in 7) - increased for better visibility
	// Zoom 6: ~20% (1 in 5)  - increased for better visibility
	// Zoom 5 and below: ~2% (1 in 50) - down from 10%
	switch zoom {
	case 7:
		return 0.15 // Send 15% of markers (skip 85%)
	case 6:
		return 0.20 // Send 20% of markers (skip 80%)
	default: // zoom <= 5
		return 0.02 // Send 2% of markers (skip 98%)
	}
}

// sampleMarkerChannel applies statistical sampling to a marker stream.
// Returns a new channel that emits only a fraction of markers based on sampleRate.
func sampleMarkerChannel(ctx context.Context, in <-chan database.Marker, sampleRate float64) <-chan database.Marker {
	out := make(chan database.Marker, 100)

	go func() {
		defer close(out)
		counter := 0

		for m := range in {
			// Deterministic sampling based on marker ID for consistency
			// Use modulo to ensure even distribution
			counter++

			// Convert sample rate to skip pattern
			// e.g., 0.25 = keep every 4th marker
			skipInterval := int(1.0 / sampleRate)
			if skipInterval < 1 {
				skipInterval = 1
			}

			// Keep marker if counter is divisible by skip interval
			if counter%skipInterval == 0 {
				select {
				case out <- m:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out
}

// streamMultipleTracks merges marker streams from multiple tracks into a single channel
func streamMultipleTracks(ctx context.Context, trackIDs []string, zoom int, minLat, minLon, maxLat, maxLon float64, dbType string) (<-chan database.Marker, <-chan error) {
	out := make(chan database.Marker)
	errOut := make(chan error, len(trackIDs))

	var wg sync.WaitGroup

	for _, trackID := range trackIDs {
		trackID = strings.TrimSpace(trackID)
		if trackID == "" {
			continue
		}

		wg.Add(1)
		go func(tid string) {
			defer wg.Done()

			markerCh, errCh := db.StreamMarkersByTrackIDZoomAndBounds(ctx, tid, zoom, minLat, minLon, maxLat, maxLon, dbType)

			for {
				select {
				case <-ctx.Done():
					return
				case err, ok := <-errCh:
					if !ok {
						return
					}
					if err != nil {
						select {
						case errOut <- err:
						case <-ctx.Done():
							return
						}
					}
				case marker, ok := <-markerCh:
					if !ok {
						return
					}
					select {
					case out <- marker:
					case <-ctx.Done():
						return
					}
				}
			}
		}(trackID)
	}

	// Close output channel when all tracks are done
	go func() {
		wg.Wait()
		close(out)
		close(errOut)
	}()

	return out, errOut
}

// streamMarkersHandler streams markers via Server-Sent Events or MessagePack binary.
// Markers are emitted as soon as they are read and aggregated.
// Use ?format=msgpack for binary encoding (~60% smaller, faster parsing).
//
// @Summary     Stream markers for viewport
// @Description Streams map markers as SSE by default, or MessagePack when requested via query/header.
// @Tags        map
// @Produce     text/event-stream
// @Param       zoom     query int    false "Requested map zoom"
// @Param       minLat   query number true  "Minimum latitude"
// @Param       minLon   query number true  "Minimum longitude"
// @Param       maxLat   query number true  "Maximum latitude"
// @Param       maxLon   query number true  "Maximum longitude"
// @Param       trackID  query string false "Single track ID filter"
// @Param       trackIDs query string false "Comma-separated track IDs filter"
// @Param       format   query string false "Set to msgpack for binary stream"
// @Success     200 {string} string "Streaming response"
// @Failure     500 {string} string "Streaming unsupported"
// @Router      /stream_markers [get]
func streamMarkersHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	zoom, _ := strconv.Atoi(q.Get("zoom"))
	minLat, _ := strconv.ParseFloat(q.Get("minLat"), 64)
	minLon, _ := strconv.ParseFloat(q.Get("minLon"), 64)
	maxLat, _ := strconv.ParseFloat(q.Get("maxLat"), 64)
	maxLon, _ := strconv.ParseFloat(q.Get("maxLon"), 64)
	trackID := q.Get("trackID")
	trackIDsParam := q.Get("trackIDs")
	showFilter := q.Get("show") // "rt" = realtime only, skip historical data

	// Check if client requests MessagePack binary format
	useMsgpack := q.Get("format") == "msgpack" ||
		strings.Contains(r.Header.Get("Accept"), "application/msgpack")

	// PERFORMANCE: Calculate density-based sampling rate for low zoom levels
	// This reduces network traffic and client-side processing for dense views
	sampleRate := calculateSampleRate(zoom, minLat, minLon, maxLat, maxLon)
	if sampleRate < 1.0 {
		log.Printf("density reduction: zoom=%d rate=%.2f (sending ~%.0f%% of markers)",
			zoom, sampleRate, sampleRate*100)
	}

	// ON-THE-FLY CLUSTERING: Always query zoom=0 (raw markers)
	// The aggregateMarkers function will cluster based on requested zoom
	const rawZoom = 0

	// Choose streaming source: either entire map, a single track, or multiple tracks
	// When showFilter=="rt", skip historical data entirely for faster RT-only views.
	ctx := r.Context()
	var (
		baseSrc <-chan database.Marker
		errCh   <-chan error
	)
	rtOnly := showFilter == "rt"
	if rtOnly && trackID == "" && trackIDsParam == "" {
		// In RT-only mode for global view, still load historical markers
		// They will be filtered client-side, but spectrum markers should be available
		baseSrc, errCh = db.StreamMarkersByZoomAndBounds(ctx, rawZoom, minLat, minLon, maxLat, maxLon, *dbType)
	} else if trackID != "" {
		baseSrc, errCh = db.StreamMarkersByTrackIDZoomAndBounds(ctx, trackID, rawZoom, minLat, minLon, maxLat, maxLon, *dbType)
	} else if trackIDsParam != "" {
		// Handle multiple tracks by merging their streams
		trackIDs := strings.Split(trackIDsParam, ",")
		baseSrc, errCh = streamMultipleTracks(ctx, trackIDs, rawZoom, minLat, minLon, maxLat, maxLon, *dbType)
	} else {
		baseSrc, errCh = db.StreamMarkersByZoomAndBounds(ctx, rawZoom, minLat, minLon, maxLat, maxLon, *dbType)
	}

	// Fetch current realtime points once so the map reflects network devices.
	// We only touch the realtime table when the dedicated flag enables it so
	// operators control the feature explicitly.
	var rtMarks []database.Marker
	if *safecastRealtimeEnabled {
		var err error
		rtMarks, err = db.GetLatestRealtimeByBounds(ctx, minLat, minLon, maxLat, maxLon, *dbType)
		if err != nil {
			log.Printf("realtime query: %v", err)
		}
		// Log bounds alongside count to help diagnose empty map tiles.
		log.Printf("realtime markers: %d lat[%f,%f] lon[%f,%f]", len(rtMarks), minLat, maxLat, minLon, maxLon)
	}

	// Apply density-based sampling if needed (zoom < 8)
	var sampledSrc <-chan database.Marker
	if sampleRate < 1.0 {
		sampledSrc = sampleMarkerChannel(ctx, baseSrc, sampleRate)
	} else {
		sampledSrc = baseSrc
	}

	agg := aggregateMarkers(ctx, sampledSrc, nil, zoom)

	// Set appropriate headers based on format
	if useMsgpack {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("X-Content-Format", "msgpack")
	} else {
		w.Header().Set("Content-Type", "text/event-stream")
	}
	w.Header().Set("Cache-Control", "no-cache")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Helper to write a msgpack frame: 4-byte length prefix + data
	writeMsgpackFrame := func(data []byte) {
		lenBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(lenBuf, uint32(len(data)))
		w.Write(lenBuf)
		w.Write(data)
	}

	// Emit realtime markers first when enabled.
	for _, m := range rtMarks {
		if useMsgpack {
			b, _ := msgpack.Marshal(m)
			writeMsgpackFrame(b)
		} else {
			b, _ := json.Marshal(m)
			fmt.Fprintf(w, "data: %s\n\n", b)
		}
	}
	flusher.Flush()

	for {
		select {
		case <-ctx.Done():
			return
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
				continue
			}
			if err != nil {
				if useMsgpack {
					// Send zero-length frame as error indicator
					w.Write([]byte{0, 0, 0, 0})
				} else {
					fmt.Fprintf(w, "event: done\ndata: %v\n\n", err)
				}
				flusher.Flush()
				return
			}
			errCh = nil
		case m, ok := <-agg:
			if !ok {
				if useMsgpack {
					// Send zero-length frame as end marker
					w.Write([]byte{0, 0, 0, 0})
				} else {
					fmt.Fprint(w, "event: done\ndata: end\n\n")
				}
				flusher.Flush()
				return
			}
			if useMsgpack {
				b, _ := msgpack.Marshal(m)
				writeMsgpackFrame(b)
			} else {
				b, _ := json.Marshal(m)
				fmt.Fprintf(w, "data: %s\n\n", b)
			}
			flusher.Flush()
		}
	}
}

// realtimePoint holds a single measurement for realtime history charts.
// Keeping the struct tiny helps when we duplicate slices for aggregation.
type realtimePoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// rangeSummary describes the plotted window and aggregation bucket.
// Returning it to the frontend lets JavaScript render friendly titles.
type rangeSummary struct {
	Start         int64 `json:"start"`
	End           int64 `json:"end"`
	BucketSeconds int64 `json:"bucketSeconds"`
}

// historyAggregate bundles the processed realtime readings so the handler can
// serialise the JSON response without juggling several parallel slices.
type historyAggregate struct {
	Series      map[string][]realtimePoint
	ExtraSeries map[string]map[string][]realtimePoint
	Extra       map[string]float64
	Ranges      map[string]rangeSummary
	DeviceName  string
	Transport   string
	Tube        string
	Country     string
}

// realtimeMeasurementPayload moves measurements between goroutines without
// sharing mutable state and keeps channel signatures consistent across helpers.
type realtimeMeasurementPayload struct {
	timestamp int64
	radiation float64
	extras    map[string]float64
	name      string
	transport string
	tube      string
	country   string
}

