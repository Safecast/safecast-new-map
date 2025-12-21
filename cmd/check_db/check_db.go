package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "safecast-new-map/pkg/database/drivers"
)

func main() {
	db, err := sql.Open("sqlite", "database-8765.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var count int
	
	// Total spectra
	err = db.QueryRow("SELECT COUNT(*) FROM spectra").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total spectra records: %d\n", count)

	// Unique marker IDs in spectra table
	err = db.QueryRow("SELECT COUNT(DISTINCT marker_id) FROM spectra").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Unique marker_ids in spectra: %d\n", count)

	// Unique lat/lon/date combinations in spectra
	err = db.QueryRow(`
		SELECT COUNT(DISTINCT m.lat || ',' || m.lon || ',' || m.date)
		FROM spectra s
		INNER JOIN markers m ON s.marker_id = m.id
	`).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Unique locations with spectra: %d\n", count)

	// Current markers with flag set
	err = db.QueryRow("SELECT COUNT(*) FROM markers WHERE has_spectrum = 1").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Markers with has_spectrum=1: %d\n", count)

	// Show sample of spectra table
	fmt.Println("\nSample spectra records:")
	rows, err := db.Query("SELECT id, marker_id, filename FROM spectra LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, markerID int64
		var filename string
		rows.Scan(&id, &markerID, &filename)
		fmt.Printf("  Spectrum ID: %d, Marker ID: %d, File: %s\n", id, markerID, filename)
	}
}
