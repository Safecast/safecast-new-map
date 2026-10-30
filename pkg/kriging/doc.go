// Package kriging estimates a local radiation field for the map overlay.
//
// The model is ordinary kriging of log(dose rate in µSv/h). The map shows the
// lognormal median exp(μ) and keeps the kriging variance on the log scale.
// Variance is uncertainty in the latent field: the measurement nugget is on the
// observation diagonal only, and it is left out of the covariance with the
// prediction target. A single measurement therefore does not produce zero
// uncertainty.
//
// Covariance is a fixed separable exponential, not a variogram fitted to drive
// tracks. Track points are highly redundant, and a fit on them would shrink the
// range. Spatial and temporal factors are separate so either can be replaced
// later. TrackID and DeviceID are kept on every observation for that later
// error term; this covariance ignores them.
//
// Practical range is the distance or age where correlation falls to about 0.05:
// 500 m and 14 days. Relative nugget is 0.05 of the partial sill. The partial
// sill is 1 on the log scale, so reported variance is already variance / sill.
// Cells whose variance exceeds 0.75 of the partial sill are marked suppressed
// and the overlay does not draw them. A blank area is unsupported, not safe.
//
// Neighborhoods overlap. Each grid point uses the observations most correlated
// with it, and the viewport query is expanded by the practical ranges. There is
// no regional model and no hard boundary between tiles. Viewports wider than
// 20 km, or windows with more than 8000 raw points, are refused instead of
// subsampled. The streamed marker layer is not an input: it is already thinned
// and, at low zoom, biased toward the hottest point in each cell.
//
// Prediction time is the end of the map date slider, or the current time when
// the slider is cleared. Observations on both sides of that time can enter the
// window. That is interpolation. A causal event mode is intentionally not
// implemented here.
//
// # Tests
//
//	go test ./pkg/kriging/
//
// # Visual check, synthetic field
//
// Restart the server after a change; map.html is embedded in the binary.
//
//	go run ./cmd/unified-server -db-type sqlite -db-path ./dev.db
//
// Open http://localhost:8765/?kriging=synthetic and turn on "Radiation estimate".
// The map zooms to a fake track. Color should follow that track, fade within a
// few hundred meters, and stay blank in the empty northeast. Pan slightly: the
// shared edge should not jump. Turn the control off and the base map is unchanged.
// The synthetic flag does not write to the database.
//
// # Visual check, stored measurements
//
// Open http://localhost:8765/ against a database that has zoom = 0 rows and zoom
// into a surveyed area until the overlay draws. Gaps between tracks stay blank.
// Markers stay clickable above the color.
package kriging
