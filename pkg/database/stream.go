package database

import (
	"context"
	"fmt"
)

// StreamMarkersForH3 streams at most `limit` markers for H3 grid aggregation.
// Uses TABLESAMPLE SYSTEM for large viewports (low zoom) so Postgres scans only
// a fraction of the table, keeping response times under ~1 second.
func (db *Database) StreamMarkersForH3(ctx context.Context, zoom int, minLat, minLon, maxLat, maxLon float64, dbType string, limit int) (<-chan Marker, <-chan error) {
	out := make(chan Marker)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		// At low zoom, the viewport covers a huge area. Use TABLESAMPLE to scan only
		// a fraction of physical pages rather than millions of rows. At high zoom the
		// bounding box is small enough that a full index scan is fast.
		samplePct := 100.0
		switch {
		case zoom <= 5:
			samplePct = 1.0
		case zoom <= 7:
			samplePct = 5.0
		case zoom <= 9:
			samplePct = 15.0
		}

		var query string
		switch dbType {
		case "pgx":
			// Exclude airborne measurements: speed > 150 m/s ≈ 540 km/h (planes ~200-270 m/s)
			speedFilter := "AND (speed IS NULL OR speed = 0 OR speed < 150)"
			if samplePct < 100 {
				query = fmt.Sprintf(`
					SELECT id, doserate, date, lon, lat, countrate, zoom, COALESCE(speed,0), trackid,
					       COALESCE(altitude,0), COALESCE(detector,''), COALESCE(radiation,''),
					       COALESCE(temperature,0), COALESCE(humidity,0), COALESCE(has_spectrum,FALSE)
					FROM markers TABLESAMPLE SYSTEM(%g)
					WHERE zoom = $1
					  AND geom && ST_MakeEnvelope($4, $2, $5, $3, 4326)
					  %s
					LIMIT %d;`, samplePct, speedFilter, limit)
			} else {
				query = fmt.Sprintf(`
					SELECT id, doserate, date, lon, lat, countrate, zoom, COALESCE(speed,0), trackid,
					       COALESCE(altitude,0), COALESCE(detector,''), COALESCE(radiation,''),
					       COALESCE(temperature,0), COALESCE(humidity,0), COALESCE(has_spectrum,FALSE)
					FROM markers
					WHERE zoom = $1
					  AND geom && ST_MakeEnvelope($4, $2, $5, $3, 4326)
					  %s
					LIMIT %d;`, speedFilter, limit)
			}
		default:
			query = fmt.Sprintf(`
				SELECT id, doseRate, date, lon, lat, countRate, zoom, COALESCE(speed,0), trackID,
				       COALESCE(altitude,0), COALESCE(detector,''), COALESCE(radiation,''),
				       COALESCE(temperature,0), COALESCE(humidity,0), COALESCE(has_spectrum,0)
				FROM markers
				WHERE zoom = ? AND lat BETWEEN ? AND ? AND lon BETWEEN ? AND ?
				LIMIT %d;`, limit)
		}

		rows, err := db.DB.QueryContext(ctx, query, zoom, minLat, maxLat, minLon, maxLon)
		if err != nil {
			errCh <- fmt.Errorf("query markers for h3: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var m Marker
			if err := rows.Scan(&m.ID, &m.DoseRate, &m.Date, &m.Lon, &m.Lat, &m.CountRate, &m.Zoom, &m.Speed, &m.TrackID,
				&m.Altitude, &m.Detector, &m.Radiation, &m.Temperature, &m.Humidity, &m.HasSpectrum); err != nil {
				errCh <- fmt.Errorf("scan marker: %w", err)
				return
			}
			select {
			case out <- m:
			case <-ctx.Done():
				return
			}
		}
		if err := rows.Err(); err != nil {
			errCh <- fmt.Errorf("iterate markers: %w", err)
		}
	}()
	return out, errCh
}

// StreamMarkersByZoomAndBounds streams markers row by row through a channel.
// It avoids loading large result sets into memory and stops when the context is done.
func (db *Database) StreamMarkersByZoomAndBounds(ctx context.Context, zoom int, minLat, minLon, maxLat, maxLon float64, dbType string) (<-chan Marker, <-chan error) {
	out := make(chan Marker)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		var query string
		switch dbType {
		case "pgx":
			// Use PostGIS spatial index with ST_Intersects and && bounding box operator
			// for optimal performance. The && operator uses the GIST index efficiently.
			// Note: PostgreSQL folds unquoted identifiers to lowercase, so use lowercase column names
			query = `
                SELECT id, doserate, date, lon, lat, countrate, zoom, COALESCE(speed, 0) as speed, trackid,
                       COALESCE(altitude, 0) as altitude,
                       COALESCE(detector, '') as detector,
                       COALESCE(radiation, '') as radiation,
                       COALESCE(temperature, 0) as temperature,
                       COALESCE(humidity, 0) as humidity,
                       COALESCE(has_spectrum, FALSE) as has_spectrum
                FROM markers
                WHERE zoom = $1
                  AND geom && ST_MakeEnvelope($4, $2, $5, $3, 4326)
                  AND ST_Intersects(geom, ST_MakeEnvelope($4, $2, $5, $3, 4326));
            `
		default:
			query = `
                SELECT id, doseRate, date, lon, lat, countRate, zoom, COALESCE(speed, 0) as speed, trackID,
                       COALESCE(altitude, 0) as altitude,
                       COALESCE(detector, '') as detector,
                       COALESCE(radiation, '') as radiation,
                       COALESCE(temperature, 0) as temperature,
                       COALESCE(humidity, 0) as humidity,
                       COALESCE(has_spectrum, 0) as has_spectrum
                FROM markers
                WHERE zoom = ? AND lat BETWEEN ? AND ? AND lon BETWEEN ? AND ?;
            `
		}

		rows, err := db.DB.QueryContext(ctx, query, zoom, minLat, maxLat, minLon, maxLon)
		if err != nil {
			errCh <- fmt.Errorf("query markers: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var m Marker
			if err := rows.Scan(&m.ID, &m.DoseRate, &m.Date, &m.Lon, &m.Lat, &m.CountRate, &m.Zoom, &m.Speed, &m.TrackID,
				&m.Altitude, &m.Detector, &m.Radiation, &m.Temperature, &m.Humidity, &m.HasSpectrum); err != nil {
				errCh <- fmt.Errorf("scan marker: %w", err)
				return
			}
			select {
			case out <- m:
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}

		if err := rows.Err(); err != nil {
			errCh <- fmt.Errorf("iterate markers: %w", err)
		}
	}()

	return out, errCh
}

// StreamMarkersByTrackIDZoomAndBounds streams markers of one track within bounds.
// This keeps memory usage low while focusing on a single track only.
func (db *Database) StreamMarkersByTrackIDZoomAndBounds(ctx context.Context, trackID string, zoom int, minLat, minLon, maxLat, maxLon float64, dbType string) (<-chan Marker, <-chan error) {
	out := make(chan Marker)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		var query string
		switch dbType {
		case "pgx":
			// Use PostGIS spatial index with ST_Intersects and && bounding box operator
			// for optimal performance. The && operator uses the GIST index efficiently.
			// Note: PostgreSQL folds unquoted identifiers to lowercase, so use lowercase column names
			query = `
                SELECT id, doserate, date, lon, lat, countrate, zoom, COALESCE(speed, 0) as speed, trackid,
                       COALESCE(altitude, 0) as altitude,
                       COALESCE(detector, '') as detector,
                       COALESCE(radiation, '') as radiation,
                       COALESCE(temperature, 0) as temperature,
                       COALESCE(humidity, 0) as humidity,
                       COALESCE(has_spectrum, FALSE) as has_spectrum
                FROM markers
                WHERE trackid = $1
                  AND zoom = $2
                  AND geom && ST_MakeEnvelope($5, $3, $6, $4, 4326)
                  AND ST_Intersects(geom, ST_MakeEnvelope($5, $3, $6, $4, 4326));
            `
		default:
			query = `
                SELECT id, doseRate, date, lon, lat, countRate, zoom, COALESCE(speed, 0) as speed, trackID,
                       COALESCE(altitude, 0) as altitude,
                       COALESCE(detector, '') as detector,
                       COALESCE(radiation, '') as radiation,
                       COALESCE(temperature, 0) as temperature,
                       COALESCE(humidity, 0) as humidity,
                       COALESCE(has_spectrum, 0) as has_spectrum
                FROM markers
                WHERE trackID = ? AND zoom = ? AND lat BETWEEN ? AND ? AND lon BETWEEN ? AND ?;
            `
		}

		rows, err := db.DB.QueryContext(ctx, query, trackID, zoom, minLat, maxLat, minLon, maxLon)
		if err != nil {
			errCh <- fmt.Errorf("query markers: %w", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var m Marker
			if err := rows.Scan(&m.ID, &m.DoseRate, &m.Date, &m.Lon, &m.Lat, &m.CountRate, &m.Zoom, &m.Speed, &m.TrackID,
				&m.Altitude, &m.Detector, &m.Radiation, &m.Temperature, &m.Humidity, &m.HasSpectrum); err != nil {
				errCh <- fmt.Errorf("scan marker: %w", err)
				return
			}
			select {
			case out <- m:
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			}
		}

		if err := rows.Err(); err != nil {
			errCh <- fmt.Errorf("iterate markers: %w", err)
		}
	}()

	return out, errCh
}

// StreamMarkersByRadius streams markers within a specified radius (in meters) from a center point.
// Uses PostGIS ST_DWithin for efficient distance-based queries.
// Only supported for PostgreSQL (pgx driver).
func (db *Database) StreamMarkersByRadius(ctx context.Context, centerLat, centerLon float64, radiusMeters int, dbType string) (<-chan Marker, <-chan error) {
out := make(chan Marker)
errCh := make(chan error, 1)

go func() {
defer close(out)
defer close(errCh)

var query string
switch dbType {
case "pgx":
// Use PostGIS ST_DWithin for efficient radius-based queries
// ST_DWithin uses the GIST index and returns points within distance in meters
// Note: ST_DWithin requires geography type for accurate distance calculations
query = `
                SELECT id, doserate, date, lon, lat, countrate, zoom, COALESCE(speed, 0) as speed, trackid,
                       COALESCE(altitude, 0) as altitude,
                       COALESCE(detector, '') as detector,
                       COALESCE(radiation, '') as radiation,
                       COALESCE(temperature, 0) as temperature,
                       COALESCE(humidity, 0) as humidity,
                       COALESCE(has_spectrum, FALSE) as has_spectrum
                FROM markers
                WHERE ST_DWithin(
                    geom::geography,
                    ST_MakePoint($2, $1)::geography,
                    $3
                )
                ORDER BY ST_Distance(geom::geography, ST_MakePoint($2, $1)::geography) ASC;
            `
default:
errCh <- fmt.Errorf("radius search only supported for PostgreSQL (pgx driver)")
return
}

rows, err := db.DB.QueryContext(ctx, query, centerLat, centerLon, radiusMeters)
if err != nil {
errCh <- fmt.Errorf("query markers by radius: %w", err)
return
}
defer rows.Close()

for rows.Next() {
var m Marker
if err := rows.Scan(&m.ID, &m.DoseRate, &m.Date, &m.Lon, &m.Lat, &m.CountRate, &m.Zoom, &m.Speed, &m.TrackID,
&m.Altitude, &m.Detector, &m.Radiation, &m.Temperature, &m.Humidity, &m.HasSpectrum); err != nil {
errCh <- fmt.Errorf("scan marker: %w", err)
return
}
select {
case out <- m:
case <-ctx.Done():
errCh <- ctx.Err()
return
}
}

if err := rows.Err(); err != nil {
errCh <- fmt.Errorf("iterate markers: %w", err)
}
}()

return out, errCh
}
