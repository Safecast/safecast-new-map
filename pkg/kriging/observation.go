package kriging

import (
	"context"
	"math"
)

// EarthRadiusM matches the radius used by the map's distance helpers.
const EarthRadiusM = 6371000

// MaxViewportMeters is the widest viewport this edition will krige.
// A wider box would force subsampling or one model for a whole region.
const MaxViewportMeters = 20000

// DefaultMaxObservations is the most raw points one window may hold.
// Above this the request is refused. The points are not thinned.
const DefaultMaxObservations = 8000

// Observation is one raw measurement. TrackID and DeviceID are retained so a
// later covariance can add track-level or device-level error without a new solver.
type Observation struct {
	Lat      float64
	Lon      float64
	TimeUnix int64
	DoseRate float64 // µSv/h. Non-positive values never enter the log model.
	TrackID  string
	DeviceID string
}

// Target is the location and time at which the latent field is predicted.
type Target struct {
	Lat      float64
	Lon      float64
	TimeUnix int64
}

// Bounds is a geographic query rectangle. MinLon > MaxLon means the rectangle
// crosses the antimeridian.
type Bounds struct {
	MinLat float64 `json:"minLat"`
	MinLon float64 `json:"minLon"`
	MaxLat float64 `json:"maxLat"`
	MaxLon float64 `json:"maxLon"`
}

// Valid reports whether the rectangle has positive area.
func (b Bounds) Valid() bool {
	if math.IsNaN(b.MinLat) || math.IsNaN(b.MaxLat) || math.IsNaN(b.MinLon) || math.IsNaN(b.MaxLon) {
		return false
	}
	if b.MinLat >= b.MaxLat || b.MinLon == b.MaxLon {
		return false
	}
	return true
}

// QueryResult is a raw window. OmittedNonPositive counts rows that were dropped
// because a log model cannot use a non-positive dose. They are not replaced
// with a stand-in value.
type QueryResult struct {
	Observations       []Observation
	OmittedNonPositive int
}

// ObservationSource retrieves raw measurements. Implementations must not cluster,
// thin, or keep only the maximum dose in a cell.
type ObservationSource interface {
	Query(ctx context.Context, bounds Bounds, t0, t1 int64) (QueryResult, error)
}

// HaversineMeters is the great-circle distance between two WGS84 points.
func HaversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	φ1 := lat1 * math.Pi / 180
	φ2 := lat2 * math.Pi / 180
	dφ := (lat2 - lat1) * math.Pi / 180
	dλ := (lon2 - lon1) * math.Pi / 180
	a := math.Sin(dφ/2)*math.Sin(dφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(dλ/2)*math.Sin(dλ/2)
	return 2 * EarthRadiusM * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func metersPerDegreeLat() float64 {
	return 2 * math.Pi * EarthRadiusM / 360
}

// Expand grows bounds by meters on every side so a prediction on the viewport
// edge still sees observations just outside it. The halo is retrieval, not a
// second model: correlation at the practical range is already about 0.05.
func Expand(b Bounds, meters float64) Bounds {
	if meters <= 0 {
		return b
	}
	dLat := meters / metersPerDegreeLat()
	mid := (b.MinLat + b.MaxLat) / 2 * math.Pi / 180
	cos := math.Cos(mid)
	if cos < 1e-6 {
		cos = 1e-6
	}
	dLon := meters / (metersPerDegreeLat() * cos)
	return Bounds{
		MinLat: b.MinLat - dLat,
		MaxLat: b.MaxLat + dLat,
		MinLon: b.MinLon - dLon,
		MaxLon: b.MaxLon + dLon,
	}
}

// SpanMeters is the longer side of the rectangle. Antimeridian crossings are
// measured the short way across ±180, which is the viewport the map is showing.
func SpanMeters(b Bounds) float64 {
	midLat := (b.MinLat + b.MaxLat) / 2
	midLon := b.MinLon
	if b.MinLon <= b.MaxLon {
		midLon = (b.MinLon + b.MaxLon) / 2
	}
	height := HaversineMeters(b.MinLat, midLon, b.MaxLat, midLon)
	var width float64
	if b.MinLon <= b.MaxLon {
		width = HaversineMeters(midLat, b.MinLon, midLat, b.MaxLon)
	} else {
		width = HaversineMeters(midLat, b.MinLon, midLat, 180) + HaversineMeters(midLat, -180, midLat, b.MaxLon)
	}
	if height > width {
		return height
	}
	return width
}

func contains(b Bounds, lat, lon float64) bool {
	if lat < b.MinLat || lat > b.MaxLat {
		return false
	}
	if b.MinLon <= b.MaxLon {
		return lon >= b.MinLon && lon <= b.MaxLon
	}
	return lon >= b.MinLon || lon <= b.MaxLon
}
