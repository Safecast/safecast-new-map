// h3_backfill.go — pre-compute H3 cell indices for existing markers rows.
//
// On startup, a background goroutine walks markers rows that have NULL h3_res5/7/9
// and fills them in batches of 5 000.  New uploads are handled inline by
// updateMarkerH3Indices (called from the upload pipeline).
//
// Column mapping:
//   h3_res5  → H3 resolution 5 (~4.5 km)  used for zoom  6–7
//   h3_res7  → H3 resolution 7 (~0.6 km)  used for zoom 10–11
//   h3_res9  → H3 resolution 9 (~70 m)    used for zoom 14–15

package main

import (
	"context"
	"log"
	"time"

	h3 "github.com/uber/h3-go/v4"
)

const h3BackfillBatch = 5_000

// h3IndicesForLatLon returns the three pre-computed H3 cell IDs as strings.
func h3IndicesForLatLon(lat, lon float64) (res5, res7, res9 string) {
	ll := h3.NewLatLng(lat, lon)
	c5, _ := h3.LatLngToCell(ll, 5)
	c7, _ := h3.LatLngToCell(ll, 7)
	c9, _ := h3.LatLngToCell(ll, 9)
	return c5.String(), c7.String(), c9.String()
}

// updateMarkerH3Indices sets the h3_res5/7/9 columns for a single marker.
// Called inline during upload processing so newly ingested rows are indexed immediately.
func updateMarkerH3Indices(ctx context.Context, markerID int64, lat, lon float64) {
	if db == nil {
		return
	}
	r5, r7, r9 := h3IndicesForLatLon(lat, lon)
	_, err := db.DB.ExecContext(ctx,
		`UPDATE markers SET h3_res5=$1, h3_res7=$2, h3_res9=$3 WHERE id=$4`,
		r5, r7, r9, markerID,
	)
	if err != nil {
		log.Printf("h3_backfill: update marker %d: %v", markerID, err)
	}
}

// backfillH3IndicesBackground runs as a goroutine after server start.
// It fills NULL h3_res* columns in small batches, sleeping 1 s between each
// to avoid overwhelming PostgreSQL during normal operation.
func backfillH3IndicesBackground(ctx context.Context) {
	// Wait a bit after startup so the rest of the server is fully initialised.
	select {
	case <-time.After(15 * time.Second):
	case <-ctx.Done():
		return
	}

	log.Println("h3_backfill: starting background H3 index fill")
	total := 0
	for {
		n, err := backfillH3Batch(ctx, h3BackfillBatch)
		if err != nil {
			log.Printf("h3_backfill: batch error: %v — retrying in 60 s", err)
			select {
			case <-time.After(60 * time.Second):
			case <-ctx.Done():
				return
			}
			continue
		}
		total += n
		if n == 0 {
			log.Printf("h3_backfill: complete — %d rows indexed", total)
			return
		}
		log.Printf("h3_backfill: indexed %d rows (total so far: %d)", n, total)
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

// backfillH3Batch fills up to batchSize markers that still have NULL h3_res5.
// Returns the number of rows updated (0 means all rows are filled).
func backfillH3Batch(ctx context.Context, batchSize int) (int, error) {
	if db == nil {
		return 0, nil
	}
	rows, err := db.DB.QueryContext(ctx,
		`SELECT id, latitude, longitude FROM markers
		 WHERE h3_res5 IS NULL AND latitude IS NOT NULL AND longitude IS NOT NULL
		 LIMIT $1`,
		batchSize,
	)
	if err != nil {
		return 0, err
	}
	type row struct {
		id  int64
		lat float64
		lon float64
	}
	var batch []row
	for rows.Next() {
		var r row
		if err := rows.Scan(&r.id, &r.lat, &r.lon); err != nil {
			rows.Close()
			return 0, err
		}
		batch = append(batch, r)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return 0, err
	}
	if len(batch) == 0 {
		return 0, nil
	}

	for _, r := range batch {
		r5, r7, r9 := h3IndicesForLatLon(r.lat, r.lon)
		if _, err := db.DB.ExecContext(ctx,
			`UPDATE markers SET h3_res5=$1, h3_res7=$2, h3_res9=$3 WHERE id=$4`,
			r5, r7, r9, r.id,
		); err != nil {
			return 0, err
		}
		if ctx.Err() != nil {
			return len(batch), nil
		}
	}
	return len(batch), nil
}
