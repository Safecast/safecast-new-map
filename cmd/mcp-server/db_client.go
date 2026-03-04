package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDB() error {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL environment variable is required")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db = pool
	return nil
}

func dbAvailable() bool {
	return db != nil
}

// queryRows executes a query and returns results as a slice of maps.
func queryRows(ctx context.Context, query string, args ...any) ([]map[string]any, error) {
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fields := rows.FieldDescriptions()
	var results []map[string]any

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}
		row := make(map[string]any, len(fields))
		for i, fd := range fields {
			row[string(fd.Name)] = values[i]
		}
		results = append(results, row)
	}

	return results, rows.Err()
}

// queryRow executes a query and returns a single row as a map.
func queryRow(ctx context.Context, query string, args ...any) (map[string]any, error) {
	rows, err := queryRows(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no rows returned")
	}
	return rows[0], nil
}
