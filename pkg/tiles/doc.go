// Package tiles delivers map images for a viewport without streaming one
// record per measurement.
//
// A tile is a Web Mercator XYZ square, the same grid Leaflet already uses for
// the base map. A processor turns the measurements inside that square into
// bytes. The aggregate processor paints the hottest dose in each pixel and
// leaves unsampled pixels transparent. Cells are a global Web Mercator grid
// whose geographic size halves at each zoom. Color is the hottest dose in
// the cell. Opacity is the observation count on a logarithmic scale, so a
// single reading stays faint. Another processor can be registered
// later without changing the cache, the HTTP routes, or the browser layer.
//
// This package does not import or copy a kriging implementation. The processor
// interface is the only seam for that work.
//
// The occupied-cell sidecar is intentionally absent. The live benchmark in
// bench_test.go (go test -tags tilebench) is what would show that a cold
// low-zoom scan misses the latency budget. This environment has no dense
// PostGIS snapshot, so that measurement was not available, and on-demand
// generation remains the only path.
package tiles
