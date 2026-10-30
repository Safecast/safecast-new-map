package kriging

import (
	"context"
	"math"
)

// SyntheticAnchor is the fixed geographic origin of the visual synthetic track.
// It matches the server's default map center so a fresh database opens nearby.
const (
	SyntheticAnchorLat = 44.08832
	SyntheticAnchorLon = 42.97577
)

// SyntheticSource returns a track-shaped sample from a smooth log-dose bump.
// The northeast of SyntheticView is left empty so the overlay can show a blank,
// unsupported quadrant. Times sit on the requested prediction time; the source
// does not read or write the database.
type SyntheticSource struct{}

// Query implements ObservationSource. Observations outside bounds are omitted
// so a viewport away from the anchor stays blank rather than painting the track
// in the wrong place.
func (SyntheticSource) Query(_ context.Context, bounds Bounds, t0, t1 int64) (QueryResult, error) {
	when := t0
	if t1 >= t0 {
		when = t0 + (t1-t0)/2
	}
	var kept []Observation
	for _, o := range SyntheticObservations(when) {
		if contains(bounds, o.Lat, o.Lon) {
			kept = append(kept, o)
		}
	}
	return QueryResult{Observations: kept}, nil
}

// SyntheticObservations builds the fake track at when. The dose is a smooth bump
// along the line, highest near the middle of the track.
func SyntheticObservations(when int64) []Observation {
	mLat := metersPerDegreeLat()
	mLon := mLat * math.Cos(SyntheticAnchorLat*math.Pi/180)
	var out []Observation
	for east := -800.0; east <= 100.0; east += 40 {
		along := east - (-800)
		peak := math.Exp(-math.Pow((along-450)/120, 2))
		dose := 0.10 * math.Exp(1.2*peak)
		out = append(out, Observation{
			Lat:      SyntheticAnchorLat,
			Lon:      SyntheticAnchorLon + east/mLon,
			TimeUnix: when,
			DoseRate: dose,
			TrackID:  "synthetic-track",
			DeviceID: "synthetic-device",
		})
	}
	return out
}

// SyntheticView is the rectangle the overlay fits to. The track runs along the
// southern part; the northeast corner is more than a kilometer from the nearest point.
func SyntheticView() Bounds {
	mLat := metersPerDegreeLat()
	mLon := mLat * math.Cos(SyntheticAnchorLat*math.Pi/180)
	return Bounds{
		MinLat: SyntheticAnchorLat + (-400)/mLat,
		MaxLat: SyntheticAnchorLat + 1200/mLat,
		MinLon: SyntheticAnchorLon + (-900)/mLon,
		MaxLon: SyntheticAnchorLon + 700/mLon,
	}
}
