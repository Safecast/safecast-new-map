package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"safecast-new-map/pkg/database"
	safecastrealtime "safecast-new-map/pkg/safecast-realtime"
)

// decodeRealtimeExtras converts the optional JSON blob with temperature or
// humidity hints into a float map.  Invalid payloads are ignored so noisy
// devices do not break chart rendering.
func decodeRealtimeExtras(raw string) map[string]float64 {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	var parsed map[string]float64
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		log.Printf("decode realtime extras: %v", err)
		return nil
	}
	clean := make(map[string]float64, len(parsed))
	for key, value := range parsed {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			continue
		}
		clean[key] = value
	}
	if len(clean) == 0 {
		return nil
	}
	return clean
}

// copyFloatMap duplicates a float map so later mutations do not affect cached
// responses.  We prefer copying over shared state to follow Go's advice of
// communicating values explicitly.
func copyFloatMap(src map[string]float64) map[string]float64 {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]float64, len(src))
	for key, value := range src {
		dst[key] = value
	}
	return dst
}

// cloneRealtimePoints allocates a fresh slice so callers can trim or
// resample without affecting other views.
func cloneRealtimePoints(points []realtimePoint) []realtimePoint {
	if len(points) == 0 {
		return nil
	}
	out := make([]realtimePoint, len(points))
	copy(out, points)
	return out
}

// filterPointsSince keeps only values with timestamps at or after the cutoff.
// The helper assumes the input slice is sorted by timestamp.
func filterPointsSince(points []realtimePoint, cutoff int64) []realtimePoint {
	if cutoff <= 0 || len(points) == 0 {
		return cloneRealtimePoints(points)
	}
	idx := sort.Search(len(points), func(i int) bool {
		return points[i].Timestamp >= cutoff
	})
	if idx >= len(points) {
		return nil
	}
	return cloneRealtimePoints(points[idx:])
}

// resampleRealtimePoints collapses samples into evenly sized buckets using the
// average value per bucket.  Returning a new slice keeps the caller's data
// immutable and mirrors the "don't mutate shared structures" guidance.
func resampleRealtimePoints(points []realtimePoint, bucket int64) []realtimePoint {
	if bucket <= 0 || len(points) == 0 {
		return cloneRealtimePoints(points)
	}
	out := make([]realtimePoint, 0, len(points))
	currentBucket := (points[0].Timestamp / bucket) * bucket
	var sum float64
	var max float64
	var count int
	for _, p := range points {
		bucketID := (p.Timestamp / bucket) * bucket
		if bucketID != currentBucket && count > 0 {
			out = append(out, realtimePoint{
				Timestamp: currentBucket + bucket/2,
				Value:     sum / float64(count),
				Max:       max,
			})
			currentBucket = bucketID
			sum = 0
			max = 0
			count = 0
		}
		sum += p.Value
		if count == 0 || p.Value > max {
			max = p.Value
		}
		count++
	}
	if count > 0 {
		out = append(out, realtimePoint{
			Timestamp: currentBucket + bucket/2,
			Value:     sum / float64(count),
			Max:       max,
		})
	}
	return out
}

// lastNRealtimePoints keeps only the newest buckets so short charts honour fixed
// divisions without redrawing hundreds of samples. Copying avoids sharing slices.
func lastNRealtimePoints(points []realtimePoint, keep int) []realtimePoint {
	if keep <= 0 || len(points) == 0 {
		return nil
	}
	if len(points) <= keep {
		return points
	}
	out := make([]realtimePoint, keep)
	copy(out, points[len(points)-keep:])
	return out
}

// alignToCeil snaps timestamps up to the next step boundary so chart grids show
// whole buckets even when the latest measurement arrives mid-interval.
func alignToCeil(t time.Time, step time.Duration) time.Time {
	if step <= 0 {
		return t.UTC()
	}
	tt := t.UTC()
	truncated := tt.Truncate(step)
	if tt.Equal(truncated) {
		return tt
	}
	return truncated.Add(step)
}

// alignMonthCeil advances to the first day of the next month so the long-term
// chart can display complete calendar months without partial buckets.
func alignMonthCeil(t time.Time) time.Time {
	tt := t.UTC()
	base := time.Date(tt.Year(), tt.Month(), 1, 0, 0, 0, 0, time.UTC)
	return base.AddDate(0, 1, 0)
}

// aggregateMonthly groups samples by calendar month using the average value per
// month. This keeps the "months" chart faithful to calendar boundaries.
func aggregateMonthly(points []realtimePoint, start time.Time, months int) []realtimePoint {
	if months <= 0 || len(points) == 0 {
		return nil
	}
	startUTC := time.Date(start.UTC().Year(), start.UTC().Month(), 1, 0, 0, 0, 0, time.UTC)
	filtered := filterPointsSince(points, startUTC.Unix())
	if len(filtered) == 0 {
		return nil
	}
	out := make([]realtimePoint, 0, months)
	idx := 0
	for i := 0; i < months; i++ {
		monthStart := startUTC.AddDate(0, i, 0)
		monthEnd := monthStart.AddDate(0, 1, 0)
		startUnix := monthStart.Unix()
		endUnix := monthEnd.Unix()
		for idx < len(filtered) && filtered[idx].Timestamp < startUnix {
			idx++
		}
		if idx >= len(filtered) {
			break
		}
		sum := 0.0
		count := 0
		j := idx
		for j < len(filtered) {
			ts := filtered[j].Timestamp
			if ts >= endUnix {
				break
			}
			sum += filtered[j].Value
			count++
			j++
		}
		if count > 0 {
			mid := startUnix + (endUnix-startUnix)/2
			out = append(out, realtimePoint{Timestamp: mid, Value: sum / float64(count)})
			idx = j
		}
	}
	return out
}

// realtimeBucketSteps enumerates pleasant aggregation steps from minutes to
// months.  The list keeps resampling predictable and friendly on charts.
var realtimeBucketSteps = []int64{
	60,
	120,
	300,
	600,
	900,
	1800,
	3600,
	7200,
	14400,
	28800,
	43200,
	86400,
	172800,
	604800,
	1209600,
	2592000,
	7776000,
}

// pickNiceBucket chooses the smallest bucket from realtimeBucketSteps that is
// greater or equal to the requested minimum.  Falling back to doubling keeps
// the function total when the data spans many years.
func pickNiceBucket(minStep int64) int64 {
	if minStep <= 1 {
		return 1
	}
	for _, step := range realtimeBucketSteps {
		if step >= minStep {
			return step
		}
	}
	return realtimeBucketSteps[len(realtimeBucketSteps)-1]
}

// nextBucket returns the next larger bucket after the provided step.
func nextBucket(step int64) int64 {
	for _, candidate := range realtimeBucketSteps {
		if candidate > step {
			return candidate
		}
	}
	if step <= 0 {
		return 1
	}
	return step * 2
}

// prepareRealtimeSeries optionally resamples a slice so the frontend receives
// at most "limit" points.  When defaultBucket is positive we always aggregate
// using that duration; otherwise the function derives a pleasant step.
func prepareRealtimeSeries(points []realtimePoint, limit int, defaultBucket int64) ([]realtimePoint, int64) {
	if len(points) == 0 {
		return nil, 0
	}
	if defaultBucket > 0 {
		bucket := defaultBucket
		resampled := resampleRealtimePoints(points, bucket)
		for len(resampled) > limit {
			bucket = nextBucket(bucket)
			resampled = resampleRealtimePoints(points, bucket)
		}
		return resampled, bucket
	}
	if len(points) <= limit {
		return cloneRealtimePoints(points), 0
	}
	span := points[len(points)-1].Timestamp - points[0].Timestamp
	if span <= 0 {
		return cloneRealtimePoints(points), 0
	}
	approx := span / int64(limit)
	if approx <= 0 {
		approx = 1
	}
	bucket := pickNiceBucket(approx)
	resampled := resampleRealtimePoints(points, bucket)
	for len(resampled) > limit {
		bucket = nextBucket(bucket)
		resampled = resampleRealtimePoints(points, bucket)
	}
	return resampled, bucket
}

// summariseRealtimeHistory processes DB rows on a background goroutine and
// returns aggregated series for day, month, and all-time windows.
func summariseRealtimeHistory(rows []database.RealtimeMeasurement, now time.Time) historyAggregate {
	input := make(chan realtimeMeasurementPayload)
	go func() {
		defer close(input)
		for _, row := range rows {
			val, ok := safecastrealtime.FromRealtime(row.Value, row.Unit)
			if !ok {
				continue
			}
			payload := realtimeMeasurementPayload{
				timestamp: row.MeasuredAt,
				radiation: safecastrealtime.ToMicroRoentgen(val),
				extras:    decodeRealtimeExtras(row.Extra),
				name:      row.DeviceName,
				transport: row.Transport,
				tube:      row.Tube,
				country:   row.Country,
			}
			input <- payload
		}
	}()

	return collectRealtimeMeasurements(input, now, len(rows))
}

// collectRealtimeMeasurements consumes the measurement stream, keeps track of
// metadata, and prepares slices for each chart timeframe.
func collectRealtimeMeasurements(input <-chan realtimeMeasurementPayload, now time.Time, expected int) historyAggregate {
	agg := historyAggregate{
		Series: map[string][]realtimePoint{
			"day":   {},
			"month": {},
			"all":   {},
		},
		ExtraSeries: map[string]map[string][]realtimePoint{
			"day":   {},
			"month": {},
			"all":   {},
		},
		Ranges: make(map[string]rangeSummary, 3),
	}

	extrasAll := make(map[string][]realtimePoint)
	if expected <= 0 {
		expected = 1024
	}
	allPoints := make([]realtimePoint, 0, expected)
	var lastTs int64

	for payload := range input {
		point := realtimePoint{Timestamp: payload.timestamp, Value: payload.radiation}
		allPoints = append(allPoints, point)
		lastTs = payload.timestamp

		if agg.DeviceName == "" && payload.name != "" {
			agg.DeviceName = payload.name
		}
		if agg.Transport == "" && payload.transport != "" {
			agg.Transport = payload.transport
		}
		if agg.Country == "" && payload.country != "" {
			agg.Country = payload.country
		}
		if agg.Tube == "" {
			if label := safecastrealtime.DetectorLabel(payload.tube, payload.transport, payload.name); label != "" {
				agg.Tube = label
			}
		}

		if len(payload.extras) > 0 {
			agg.Extra = copyFloatMap(payload.extras)
			for key, value := range payload.extras {
				extrasAll[key] = append(extrasAll[key], realtimePoint{Timestamp: payload.timestamp, Value: value})
			}
		}
	}

	reference := now.UTC()
	if lastTs > 0 {
		lastSeen := time.Unix(lastTs, 0).UTC()
		if lastSeen.After(reference) {
			reference = lastSeen
		}
	}

	const (
		hourlySegments  = 24
		dailySegments   = 24
		monthlySegments = 24
	)
	const dayDuration = 24 * time.Hour

	dayEnd := alignToCeil(reference, time.Hour)
	dayStart := dayEnd.Add(-time.Duration(hourlySegments) * time.Hour)
	monthEnd := alignToCeil(reference, dayDuration)
	monthStart := monthEnd.Add(-time.Duration(dailySegments) * dayDuration)
	allEnd := alignMonthCeil(reference)
	allStart := allEnd.AddDate(0, -monthlySegments, 0)

	dayCutoff := dayStart.Unix()
	monthCutoff := monthStart.Unix()
	allCutoff := allStart.Unix()

	hourBucket := int64(time.Hour / time.Second)
	daySeries := lastNRealtimePoints(resampleRealtimePoints(filterPointsSince(allPoints, dayCutoff), hourBucket), hourlySegments)

	dayBucketSeconds := int64(dayDuration / time.Second)
	monthSeries := lastNRealtimePoints(resampleRealtimePoints(filterPointsSince(allPoints, monthCutoff), dayBucketSeconds), dailySegments)

	allSeries := aggregateMonthly(allPoints, allStart, monthlySegments)
	avgBucket := (allEnd.Sub(allStart) / time.Duration(monthlySegments)) / time.Second
	allBucketSeconds := int64(avgBucket)
	if allBucketSeconds <= 0 {
		allBucketSeconds = int64(30 * dayDuration / time.Second)
	}

	agg.Series["day"] = daySeries
	agg.Series["month"] = monthSeries
	agg.Series["all"] = allSeries
	agg.Ranges["day"] = rangeSummary{Start: dayCutoff, End: dayEnd.Unix(), BucketSeconds: hourBucket}
	agg.Ranges["month"] = rangeSummary{Start: monthCutoff, End: monthEnd.Unix(), BucketSeconds: dayBucketSeconds}
	agg.Ranges["all"] = rangeSummary{Start: allCutoff, End: allEnd.Unix(), BucketSeconds: allBucketSeconds}

	buildExtras := func(target string, cutoff int64, bucket int64) {
		if len(extrasAll) == 0 {
			return
		}
		series := make(map[string][]realtimePoint)
		for key, pts := range extrasAll {
			filtered := filterPointsSince(pts, cutoff)
			if len(filtered) == 0 {
				continue
			}
			if bucket > 0 {
				series[key] = resampleRealtimePoints(filtered, bucket)
			} else {
				series[key] = filtered
			}
		}
		if len(series) > 0 {
			agg.ExtraSeries[target] = series
		}
	}

	buildMonthlyExtras := func(target string, start time.Time, months int) {
		if len(extrasAll) == 0 {
			return
		}
		series := make(map[string][]realtimePoint)
		for key, pts := range extrasAll {
			aggregated := aggregateMonthly(pts, start, months)
			if len(aggregated) == 0 {
				continue
			}
			series[key] = aggregated
		}
		if len(series) > 0 {
			agg.ExtraSeries[target] = series
		}
	}

	buildExtras("day", dayCutoff, hourBucket)
	buildExtras("month", monthCutoff, dayBucketSeconds)
	buildMonthlyExtras("all", allStart, monthlySegments)

	if agg.Series["day"] == nil {
		agg.Series["day"] = []realtimePoint{}
	}
	if agg.Series["month"] == nil {
		agg.Series["month"] = []realtimePoint{}
	}
	if agg.Series["all"] == nil {
		agg.Series["all"] = []realtimePoint{}
	}

	return agg
}

// realtimeHistoryHandler returns one year of realtime measurements for a device.
// The handler keeps the response lightweight so the frontend can draw Grafana-style
// charts without shipping a dedicated dashboard backend.
//
// @Summary     Get realtime device history
// @Description Returns aggregated realtime history and summary ranges for a device.
// @Tags        map
// @Produce     json
// @Param       device query string true "Realtime device ID"
// @Success     200 {object} map[string]interface{} "Realtime history payload"
// @Failure     400 {string} string "Missing device"
// @Failure     404 {string} string "Realtime feature disabled"
// @Failure     500 {string} string "Server error"
// @Router      /realtime_history [get]
func realtimeHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if !*safecastRealtimeEnabled {
		http.NotFound(w, r)
		return
	}

	device := strings.TrimSpace(r.URL.Query().Get("device"))
	if device == "" {
		http.Error(w, "missing device", http.StatusBadRequest)
		return
	}
	if strings.HasPrefix(device, "live:") {
		device = strings.TrimPrefix(device, "live:")
	}

	now := time.Now()
	rows, err := db.GetRealtimeHistory(device, 0, *dbType)
	if err != nil {
		http.Error(w, "history error", http.StatusInternalServerError)
		return
	}

	historyCh := make(chan historyAggregate, 1)
	go func() {
		historyCh <- summariseRealtimeHistory(rows, now)
	}()

	var agg historyAggregate
	select {
	case <-r.Context().Done():
		http.Error(w, "request cancelled", http.StatusRequestTimeout)
		return
	case agg = <-historyCh:
	}

	resp := struct {
		DeviceID    string                                `json:"deviceID"`
		DeviceName  string                                `json:"deviceName,omitempty"`
		Transport   string                                `json:"transport,omitempty"`
		Tube        string                                `json:"tube,omitempty"`
		Country     string                                `json:"country,omitempty"`
		Series      map[string][]realtimePoint            `json:"series"`
		Extra       map[string]float64                    `json:"extra,omitempty"`
		ExtraSeries map[string]map[string][]realtimePoint `json:"extraSeries,omitempty"`
		Ranges      map[string]rangeSummary               `json:"ranges,omitempty"`
	}{
		DeviceID:    device,
		DeviceName:  agg.DeviceName,
		Transport:   agg.Transport,
		Tube:        agg.Tube,
		Country:     agg.Country,
		Series:      agg.Series,
		Extra:       agg.Extra,
		ExtraSeries: agg.ExtraSeries,
		Ranges:      agg.Ranges,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// Gzip compression for API responses is in pkg/httpapi (Server.gzipWrap).

// =====================
// MAIN
// =====================

// main parses flags, initialises the DB & routes, then either
// (a) serves plain HTTP on a custom port, or
// (b) if -domain is given, serves ACME-backed HTTPS on 443 plus
//     an ACME/redirect helper on 80.
//
// If any web-server returns an error it is only logged – the
// application continues running.  A final `select{}` keeps the
// main goroutine alive without resorting to mutexes.

// main: парсинг флагов, инициализация БД и запуск веб-серверов.
// Добавлен withServerHeader для всех запросов.
// =====================
// MAIN
// =====================
