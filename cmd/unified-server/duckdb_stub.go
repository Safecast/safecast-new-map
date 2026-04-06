//go:build !duckdb

package main

import (
	"database/sql"
	"log"
)

// duckDB is nil when the duckdb build tag is not provided.
var duckDB *sql.DB

// Stub for DuckDB when not enabled in build tags
func initDuckDBAnalytics() error {
	log.Println("DuckDB analytics disabled (duckdb build tag not provided)")
	return nil
}

func duckDBAvailable() bool {
	return false
}
