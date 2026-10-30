package kriging

import "math"

const (
	// SpatialPracticalRangeM is where the spatial factor falls to about exp(-3).
	SpatialPracticalRangeM float64 = 500
	// TemporalPracticalRangeSec is 14 days, where the temporal factor falls to about exp(-3).
	TemporalPracticalRangeSec int64 = 14 * 24 * 60 * 60
	// DefaultPartialSill is the variance of the latent log-dose field.
	// Weights are invariant to scaling it; the reported variance is in these units.
	DefaultPartialSill = 1
	// RelativeNugget is measurement-error variance as a fraction of the partial sill.
	RelativeNugget = 0.05
	// SuppressionFraction is the variance, relative to the partial sill, above which
	// a cell is not drawn. At the sill the neighborhood has told us almost nothing.
	SuppressionFraction = 0.75
	// ModelName identifies this fixed covariance in API responses.
	ModelName = "ordinary_kriging_log_separable_exponential"
)

// CovarianceModel calculates covariance. Pair is the latent-field covariance and
// does not include measurement error. MeasurementNugget is added only on the
// observation diagonal. Distance uses the same separable metric as Pair, so the
// neighborhood keeps the most correlated observations.
//
// TrackID and DeviceID are visible on the observations. A later model can use
// them for a track or device term without changing the solver.
type CovarianceModel interface {
	Pair(a, b Observation) float64
	MeasurementNugget() float64
	FieldVariance() float64
	Distance(target Target, o Observation) float64
}

// ExponentialModel is a separable exponential covariance with fixed ranges.
type ExponentialModel struct {
	PartialSill               float64
	Nugget                    float64
	SpatialPracticalRangeM    float64
	TemporalPracticalRangeSec int64
}

// DefaultModel returns the prior used by the map layer and the tests.
func DefaultModel() ExponentialModel {
	return ExponentialModel{
		PartialSill:               DefaultPartialSill,
		Nugget:                    RelativeNugget * DefaultPartialSill,
		SpatialPracticalRangeM:    SpatialPracticalRangeM,
		TemporalPracticalRangeSec: TemporalPracticalRangeSec,
	}
}

// Pair is σ² exp(-3 d / range_space) exp(-3 |Δt| / range_time).
func (m ExponentialModel) Pair(a, b Observation) float64 {
	return m.PartialSill * math.Exp(-m.distanceObs(a, b))
}

// MeasurementNugget is the log-scale measurement variance.
func (m ExponentialModel) MeasurementNugget() float64 { return m.Nugget }

// FieldVariance is the partial sill of the latent log-dose field.
func (m ExponentialModel) FieldVariance() float64 { return m.PartialSill }

// Distance is the dimensionless separable distance from a prediction target.
func (m ExponentialModel) Distance(target Target, o Observation) float64 {
	return m.separable(target.Lat, target.Lon, target.TimeUnix, o.Lat, o.Lon, o.TimeUnix)
}

func (m ExponentialModel) distanceObs(a, b Observation) float64 {
	return m.separable(a.Lat, a.Lon, a.TimeUnix, b.Lat, b.Lon, b.TimeUnix)
}

func (m ExponentialModel) separable(lat1, lon1 float64, t1 int64, lat2, lon2 float64, t2 int64) float64 {
	spatial := HaversineMeters(lat1, lon1, lat2, lon2) / m.SpatialPracticalRangeM
	dt := t1 - t2
	if dt < 0 {
		dt = -dt
	}
	temporal := float64(dt) / float64(m.TemporalPracticalRangeSec)
	return 3*spatial + 3*temporal
}

func (m ExponentialModel) usable() bool {
	return m.PartialSill > 0 && m.Nugget >= 0 && m.SpatialPracticalRangeM > 0 && m.TemporalPracticalRangeSec > 0
}

// Suppress reports whether a latent variance is too large to draw.
func Suppress(fieldVariance, logVariance float64) bool {
	if math.IsNaN(logVariance) || math.IsInf(logVariance, 0) {
		return true
	}
	return logVariance > SuppressionFraction*fieldVariance
}
