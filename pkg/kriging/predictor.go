package kriging

import (
	"context"
	"errors"
	"math"
)

// ErrViewportTooWide means the rectangle is larger than one local neighborhood
// can honestly cover. Callers must not thin the points to make it fit.
var ErrViewportTooWide = errors.New("kriging: viewport too wide")

// ErrTooManyObservations means the raw window is larger than the solver cap.
// The extra points are not dropped.
var ErrTooManyObservations = errors.New("kriging: too many observations")

// ErrInvalidBounds means the rectangle has no area.
var ErrInvalidBounds = errors.New("kriging: invalid bounds")

const (
	defaultGrid = 32
	maxGrid     = 48
)

// Cell is one prediction. MedianUSvH is nil when the system did not solve.
// Suppressed cells must not be drawn, even when a median is present: ordinary
// kriging still returns a number where the observations are too far to trust.
type Cell struct {
	Lat         float64  `json:"lat"`
	Lon         float64  `json:"lon"`
	MedianUSvH  *float64 `json:"medianUSvH"`
	LogVariance float64  `json:"logVariance"`
	Suppressed  bool     `json:"suppressed"`
}

// GridRequest is one viewport. Cols and Rows are the prediction lattice, not map tiles.
type GridRequest struct {
	Bounds   Bounds
	TimeUnix int64
	Cols     int
	Rows     int
}

// Grid is a lattice of predictions that share one candidate window.
type Grid struct {
	Bounds             Bounds
	Cols               int
	Rows               int
	PredictionTime     int64
	Cells              []Cell
	CandidateCount     int
	OmittedNonPositive int
}

// Predictor coordinates retrieval, the neighborhood, and the solver.
// Map code should depend on Predictor rather than on the covariance or the database.
type Predictor struct {
	Source          ObservationSource
	Model           CovarianceModel
	Neighbors       int
	MaxObservations int
}

// NewPredictor wires a source to the given covariance. A zero model is replaced
// with DefaultModel. The neighborhood uses that same model, so the points it
// keeps are the ones the covariance treats as correlated.
func NewPredictor(src ObservationSource, model CovarianceModel) *Predictor {
	if model == nil {
		model = DefaultModel()
	}
	return &Predictor{
		Source:          src,
		Model:           model,
		Neighbors:       DefaultNeighbors,
		MaxObservations: DefaultMaxObservations,
	}
}

// PredictGrid predicts every cell center in the viewport.
// The observation query is the viewport expanded by the spatial practical range,
// and the time window is one temporal practical range on either side of the
// prediction time. Points outside that window have correlation below about 0.05
// and are not fetched. That cutoff is not a boundary between regional models.
func (p *Predictor) PredictGrid(ctx context.Context, req GridRequest) (Grid, error) {
	if p == nil || p.Source == nil || p.Model == nil {
		return Grid{}, errors.New("kriging: predictor is not configured")
	}
	if !req.Bounds.Valid() {
		return Grid{}, ErrInvalidBounds
	}
	if SpanMeters(req.Bounds) > MaxViewportMeters {
		return Grid{}, ErrViewportTooWide
	}
	cols, rows := normalizeGrid(req.Cols, req.Rows)
	haloM := SpatialPracticalRangeM
	window := TemporalPracticalRangeSec
	if exp, ok := p.Model.(ExponentialModel); ok {
		haloM = exp.SpatialPracticalRangeM
		window = exp.TemporalPracticalRangeSec
	}
	q, err := p.Source.Query(ctx, Expand(req.Bounds, haloM), req.TimeUnix-window, req.TimeUnix+window)
	if err != nil {
		return Grid{}, err
	}
	limit := p.MaxObservations
	if limit <= 0 {
		limit = DefaultMaxObservations
	}
	if len(q.Observations) > limit {
		return Grid{}, ErrTooManyObservations
	}

	lats := cellCoords(req.Bounds.MinLat, req.Bounds.MaxLat, rows)
	lons := cellLons(req.Bounds, cols)
	cells := make([]Cell, 0, rows*cols)
	for _, lat := range lats {
		for _, lon := range lons {
			cells = append(cells, p.predict(Target{Lat: lat, Lon: lon, TimeUnix: req.TimeUnix}, q.Observations))
		}
	}
	return Grid{
		Bounds:             req.Bounds,
		Cols:               cols,
		Rows:               rows,
		PredictionTime:     req.TimeUnix,
		Cells:              cells,
		CandidateCount:     len(q.Observations),
		OmittedNonPositive: q.OmittedNonPositive,
	}, nil
}

func (p *Predictor) predict(target Target, candidates []Observation) Cell {
	cell := Cell{
		Lat:         target.Lat,
		Lon:         target.Lon,
		LogVariance: uninformedVariance(p.Model),
		Suppressed:  true,
	}
	k := p.Neighbors
	if k <= 0 {
		k = DefaultNeighbors
	}
	neigh := (NeighborhoodSelector{K: k, Model: p.Model}).Select(target, candidates)
	if len(neigh) == 0 {
		return cell
	}
	sol, ok := Solve(p.Model, target, neigh)
	if !ok {
		return cell
	}
	med := math.Exp(sol.MeanLog)
	if math.IsNaN(med) || math.IsInf(med, 0) || med <= 0 {
		return cell
	}
	cell.LogVariance = sol.Variance
	cell.MedianUSvH = &med
	cell.Suppressed = Suppress(p.Model.FieldVariance(), sol.Variance)
	return cell
}

func uninformedVariance(model CovarianceModel) float64 {
	if model == nil || model.FieldVariance() <= 0 {
		return 1
	}
	// Far from data, ordinary kriging variance is about two partial sills
	// plus the nugget. A finite stand-in keeps JSON finite and the cell suppressed.
	return 2*model.FieldVariance() + model.MeasurementNugget()
}

func normalizeGrid(cols, rows int) (int, int) {
	if cols <= 0 {
		cols = defaultGrid
	}
	if rows <= 0 {
		rows = defaultGrid
	}
	if cols > maxGrid {
		cols = maxGrid
	}
	if rows > maxGrid {
		rows = maxGrid
	}
	return cols, rows
}

func cellCoords(min, max float64, n int) []float64 {
	step := (max - min) / float64(n)
	out := make([]float64, n)
	for i := 0; i < n; i++ {
		out[i] = min + (float64(i)+0.5)*step
	}
	return out
}

func cellLons(b Bounds, cols int) []float64 {
	if b.MinLon <= b.MaxLon {
		return cellCoords(b.MinLon, b.MaxLon, cols)
	}
	span := (180 - b.MinLon) + (b.MaxLon + 180)
	step := span / float64(cols)
	out := make([]float64, cols)
	for i := 0; i < cols; i++ {
		lon := b.MinLon + (float64(i)+0.5)*step
		for lon > 180 {
			lon -= 360
		}
		out[i] = lon
	}
	return out
}
