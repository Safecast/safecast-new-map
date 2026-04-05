//go:build !duckdb

package main

import (
	"log"
)

// Stub for DuckDB when not enabled in build tags
func initDuckDBAnalytics() error {
	log.Println("DuckDB analytics disabled (duckdb build tag not provided)")
	return nil
}

func duckDBAvailable() bool {
	return false
}
