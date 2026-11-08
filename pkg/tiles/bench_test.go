//go:build tilebench

package tiles

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"safecast-new-map/pkg/database"
	_ "safecast-new-map/pkg/database/drivers"
)

// TestLiveScan separates database time from JSON size and PNG size.
// It only reads. Without a populated database it skips:
//
//	TILEBENCH_DB_TYPE=pgx TILEBENCH_DSN='postgres://...' go test -tags tilebench ./pkg/tiles -run TestLiveScan -v
//
// SQLite and DuckDB use TILEBENCH_DB_PATH instead of TILEBENCH_DSN.
// A cold zoom-8 city tile that stays above about a second is the signal
// to add an occupied-cell sidecar. This test records that time and does not
// add the sidecar itself.
func TestLiveScan(t *testing.T) {
	db := openBenchDB(t)
	defer db.Close()
	driver := os.Getenv("TILEBENCH_DB_TYPE")

	tokyo := Bounds{MinLat: 35.53, MinLon: 139.56, MaxLat: 35.82, MaxLon: 139.93}
	t.Log(zoomHistogram(t, db, driver, tokyo))

	reportScan(t, db, driver, "tokyo-viewport-btree-narrow", tokyo, false, false)
	reportScan(t, db, driver, "tokyo-viewport-btree-full-json", tokyo, false, true)
	if strings.EqualFold(driver, "pgx") {
		reportScan(t, db, driver, "tokyo-viewport-gist-narrow", tokyo, true, false)
	}

	agg := &Aggregate{DB: db, DBType: driver}
	reportTile(t, agg, "tokyo-z12", mustTile(t, 35.675, 139.745, 12))
	reportTile(t, agg, "tokyo-z8", mustTile(t, 35.675, 139.745, 8))
	reportTile(t, agg, "sparse-ocean-z8", mustTile(t, -40, -30, 8))

	trackID, corridor := busiestTrack(t, db, driver, tokyo)
	if trackID == "" {
		t.Log("no track in the Tokyo viewport; drive-survey arm skipped")
		return
	}
	t.Logf("busiest track %s corridor lat %.5f..%.5f lon %.5f..%.5f", trackID, corridor.MinLat, corridor.MaxLat, corridor.MinLon, corridor.MaxLon)
	reportScan(t, db, driver, "drive-track-btree-narrow", corridor, false, false)
}

func openBenchDB(t *testing.T) *database.Database {
	t.Helper()
	driver := strings.TrimSpace(os.Getenv("TILEBENCH_DB_TYPE"))
	if driver == "" {
		t.Skip("set TILEBENCH_DB_TYPE and TILEBENCH_DSN or TILEBENCH_DB_PATH to measure a real database")
	}
	db, err := database.NewDatabase(database.Config{
		DBType: driver,
		DBConn: os.Getenv("TILEBENCH_DSN"),
		DBPath: os.Getenv("TILEBENCH_DB_PATH"),
		Port:   1,
	})
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func mustTile(t *testing.T, lat, lon float64, z int) TileID {
	t.Helper()
	id, ok := TileAt(lat, lon, z)
	if !ok {
		t.Fatalf("tile at %f,%f z%d", lat, lon, z)
	}
	return id
}

func zoomHistogram(t *testing.T, db *database.Database, driver string, b Bounds) string {
	t.Helper()
	q := "SELECT zoom, COUNT(*) FROM markers WHERE lat BETWEEN ? AND ? AND lon BETWEEN ? AND ? GROUP BY zoom ORDER BY zoom"
	args := []any{b.MinLat, b.MaxLat, b.MinLon, b.MaxLon}
	if strings.EqualFold(driver, "pgx") {
		q = "SELECT zoom, COUNT(*) FROM markers WHERE lat BETWEEN $1 AND $2 AND lon BETWEEN $3 AND $4 GROUP BY zoom ORDER BY zoom"
	}
	var lines []string
	err := db.WithWebRead(context.Background(), func(ctx context.Context, conn *sql.DB) error {
		rows, err := conn.QueryContext(ctx, q, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var zoom, n int
			if err := rows.Scan(&zoom, &n); err != nil {
				return err
			}
			lines = append(lines, fmt.Sprintf("zoom %d: %d rows", zoom, n))
		}
		return rows.Err()
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) == 0 {
		return "zoom histogram: no rows"
	}
	return "zoom histogram:\n" + strings.Join(lines, "\n")
}

func reportScan(t *testing.T, db *database.Database, driver, name string, b Bounds, gist, full bool) {
	t.Helper()
	q, args := scanQuery(driver, b, gist, full)
	start := time.Now()
	var first time.Duration
	var rowsN int
	var jsonBytes int
	err := db.WithWebRead(context.Background(), func(ctx context.Context, conn *sql.DB) error {
		rows, err := conn.QueryContext(ctx, q, args...)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			if rowsN == 0 {
				first = time.Since(start)
			}
			if full {
				var rec fullRow
				var has any
				if err := rows.Scan(&rec.ID, &rec.Dose, &rec.Date, &rec.Lon, &rec.Lat, &rec.Count, &rec.Zoom, &rec.Speed, &rec.TrackID, &rec.Detector, &has); err != nil {
					return err
				}
				encoded, err := json.Marshal(rec)
				if err != nil {
					return err
				}
				jsonBytes += len(encoded) + 1
			} else {
				var p Point
				var has any
				if err := rows.Scan(&p.ID, &p.Lat, &p.Lon, &p.Dose, &p.Date, &p.TrackID, &has); err != nil {
					return err
				}
			}
			rowsN++
		}
		return rows.Err()
	})
	if err != nil {
		t.Fatalf("%s: %v", name, err)
	}
	t.Logf("%s rows=%d scan=%s first=%s json=%dB", name, rowsN, time.Since(start).Round(time.Millisecond), first.Round(time.Millisecond), jsonBytes)
}

type fullRow struct {
	ID       int64   `json:"id"`
	Dose     float64 `json:"doseRate"`
	Date     int64   `json:"date"`
	Lon      float64 `json:"lon"`
	Lat      float64 `json:"lat"`
	Count    float64 `json:"countRate"`
	Zoom     int     `json:"zoom"`
	Speed    float64 `json:"speed"`
	TrackID  string  `json:"trackID"`
	Detector string  `json:"detector"`
}

func scanQuery(driver string, b Bounds, gist, full bool) (string, []any) {
	pg := strings.EqualFold(driver, "pgx")
	cols := "id, lat, lon, doseRate, date, COALESCE(trackID, ''), COALESCE(has_spectrum, 0)"
	if full && pg {
		cols = "id, doserate, date, lon, lat, countrate, zoom, COALESCE(speed, 0), trackid, COALESCE(detector, ''), COALESCE(has_spectrum, FALSE)"
	} else if full {
		cols = "id, doseRate, date, lon, lat, countRate, zoom, COALESCE(speed, 0), trackID, COALESCE(detector, ''), COALESCE(has_spectrum, 0)"
	} else if pg {
		cols = "id, lat, lon, doserate, date, COALESCE(trackid, ''), COALESCE(has_spectrum, FALSE)"
	}
	if gist {
		return `SELECT ` + cols + ` FROM markers WHERE zoom = $1 AND geom && ST_MakeEnvelope($2, $3, $4, $5, 4326) AND ST_Intersects(geom, ST_MakeEnvelope($2, $3, $4, $5, 4326))`,
			[]any{0, b.MinLon, b.MinLat, b.MaxLon, b.MaxLat}
	}
	if pg {
		return `SELECT ` + cols + ` FROM markers WHERE zoom = $1 AND lat BETWEEN $2 AND $3 AND lon BETWEEN $4 AND $5`,
			[]any{0, b.MinLat, b.MaxLat, b.MinLon, b.MaxLon}
	}
	return `SELECT ` + cols + ` FROM markers WHERE zoom = ? AND lat BETWEEN ? AND ? AND lon BETWEEN ? AND ?`,
		[]any{0, b.MinLat, b.MaxLat, b.MinLon, b.MaxLon}
}

func reportTile(t *testing.T, agg *Aggregate, name string, id TileID) {
	t.Helper()
	bounds, ok := TileBounds(id)
	if !ok {
		t.Fatalf("%s bounds", name)
	}
	start := time.Now()
	points, err := agg.Load(context.Background(), bounds, 0, 0, "all")
	if err != nil {
		t.Fatalf("%s load: %v", name, err)
	}
	loaded := time.Since(start)
	binStart := time.Now()
	grid := Bin(points, id)
	binned := time.Since(binStart)
	body, err := EncodePNG(grid)
	if err != nil {
		t.Fatal(err)
	}
	occupied := 0
	for _, c := range grid.Cells {
		if c.Occupied {
			occupied++
		}
	}
	t.Logf("%s points=%d occupied=%d load=%s bin=%s png=%dB", name, len(points), occupied, loaded.Round(time.Millisecond), binned.Round(time.Millisecond), len(body))
	if id.Z <= 8 && loaded > time.Second {
		t.Logf("%s cold scan exceeded 1s; an occupied-cell sidecar would be justified for this zoom", name)
	}
}

func busiestTrack(t *testing.T, db *database.Database, driver string, b Bounds) (string, Bounds) {
	t.Helper()
	q := `SELECT trackID, COUNT(*) FROM markers WHERE zoom = 0 AND lat BETWEEN ? AND ? AND lon BETWEEN ? AND ? AND trackID <> '' GROUP BY trackID ORDER BY COUNT(*) DESC LIMIT 1`
	args := []any{b.MinLat, b.MaxLat, b.MinLon, b.MaxLon}
	if strings.EqualFold(driver, "pgx") {
		q = `SELECT trackid, COUNT(*) FROM markers WHERE zoom = 0 AND lat BETWEEN $1 AND $2 AND lon BETWEEN $3 AND $4 AND trackid <> '' GROUP BY trackid ORDER BY COUNT(*) DESC LIMIT 1`
	}
	var track string
	var n int
	err := db.WithWebRead(context.Background(), func(ctx context.Context, conn *sql.DB) error {
		return conn.QueryRowContext(ctx, q, args...).Scan(&track, &n)
	})
	if err == sql.ErrNoRows || track == "" {
		return "", Bounds{}
	}
	if err != nil {
		t.Fatal(err)
	}
	rangeQ := `SELECT MIN(lat), MAX(lat), MIN(lon), MAX(lon) FROM markers WHERE zoom = 0 AND trackID = ?`
	rangeArgs := []any{track}
	if strings.EqualFold(driver, "pgx") {
		rangeQ = `SELECT MIN(lat), MAX(lat), MIN(lon), MAX(lon) FROM markers WHERE zoom = 0 AND trackid = $1`
	}
	var out Bounds
	err = db.WithWebRead(context.Background(), func(ctx context.Context, conn *sql.DB) error {
		return conn.QueryRowContext(ctx, rangeQ, rangeArgs...).Scan(&out.MinLat, &out.MaxLat, &out.MinLon, &out.MaxLon)
	})
	if err != nil {
		t.Fatal(err)
	}
	return track, out
}
