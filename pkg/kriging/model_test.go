package kriging

import (
	"context"
	"math"
	"testing"
)

func atMeters(east, north, dose float64, when int64, track string) Observation {
	m := metersPerDegreeLat()
	return Observation{
		Lat:      north / m,
		Lon:      east / m,
		TimeUnix: when,
		DoseRate: dose,
		TrackID:  track,
		DeviceID: "meter",
	}
}

func targetMeters(east, north float64, when int64) Target {
	m := metersPerDegreeLat()
	return Target{Lat: north / m, Lon: east / m, TimeUnix: when}
}

func mustSolve(t *testing.T, target Target, obs []Observation) Solution {
	t.Helper()
	sol, ok := Solve(DefaultModel(), target, obs)
	if !ok {
		t.Fatal("solve failed")
	}
	return sol
}

func TestConstantFieldAndWeights(t *testing.T) {
	const when int64 = 1_700_000_000
	obs := []Observation{
		atMeters(0, 0, 0.2, when, "a"),
		atMeters(80, 20, 0.2, when, "b"),
		atMeters(-40, 60, 0.2, when, "c"),
	}
	sol := mustSolve(t, targetMeters(20, 15, when), obs)
	sum := 0.0
	for _, w := range sol.Weights {
		sum += w
	}
	if math.Abs(sum-1) > 1e-8 {
		t.Fatalf("weights sum to %g, want 1", sum)
	}
	if math.Abs(sol.MeanLog-math.Log(0.2)) > 1e-8 {
		t.Fatalf("log mean %g, want %g", sol.MeanLog, math.Log(0.2))
	}
	if math.Exp(sol.MeanLog) <= 0 {
		t.Fatal("median is not positive")
	}
}

func TestLatentVarianceAtObservationEqualsNugget(t *testing.T) {
	const when int64 = 1_700_000_000
	obs := atMeters(10, -5, 0.15, when, "here")
	sol := mustSolve(t, targetMeters(10, -5, when), []Observation{obs})
	nugget := DefaultModel().Nugget
	if math.Abs(sol.Variance-nugget) > 1e-8 {
		t.Fatalf("variance %g, want nugget %g", sol.Variance, nugget)
	}
	if math.Abs(math.Exp(sol.MeanLog)-0.15) > 1e-8 {
		t.Fatalf("median %g, want 0.15", math.Exp(sol.MeanLog))
	}
	if Suppress(DefaultPartialSill, sol.Variance) {
		t.Fatal("a measurement at the prediction point should be drawn")
	}
}

func TestVarianceGrowsWithDistanceAndIsSuppressed(t *testing.T) {
	const when int64 = 1_700_000_000
	obs := []Observation{atMeters(0, 0, 0.12, when, "origin")}
	near := mustSolve(t, targetMeters(0, 0, when), obs)
	mid := mustSolve(t, targetMeters(200, 0, when), obs)
	far := mustSolve(t, targetMeters(3000, 0, when), obs)
	if !(near.Variance < mid.Variance && mid.Variance < far.Variance) {
		t.Fatalf("variance did not grow: near %g mid %g far %g", near.Variance, mid.Variance, far.Variance)
	}
	if Suppress(DefaultPartialSill, near.Variance) {
		t.Fatal("prediction on the observation was suppressed")
	}
	if !Suppress(DefaultPartialSill, far.Variance) {
		t.Fatalf("far variance %g was not suppressed", far.Variance)
	}
	cell := NewPredictor(memorySource{obs: obs}, DefaultModel()).predict(targetMeters(3000, 0, when), obs)
	if !cell.Suppressed {
		t.Fatal("far cell would be drawn")
	}
	if cell.MedianUSvH != nil && *cell.MedianUSvH <= 0 {
		t.Fatal("median must stay positive when it is present")
	}
}

func TestClusterDoesNotActLikeIndependentSamples(t *testing.T) {
	const when int64 = 1_700_000_000
	cluster := make([]Observation, 40)
	for i := range cluster {
		// A few tens of centimeters and a few seconds: one passage, not 40 places.
		cluster[i] = atMeters(200+float64(i%5)*0.25, float64(i/5)*0.25, 0.3, when+int64(i), "drive")
	}
	spread := make([]Observation, 0, 40)
	for i := 0; i < 8; i++ {
		for j := 0; j < 5; j++ {
			spread = append(spread, atMeters(-180+float64(i)*50, -100+float64(j)*50, 0.3, when, "spread"))
		}
	}
	single := []Observation{atMeters(200, 0, 0.3, when, "one")}
	target := targetMeters(0, 0, when)

	vSingle := mustSolve(t, target, single).Variance
	vCluster := mustSolve(t, target, (NeighborhoodSelector{K: DefaultNeighbors, Model: DefaultModel()}).Select(target, cluster)).Variance
	vSpread := mustSolve(t, target, (NeighborhoodSelector{K: DefaultNeighbors, Model: DefaultModel()}).Select(target, spread)).Variance

	if vCluster < vSingle*0.5 {
		t.Fatalf("cluster variance %g collapsed relative to one observation %g", vCluster, vSingle)
	}
	if vSpread >= vCluster*0.5 {
		t.Fatalf("spread variance %g is not well below cluster variance %g (single %g)", vSpread, vCluster, vSingle)
	}

	sol := mustSolve(t, target, (NeighborhoodSelector{K: DefaultNeighbors, Model: DefaultModel()}).Select(target, cluster))
	sum := 0.0
	maxW := 0.0
	for _, w := range sol.Weights {
		sum += w
		if w > maxW {
			maxW = w
		}
	}
	if math.Abs(sum-1) > 1e-6 {
		t.Fatalf("cluster weights sum to %g", sum)
	}
	if maxW >= 0.1 {
		t.Fatalf("cluster weight %g was not split across the track", maxW)
	}
}

func TestOverlappingViewportsAgree(t *testing.T) {
	const when int64 = 1_700_000_000
	obs := []Observation{
		{Lat: 0.003, Lon: 0.003, TimeUnix: when, DoseRate: 0.2, TrackID: "inside", DeviceID: "d"},
		// Outside the western viewport, inside the eastern one, and inside the western halo.
		{Lat: 0.0055, Lon: 0.003, TimeUnix: when, DoseRate: 2, TrackID: "halo", DeviceID: "d"},
	}
	src := memorySource{obs: obs}
	p := NewPredictor(src, DefaultModel())
	west := Bounds{MinLat: 0, MinLon: 0, MaxLat: 0.004, MaxLon: 0.004}
	east := Bounds{MinLat: 0.002, MinLon: 0.002, MaxLat: 0.006, MaxLon: 0.006}
	g1, err := p.PredictGrid(context.Background(), GridRequest{Bounds: west, TimeUnix: when, Cols: 2, Rows: 2})
	if err != nil {
		t.Fatal(err)
	}
	g2, err := p.PredictGrid(context.Background(), GridRequest{Bounds: east, TimeUnix: when, Cols: 2, Rows: 2})
	if err != nil {
		t.Fatal(err)
	}
	shared := 0
	for _, a := range g1.Cells {
		for _, b := range g2.Cells {
			if math.Abs(a.Lat-b.Lat) > 1e-9 || math.Abs(a.Lon-b.Lon) > 1e-9 {
				continue
			}
			shared++
			if a.Suppressed != b.Suppressed {
				t.Fatalf("suppression differs at %g,%g", a.Lat, a.Lon)
			}
			if a.MedianUSvH == nil || b.MedianUSvH == nil {
				t.Fatalf("missing median at shared cell %g,%g", a.Lat, a.Lon)
			}
			if math.Abs(*a.MedianUSvH-*b.MedianUSvH) > 1e-9 {
				t.Fatalf("median %g vs %g at %g,%g", *a.MedianUSvH, *b.MedianUSvH, a.Lat, a.Lon)
			}
			if math.Abs(a.LogVariance-b.LogVariance) > 1e-9 {
				t.Fatalf("variance %g vs %g", a.LogVariance, b.LogVariance)
			}
			if a.Suppressed {
				t.Fatal("shared cell on top of a measurement was suppressed")
			}
		}
	}
	if shared == 0 {
		t.Fatal("viewports did not share a cell center")
	}
}

func TestTemporalDecayAndHistoricalPrediction(t *testing.T) {
	const past int64 = 1_300_000_000
	contemporary := Observation{Lat: 0, Lon: 0, TimeUnix: past, DoseRate: 0.2, TrackID: "then", DeviceID: "d"}
	future := Observation{Lat: 0, Lon: 0, TimeUnix: past + 60*24*3600, DoseRate: 8, TrackID: "later", DeviceID: "d"}
	old := Observation{Lat: 0, Lon: 0, TimeUnix: past - 60*24*3600, DoseRate: 8, TrackID: "older", DeviceID: "d"}

	sol, ok := Solve(DefaultModel(), Target{TimeUnix: past}, []Observation{contemporary, future})
	if !ok {
		t.Fatal("solve failed")
	}
	med := math.Exp(sol.MeanLog)
	if med < 0.15 || med > 0.35 {
		t.Fatalf("future point was treated as contemporary, median %g", med)
	}

	p := NewPredictor(memorySource{obs: []Observation{contemporary, future, old}}, DefaultModel())
	grid, err := p.PredictGrid(context.Background(), GridRequest{
		Bounds:   Bounds{MinLat: -0.001, MinLon: -0.001, MaxLat: 0.001, MaxLon: 0.001},
		TimeUnix: past,
		Cols:     2,
		Rows:     2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if grid.CandidateCount != 1 {
		t.Fatalf("window kept %d observations, want the contemporary one", grid.CandidateCount)
	}
	for _, cell := range grid.Cells {
		if cell.MedianUSvH == nil || math.Abs(*cell.MedianUSvH-0.2) > 1e-6 {
			t.Fatalf("historical median %#v", cell.MedianUSvH)
		}
	}

	onlyOld, err := NewPredictor(memorySource{obs: []Observation{old}}, DefaultModel()).PredictGrid(context.Background(), GridRequest{
		Bounds:   Bounds{MinLat: -0.001, MinLon: -0.001, MaxLat: 0.001, MaxLon: 0.001},
		TimeUnix: past,
		Cols:     1,
		Rows:     1,
	})
	if err != nil {
		t.Fatal(err)
	}
	if onlyOld.CandidateCount != 0 {
		t.Fatalf("old observation was fetched, count %d", onlyOld.CandidateCount)
	}
	if !onlyOld.Cells[0].Suppressed {
		t.Fatal("empty historical window was drawn")
	}
}

func TestDuplicateNeighborhoodIsStable(t *testing.T) {
	const when int64 = 1_700_000_000
	dupes := make([]Observation, 8)
	for i := range dupes {
		dupes[i] = atMeters(0, 0, 0.4, when, "same")
	}
	sol := mustSolve(t, targetMeters(0, 0, when), dupes)
	if math.Abs(math.Exp(sol.MeanLog)-0.4) > 1e-6 {
		t.Fatalf("median %g", math.Exp(sol.MeanLog))
	}
	if sol.Variance < 0 || math.IsNaN(sol.Variance) {
		t.Fatalf("variance %g", sol.Variance)
	}
	// Repeated measurements of one location reduce measurement error, not as n independent places.
	want := DefaultModel().Nugget / float64(len(dupes))
	if math.Abs(sol.Variance-want) > 1e-6 {
		t.Fatalf("variance %g, want %g", sol.Variance, want)
	}

	singular := DefaultModel()
	singular.Nugget = 0
	if _, ok := Solve(singular, targetMeters(0, 0, when), dupes[:2]); ok {
		t.Fatal("zero nugget with identical points should refuse the singular system")
	}
}

func TestGridRefusesWideViewportAndHugeSample(t *testing.T) {
	p := NewPredictor(memorySource{}, DefaultModel())
	_, err := p.PredictGrid(context.Background(), GridRequest{
		Bounds:   Bounds{MinLat: 0, MinLon: 0, MaxLat: 1, MaxLon: 1},
		TimeUnix: 1,
		Cols:     2,
		Rows:     2,
	})
	if err != ErrViewportTooWide {
		t.Fatalf("wide viewport err %v", err)
	}
	_, err = p.PredictGrid(context.Background(), GridRequest{
		Bounds: Bounds{MinLat: 1, MinLon: 0, MaxLat: 1, MaxLon: 1},
	})
	if err != ErrInvalidBounds {
		t.Fatalf("empty viewport err %v", err)
	}

	obs := make([]Observation, DefaultMaxObservations+1)
	for i := range obs {
		obs[i] = atMeters(float64(i%10), float64(i/10), 0.1, 10, "t")
	}
	_, err = NewPredictor(memorySource{obs: obs}, DefaultModel()).PredictGrid(context.Background(), GridRequest{
		Bounds:   Bounds{MinLat: -0.01, MinLon: -0.01, MaxLat: 0.01, MaxLon: 0.01},
		TimeUnix: 10,
		Cols:     2,
		Rows:     2,
	})
	if err != ErrTooManyObservations {
		t.Fatalf("too many err %v", err)
	}
}

func TestHaloIsQueried(t *testing.T) {
	const when int64 = 1_700_000_000
	rec := &recordingSource{obs: []Observation{atMeters(0, 0, 0.1, when, "t")}}
	bounds := Bounds{MinLat: 0, MinLon: 0, MaxLat: 0.002, MaxLon: 0.002}
	_, err := NewPredictor(rec, DefaultModel()).PredictGrid(context.Background(), GridRequest{
		Bounds: bounds, TimeUnix: when, Cols: 2, Rows: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if rec.bounds.MinLat >= bounds.MinLat || rec.bounds.MaxLat <= bounds.MaxLat {
		t.Fatalf("query was not expanded: %+v", rec.bounds)
	}
	if rec.t1-rec.t0 != 2*TemporalPracticalRangeSec {
		t.Fatalf("time window %d", rec.t1-rec.t0)
	}
}

type memorySource struct {
	obs     []Observation
	omitted int
}

func (m memorySource) Query(_ context.Context, bounds Bounds, t0, t1 int64) (QueryResult, error) {
	var kept []Observation
	omitted := m.omitted
	for _, o := range m.obs {
		if !contains(bounds, o.Lat, o.Lon) || o.TimeUnix < t0 || o.TimeUnix > t1 {
			continue
		}
		if o.DoseRate <= 0 {
			omitted++
			continue
		}
		kept = append(kept, o)
	}
	return QueryResult{Observations: kept, OmittedNonPositive: omitted}, nil
}

type recordingSource struct {
	obs    []Observation
	bounds Bounds
	t0, t1 int64
}

func (r *recordingSource) Query(ctx context.Context, bounds Bounds, t0, t1 int64) (QueryResult, error) {
	r.bounds = bounds
	r.t0 = t0
	r.t1 = t1
	return memorySource{obs: r.obs}.Query(ctx, bounds, t0, t1)
}
