package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "database-8765.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM spectra").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total spectra: %d\n", count)

	err = db.QueryRow("SELECT COUNT(*) FROM markers WHERE has_spectrum = 1").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Markers with has_spectrum=1: %d\n", count)
}
