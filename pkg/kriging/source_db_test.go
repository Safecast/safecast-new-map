package kriging

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"safecast-new-map/pkg/database"
	_ "safecast-new-map/pkg/database/drivers"
)

func TestDBSourceUsesRawPointsOnly(t *testing.T) {
	dir := t.TempDir()
	cfg := database.Config{DBType: "sqlite", DBPath: filepath.Join(dir, "kriging.db"), Port: 1}
	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	if err := db.InitSchema(cfg); err != nil {
		t.Fatalf("InitSchema: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	const when = int64(1_700_000_000)
	inserts := []struct {
		id    int
		dose  float64
		zoom  int
		lat   float64
		when  int64
		track string
		dev   string
	}{
		{1, 0.2, 0, 35.0, when, "raw-track", "device-a"},
		{2, 0, 0, 35.0, when, "raw-track", "device-a"},
		{3, -1, 0, 35.0, when, "raw-track", "device-a"},
		{4, 99, 5, 35.0, when, "clustered", "device-a"},
		{5, 0.4, 0, 35.0, when - 60*24*3600, "old-track", "device-b"},
		{6, 0.5, 0, 50.0, when, "elsewhere", "device-c"},
	}
	for _, row := range inserts {
		_, err := db.DB.ExecContext(context.Background(),
			`INSERT INTO markers (id, doseRate, date, lon, lat, countRate, zoom, speed, trackID, device_id)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			row.id, row.dose, row.when, 139.0, row.lat, 1.0, row.zoom, 0.0, row.track, row.dev)
		if err != nil {
			t.Fatalf("insert %d: %v", row.id, err)
		}
	}

	src := NewDBSource(db)
	got, err := src.Query(context.Background(), Bounds{
		MinLat: 34, MinLon: 138, MaxLat: 36, MaxLon: 140,
	}, when-TemporalPracticalRangeSec, when+TemporalPracticalRangeSec)
	if err != nil {
		t.Fatal(err)
	}
	if got.OmittedNonPositive != 2 {
		t.Fatalf("omitted %d, want 2", got.OmittedNonPositive)
	}
	if len(got.Observations) != 1 {
		t.Fatalf("got %d observations, want the one raw positive point in the window", len(got.Observations))
	}
	obs := got.Observations[0]
	if obs.DoseRate != 0.2 || obs.TrackID != "raw-track" || obs.DeviceID != "device-a" {
		t.Fatalf("observation %+v", obs)
	}

	src.MaxObservations = 1
	_, err = src.Query(context.Background(), Bounds{
		MinLat: -90, MinLon: -180, MaxLat: 90, MaxLon: 180,
	}, 0, time.Now().Unix()+10_000)
	if err != ErrTooManyObservations {
		t.Fatalf("cap err %v", err)
	}
}
