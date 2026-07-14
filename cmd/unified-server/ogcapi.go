package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// OGC API - Features (a.k.a. WFS 3.0) endpoint for QGIS / GIS clients.
//
// Rooted at /ogc. Minimal conformant set QGIS needs to auto-discover layers:
//
//	GET /ogc                              landing page
//	GET /ogc/conformance                  conformance declaration
//	GET /ogc/collections                  list of feature collections
//	GET /ogc/collections/{id}             collection metadata
//	GET /ogc/collections/{id}/items       GeoJSON FeatureCollection (bbox,limit,offset)
//
// In QGIS: Layer → Add Layer → WFS / OGC API - Features, and use the landing
// URL https://simplemap.safecast.org/ogc as the connection URL.
//
// v1 exposes one collection, "sensors" (live fixed-sensor network). The
// historical bGeigie "measurements" collection is a planned v2 — add a case to
// ogcCollections() and ogcItems() when its paged query lands.

const (
	ogcBase          = "/ogc"
	ogcDefaultLimit  = 1000
	ogcMaxLimit      = 10000
	ogcSensorsCollID = "sensors"
)

// ogcBaseURL reconstructs the externally visible base URL (scheme + host),
// honouring the reverse proxy's X-Forwarded-Proto so links are correct behind
// CloudFront / nginx.
func ogcBaseURL(r *http.Request) string {
	scheme := "http"
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	} else if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}

func ogcWriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type ogcLink struct {
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Type  string `json:"type,omitempty"`
	Title string `json:"title,omitempty"`
}

// handleOGCAPI dispatches all /ogc* requests.
func (h *RESTHandler) handleOGCAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ogcWriteJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, ogcBase)
	path = strings.TrimSuffix(path, "/")

	switch {
	case path == "" || path == "/":
		h.ogcLanding(w, r)
	case path == "/api":
		h.ogcAPIDoc(w, r)
	case path == "/conformance":
		h.ogcConformance(w, r)
	case path == "/collections":
		h.ogcCollections(w, r)
	case strings.HasPrefix(path, "/collections/"):
		rest := strings.TrimPrefix(path, "/collections/")
		if id, ok := strings.CutSuffix(rest, "/items"); ok {
			h.ogcItems(w, r, id)
			return
		}
		h.ogcCollection(w, r, rest)
	default:
		ogcWriteJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
	}
}

func (h *RESTHandler) ogcLanding(w http.ResponseWriter, r *http.Request) {
	base := ogcBaseURL(r) + ogcBase
	ogcWriteJSON(w, http.StatusOK, map[string]any{
		"title":       "Safecast OGC API - Features",
		"description": "Safecast radiation data as OGC API - Features (WFS) collections for GIS clients.",
		"links": []ogcLink{
			{Href: base, Rel: "self", Type: "application/json", Title: "This document"},
			{Href: base + "/api", Rel: "service-desc", Type: "application/vnd.oai.openapi+json;version=3.0", Title: "API definition"},
			{Href: base + "/conformance", Rel: "conformance", Type: "application/json", Title: "Conformance"},
			{Href: base + "/collections", Rel: "data", Type: "application/json", Title: "Feature collections"},
		},
	})
}

// ogcAPIDoc serves a minimal OpenAPI 3.0 definition. QGIS requires the
// service-desc link and reads the items `limit` schema to size its paging.
func (h *RESTHandler) ogcAPIDoc(w http.ResponseWriter, r *http.Request) {
	base := ogcBaseURL(r) + ogcBase
	ogcWriteJSON(w, http.StatusOK, map[string]any{
		"openapi": "3.0.2",
		"info": map[string]any{
			"title":       "Safecast OGC API - Features",
			"description": "Safecast radiation data as OGC API - Features collections.",
			"version":     "1.0.0",
		},
		"servers": []map[string]any{{"url": base}},
		"paths": map[string]any{
			"/collections/{collectionId}/items": map[string]any{
				"get": map[string]any{
					"summary":     "Retrieve features of a collection",
					"operationId": "getFeatures",
					"parameters": []map[string]any{
						{
							"name": "collectionId", "in": "path", "required": true,
							"schema": map[string]any{"type": "string"},
						},
						{
							"name": "bbox", "in": "query", "required": false,
							"schema": map[string]any{
								"type": "array", "minItems": 4, "maxItems": 4,
								"items": map[string]any{"type": "number"},
							},
						},
						{
							"name": "limit", "in": "query", "required": false,
							"schema": map[string]any{
								"type": "integer", "minimum": 1,
								"maximum": ogcMaxLimit, "default": ogcDefaultLimit,
							},
						},
						{
							"name": "offset", "in": "query", "required": false,
							"schema": map[string]any{"type": "integer", "minimum": 0, "default": 0},
						},
					},
					"responses": map[string]any{
						"200": map[string]any{"description": "GeoJSON FeatureCollection"},
					},
				},
			},
		},
	})
}

func (h *RESTHandler) ogcConformance(w http.ResponseWriter, r *http.Request) {
	ogcWriteJSON(w, http.StatusOK, map[string]any{
		"conformsTo": []string{
			"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/core",
			"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/oas30",
			"http://www.opengis.net/spec/ogcapi-features-1/1.0/conf/geojson",
		},
	})
}

// ogcCollectionMeta builds the metadata object for a single collection.
func ogcCollectionMeta(base, id string) map[string]any {
	titles := map[string]string{
		ogcSensorsCollID: "Safecast fixed sensors (live)",
	}
	descs := map[string]string{
		ogcSensorsCollID: "Most recent reading from each fixed real-time sensor (bGeigieZen, Pointcast, Solarcast, nGeigie, Notehub, …).",
	}
	return map[string]any{
		"id":          id,
		"title":       titles[id],
		"description": descs[id],
		"itemType":    "feature",
		"crs":         []string{"http://www.opengis.net/def/crs/OGC/1.3/CRS84"},
		"extent": map[string]any{
			"spatial": map[string]any{
				"bbox": [][]float64{{-180, -90, 180, 90}},
				"crs":  "http://www.opengis.net/def/crs/OGC/1.3/CRS84",
			},
		},
		"links": []ogcLink{
			{Href: base + "/collections/" + id, Rel: "self", Type: "application/json"},
			{Href: base + "/collections/" + id + "/items", Rel: "items", Type: "application/geo+json", Title: titles[id]},
		},
	}
}

func ogcKnownCollection(id string) bool {
	return id == ogcSensorsCollID
}

func (h *RESTHandler) ogcCollections(w http.ResponseWriter, r *http.Request) {
	base := ogcBaseURL(r) + ogcBase
	ogcWriteJSON(w, http.StatusOK, map[string]any{
		"links": []ogcLink{
			{Href: base + "/collections", Rel: "self", Type: "application/json"},
		},
		"collections": []map[string]any{
			ogcCollectionMeta(base, ogcSensorsCollID),
		},
	})
}

func (h *RESTHandler) ogcCollection(w http.ResponseWriter, r *http.Request, id string) {
	if !ogcKnownCollection(id) {
		ogcWriteJSON(w, http.StatusNotFound, map[string]string{"error": "unknown collection: " + id})
		return
	}
	ogcWriteJSON(w, http.StatusOK, ogcCollectionMeta(ogcBaseURL(r)+ogcBase, id))
}

// ogcItems serves a GeoJSON FeatureCollection for the given collection.
func (h *RESTHandler) ogcItems(w http.ResponseWriter, r *http.Request, id string) {
	if !ogcKnownCollection(id) {
		ogcWriteJSON(w, http.StatusNotFound, map[string]string{"error": "unknown collection: " + id})
		return
	}
	if !dbAvailable() {
		ogcWriteJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "database unavailable"})
		return
	}

	q := r.URL.Query()

	// bbox=minLon,minLat,maxLon,maxLat (CRS84 / lon,lat order — OGC standard).
	minLon, minLat, maxLon, maxLat := -180.0, -90.0, 180.0, 90.0
	if bs := q.Get("bbox"); bs != "" {
		parts := strings.Split(bs, ",")
		if len(parts) != 4 {
			ogcWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "bbox must be minLon,minLat,maxLon,maxLat"})
			return
		}
		var err error
		if minLon, err = strconv.ParseFloat(strings.TrimSpace(parts[0]), 64); err != nil {
			ogcWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid bbox"})
			return
		}
		if minLat, err = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err != nil {
			ogcWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid bbox"})
			return
		}
		if maxLon, err = strconv.ParseFloat(strings.TrimSpace(parts[2]), 64); err != nil {
			ogcWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid bbox"})
			return
		}
		if maxLat, err = strconv.ParseFloat(strings.TrimSpace(parts[3]), 64); err != nil {
			ogcWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid bbox"})
			return
		}
	}

	limit := ogcDefaultLimit
	if s := q.Get("limit"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 1 {
			ogcWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid limit"})
			return
		}
		limit = v
	}
	if limit > ogcMaxLimit {
		limit = ogcMaxLimit
	}
	offset := 0
	if s := q.Get("offset"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil || v < 0 {
			ogcWriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid offset"})
			return
		}
		offset = v
	}

	features, matched, err := ogcSensorFeatures(r.Context(), minLat, maxLat, minLon, maxLon, limit, offset)
	if err != nil {
		ogcWriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	base := ogcBaseURL(r) + ogcBase
	links := []ogcLink{
		{Href: base + "/collections/" + id + "/items", Rel: "self", Type: "application/geo+json"},
		{Href: base + "/collections/" + id, Rel: "collection", Type: "application/json"},
	}
	// Advertise a next page when more rows remain (paging QGIS follows).
	if offset+len(features) < matched {
		next := fmt.Sprintf("%s/collections/%s/items?limit=%d&offset=%d", base, id, limit, offset+limit)
		if bs := q.Get("bbox"); bs != "" {
			next += "&bbox=" + bs
		}
		links = append(links, ogcLink{Href: next, Rel: "next", Type: "application/geo+json"})
	}

	w.Header().Set("Content-Type", "application/geo+json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"type":           "FeatureCollection",
		"features":       features,
		"numberMatched":  matched,
		"numberReturned": len(features),
		"links":          links,
	})
}

// ogcSensorFeatures returns the latest reading per fixed sensor inside the bbox
// as GeoJSON features, plus the total match count (independent of limit/offset).
func ogcSensorFeatures(ctx context.Context, minLat, maxLat, minLon, maxLon float64, limit, offset int) ([]map[string]any, int, error) {
	realtimeTable, _, err := findRealtimeTable(ctx)
	if err != nil {
		return nil, 0, err
	}
	if realtimeTable == "" {
		return nil, 0, fmt.Errorf("real-time sensor table not found")
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT device_id) AS total
		FROM %s
		WHERE lat >= $1 AND lat <= $2 AND lon >= $3 AND lon <= $4`, realtimeTable)
	matched := 0
	if row, err := queryRow(ctx, countQuery, minLat, maxLat, minLon, maxLon); err == nil {
		switch v := row["total"].(type) {
		case int64:
			matched = int(v)
		case int32:
			matched = int(v)
		case int:
			matched = v
		}
	}

	query := fmt.Sprintf(`
		SELECT
			rm.device_id,
			COALESCE(rm.device_name, rm.device_id) AS device_name,
			COALESCE(rm.transport, '') AS transport,
			rm.lat AS latitude,
			rm.lon AS longitude,
			rm.value AS value,
			COALESCE(rm.unit, '') AS unit,
			to_timestamp(rm.measured_at) AS last_reading_at
		FROM %s rm
		INNER JOIN (
			SELECT device_id, MAX(measured_at) AS max_measured_at
			FROM %s
			WHERE lat >= $1 AND lat <= $2 AND lon >= $3 AND lon <= $4
			GROUP BY device_id
		) latest ON rm.device_id = latest.device_id AND rm.measured_at = latest.max_measured_at
		WHERE rm.lat >= $1 AND rm.lat <= $2 AND rm.lon >= $3 AND rm.lon <= $4
		ORDER BY rm.measured_at DESC
		LIMIT $5 OFFSET $6`, realtimeTable, realtimeTable)

	rows, err := queryRows(ctx, query, minLat, maxLat, minLon, maxLon, limit, offset)
	if err != nil {
		return nil, matched, fmt.Errorf("error querying %s: %w", realtimeTable, err)
	}

	features := make([]map[string]any, 0, len(rows))
	for _, r := range rows {
		lat, latOK := toFloat(r["latitude"])
		lon, lonOK := toFloat(r["longitude"])
		if !latOK || !lonOK {
			continue
		}
		features = append(features, map[string]any{
			"type":     "Feature",
			"id":       fmt.Sprintf("%v", r["device_id"]),
			"geometry": map[string]any{"type": "Point", "coordinates": []float64{lon, lat}},
			"properties": map[string]any{
				"device_id":       r["device_id"],
				"device_name":     r["device_name"],
				"type":            r["transport"],
				"value":           r["value"],
				"unit":            r["unit"],
				"last_reading_at": r["last_reading_at"],
			},
		})
	}
	return features, matched, nil
}
