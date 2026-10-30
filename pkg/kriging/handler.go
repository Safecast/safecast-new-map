package kriging

import (
	_ "embed"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"safecast-new-map/pkg/database"
)

//go:embed overlay.js
var overlayJS []byte

// ResponseNote tells API clients why a cell can be blank.
const ResponseNote = "Suppressed cells have too little independent support to draw. A blank area is not a measurement of safety."

// Handler serves the prediction grid and the map overlay script.
type Handler struct {
	real      *Predictor
	synthetic *Predictor
}

// Register attaches GET /api/kriging/grid and GET /kriging/overlay.js.
// db may be nil; the synthetic source still answers.
func Register(mux *http.ServeMux, db *database.Database) {
	if mux == nil {
		return
	}
	h := &Handler{
		real:      NewPredictor(NewDBSource(db), DefaultModel()),
		synthetic: NewPredictor(SyntheticSource{}, DefaultModel()),
	}
	mux.HandleFunc("/api/kriging/grid", h.ServeGrid)
	mux.HandleFunc("/kriging/overlay.js", h.ServeOverlay)
}

// ServeOverlay serves the Leaflet control. It is embedded so the algorithm and
// the layer that draws it stay in this package.
func (h *Handler) ServeOverlay(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	_, _ = w.Write(overlayJS)
}

// ServeGrid predicts a local log-dose field for one viewport.
//
// @Summary     Local radiation kriging grid
// @Description Ordinary kriging of log dose rate for a small viewport. Returns a median in µSv/h and the latent log-variance. Suppressed cells must not be drawn. Viewports wider than 20 km are refused instead of subsampled. synthetic=1 uses a fixed fake track and does not read the database.
// @Tags        map
// @Produce     json
// @Param       minLat    query number true  "Southern boundary"
// @Param       minLon    query number true  "Western boundary"
// @Param       maxLat    query number true  "Northern boundary"
// @Param       maxLon    query number true  "Eastern boundary"
// @Param       time      query integer false "Prediction time, unix seconds. Defaults to now."
// @Param       cols      query integer false "Grid columns, 1 to 48"
// @Param       rows      query integer false "Grid rows, 1 to 48"
// @Param       synthetic query boolean false "Use the built-in synthetic track"
// @Success     200 {object} map[string]interface{} "Prediction grid"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     422 {object} map[string]string "Viewport too wide or too many raw observations"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /api/kriging/grid [get]
func (h *Handler) ServeGrid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	q := r.URL.Query()
	bounds, err := parseBounds(q)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	when, err := parseTime(q.Get("time"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	cols, err := parseGridCount(q.Get("cols"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid cols"})
		return
	}
	rows, err := parseGridCount(q.Get("rows"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid rows"})
		return
	}
	synthetic := q.Get("synthetic") == "1" || q.Get("synthetic") == "true"
	predictor := h.real
	if synthetic {
		predictor = h.synthetic
	}
	if predictor == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "kriging is not configured"})
		return
	}

	grid, err := predictor.PredictGrid(r.Context(), GridRequest{
		Bounds:   bounds,
		TimeUnix: when,
		Cols:     cols,
		Rows:     rows,
	})
	if err != nil {
		h.writePredictError(w, err, synthetic)
		return
	}
	writeJSON(w, http.StatusOK, gridPayload(grid, predictor.Model, synthetic))
}

func (h *Handler) writePredictError(w http.ResponseWriter, err error, synthetic bool) {
	body := map[string]any{"error": err.Error()}
	if synthetic {
		view := SyntheticView()
		body["syntheticView"] = view
	}
	switch {
	case errors.Is(err, ErrInvalidBounds):
		writeJSON(w, http.StatusBadRequest, body)
	case errors.Is(err, ErrViewportTooWide):
		body["refusal"] = "viewport_too_wide"
		writeJSON(w, http.StatusUnprocessableEntity, body)
	case errors.Is(err, ErrTooManyObservations):
		body["refusal"] = "too_many_observations"
		writeJSON(w, http.StatusUnprocessableEntity, body)
	default:
		if errors.Is(err, ErrDatabaseUnavailable) {
			writeJSON(w, http.StatusServiceUnavailable, body)
			return
		}
		writeJSON(w, http.StatusInternalServerError, body)
	}
}

func gridPayload(grid Grid, model CovarianceModel, synthetic bool) map[string]any {
	name, spatial, temporal, sill, nugget := modelMeta(model)
	payload := map[string]any{
		"predictionTime":            grid.PredictionTime,
		"valueUnit":                 "uSv/h",
		"varianceScale":             "log",
		"model":                     name,
		"spatialPracticalRangeM":    spatial,
		"temporalPracticalRangeSec": temporal,
		"partialSill":               sill,
		"nugget":                    nugget,
		"suppressionThreshold":      SuppressionFraction * sill,
		"candidateCount":            grid.CandidateCount,
		"omittedNonPositive":        grid.OmittedNonPositive,
		"note":                      ResponseNote,
		"minLat":                    grid.Bounds.MinLat,
		"minLon":                    grid.Bounds.MinLon,
		"maxLat":                    grid.Bounds.MaxLat,
		"maxLon":                    grid.Bounds.MaxLon,
		"cols":                      grid.Cols,
		"rows":                      grid.Rows,
		"cells":                     grid.Cells,
	}
	if synthetic {
		payload["syntheticView"] = SyntheticView()
	}
	return payload
}

func modelMeta(model CovarianceModel) (name string, spatial float64, temporal int64, sill, nugget float64) {
	sill = model.FieldVariance()
	nugget = model.MeasurementNugget()
	if exp, ok := model.(ExponentialModel); ok {
		return ModelName, exp.SpatialPracticalRangeM, exp.TemporalPracticalRangeSec, sill, nugget
	}
	return "covariance", 0, 0, sill, nugget
}

func parseBounds(q interface{ Get(string) string }) (Bounds, error) {
	minLat, err := parseCoord(q, "minLat")
	if err != nil {
		return Bounds{}, err
	}
	minLon, err := parseCoord(q, "minLon")
	if err != nil {
		return Bounds{}, err
	}
	maxLat, err := parseCoord(q, "maxLat")
	if err != nil {
		return Bounds{}, err
	}
	maxLon, err := parseCoord(q, "maxLon")
	if err != nil {
		return Bounds{}, err
	}
	if minLat < -90 || maxLat > 90 || minLat > 90 || maxLat < -90 {
		return Bounds{}, errors.New("latitude out of range")
	}
	b := Bounds{MinLat: minLat, MinLon: minLon, MaxLat: maxLat, MaxLon: maxLon}
	if !b.Valid() {
		return Bounds{}, ErrInvalidBounds
	}
	return b, nil
}

func parseCoord(q interface{ Get(string) string }, key string) (float64, error) {
	s := q.Get(key)
	if s == "" {
		return 0, errors.New("missing " + key)
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) {
		return 0, errors.New("invalid " + key)
	}
	return v, nil
}

func parseTime(s string) (int64, error) {
	if s == "" {
		return nowUnix(), nil
	}
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.New("invalid time")
	}
	return v, nil
}

func parseGridCount(s string) (int, error) {
	if s == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return 0, errors.New("invalid grid count")
	}
	return v, nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// nowUnix is a variable so tests can pin the default prediction time.
var nowUnix = func() int64 { return time.Now().Unix() }
