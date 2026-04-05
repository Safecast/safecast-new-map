package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	h3 "github.com/uber/h3-go/v4"
)

// H3Cell is the aggregated result for a single H3 hexagon cell.
type H3Cell struct {
	CellID  string  `json:"cell"`
	Lat     float64 `json:"lat"`      // center latitude
	Lon     float64 `json:"lon"`      // center longitude
	AvgDose float64 `json:"avg_dose"` // µSv/h average
	MaxDose float64 `json:"max_dose"` // µSv/h peak
	Count   int     `json:"count"`    // number of measurements
}

// zoomToH3Resolution maps a Leaflet zoom level to a sensible H3 resolution.
// Target: 200–2000 visible hex cells per viewport.
func zoomToH3Resolution(zoom int) int {
	// res 3 = ~33km, res 4 = ~12km, res 5 = ~4.5km, res 6 = ~1.7km,
	// res 7 = ~0.6km, res 8 = ~200m, res 9 = ~70m, res 10 = ~25m, res 11 = ~9m
	switch {
	case zoom <= 3:
		return 3
	case zoom <= 5:
		return 4
	case zoom <= 7:
		return 5
	case zoom <= 9:
		return 6
	case zoom <= 11:
		return 7
	case zoom <= 13:
		return 8
	case zoom <= 15:
		return 9
	case zoom <= 17:
		return 10
	default:
		return 11
	}
}

// h3ResColumn maps an H3 resolution to the pre-computed column name, or "" if none.
// Resolutions 5, 7, 9 have dedicated columns that are pre-filled by h3_backfill.go.
func h3ResColumn(res int) string {
	switch res {
	case 5:
		return "h3_res5"
	case 7:
		return "h3_res7"
	case 9:
		return "h3_res9"
	default:
		return ""
	}
}

// handleH3Grid handles GET /api/h3grid
//
// @Summary     Aggregate radiation measurements into H3 hexagonal grid cells
// @Description Returns radiation data aggregated into H3 hex cells for the given viewport.
//
//	Resolution is chosen automatically based on zoom level (zoom→H3 res: 0-3→3, 4-5→4, …, 18+→11).
//	For resolutions 5, 7, 9 the query uses pre-computed h3_res* columns (SQL GROUP BY — fast).
//	For all other resolutions it falls back to Go-side aggregation over a streamed row set.
//	Each cell returns the average dose rate, peak dose rate, and measurement count.
//
// @Tags        historical
// @Produce     json
// @Param       zoom   query  integer true  "Leaflet zoom level (0–18)"
// @Param       minLat query  number  true  "Bounding box min latitude"
// @Param       minLon query  number  true  "Bounding box min longitude"
// @Param       maxLat query  number  true  "Bounding box max latitude"
// @Param       maxLon query  number  true  "Bounding box max longitude"
// @Success     200 {array}  H3Cell "Array of H3 cells with aggregated dose rates"
// @Failure     400 {object} map[string]string "Invalid parameters"
// @Failure     503 {object} map[string]string "Database unavailable"
// @Router      /h3grid [get]
func (h *RESTHandler) handleH3Grid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if !dbAvailable() {
		writeError(w, http.StatusServiceUnavailable, "database not available")
		return
	}

	q := r.URL.Query()
	zoom, err := strconv.Atoi(q.Get("zoom"))
	if err != nil || zoom < 0 || zoom > 18 {
		writeError(w, http.StatusBadRequest, "zoom must be 0–18")
		return
	}
	minLat, err := strconv.ParseFloat(q.Get("minLat"), 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid minLat")
		return
	}
	minLon, err := strconv.ParseFloat(q.Get("minLon"), 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid minLon")
		return
	}
	maxLat, err := strconv.ParseFloat(q.Get("maxLat"), 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid maxLat")
		return
	}
	maxLon, err := strconv.ParseFloat(q.Get("maxLon"), 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid maxLon")
		return
	}

	res := zoomToH3Resolution(zoom)
	col := h3ResColumn(res)

	var result []H3Cell

	if col != "" {
		// Fast path: SQL GROUP BY on pre-computed column.
		result, err = queryH3ByColumn(r.Context(), col, minLat, minLon, maxLat, maxLon)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		// Fallback: stream rows and aggregate in Go.
		result, err = queryH3ByStreaming(r, res, minLat, minLon, maxLat, maxLon)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=30")
	json.NewEncoder(w).Encode(result) //nolint:errcheck
}

// queryH3ByColumn uses a SQL GROUP BY on a pre-computed h3_res* column.
// This is dramatically faster than streaming — the DB does all the work.
func queryH3ByColumn(ctx context.Context, col string, minLat, minLon, maxLat, maxLon float64) ([]H3Cell, error) {
	// col is always one of h3_res5 / h3_res7 / h3_res9 — no user input, safe to interpolate.
	query := `SELECT ` + col + `, AVG(doserate), MAX(doserate), COUNT(*)
	          FROM markers
	          WHERE latitude  BETWEEN $1 AND $2
	            AND longitude BETWEEN $3 AND $4
	            AND doserate  > 0
	            AND (speed IS NULL OR speed = 0 OR speed < 45)
	            AND ` + col + ` IS NOT NULL
	          GROUP BY ` + col + `
	          LIMIT 20000`

	rows, err := db.DB.QueryContext(ctx, query, minLat, maxLat, minLon, maxLon)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []H3Cell
	for rows.Next() {
		var cellStr string
		var avg, maxD float64
		var count int
		if err := rows.Scan(&cellStr, &avg, &maxD, &count); err != nil {
			return nil, err
		}
		cell := h3.CellFromString(cellStr)
		center, _ := h3.CellToLatLng(cell)
		result = append(result, H3Cell{
			CellID:  cellStr,
			Lat:     center.Lat,
			Lon:     center.Lng,
			AvgDose: avg,
			MaxDose: maxD,
			Count:   count,
		})
	}
	return result, rows.Err()
}

// queryH3ByStreaming is the fallback used for resolutions without a pre-computed column.
// It streams up to 200 000 markers and aggregates in Go (same logic as the original handler).
func queryH3ByStreaming(r *http.Request, res int, minLat, minLon, maxLat, maxLon float64) ([]H3Cell, error) {
	const maxMarkers = 200_000
	markers, errCh := db.StreamMarkersForH3(r.Context(), 0, minLat, minLon, maxLat, maxLon, *dbType, maxMarkers)

	type cellAcc struct {
		sumDose float64
		maxDose float64
		count   int
	}
	cells := make(map[h3.Cell]*cellAcc)

	for m := range markers {
		if m.DoseRate <= 0 {
			continue
		}
		cell, _ := h3.LatLngToCell(h3.NewLatLng(m.Lat, m.Lon), res)
		acc := cells[cell]
		if acc == nil {
			acc = &cellAcc{}
			cells[cell] = acc
		}
		acc.sumDose += m.DoseRate
		acc.count++
		if m.DoseRate > acc.maxDose {
			acc.maxDose = m.DoseRate
		}
	}

	if e := <-errCh; e != nil && r.Context().Err() == nil {
		return nil, e
	}

	result := make([]H3Cell, 0, len(cells))
	for cell, acc := range cells {
		center, _ := h3.CellToLatLng(cell)
		result = append(result, H3Cell{
			CellID:  cell.String(),
			Lat:     center.Lat,
			Lon:     center.Lng,
			AvgDose: acc.sumDose / float64(acc.count),
			MaxDose: acc.maxDose,
			Count:   acc.count,
		})
	}
	return result, nil
}
