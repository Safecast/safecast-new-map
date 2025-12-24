package database

import (
	"context"
	"fmt"
)

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
