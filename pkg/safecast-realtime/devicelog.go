// Copyright 2026 Safecast.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Device-log backfill: ttserve's /devices feed only exposes each device's
// latest value (the last reading of each Notehub upload batch), so transient
// peaks never reach us. ttserve also serves the *full* per-device monthly log
// at /device-log/<YYYY-MM>$<urn-with-dashes>.json — every ~15-minute reading,
// peaks included. This file pulls and parses those logs so the history table
// reflects real peaks rather than sparse snapshots.
package safecastrealtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"safecast-new-map/pkg/countryresolver"
	"safecast-new-map/pkg/database"
)

// deviceLogBaseURL is the ttserve endpoint that serves full per-device history.
const deviceLogBaseURL = "https://tt.safecast.org/device-log/"

// deviceLogClient bounds device-log downloads so a slow-but-alive ttserve can't
// hang a request goroutine indefinitely.
var deviceLogClient = &http.Client{Timeout: 30 * time.Second}

// deviceLogFilename builds a monthly log filename for a device URN, mirroring
// ttserve's DeviceUIDFilename (colons and dots become dashes).
// e.g. ("note:dev:867648049106527", June 2026) -> "2026-06$note-dev-867648049106527.json"
func deviceLogFilename(urn string, month time.Time) string {
	clean := strings.NewReplacer(":", "-", ".", "-").Replace(urn)
	return fmt.Sprintf("%s$%s.json", month.UTC().Format("2006-01"), clean)
}

// FetchDeviceLog downloads and parses one device's monthly history log from
// ttserve, returning radiation measurements ready to upsert. A missing month
// (HTTP 404 / "no such file") is not an error — it returns an empty slice.
func FetchDeviceLog(ctx context.Context, urn string, month time.Time) ([]database.RealtimeMeasurement, error) {
	url := deviceLogBaseURL + deviceLogFilename(urn, month)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := deviceLogClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	// The log is a stream of devicePayload-compatible JSON objects separated by
	// "}\n,\n{" with a trailing comma. Normalise it into a JSON array so the
	// existing devicePayload decoder handles each record.
	trimmed := bytes.TrimRight(bytes.TrimSpace(body), ",\n\r\t ")
	if len(trimmed) == 0 {
		return nil, nil
	}
	if trimmed[0] != '{' {
		// Not a log file (e.g. an error string like "no such file or directory").
		return nil, nil
	}
	wrapped := make([]byte, 0, len(trimmed)+2)
	wrapped = append(wrapped, '[')
	wrapped = append(wrapped, trimmed...)
	wrapped = append(wrapped, ']')

	var payloads []devicePayload
	if err := json.Unmarshal(wrapped, &payloads); err != nil {
		return nil, fmt.Errorf("parse device log %s: %w", url, err)
	}

	fetchedAt := time.Now().Unix()
	out := make([]database.RealtimeMeasurement, 0, len(payloads))
	for _, d := range payloads {
		if d.ID == "" || d.Time == 0 {
			continue
		}
		if _, ok := convertIfRadiation(d); !ok {
			continue
		}
		if d.Lat == 0 && d.Lon == 0 {
			continue
		}

		country, _ := countryresolver.Resolve(d.Lat, d.Lon)
		if country == "" {
			country = strings.ToUpper(strings.TrimSpace(d.Country))
		}

		out = append(out, database.RealtimeMeasurement{
			DeviceID:   d.ID,
			Transport:  d.Type,
			DeviceName: d.Name,
			Tube:       DetectorLabel(d.Tube, d.Type, d.Name),
			Country:    country,
			Value:      d.Value,
			Unit:       d.Unit,
			Lat:        d.Lat,
			Lon:        d.Lon,
			MeasuredAt: d.Time,
			FetchedAt:  fetchedAt,
			Extra:      encodeMetrics(d.Metrics),
		})
	}
	return out, nil
}

// BackfillDeviceHistory fetches the current month's device log (and the
// previous month when we're early enough that it still matters for a 30-day
// window) and upserts every reading. Duplicates are skipped by the insert's
// unique constraint, so repeated calls are cheap and idempotent. It returns the
// number of rows stored.
func BackfillDeviceHistory(ctx context.Context, db *database.Database, dbType, urn string, now time.Time) (int, error) {
	months := []time.Time{now.UTC()}
	// Include the previous month so the 30-day chart is complete near month start.
	if now.UTC().Day() <= 30 {
		months = append(months, now.UTC().AddDate(0, -1, 0))
	}

	var all []database.RealtimeMeasurement
	for _, m := range months {
		measurements, err := FetchDeviceLog(ctx, urn, m)
		if err != nil {
			return 0, err
		}
		all = append(all, measurements...)
	}
	if err := db.InsertRealtimeMeasurements(all, dbType); err != nil {
		return 0, err
	}
	return len(all), nil
}
