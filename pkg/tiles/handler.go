package tiles

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"safecast-new-map/pkg/database"
)

// active is the service uploads invalidate. The map process sets it once.
var (
	activeMu sync.RWMutex
	active   *Service
)

// Use selects the service InvalidateBounds talks to.
func Use(s *Service) {
	activeMu.Lock()
	active = s
	activeMu.Unlock()
}

// InvalidateBounds drops overview tiles that meet b on the active service.
// A process that has not called Use does nothing, so tests and tools stay quiet.
func InvalidateBounds(b Bounds) error {
	activeMu.RLock()
	s := active
	activeMu.RUnlock()
	if s == nil {
		return nil
	}
	return s.InvalidateBounds(b)
}

// Register mounts GET /tiles/... on mux.
func (s *Service) Register(mux *http.ServeMux) {
	if mux == nil || s == nil {
		return
	}
	mux.HandleFunc("GET /tiles/{path...}", s.serve)
}

// serve dispatches tile images, the click probe, and the live-sensor overlay.
//
// @Summary     Radiation map tile
// @Description Serves an aggregated PNG for one XYZ tile, a probe of the winning measurement, or the current realtime markers in a viewport.
// @Tags        map
// @Produce     image/png
// @Param       path path string true "aggregate/{z}/{x}/{y}.png, aggregate/probe, or live"
// @Param       from query int false "Start unix time, quantized to the UTC day"
// @Param       to query int false "End unix time, quantized to the UTC day"
// @Param       speed query string false "Speed bands: a plane, c car, w walk, combined as cw"
// @Success     200 {string} string "PNG tile, probe JSON, or realtime JSON"
// @Failure     400 {string} string "Bad coordinates or time window"
// @Failure     404 {string} string "Unknown processor"
// @Failure     500 {string} string "Render failed"
// @Router      /tiles/{path} [get]
func (s *Service) serve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/tiles/")
	parts := strings.Split(path, "/")
	if len(parts) == 1 && parts[0] == "live" {
		s.serveLive(w, r)
		return
	}
	if len(parts) == 2 && parts[1] == "probe" {
		s.serveProbe(w, r, parts[0])
		return
	}
	if len(parts) == 4 && strings.HasSuffix(parts[3], ".png") {
		yPart := strings.TrimSuffix(parts[3], ".png")
		s.serveTile(w, r, parts[0], parts[1], parts[2], yPart)
		return
	}
	http.NotFound(w, r)
}

func (s *Service) serveTile(w http.ResponseWriter, r *http.Request, name, zs, xs, ys string) {
	id, err := parseTile(zs, xs, ys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	from, to, speed, err := filters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := s.Render(r.Context(), name, Request{Tile: id, From: from, To: to, Speed: speed})
	if err != nil {
		writeRenderError(w, err)
		return
	}
	w.Header().Set("Content-Type", res.ContentType)
	w.Header().Set("Cache-Control", "public, max-age=60")
	_, _ = w.Write(res.Body)
}

func (s *Service) serveProbe(w http.ResponseWriter, r *http.Request, name string) {
	if name != "aggregate" || s.agg == nil {
		http.NotFound(w, r)
		return
	}
	q := r.URL.Query()
	lat, err1 := strconv.ParseFloat(q.Get("lat"), 64)
	lon, err2 := strconv.ParseFloat(q.Get("lon"), 64)
	z, err3 := strconv.Atoi(q.Get("z"))
	if err1 != nil || err2 != nil || err3 != nil {
		http.Error(w, "lat, lon, and z are required", http.StatusBadRequest)
		return
	}
	from, to, speed, err := filters(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hit, err := s.agg.Probe(r.Context(), lat, lon, z, QuantizeFrom(from), QuantizeTo(to), speed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(hit)
}

func (s *Service) serveLive(w http.ResponseWriter, r *http.Request) {
	if s.agg == nil || s.agg.DB == nil {
		http.Error(w, "database not available", http.StatusServiceUnavailable)
		return
	}
	q := r.URL.Query()
	minLat, err1 := strconv.ParseFloat(q.Get("minLat"), 64)
	minLon, err2 := strconv.ParseFloat(q.Get("minLon"), 64)
	maxLat, err3 := strconv.ParseFloat(q.Get("maxLat"), 64)
	maxLon, err4 := strconv.ParseFloat(q.Get("maxLon"), 64)
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		http.Error(w, "minLat, minLon, maxLat, and maxLon are required", http.StatusBadRequest)
		return
	}
	markers, err := s.agg.DB.GetLatestRealtimeByBounds(r.Context(), minLat, minLon, maxLat, maxLon, s.agg.DBType)
	if err != nil {
		http.Error(w, "realtime query failed", http.StatusInternalServerError)
		return
	}
	if markers == nil {
		markers = []database.Marker{}
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	_ = json.NewEncoder(w).Encode(markers)
}

func parseTile(zs, xs, ys string) (TileID, error) {
	z, err1 := strconv.Atoi(zs)
	x, err2 := strconv.Atoi(xs)
	y, err3 := strconv.Atoi(ys)
	if err1 != nil || err2 != nil || err3 != nil || z < 0 || z > MaxOverviewZoom {
		return TileID{}, errBadTile
	}
	if _, ok := TileBounds(TileID{Z: z, X: x, Y: y}); !ok {
		return TileID{}, errBadTile
	}
	return TileID{Z: z, X: x, Y: y}, nil
}

var errBadTile = errors.New("invalid tile coordinate")

func filters(r *http.Request) (from, to int64, speed string, err error) {
	q := r.URL.Query()
	if s := q.Get("from"); s != "" {
		from, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, 0, "", errors.New("invalid from")
		}
	}
	if s := q.Get("to"); s != "" {
		to, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0, 0, "", errors.New("invalid to")
		}
	}
	if from > 0 && to > 0 && QuantizeFrom(from) > QuantizeTo(to) {
		return 0, 0, "", errors.New("from is after to")
	}
	return from, to, q.Get("speed"), nil
}

func writeRenderError(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrUnknownProcessor) {
		http.Error(w, "unknown processor", http.StatusNotFound)
		return
	}
	http.Error(w, "tile render failed", http.StatusInternalServerError)
}
