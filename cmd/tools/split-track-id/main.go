// split-track-id repairs a trackID that DetectExistingTrackID incorrectly
// merged across unrelated uploads (see pkg/database/database.go
// DetectExistingTrackID). It splits the markers of one corrupted trackID
// back into one trackID per distinct drive, using each upload's filename
// timestamp as the clustering key.
//
// uploads.recording_date is NOT used: it's computed as
// `MIN(date) FROM markers WHERE trackID = ...` (handlers_upload.go), so once
// a trackID is wrongly merged, every later upload into that same trackID
// inherits the *earliest* upload's recording_date instead of its own. The
// bGeigieZen filename (YYYY-MM-DD_HHMM.log, the drive start time) is the
// only reliable signal left; raw uploaded files are not retained so exact
// per-upload boundaries can't be recovered from content. This is the best
// available proxy.
//
// Usage:
//
//	go run ./cmd/tools/split-track-id -track-id U68RVV -db-conn "postgres://..." [-apply]
//
// Without -apply it only prints the plan (dry run).
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var filenameDateRe = regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})_(\d{2})(\d{2})`)

func dateFromFilename(name string) (time.Time, bool) {
	m := filenameDateRe.FindStringSubmatch(name)
	if m == nil {
		return time.Time{}, false
	}
	t, err := time.Parse("2006-01-02_1504", fmt.Sprintf("%s-%s-%s_%s%s", m[1], m[2], m[3], m[4], m[5]))
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func newTrackID(seed int64) string {
	const maxLength = 6
	timestamp := uint64(seed)
	encoded := ""
	base := uint64(len(base62Chars))
	for timestamp > 0 && len(encoded) < maxLength {
		encoded = string(base62Chars[timestamp%base]) + encoded
		timestamp /= base
	}
	r := rand.New(rand.NewSource(seed))
	for len(encoded) < maxLength {
		encoded += string(base62Chars[r.Intn(len(base62Chars))])
	}
	return encoded
}

type uploadRow struct {
	id            int64
	filename      string
	ts            time.Time
	tsFromComment string // "filename" or "recording_date (fallback)"
}

// cluster groups upload rows that share the same drive (identical filename
// timestamp) — they get reassigned to the same new trackID.
type cluster struct {
	ts      time.Time
	uploads []uploadRow
}

func main() {
	dbConn := flag.String("db-conn", "postgres://postgres@127.0.0.1:5432/safecast", "Postgres connection string")
	trackID := flag.String("track-id", "", "corrupted trackID to split (required)")
	apply := flag.Bool("apply", false, "actually write changes; without this flag only prints the plan")
	flag.Parse()

	if *trackID == "" {
		log.Fatal("-track-id is required")
	}

	db, err := sql.Open("pgx", *dbConn)
	if err != nil {
		log.Fatalf("connecting to database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("pinging database: %v", err)
	}
	ctx := context.Background()

	uploads, err := loadUploads(ctx, db, *trackID)
	if err != nil {
		log.Fatalf("loading uploads: %v", err)
	}
	if len(uploads) < 2 {
		log.Fatalf("trackID %s has %d upload(s); nothing to split", *trackID, len(uploads))
	}

	clusters := groupByTimestamp(uploads)
	sort.Slice(clusters, func(i, j int) bool { return clusters[i].ts.Before(clusters[j].ts) })

	// Boundaries are midpoints between consecutive clusters' timestamps.
	// Marker dates on or after boundary[i] and before boundary[i+1] belong
	// to clusters[i].
	boundaries := make([]int64, len(clusters)+1)
	boundaries[0] = 0
	boundaries[len(clusters)] = 1 << 62
	for i := 0; i < len(clusters)-1; i++ {
		boundaries[i+1] = (clusters[i].ts.Unix() + clusters[i+1].ts.Unix()) / 2
	}

	newIDs := make([]string, len(clusters))
	used := map[string]bool{}
	for i, c := range clusters {
		id := newTrackID(c.ts.UnixNano() + int64(i))
		for used[id] {
			id = newTrackID(time.Now().UnixNano())
		}
		used[id] = true
		newIDs[i] = id
	}

	fmt.Printf("Splitting trackID %s into %d new tracks:\n", *trackID, len(clusters))
	for i, c := range clusters {
		fmt.Printf("  cluster ts=%s  window=[%d,%d) -> new trackID %s\n",
			c.ts.Format(time.RFC3339), boundaries[i], boundaries[i+1], newIDs[i])
		for _, u := range c.uploads {
			fmt.Printf("      upload #%d %-30s (ts source: %s)\n", u.id, u.filename, u.tsFromComment)
		}
	}

	totalMarkers, err := countMarkers(ctx, db, *trackID)
	if err != nil {
		log.Fatalf("counting markers: %v", err)
	}
	fmt.Printf("Total markers currently under %s: %d\n", *trackID, totalMarkers)

	if !*apply {
		fmt.Println("\nDry run only. Re-run with -apply to write these changes.")
		return
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("beginning transaction: %v", err)
	}
	defer tx.Rollback()

	var reassigned int64
	for i, c := range clusters {
		res, err := tx.ExecContext(ctx,
			`UPDATE markers SET trackID = $1
			 WHERE trackID = $2 AND date >= $3 AND date < $4`,
			newIDs[i], *trackID, boundaries[i], boundaries[i+1])
		if err != nil {
			log.Fatalf("updating markers for cluster %s: %v", c.ts, err)
		}
		n, _ := res.RowsAffected()
		reassigned += n

		for _, u := range c.uploads {
			if _, err := tx.ExecContext(ctx,
				`UPDATE uploads SET track_id = $1 WHERE id = $2`,
				newIDs[i], u.id); err != nil {
				log.Fatalf("updating upload #%d: %v", u.id, err)
			}
		}
	}

	if reassigned != totalMarkers {
		log.Fatalf("safety check failed: reassigned %d markers but %s had %d — rolling back",
			reassigned, *trackID, totalMarkers)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("committing: %v", err)
	}
	fmt.Printf("Done. Reassigned %d markers across %d new trackIDs.\n", reassigned, len(clusters))
}

func groupByTimestamp(uploads []uploadRow) []cluster {
	byTS := map[int64]*cluster{}
	var order []int64
	for _, u := range uploads {
		key := u.ts.Unix()
		c, ok := byTS[key]
		if !ok {
			c = &cluster{ts: u.ts}
			byTS[key] = c
			order = append(order, key)
		}
		c.uploads = append(c.uploads, u)
	}
	out := make([]cluster, 0, len(order))
	for _, k := range order {
		out = append(out, *byTS[k])
	}
	return out
}

func loadUploads(ctx context.Context, db *sql.DB, trackID string) ([]uploadRow, error) {
	rows, err := db.QueryContext(ctx,
		`SELECT id, filename, recording_date FROM uploads WHERE track_id = $1`,
		trackID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []uploadRow
	for rows.Next() {
		var (
			u             uploadRow
			recordingDate sql.NullTime
		)
		if err := rows.Scan(&u.id, &u.filename, &recordingDate); err != nil {
			return nil, err
		}
		if ts, ok := dateFromFilename(u.filename); ok {
			u.ts = ts
			u.tsFromComment = "filename"
		} else if recordingDate.Valid {
			u.ts = recordingDate.Time
			u.tsFromComment = "recording_date (fallback)"
		} else {
			return nil, fmt.Errorf("upload #%d (%s) has no usable timestamp (filename unparseable, recording_date NULL)", u.id, u.filename)
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

func countMarkers(ctx context.Context, db *sql.DB, trackID string) (int64, error) {
	var n int64
	err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM markers WHERE trackID = $1`, trackID).Scan(&n)
	return n, err
}
