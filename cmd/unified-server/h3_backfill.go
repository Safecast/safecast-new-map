// h3_backfill.go — pre-compute H3 cell indices for existing markers rows.
//
// Strategy: divide the markers ID range into N worker chunks, each worker
// fetches its own batch and bulk-UPDATEs via unnest — no row-by-row round trips.
//
// Column mapping:
//   h3_res5  → H3 resolution 5 (~4.5 km)  used for zoom  6–7
//   h3_res7  → H3 resolution 7 (~0.6 km)  used for zoom 10–11
//   h3_res9  → H3 resolution 9 (~70 m)    used for zoom 14–15

package main

import (
	"context"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	h3 "github.com/uber/h3-go/v4"
)

const h3BackfillBatch = 10_000 // rows per worker per iteration

// h3IndicesForLatLon returns the three pre-computed H3 cell IDs as strings.
func h3IndicesForLatLon(lat, lon float64) (res5, res7, res9 string) {
	ll := h3.NewLatLng(lat, lon)
	c5, _ := h3.LatLngToCell(ll, 5)
	c7, _ := h3.LatLngToCell(ll, 7)
	c9, _ := h3.LatLngToCell(ll, 9)
	return c5.String(), c7.String(), c9.String()
}

// updateMarkerH3Indices sets h3_res5/7/9 for a single marker (called on new uploads).
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
// It spawns numCPU/2 workers (min 4, max 12) that each process independent
// ID-range slices in parallel.
func backfillH3IndicesBackground(ctx context.Context) {
	select {
	case <-time.After(15 * time.Second):
	case <-ctx.Done():
		return
	}

	if db == nil {
		return
	}

	// Find min/max IDs with NULL h3_res5.
	var minID, maxID int64
	err := db.DB.QueryRowContext(ctx,
		`SELECT COALESCE(MIN(id),0), COALESCE(MAX(id),0) FROM markers WHERE h3_res5 IS NULL AND lat IS NOT NULL`,
	).Scan(&minID, &maxID)
	if err != nil || maxID == 0 {
		log.Printf("h3_backfill: nothing to fill (err=%v)", err)
		return
	}

	nWorkers := runtime.NumCPU() / 2
	if nWorkers < 4 {
		nWorkers = 4
	}
	if nWorkers > 12 {
		nWorkers = 12
	}

	log.Printf("h3_backfill: starting — IDs %d–%d, %d parallel workers", minID, maxID, nWorkers)

	// Work queue: each item is a (fromID, toID) range slice.
	type workItem struct{ from, to int64 }
	sliceSize := (maxID - minID + 1) / int64(nWorkers)
	if sliceSize < int64(h3BackfillBatch) {
		sliceSize = int64(h3BackfillBatch)
	}

	queue := make(chan workItem, nWorkers*4)
	go func() {
		for lo := minID; lo <= maxID; lo += sliceSize {
			hi := lo + sliceSize - 1
			if hi > maxID {
				hi = maxID
			}
			select {
			case queue <- workItem{lo, hi}:
			case <-ctx.Done():
				close(queue)
				return
			}
		}
		close(queue)
	}()

	var total int64
	var wg sync.WaitGroup
	for i := 0; i < nWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for item := range queue {
				n, err := h3WorkerProcess(ctx, item.from, item.to)
				if err != nil && ctx.Err() == nil {
					log.Printf("h3_backfill worker %d: %v", workerID, err)
				}
				if n > 0 {
					newTotal := atomic.AddInt64(&total, int64(n))
					log.Printf("h3_backfill: indexed %d rows (total: %d)", n, newTotal)
				}
				if ctx.Err() != nil {
					return
				}
			}
		}(i)
	}

	wg.Wait()
	log.Printf("h3_backfill: complete — %d rows indexed", atomic.LoadInt64(&total))
}

// h3WorkerProcess fills h3_res* for all NULL rows in [fromID, toID].
// Uses bulk unnest UPDATE to minimise round trips.
func h3WorkerProcess(ctx context.Context, fromID, toID int64) (int, error) {
	total := 0
	cursor := fromID

	for {
		rows, err := db.DB.QueryContext(ctx,
			`SELECT id, lat, lon FROM markers
			 WHERE id BETWEEN $1 AND $2
			   AND h3_res5 IS NULL
			   AND lat IS NOT NULL AND lon IS NOT NULL
			 ORDER BY id
			 LIMIT $3`,
			cursor, toID, h3BackfillBatch,
		)
		if err != nil {
			return total, err
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
				return total, err
			}
			batch = append(batch, r)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return total, err
		}
		if len(batch) == 0 {
			break
		}

		// Compute H3 indices in Go (CPU-bound, fast) then bulk-write in one transaction.
		tx, err := db.DB.BeginTx(ctx, nil)
		if err != nil {
			return total, err
		}
		stmt, err := tx.PrepareContext(ctx,
			`UPDATE markers SET h3_res5=$1, h3_res7=$2, h3_res9=$3 WHERE id=$4`)
		if err != nil {
			tx.Rollback() //nolint:errcheck
			return total, err
		}
		for _, r := range batch {
			r5, r7, r9 := h3IndicesForLatLon(r.lat, r.lon)
			if _, err = stmt.ExecContext(ctx, r5, r7, r9, r.id); err != nil {
				stmt.Close()
				tx.Rollback() //nolint:errcheck
				return total, err
			}
		}
		stmt.Close()
		if err = tx.Commit(); err != nil {
			return total, err
		}

		total += len(batch)
		cursor = batch[len(batch)-1].id + 1
		if cursor > toID {
			break
		}

		if ctx.Err() != nil {
			return total, nil
		}
	}
	return total, nil
}
