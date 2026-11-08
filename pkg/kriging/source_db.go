package kriging

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	"safecast-new-map/pkg/database"
)

// ErrDatabaseUnavailable means the map has no database to read raw markers from.
var ErrDatabaseUnavailable = errors.New("kriging: database unavailable")

// DBSource reads raw zoom = 0 marker rows. It does not use the streamed or
// clustered marker feed: those paths sample the archive and, below zoom 13,
// keep the hottest point in a cell, which would bias the field high.
//
// The query uses database/sql directly. The single-writer pipeline in
// pkg/database is unexported, and database/sql already serializes the one
// connection those engines open. The 20 km window and the row cap bound how
// long that read can delay another statement.
type DBSource struct {
	DB *database.Database
	// MaxObservations caps raw rows in one window. Zero uses DefaultMaxObservations.
	MaxObservations int
}

// NewDBSource adapts a database. A nil database fails at query time.
func NewDBSource(db *database.Database) *DBSource {
	return &DBSource{DB: db, MaxObservations: DefaultMaxObservations}
}

// Query returns raw observations inside bounds and [t0, t1]. Non-positive doses
// are counted and omitted. Crossing the row cap returns ErrTooManyObservations
// without a truncated sample.
func (s *DBSource) Query(ctx context.Context, bounds Bounds, t0, t1 int64) (QueryResult, error) {
	if s == nil || s.DB == nil || s.DB.DB == nil {
		return QueryResult{}, ErrDatabaseUnavailable
	}
	limit := s.MaxObservations
	if limit <= 0 {
		limit = DefaultMaxObservations
	}
	var (
		out     []Observation
		omitted int
	)
	// A point on the antimeridian can match both halves of a split window.
	// Counting it twice would treat one measurement as two.
	seen := make(map[string]struct{})
	for _, box := range queryBoxes(bounds) {
		rows, err := s.queryBox(ctx, box, t0, t1)
		if err != nil {
			return QueryResult{}, err
		}
		for _, row := range rows {
			key := rowKey(row)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			if len(seen) > limit {
				return QueryResult{}, ErrTooManyObservations
			}
			if !row.dose.Valid || row.dose.Float64 <= 0 {
				omitted++
				continue
			}
			out = append(out, Observation{
				Lat:      row.lat,
				Lon:      row.lon,
				TimeUnix: row.when,
				DoseRate: row.dose.Float64,
				TrackID:  row.trackID,
				DeviceID: row.deviceID,
			})
		}
	}
	return QueryResult{Observations: out, OmittedNonPositive: omitted}, nil
}

func rowKey(row markerRow) string {
	dose := math.NaN()
	if row.dose.Valid {
		dose = row.dose.Float64
	}
	return fmt.Sprintf("%.7f|%.7f|%d|%.8f|%s|%s", row.lat, row.lon, row.when, dose, row.trackID, row.deviceID)
}

type markerRow struct {
	lat, lon float64
	dose     sql.NullFloat64
	when     int64
	trackID  string
	deviceID string
}

func (s *DBSource) queryBox(ctx context.Context, box Bounds, t0, t1 int64) ([]markerRow, error) {
	driver := strings.ToLower(s.DB.Driver)
	var (
		query string
		args  []any
	)
	if driver == "pgx" {
		query = `
			SELECT lat, lon, doserate, date, COALESCE(trackid, ''), COALESCE(device_id, '')
			FROM markers
			WHERE zoom = $1
			  AND date BETWEEN $2 AND $3
			  AND geom && ST_MakeEnvelope($4, $5, $6, $7, 4326)
			  AND ST_Intersects(geom, ST_MakeEnvelope($4, $5, $6, $7, 4326))`
		args = []any{0, t0, t1, box.MinLon, box.MinLat, box.MaxLon, box.MaxLat}
	} else {
		query = `
			SELECT lat, lon, doseRate, date, COALESCE(trackID, ''), COALESCE(device_id, '')
			FROM markers
			WHERE zoom = ?
			  AND date BETWEEN ? AND ?
			  AND lat BETWEEN ? AND ?
			  AND lon BETWEEN ? AND ?`
		args = []any{0, t0, t1, box.MinLat, box.MaxLat, box.MinLon, box.MaxLon}
	}
	rows, err := s.DB.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query kriging observations: %w", err)
	}
	defer rows.Close()

	var out []markerRow
	for rows.Next() {
		var row markerRow
		if err := rows.Scan(&row.lat, &row.lon, &row.dose, &row.when, &row.trackID, &row.deviceID); err != nil {
			return nil, fmt.Errorf("scan kriging observation: %w", err)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate kriging observations: %w", err)
	}
	return out, nil
}

func queryBoxes(b Bounds) []Bounds {
	if b.MinLon > b.MaxLon {
		return []Bounds{
			{MinLat: b.MinLat, MaxLat: b.MaxLat, MinLon: b.MinLon, MaxLon: 180},
			{MinLat: b.MinLat, MaxLat: b.MaxLat, MinLon: -180, MaxLon: b.MaxLon},
		}
	}
	if b.MinLon >= -180 && b.MaxLon <= 180 {
		return []Bounds{b}
	}
	// The halo pushed the rectangle past ±180. Query both sides rather than
	// dropping the observations that wrapped.
	var boxes []Bounds
	if b.MinLon < -180 {
		boxes = append(boxes, Bounds{MinLat: b.MinLat, MaxLat: b.MaxLat, MinLon: b.MinLon + 360, MaxLon: 180})
		boxes = append(boxes, Bounds{MinLat: b.MinLat, MaxLat: b.MaxLat, MinLon: -180, MaxLon: math.Min(b.MaxLon, 180)})
	}
	if b.MaxLon > 180 {
		boxes = append(boxes, Bounds{MinLat: b.MinLat, MaxLat: b.MaxLat, MinLon: math.Max(b.MinLon, -180), MaxLon: 180})
		boxes = append(boxes, Bounds{MinLat: b.MinLat, MaxLat: b.MaxLat, MinLon: -180, MaxLon: b.MaxLon - 360})
	}
	if len(boxes) == 0 {
		return []Bounds{b}
	}
	return boxes
}
