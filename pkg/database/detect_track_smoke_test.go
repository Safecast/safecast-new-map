package database

import "testing"

// Regression test for the arg/placeholder count bug that shipped alongside
// the DetectExistingTrackID distinct-point fix: every upload started failing
// with "expected 128 arguments, got 256" because the CASE column reused the
// same $N placeholders as the WHERE clause but still appended a second,
// duplicate set of args for them.
func TestDetectExistingTrackIDArgCount(t *testing.T) {
	db, err := NewDatabase(Config{DBType: "pgx", DBConn: "postgres://postgres@127.0.0.1:5432/safecast?sslmode=disable"})
	if err != nil {
		t.Skipf("local postgres not available: %v", err)
	}
	defer db.DB.Close()

	markers := make([]Marker, 0, 40)
	for i := 0; i < 40; i++ {
		markers = append(markers, Marker{
			Lat: float64(i) + 0.123456, Lon: float64(i) + 0.654321, Date: int64(1700000000 + i), DoseRate: float64(i) + 0.1,
		})
	}

	if _, err := db.DetectExistingTrackID(markers, 10, "pgx"); err != nil {
		t.Fatalf("DetectExistingTrackID: %v", err)
	}
}
