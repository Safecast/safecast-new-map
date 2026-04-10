package main

import (
	"math"
	"sort"
	"time"

	"safecast-new-map/pkg/database"
)

// ==========
// Константы для слияния маркеров
// ==========
const (
	markerRadiusPx = 10.0       // радиус кружка в пикселях
	minValidTS     = 1262304000 // 2010-01-01 00:00:00 UTC
)

// microRoentgenPerMicroSievert keeps conversion logic explicit so both the API
// exporter and the JSON importer agree on the units we advertise publicly.
const microRoentgenPerMicroSievert = 100.0

type SpeedRange struct{ Min, Max float64 }

// convertRhToSv и convertSvToRh - вспомогательные функции перевода
func convertRhToSv(markers []database.Marker) []database.Marker {
	filteredMarkers := []database.Marker{}
	const conversionFactor = 0.01 // 1 Rh = 0.01 Sv

	for _, newMarker := range markers {
		newMarker.DoseRate = newMarker.DoseRate * conversionFactor
		filteredMarkers = append(filteredMarkers, newMarker)
	}
	return filteredMarkers
}

// filterZeroMarkers убирает маркеры с нулевым значением дозы
func filterZeroMarkers(markers []database.Marker) []database.Marker {
	filteredMarkers := []database.Marker{}
	for _, m := range markers {
		if m.DoseRate == 0 {
			continue
		}
		filteredMarkers = append(filteredMarkers, m)
	}
	return filteredMarkers
}

// NEW ────────────────
func isValidDate(ts int64) bool {
	// допустимо «сегодня плюс сутки» с учётом часовых поясов
	max := time.Now().Add(24 * time.Hour).Unix()
	return ts >= minValidTS && ts <= max
}

func filterInvalidDateMarkers(markers []database.Marker) []database.Marker {
	out := markers[:0]
	for _, m := range markers {
		if isValidDate(m.Date) {
			out = append(out, m)
		}
	}
	return out
}

// Проекция Web Mercator приблизительно переводит широту/долготу в "метры".
// Формулы стандартные, здесь используется для перевода в пиксельные координаты.
func latLonToWebMercator(lat, lon float64) (x, y float64) {
	// const радиус Земли для WebMercator
	const originShift = 2.0 * math.Pi * 6378137.0 / 2.0

	x = lon * originShift / 180.0
	y = math.Log(math.Tan((90.0+lat)*math.Pi/360.0)) / (math.Pi / 180.0)
	y = y * originShift / 180.0
	return x, y
}

// webMercatorToPixel переводит Web Mercator координаты (x,y) в пиксели на данном зуме.
func webMercatorToPixel(x, y float64, zoom int) (px, py float64) {
	// тайл 256x256, увеличиваем в 2^zoom
	scale := math.Exp2(float64(zoom))
	px = (x + 2.0*math.Pi*6378137.0/2.0) / (2.0 * math.Pi * 6378137.0) * 256.0 * scale
	py = (2.0*math.Pi*6378137.0/2.0 - y) / (2.0 * math.Pi * 6378137.0) * 256.0 * scale
	return
}

// latLonToPixel - удобная обёртка
func latLonToPixel(lat, lon float64, zoom int) (px, py float64) {
	x, y := latLonToWebMercator(lat, lon)
	return webMercatorToPixel(x, y, zoom)
}

// fastMergeMarkersByZoom группирует маркеры в «ячейку» сетки
// (диаметр = 2*radiusPx) и усредняет данные кластера.
// • O(N) • без мьютексов • подходит для любых зумов.
func fastMergeMarkersByZoom(markers []database.Marker, zoom int, radiusPx float64) []database.Marker {
	if len(markers) == 0 {
		return nil
	}

	cell := 2*radiusPx + 1 // px
	type acc struct {
		sumLat, sumLon, sumDose, sumCnt, sumSp float64
		sumAlt, sumTemp, sumHum                float64
		altCount, tempCount, humCount          int
		detector, radiation                    string
		latest                                 int64
		n                                      int
	}
	cl := make(map[int64]*acc) // key := cx<<32 | cy

	for _, m := range markers {
		px, py := latLonToPixel(m.Lat, m.Lon, zoom)
		key := int64(int(px/cell))<<32 | int64(int32(py/cell))
		a := cl[key]
		if a == nil {
			a = &acc{}
			cl[key] = a
		}
		a.sumLat += m.Lat
		a.sumLon += m.Lon
		a.sumDose += m.DoseRate
		a.sumCnt += m.CountRate
		a.sumSp += m.Speed
		if m.AltitudeValid {
			a.sumAlt += m.Altitude
			a.altCount++
		}
		if m.TemperatureValid {
			a.sumTemp += m.Temperature
			a.tempCount++
		}
		if m.HumidityValid {
			a.sumHum += m.Humidity
			a.humCount++
		}
		if m.Date > a.latest {
			a.latest = m.Date
		}
		if a.detector == "" && m.Detector != "" {
			a.detector = m.Detector
		}
		if a.radiation == "" && m.Radiation != "" {
			a.radiation = m.Radiation
		}
		a.n++
	}

	out := make([]database.Marker, 0, len(cl))
	for _, c := range cl {
		n := float64(c.n)
		var (
			altitude float64
			temp     float64
			hum      float64
		)
		var (
			altValid  bool
			tempValid bool
			humValid  bool
		)
		if c.altCount > 0 {
			altitude = c.sumAlt / float64(c.altCount)
			altValid = true
		}
		if c.tempCount > 0 {
			temp = c.sumTemp / float64(c.tempCount)
			tempValid = true
		}
		if c.humCount > 0 {
			hum = c.sumHum / float64(c.humCount)
			humValid = true
		}
		out = append(out, database.Marker{
			Lat:              c.sumLat / n,
			Lon:              c.sumLon / n,
			DoseRate:         c.sumDose / n,
			CountRate:        c.sumCnt / n,
			Speed:            c.sumSp / n,
			Altitude:         altitude,
			Temperature:      temp,
			Humidity:         hum,
			Detector:         c.detector,
			Radiation:        c.radiation,
			Date:             c.latest,
			Zoom:             zoom,
			TrackID:          markers[0].TrackID,
			AltitudeValid:    altValid,
			TemperatureValid: tempValid,
			HumidityValid:    humValid,
		})
	}
	return out
}

// mergeMarkersByZoom "сливает" (усредняет) маркеры, которые пересекаются в пиксельных координатах
// на текущем зуме. Если расстояние между центрами меньше 2*markerRadiusPx (плюс 1px "запас"), то объединяем.
// deprecated
func mergeMarkersByZoom(markers []database.Marker, zoom int, radiusPx float64) []database.Marker {
	if len(markers) == 0 {
		return nil
	}

	// Сначала готовим структуру с пиксельными координатами
	type markerPixel struct {
		Marker    database.Marker
		Px, Py    float64
		MergedIdx int // -1, если ни с кем ещё не сливался
	}

	mPixels := make([]markerPixel, len(markers))
	for i, m := range markers {
		px, py := latLonToPixel(m.Lat, m.Lon, zoom)
		mPixels[i] = markerPixel{
			Marker:    m,
			Px:        px,
			Py:        py,
			MergedIdx: -1,
		}
	}

	var result []database.Marker

	// Жадно идём по списку, сливаем близкие друг к другу
	for i := 0; i < len(mPixels); i++ {
		if mPixels[i].MergedIdx != -1 {
			// уже слит с кем-то
			continue
		}
		// начинаем новый кластер
		cluster := []markerPixel{mPixels[i]}
		mPixels[i].MergedIdx = i

		// проверяем всех последующих
		for j := i + 1; j < len(mPixels); j++ {
			if mPixels[j].MergedIdx != -1 {
				continue
			}
			dist := math.Hypot(mPixels[i].Px-mPixels[j].Px, mPixels[i].Py-mPixels[j].Py)
			if dist < 2.0*radiusPx {
				// Сливаем
				cluster = append(cluster, mPixels[j])
				mPixels[j].MergedIdx = i // значит, слит к кластеру i
			}
		}

		// Усредняем данные кластера
		var sumLat, sumLon, sumDose, sumCount float64
		var sumAlt, sumTemp, sumHum float64
		var altCount, tempCount, humCount int
		var latestDate int64
		detector := ""
		radiation := ""
		for _, c := range cluster {
			sumLat += c.Marker.Lat
			sumLon += c.Marker.Lon
			sumDose += c.Marker.DoseRate
			sumCount += c.Marker.CountRate
			if c.Marker.AltitudeValid {
				sumAlt += c.Marker.Altitude
				altCount++
			}
			if c.Marker.TemperatureValid {
				sumTemp += c.Marker.Temperature
				tempCount++
			}
			if c.Marker.HumidityValid {
				sumHum += c.Marker.Humidity
				humCount++
			}
			if detector == "" && c.Marker.Detector != "" {
				detector = c.Marker.Detector
			}
			if radiation == "" && c.Marker.Radiation != "" {
				radiation = c.Marker.Radiation
			}
			// возьмём дату последнего
			if c.Marker.Date > latestDate {
				latestDate = c.Marker.Date
			}
		}
		n := float64(len(cluster))
		avgLat := sumLat / n
		avgLon := sumLon / n
		avgDose := sumDose / n
		avgCount := sumCount / n

		var sumSpeed float64
		for _, c := range cluster {
			sumSpeed += c.Marker.Speed
		}
		avgSpeed := sumSpeed / n

		var altitude float64
		var temp float64
		var hum float64
		var altValid bool
		var tempValid bool
		var humValid bool
		if altCount > 0 {
			altitude = sumAlt / float64(altCount)
			altValid = true
		}
		if tempCount > 0 {
			temp = sumTemp / float64(tempCount)
			tempValid = true
		}
		if humCount > 0 {
			hum = sumHum / float64(humCount)
			humValid = true
		}

		// Создаём новый слитый маркер
		newMarker := database.Marker{
			Lat:              avgLat,
			Lon:              avgLon,
			DoseRate:         avgDose,
			CountRate:        avgCount,
			Altitude:         altitude,
			Temperature:      temp,
			Humidity:         hum,
			Detector:         detector,
			Radiation:        radiation,
			Date:             latestDate,
			Speed:            avgSpeed,
			Zoom:             zoom,
			TrackID:          cluster[0].Marker.TrackID, // берем хотя бы у первого
			AltitudeValid:    altValid,
			TemperatureValid: tempValid,
			HumidityValid:    humValid,
		}
		result = append(result, newMarker)
	}

	return result
}

// pickIdentityProbe returns up to 'limit' evenly spaced, non-zero markers
// to cheaply "probe" the DB for an existing track. This avoids thousands
// of random point-lookups on huge tables.
// • No mutexes: pure functional slice logic.
// • Streaming friendly: does not allocate more than needed.
func pickIdentityProbe(src []database.Marker, limit int) []database.Marker {
	if limit <= 0 || len(src) == 0 {
		return nil
	}
	// 1) filter out zero-dose points (they are common and uninformative)
	tmp := make([]database.Marker, 0, min(len(src), limit*2))
	for _, m := range src {
		if m.DoseRate != 0 || m.CountRate != 0 {
			tmp = append(tmp, m)
		}
	}
	if len(tmp) == 0 {
		// fall back to original src if everything was zero
		tmp = src
	}
	// 2) take evenly spaced sample up to 'limit'
	n := len(tmp)
	if n <= limit {
		out := make([]database.Marker, n)
		copy(out, tmp)
		return out
	}
	out := make([]database.Marker, 0, limit)
	stride := n / limit
	if stride <= 0 {
		stride = 1
	}
	for i := 0; i < n && len(out) < limit; i += stride {
		out = append(out, tmp[i])
	}
	return out
}

// calculateSpeedForMarkers recomputes Speed (m/s) for all markers,
// normalizing timestamp units (ms → s when needed). We ignore any
// prefilled speeds and derive velocity from geodesic distance / Δt.
//
// Complexity: O(N).
func calculateSpeedForMarkers(markers []database.Marker) []database.Marker {
	if len(markers) == 0 {
		return markers
	}

	// 1) Chronological order to keep Δt positive and stable.
	sort.Slice(markers, func(i, j int) bool { return markers[i].Date < markers[j].Date })

	// 2) Decide the epoch units once per track:
	//    ~1e9 → seconds (Unix s), ~1e12 → milliseconds (Unix ms).
	//    Check both ends to be safe with mixed sources.
	scale := 1.0 // seconds by default
	if markers[0].Date > 1_000_000_000_000 || markers[len(markers)-1].Date > 1_000_000_000_000 {
		scale = 1000.0 // timestamps are in ms → convert Δt to seconds
	}

	// helper to get Δt in seconds
	dtSec := func(prev, curr int64) float64 {
		if curr <= prev {
			return 0
		}
		return float64(curr-prev) / scale
	}

	const maxSpeed = 1000.0 // m/s, sanity cap for aircraft

	// 3) Recompute pairwise speeds from distance / Δt.
	for i := 1; i < len(markers); i++ {
		dt := dtSec(markers[i-1].Date, markers[i].Date)
		if dt <= 0 {
			continue // duplicate or invalid timestamp
		}
		dist := haversineDistance(
			markers[i-1].Lat, markers[i-1].Lon,
			markers[i].Lat, markers[i].Lon,
		)
		v := dist / dt // m/s
		if v >= 0 && v <= maxSpeed {
			markers[i].Speed = v
		} else {
			// Leave zero if insane (spikes/outliers)
			markers[i].Speed = 0
		}
	}

	// 4) Seed the very first point if needed.
	if len(markers) > 1 && markers[0].Speed == 0 {
		markers[0].Speed = markers[1].Speed
	}

	// 5) Fill zero-speed gaps by borrowing from neighbours.
	lastWithSpeed := -1
	for i := 0; i < len(markers); {
		if markers[i].Speed > 0 {
			lastWithSpeed = i
			i++
			continue
		}
		// zero-run [gapStart..gapEnd]
		gapStart := i
		for i < len(markers) && markers[i].Speed == 0 {
			i++
		}
		gapEnd := i - 1

		// right anchor (if any)
		nextWithSpeed := -1
		if i < len(markers) && markers[i].Speed > 0 {
			nextWithSpeed = i
		}

		var fill float64
		switch {
		case lastWithSpeed != -1 && nextWithSpeed != -1:
			// Prefer average speed derived from anchors distance/time.
			dt := dtSec(markers[lastWithSpeed].Date, markers[nextWithSpeed].Date)
			if dt > 0 {
				dist := haversineDistance(
					markers[lastWithSpeed].Lat, markers[lastWithSpeed].Lon,
					markers[nextWithSpeed].Lat, markers[nextWithSpeed].Lon,
				)
				fill = dist / dt
			}
		case lastWithSpeed != -1:
			fill = markers[lastWithSpeed].Speed
		case nextWithSpeed != -1:
			fill = markers[nextWithSpeed].Speed
		}

		if fill > 0 && fill <= maxSpeed {
			for j := gapStart; j <= gapEnd; j++ {
				markers[j].Speed = fill
			}
		}
	}

	// 6) Global fallback: if anything is still zero, use total distance / total time.
	needFallback := false
	for _, m := range markers {
		if m.Speed == 0 {
			needFallback = true
			break
		}
	}
	if needFallback && len(markers) >= 2 {
		totalDt := dtSec(markers[0].Date, markers[len(markers)-1].Date)
		if totalDt > 0 {
			dist := haversineDistance(
				markers[0].Lat, markers[0].Lon,
				markers[len(markers)-1].Lat, markers[len(markers)-1].Lon,
			)
			v := dist / totalDt
			if v >= 0 && v <= maxSpeed {
				for k := range markers {
					if markers[k].Speed == 0 {
						markers[k].Speed = v
					}
				}
			}
		}
	}

	return markers
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	phi1, phi2 := lat1*math.Pi/180, lat2*math.Pi/180
	dPhi, dLambda := (lat2-lat1)*math.Pi/180, (lon2-lon1)*math.Pi/180
	a := math.Sin(dPhi/2)*math.Sin(dPhi/2) + math.Cos(phi1)*math.Cos(phi2)*math.Sin(dLambda/2)*math.Sin(dLambda/2)
	return 2 * R * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func radiusForZoom(zoom int) float64 {
	// линейная шкала: z=20 → 10 px, z=10 → 5 px, z=5 → 2.5 px …
	return markerRadiusPx * float64(zoom) / 20.0
}

// clusterMarkersForZoom applies on-the-fly clustering to raw markers (zoom=0)
// for display at a specific zoom level. This enables storing only raw markers
// in the database while still providing clustered views for lower zoom levels.
// At zoom >= 15, returns markers unchanged for maximum detail.
func clusterMarkersForZoom(markers []database.Marker, requestedZoom int) []database.Marker {
	if len(markers) == 0 {
		return markers
	}
	// At high zoom levels (15+), show raw markers for maximum detail
	if requestedZoom >= 15 {
		return markers
	}
	// For lower zoom levels, cluster markers to reduce visual clutter
	return fastMergeMarkersByZoom(markers, requestedZoom, radiusForZoom(requestedZoom))
}

// precomputeMarkersForAllZoomLevels creates aggregates for z=1…20
// DEPRECATED: This is kept for backward compatibility but no longer used.
// New uploads store only zoom=0 markers and cluster on-the-fly.
func precomputeMarkersForAllZoomLevels(src []database.Marker) []database.Marker {
	type job struct {
		z   int
		out []database.Marker
	}
	ch := make(chan job, 20)

	for z := 1; z <= 20; z++ {
		go func(zoom int) {
			merged := fastMergeMarkersByZoom(src, zoom, radiusForZoom(zoom))
			ch <- job{z: zoom, out: merged}
		}(z)
	}

	var res []database.Marker
	for i := 0; i < 20; i++ {
		res = append(res, (<-ch).out...)
	}
	return res
}
