package database

import (
	"context"
	"fmt"
)

// CountryStats holds aggregated measurement statistics for a single country.
type CountryStats struct {
	Country      string  `json:"country"`
	Measurements int64   `json:"measurements"`
	AvgDoseRate  float64 `json:"avgDoseRate"`
	FirstSeen    int64   `json:"firstSeen"`
	LastSeen     int64   `json:"lastSeen"`
}

// QueryCountryStats returns per-country aggregated statistics from the markers table.
func (db *Database) QueryCountryStats(ctx context.Context) ([]CountryStats, error) {
	if db.DB == nil {
		return nil, fmt.Errorf("database unavailable")
	}

	query := `
SELECT country,
       COUNT(*) AS measurements,
       AVG(doserate) AS avg_dose_rate,
       MIN(date) AS first_seen,
       MAX(date) AS last_seen
FROM markers
WHERE country IS NOT NULL AND country != ''
GROUP BY country
ORDER BY measurements DESC`

	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query country stats: %w", err)
	}
	defer rows.Close()

	var stats []CountryStats
	for rows.Next() {
		var s CountryStats
		if err := rows.Scan(&s.Country, &s.Measurements, &s.AvgDoseRate, &s.FirstSeen, &s.LastSeen); err != nil {
			return nil, fmt.Errorf("scan country stats: %w", err)
		}
		stats = append(stats, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate country stats: %w", err)
	}
	return stats, nil
}
