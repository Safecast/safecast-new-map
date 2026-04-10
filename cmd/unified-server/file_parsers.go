package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"safecast-new-map/pkg/database"
	"safecast-new-map/pkg/mapimport"
	"safecast-new-map/pkg/spectrum"
)

func processBGeigieZenFile(
	file multipart.File,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {
	logT(trackID, "BGEIGIE", "▶ start (stream)")

	sc := bufio.NewScanner(file)
	sc.Buffer(make([]byte, 0, 64*1024), 2*1024*1024)

	const cpmPerMicroSv = 334.0
	markers := make([]database.Marker, 0, 4096)

	// Extract header information for detector field
	var formatVersion string
	var deviceID string
	var deviceName string

	parsed := 0
	skipped := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())

		// Parse header lines (# comments at start of file)
		if strings.HasPrefix(line, "#") {
			// Extract format version: # format=1.0.3nano or # format=2.0.0CzechRad
			if strings.HasPrefix(line, "# format=") {
				formatVersion = strings.TrimPrefix(line, "# format=")
				logT(trackID, "BGEIGIE", "detected format: %s", formatVersion)
			}
			// Extract device name from nm= field: # nm=PieterFranken,tz=9,cn=JPN,...
			if strings.Contains(line, "nm=") {
				parts := strings.Split(line, ",")
				for _, part := range parts {
					part = strings.TrimPrefix(part, "# ")
					if strings.HasPrefix(part, "nm=") {
						deviceName = strings.TrimPrefix(part, "nm=")
						logT(trackID, "BGEIGIE", "detected device name: %s", deviceName)
						break
					}
				}
			}
			skipped++
			continue
		}

		// Accept all bGeigie format variants: $BNRDD, $BMRDD, $BNXRDD, $CZRDD, $CZRA1, etc.
		if line == "" || !strings.HasPrefix(line, "$") {
			skipped++
			continue
		}
		// Accept lines with RDD (bGeigie standard) or CZR (CzechRad)
		if !strings.Contains(line, "RDD") && !strings.Contains(line, "CZR") {
			skipped++
			continue
		}
		if i := strings.IndexByte(line, '*'); i != -1 {
			line = line[:i]
		}
		p := strings.Split(line, ",")
		if len(p) < 6 { // need at least up to CPMvalid
			skipped++
			continue
		}

		// Extract device ID from first data line (field 1): $BNRDD,0208,... or $CZRDD,0486,...
		if deviceID == "" && len(p) >= 2 {
			deviceID = strings.TrimSpace(p[1])
			logT(trackID, "BGEIGIE", "detected device ID: %s", deviceID)
		}

		var (
			ts  int64
			cps float64
			cpm float64
			lat float64
			lon float64
		)

		// Zen variant: 0:$BNRDD 1:ver 2:ISO8601 3:CPM 4:CPS 5:TotalCounts 6:fix 7:LATdmm 8:N/S 9:LONdmm 10:E/W ...
		// We rely on CPM for the µSv/h conversion because CPS is instantaneous and noisy.
		if len(p) >= 11 && strings.Contains(p[2], "T") {
			if t, err := time.Parse(time.RFC3339, strings.TrimSpace(p[2])); err == nil {
				ts = t.Unix()
			}
			cpm = parseFloat(p[3])
			cps = parseFloat(p[4])
			lat = parseDMM(p[7], p[8], 2)
			lon = parseDMM(p[9], p[10], 3)
		} else if len(p) >= 8 { // legacy fallback: decimals (+ optional suffix)
			// We only accept if date/time parse succeeds via known helper; otherwise skip silently.
			// If parseBGeigieDateTime isn't present, ts remains 0 and entry is skipped.
			ts = 0
			// try compact forms if helper exists in build
			// cps/cpm & coords
			cpm = parseFloat(p[3])
			cps = parseFloat(p[4])
			lat = parseBGeigieCoord(p[6])
			lon = parseBGeigieCoord(p[7])
		}

		if ts == 0 || (lat == 0 && lon == 0) {
			skipped++
			continue
		}

		dose := 0.0
		if cpm > 0 {
			dose = cpm / cpmPerMicroSv
		} else if cps > 0 {
			dose = (cps * 60.0) / cpmPerMicroSv
		}
		// Allow zero-dose records: GPS tracks are valuable even without radiation data

		countRate := cps
		if countRate == 0 && cpm > 0 {
			countRate = cpm / 60.0
		}

		// Build detector string from header info: "bGeigie-<deviceID> (<format>) [<name>]"
		// Examples: "bGeigie-0208 (1.0.3nano) [PieterFranken]" or "bGeigie-0486 (2.0.0CzechRad)"
		detector := ""
		if deviceID != "" {
			detector = "bGeigie-" + deviceID
			if formatVersion != "" {
				detector += " (" + formatVersion + ")"
			}
			if deviceName != "" {
				detector += " [" + deviceName + "]"
			}
		}

		markers = append(markers, database.Marker{
			DoseRate:  dose,
			Date:      ts,
			Lon:       lon,
			Lat:       lat,
			CountRate: countRate,
			Zoom:      0,
			Speed:     0,
			TrackID:   trackID,
			Detector:  detector,
		})
		parsed++
	}
	if err := sc.Err(); err != nil {
		return database.Bounds{}, trackID, err
	}
	if len(markers) == 0 {
		return database.Bounds{}, trackID, fmt.Errorf("no valid bGeigie data points found (parsed=%d skipped=%d)", parsed, skipped)
	}

	bbox, trackID, err := processAndStoreMarkers(markers, trackID, db, dbType)
	if err != nil {
		return bbox, trackID, err
	}
	logT(trackID, "BGEIGIE", "✔ done (parsed=%d)", parsed)
	return bbox, trackID, nil
}

var speedCatalog = map[string]SpeedRange{
	"ped":   {0, 7},     // < 7 м/с   (~0-25 км/ч)
	"car":   {7, 70},    // 7–70 м/с  (~25-250 км/ч)
	"plane": {70, 1000}, // > 70 м/с  (~250-1800 км/ч)
}

// withServerHeader оборачивает любой http.Handler, добавляя
// заголовок "Server: safecast-new-map/<CompileVersion>".
//
// На запрос HEAD к “/” сразу отвечает 200 OK без тела, чтобы
// показать, что сервис жив.


func parseFloat(value string) float64 {
	parsedValue, _ := strconv.ParseFloat(value, 64)
	return parsedValue
}

func getTimeZoneByLongitude(lon float64) *time.Location {
	switch {
	case lon >= -10 && lon <= 0:
		loc, _ := time.LoadLocation("Europe/London")
		return loc
	case lon > 0 && lon <= 15:
		loc, _ := time.LoadLocation("Europe/Berlin")
		return loc
	case lon > 15 && lon <= 30:
		loc, _ := time.LoadLocation("Europe/Kiev")
		return loc
	case lon > 30 && lon <= 45:
		loc, _ := time.LoadLocation("Europe/Moscow")
		return loc
	case lon > 45 && lon <= 60:
		loc, _ := time.LoadLocation("Asia/Yekaterinburg")
		return loc
	case lon > 60 && lon <= 90:
		loc, _ := time.LoadLocation("Asia/Novosibirsk")
		return loc
	case lon > 90 && lon <= 120:
		loc, _ := time.LoadLocation("Asia/Irkutsk")
		return loc
	case lon > 120 && lon <= 135:
		loc, _ := time.LoadLocation("Asia/Yakutsk")
		return loc
	case lon > 135 && lon <= 180:
		loc, _ := time.LoadLocation("Asia/Vladivostok")
		return loc

	case lon >= -180 && lon < -150:
		loc, _ := time.LoadLocation("America/Anchorage")
		return loc
	case lon >= -150 && lon < -120:
		loc, _ := time.LoadLocation("America/Los_Angeles")
		return loc
	case lon >= -120 && lon < -90:
		loc, _ := time.LoadLocation("America/Denver")
		return loc
	case lon >= -90 && lon < -60:
		loc, _ := time.LoadLocation("America/Chicago")
		return loc
	case lon >= -60 && lon < -30:
		loc, _ := time.LoadLocation("America/New_York")
		return loc
	case lon >= -30 && lon < 0:
		loc, _ := time.LoadLocation("America/Halifax")
		return loc

	case lon >= 60 && lon < 75:
		loc, _ := time.LoadLocation("Asia/Karachi")
		return loc
	case lon >= 75 && lon < 90:
		loc, _ := time.LoadLocation("Asia/Kolkata")
		return loc
	case lon >= 90 && lon < 105:
		loc, _ := time.LoadLocation("Asia/Dhaka")
		return loc
	case lon >= 105 && lon < 120:
		loc, _ := time.LoadLocation("Asia/Bangkok")
		return loc
	case lon >= 120 && lon < 135:
		loc, _ := time.LoadLocation("Asia/Shanghai")
		return loc
	case lon >= 135 && lon < 150:
		loc, _ := time.LoadLocation("Asia/Tokyo")
		return loc
	case lon >= 150 && lon <= 180:
		loc, _ := time.LoadLocation("Australia/Sydney")
		return loc

	default:
		loc, _ := time.LoadLocation("UTC")
		return loc
	}
}

// -----------------------------------------------------------------------------
// extractDoseRate — extracts the dose rate from an arbitrary text fragment.
//
//   - «12.3 µR/h»  → 0.123 µSv/h      (1 µR/h ≈ 0.01 µSv/h, legacy iPhone dump)
//   - «0.136 uSv/h»→ 0.136 µSv/h      (Safecast)
//   - «0.29 мкЗв/ч»→ 0.29  µSv/h      (Radiacode-101 Android, RU locale)
//
// -----------------------------------------------------------------------------
func extractDoseRate(s string) float64 {
	// block: legacy µR/h → µSv/h
	reMicroRh := regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*µ?R/h`)
	if m := reMicroRh.FindStringSubmatch(s); len(m) > 0 {
		return parseFloat(m[1]) / 100.0 // convert µR/h → µSv/h
	}

	// block: standard uSv/h (Safecast, iPhone)
	reMicroSv := regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*uSv/h`)
	if m := reMicroSv.FindStringSubmatch(s); len(m) > 0 {
		return parseFloat(m[1]) // already in µSv/h
	}

	// block: russian «мкЗв/ч» (μSv/h in Cyrillic)
	reRuMicroSv := regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*мк?з?в/ч`)
	if m := reRuMicroSv.FindStringSubmatch(s); len(m) > 0 {
		return parseFloat(m[1]) // already in µSv/h
	}
	return 0
}

// -----------------------------------------------------------------------------
// extractCountRate — searches for count rate and normalises it to cps.
//
//   - «24 cps»      → 24
//   - «1500 CPM»    → 25  (1 min → sec)
//   - «24.7 имп/c»  → 24.7 (Radiacode-101 Android, RU locale)
//
// -----------------------------------------------------------------------------
func extractCountRate(s string) float64 {
	// block: cps (all locales)
	reCPS := regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*cps`)
	if m := reCPS.FindStringSubmatch(s); len(m) > 0 {
		return parseFloat(m[1])
	}

	// block: CPM (Safecast CSV)
	reCPM := regexp.MustCompile(`(?i)CPM\s*Value\s*=\s*(\d+(?:\.\d+)?)`)
	if m := reCPM.FindStringSubmatch(s); len(m) > 0 {
		return parseFloat(m[1]) / 60.0 // 1 minute → seconds
	}

	// block: russian «имп/с» or «имп/c» (cyrillic / latin 'c')
	reRU := regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*имп\s*/\s*[cс]`)
	if m := reRU.FindStringSubmatch(s); len(m) > 0 {
		return parseFloat(m[1])
	}
	return 0
}

// -----------------------------------------------------------------------------
// parseDate — recognises three date formats:
//
//   - «May 23, 2012 04:10:08»      (old .rctrk / AtomFast KML)
//   - «2012/05/23 04:10:08»        (Safecast)
//   - «26 июл 2025 11:29:54»       (Radiacode-101 Android, RU locale)
//
// loc — time-zone calculated from longitude (nil → UTC).
// -----------------------------------------------------------------------------
func parseDate(s string, loc *time.Location) int64 {
	if loc == nil {
		loc = time.UTC
	}

	// block: English «Jan 2, 2006 …»
	if m := regexp.MustCompile(`([A-Za-z]{3} \d{1,2}, \d{4} \d{2}:\d{2}:\d{2})`).FindStringSubmatch(s); len(m) > 0 {
		const layout = "Jan 2, 2006 15:04:05"
		if t, err := time.ParseInLocation(layout, m[1], loc); err == nil {
			return t.Unix()
		}
	}

	// block: ISO-ish «2006/01/02 …»
	if m := regexp.MustCompile(`(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2})`).FindStringSubmatch(s); len(m) > 0 {
		const layout = "2006/01/02 15:04:05"
		if t, err := time.ParseInLocation(layout, m[1], loc); err == nil {
			return t.Unix()
		}
	}

	// block: Russian «02 янв 2006 …»
	reRu := regexp.MustCompile(`(\d{1,2})\s+([А-Яа-я]{3})\s+(\d{4})\s+(\d{2}:\d{2}:\d{2})`)
	if m := reRu.FindStringSubmatch(s); len(m) > 0 {
		// map short russian month → number
		ruMon := map[string]string{
			"янв": "01", "фев": "02", "мар": "03", "апр": "04",
			"май": "05", "июн": "06", "июл": "07", "авг": "08",
			"сен": "09", "окт": "10", "ноя": "11", "дек": "12",
		}
		monNum, ok := ruMon[strings.ToLower(m[2])]
		if !ok {
			return 0
		}

		// build ISO-like string and parse
		dateStr := fmt.Sprintf("%s-%s-%02s %s", m[3], monNum, m[1], m[4])
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, loc); err == nil {
			return t.Unix()
		}
	}
	return 0
}

// =====================================================================================
// parseGPX.go — потоковый парсер GPX 1.1 (AtomSwift)
// =====================================================================================
//
// Формат AtomSwift:
//
//	<trkpt lat="…" lon="…">
//	  <time>2025-04-19T14:57:46Z</time>
//	  …
//	  <extensions>
//	    <atom:marker>
//	       <atom:doserate>0.018526316</atom:doserate>  <!-- µSv/h -->
//	       <atom:cp2s>1.0</atom:cp2s>                   <!-- counts / 2 s -->
//	       <atom:speed>0.41898388</atom:speed>          <!-- m/s -->
//	    </atom:marker>
//	  </extensions>
//
// Все интересующие поля находятся внутри <trkpt>.  Парсим потоково без
// дополнительного выделения памяти, никаких mutex – только канал результатов
// внутри ф-ции (go-routine → main goroutine).
//
// =====================================================================================
// parseGPX (stream) — token-driven GPX 1.1 parser (AtomSwift).
// Uses xml.Decoder directly on io.Reader, so we do *zero* extra allocations.
// =====================================================================================
func parseGPX(trackID string, r io.Reader) ([]database.Marker, error) {
	logT(trackID, "GPX", "parser start (stream)")

	type result struct {
		marker database.Marker
		err    error
	}

	out := make(chan result)
	go func() { // parser goroutine
		defer close(out)

		dec := xml.NewDecoder(r)
		var (
			inTrkpt       bool
			lat, lon      float64
			tUnix, doseSv float64
			count, speed  float64
		)

		for {
			tok, err := dec.Token()
			if err == io.EOF {
				return
			}
			if err != nil {
				out <- result{err: fmt.Errorf("XML decode: %w", err)}
				return
			}

			switch el := tok.(type) {
			case xml.StartElement:
				switch el.Name.Local {
				case "trkpt":
					inTrkpt = true
					lat, lon, tUnix, doseSv, count, speed = 0, 0, 0, 0, 0, 0
					for _, a := range el.Attr {
						if a.Name.Local == "lat" {
							lat = parseFloat(a.Value)
						} else if a.Name.Local == "lon" {
							lon = parseFloat(a.Value)
						}
					}
				case "time":
					if inTrkpt {
						var ts string
						_ = dec.DecodeElement(&ts, &el)
						if tt, err := time.Parse(time.RFC3339, ts); err == nil {
							tUnix = float64(tt.Unix())
						}
					}
				case "doserate":
					if inTrkpt {
						var s string
						_ = dec.DecodeElement(&s, &el)
						doseSv = parseFloat(s)
					}
				case "cp2s":
					if inTrkpt {
						var s string
						_ = dec.DecodeElement(&s, &el)
						count = parseFloat(s) / 2.0
					}
				case "speed":
					if inTrkpt {
						var s string
						_ = dec.DecodeElement(&s, &el)
						speed = parseFloat(s)
					}
				}
			case xml.EndElement:
				if el.Name.Local == "trkpt" && inTrkpt {
					inTrkpt = false
					if doseSv == 0 && count == 0 {
						continue
					}
					out <- result{marker: database.Marker{
						Lat:       lat,
						Lon:       lon,
						Date:      int64(tUnix),
						DoseRate:  doseSv,
						CountRate: count,
						Speed:     speed,
					}}
				}
			}
		}
	}()

	var markers []database.Marker
	for r := range out {
		if r.err != nil {
			logT(trackID, "GPX", "✖ %v", r.err)
			return nil, r.err
		}
		markers = append(markers, r.marker)
	}
	logT(trackID, "GPX", "parser done, parsed=%d markers", len(markers))
	if len(markers) == 0 {
		return nil, fmt.Errorf("no <trkpt> with numeric data found")
	}
	return markers, nil
}

// parseKML (stream) — SAX-style KML parser with *constant* time-zone
// for the whole file.  Fixes wrong speeds on tracks that cross
// several time-zones (e.g. airplanes).
func parseKML(trackID string, r io.Reader) ([]database.Marker, error) {
	logT(trackID, "KML", "parser start (stream)")

	dec := xml.NewDecoder(r)

	var (
		inPlacemark bool
		lat, lon    float64
		name, desc  string
		markers     []database.Marker
		tz          *time.Location // ← NEW: chosen once
		tzLocked    bool           // ←   and then locked
	)

	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("XML decode: %v", err)
		}

		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name.Local {
			case "Placemark":
				inPlacemark, lat, lon, name, desc = true, 0, 0, "", ""
			case "name":
				if inPlacemark {
					_ = dec.DecodeElement(&name, &el)
				}
			case "description":
				if inPlacemark {
					_ = dec.DecodeElement(&desc, &el)
				}
			case "coordinates":
				if inPlacemark {
					var coord string
					_ = dec.DecodeElement(&coord, &el)
					parts := strings.Split(coord, ",")
					if len(parts) >= 2 {
						lon = parseFloat(parts[0])
						lat = parseFloat(parts[1])
					}
					// ── выбираем TZ только *один раз* ─────────────
					if !tzLocked {
						tz = getTimeZoneByLongitude(lon)
						tzLocked = true
					}
				}
			}
		case xml.EndElement:
			if el.Name.Local == "Placemark" && inPlacemark {
				inPlacemark = false
				dose := extractDoseRate(name)
				if dose == 0 {
					dose = extractDoseRate(desc)
				}
				count := extractCountRate(desc)
				date := parseDate(desc, tz) // ← используем ЕДИНЫЙ TZ
				if dose == 0 && count == 0 {
					continue
				}
				markers = append(markers, database.Marker{
					DoseRate:  dose,
					CountRate: count,
					Lat:       lat,
					Lon:       lon,
					Date:      date,
				})
			}
		}
	}

	logT(trackID, "KML", "parser done, parsed=%d markers", len(markers))
	if len(markers) == 0 {
		return nil, fmt.Errorf("no valid <Placemark> with numeric data found")
	}
	return markers, nil
}

// =====================================================================================
// parseTextRCTRK.go  — теперь принимает trackID
// =====================================================================================
func parseTextRCTRK(trackID string, data []byte) ([]database.Marker, error) {
	logT(trackID, "RCTRK", "text parser start")

	var markers []database.Marker
	lines := strings.Split(string(data), "\n")

	for idx, line := range lines {
		if idx == 0 || strings.HasPrefix(line, "Timestamp") || strings.TrimSpace(line) == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 8 {
			logT(trackID, "RCTRK", "skip line %d: insufficient fields (%d)", idx+1, len(fields))
			continue
		}

		tsStr := fields[1] + " " + fields[2]
		t, err := time.Parse("2006-01-02 15:04:05", tsStr)
		if err != nil {
			logT(trackID, "RCTRK", "skip line %d: time parse error: %v", idx+1, err)
			continue
		}

		lat, lon := parseFloat(fields[3]), parseFloat(fields[4])
		if lat == 0 || lon == 0 {
			logT(trackID, "RCTRK", "skip line %d: invalid coords (%.6f,%.6f)", idx+1, lat, lon)
			continue
		}

		doseRaw, countRaw := parseFloat(fields[6]), parseFloat(fields[7])
		if doseRaw < 0 || countRaw < 0 {
			logT(trackID, "RCTRK", "skip line %d: negative dose/count", idx+1)
			continue
		}

		markers = append(markers, database.Marker{
			DoseRate:  doseRaw / 100.0,
			CountRate: countRaw,
			Lat:       lat,
			Lon:       lon,
			Date:      t.Unix(),
		})
	}

	logT(trackID, "RCTRK", "text parser done, parsed=%d markers", len(markers))
	return markers, nil
}

// =====================================================================================
// parseRadiacodeCSV — parses Radiacode 103 TSV/CSV export format
// Format: Timestamp	Time	Latitude	Longitude	Accuracy	DoseRate	CountRate	Comment
// =====================================================================================
func parseRadiacodeCSV(trackID string, r io.Reader) ([]database.Marker, error) {
	logT(trackID, "CSV", "Radiacode parser start")

	br := bufio.NewReaderSize(r, 512*1024)
	cr := csv.NewReader(br)
	cr.Comma = '\t' // Tab-separated
	cr.FieldsPerRecord = -1
	cr.LazyQuotes = true

	// Read first line - might be metadata or header
	firstLine, err := cr.Read()
	if err != nil {
		return nil, fmt.Errorf("CSV read error: %v", err)
	}

	// Extract device model from metadata line if present
	// Format: "Track: 2025-12-04 16-57-03\tRC-103-012737\tEC"
	deviceModel := ""
	var header []string
	if len(firstLine) > 0 && strings.HasPrefix(firstLine[0], "Track:") {
		// Extract device model from second column (e.g., "RC-103-012737")
		if len(firstLine) > 1 {
			deviceModel = strings.TrimSpace(firstLine[1])
			logT(trackID, "CSV", "detected device: %s", deviceModel)
		}
		logT(trackID, "CSV", "skipping metadata line: %s", firstLine[0])
		header, err = cr.Read()
		if err != nil {
			return nil, fmt.Errorf("CSV header: %v", err)
		}
	} else {
		header = firstLine
	}

	// Verify this is a Radiacode format
	if len(header) < 7 || !strings.Contains(strings.ToLower(header[0]), "timestamp") {
		return nil, fmt.Errorf("not a Radiacode CSV format")
	}

	markers := make([]database.Marker, 0, 4096)
	rowN := 1

	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logT(trackID, "CSV", "row %d: %v", rowN+1, err)
			rowN++
			continue
		}
		rowN++

		if len(rec) < 7 {
			logT(trackID, "CSV", "skip row %d: insufficient fields (%d)", rowN, len(rec))
			continue
		}

		// Parse timestamp from Time column (index 1: "2025-12-04 21:57:06")
		tsStr := strings.TrimSpace(rec[1])

		lat := parseFloat(rec[2])
		lon := parseFloat(rec[3])
		// rec[4] is Accuracy - we could store this but currently not used
		doseRate := parseFloat(rec[5])
		countRate := parseFloat(rec[6])

		if lat == 0 && lon == 0 {
			logT(trackID, "CSV", "skip row %d: invalid coords (0,0)", rowN)
			continue
		}

		if doseRate < 0 || countRate < 0 {
			logT(trackID, "CSV", "skip row %d: negative dose/count", rowN)
			continue
		}

		// Infer timezone from longitude and parse timestamp in that timezone
		loc := getTimeZoneByLongitude(lon)
		t, err := time.ParseInLocation("2006-01-02 15:04:05", tsStr, loc)
		if err != nil {
			logT(trackID, "CSV", "skip row %d: time parse error: %v", rowN, err)
			continue
		}

		// Radiacode CSV exports DoseRate in units that are ~50x higher than µSv/h
		// Divide by 50 to convert to standard µSv/h for consistency with Safecast data
		markers = append(markers, database.Marker{
			DoseRate:  doseRate / 50.0,
			CountRate: countRate,
			Lat:       lat,
			Lon:       lon,
			Date:      t.Unix(),
			Detector:  deviceModel,
		})
	}

	if len(markers) == 0 {
		return nil, fmt.Errorf("no valid data rows found")
	}
	logT(trackID, "CSV", "Radiacode parser done, parsed=%d markers", len(markers))
	return markers, nil
}

// =====================================================================================
// parseAtomSwiftCSV (stream) — parses huge .csv produced by AtomSwift logger
// fast & memory-friendly: no ReadAll(), we read record-by-record through bufio.Reader.
// =====================================================================================
func parseAtomSwiftCSV(trackID string, r io.Reader) ([]database.Marker, error) {
	logT(trackID, "CSV", "AtomSwift parser start (stream)")

	br := bufio.NewReaderSize(r, 512*1024) // 512 KiB read-ahead buffer
	cr := csv.NewReader(br)
	cr.Comma = ';'
	cr.FieldsPerRecord = -1 // keep tolerant

	// skip header -----------------------------------------------------------
	if _, err := cr.Read(); err != nil {
		return nil, fmt.Errorf("CSV header: %v", err)
	}

	markers := make([]database.Marker, 0, 4096) // pre-allocate reasonable cap
	rowN := 1
	for {
		rec, err := cr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logT(trackID, "CSV", "row %d: %v", rowN+1, err)
			continue
		}
		rowN++

		if len(rec) < 7 {
			logT(trackID, "CSV", "skip row %d: insufficient fields (%d)", rowN, len(rec))
			continue
		}

		ts, err := strconv.ParseInt(strings.TrimSpace(rec[0]), 10, 64)
		if err != nil {
			logT(trackID, "CSV", "skip row %d: bad timestamp", rowN)
			continue
		}

		dose := parseFloat(rec[1]) // µSv/h
		lat := parseFloat(rec[2])
		lon := parseFloat(rec[3])
		speed := parseFloat(rec[5]) // m/s
		cps := parseFloat(rec[6])

		if lat == 0 || lon == 0 || dose == 0 {
			continue
		}

		markers = append(markers, database.Marker{
			DoseRate:  dose,
			CountRate: cps,
			Lat:       lat,
			Lon:       lon,
			Date:      ts,
			Speed:     speed,
		})
	}

	if len(markers) == 0 {
		return nil, fmt.Errorf("no valid data rows found")
	}
	logT(trackID, "CSV", "AtomSwift parser done, parsed=%d markers", len(markers))
	return markers, nil
}

// processCSVFile handles *.csv uploads - auto-detects format (AtomSwift or Radiacode)
func processCSVFile(
	file multipart.File,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {

	logT(trackID, "CSV", "▶ start - detecting format")

	// Read first line to detect format
	peek := make([]byte, 1024)
	n, _ := file.Read(peek)
	header := string(peek[:n])

	// Reset file pointer to beginning
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, 0)
	} else {
		return database.Bounds{}, trackID, fmt.Errorf("cannot seek in file to auto-detect CSV format")
	}

	var markers []database.Marker
	var err error

	// Detect Radiacode format (tab-separated with "Timestamp" header)
	if strings.Contains(header, "Timestamp\t") && strings.Contains(header, "DoseRate") && strings.Contains(header, "CountRate") {
		logT(trackID, "CSV", "detected Radiacode 103 format")
		markers, err = parseRadiacodeCSV(trackID, file)
	} else {
		// Default to AtomSwift format (semicolon-separated)
		logT(trackID, "CSV", "detected AtomSwift format")
		markers, err = parseAtomSwiftCSV(trackID, file)
	}

	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("parse CSV: %w", err)
	}
	logT(trackID, "CSV", "parsed %d markers", len(markers))

	bbox, trackID, err := processAndStoreMarkers(markers, trackID, db, dbType)
	if err != nil {
		return bbox, trackID, err
	}
	logT(trackID, "CSV", "✔ done")
	return bbox, trackID, nil
}

// processGPXFile handles plain *.gpx uploads in streaming mode.
func processGPXFile(
	file multipart.File,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {
	logT(trackID, "GPX", "▶ start (stream)")

	markers, err := parseGPX(trackID, file)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("parse GPX: %w", err)
	}
	logT(trackID, "GPX", "parsed %d markers", len(markers))

	bbox, trackID, err := processAndStoreMarkers(markers, trackID, db, dbType)
	if err != nil {
		return bbox, trackID, err
	}

	logT(trackID, "GPX", "✔ done")
	return bbox, trackID, nil
}

// processKMLFile handles plain *.kml uploads in streaming mode.
func processKMLFile(
	file multipart.File,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {
	logT(trackID, "KML", "▶ start (stream)")

	markers, err := parseKML(trackID, file)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("parse KML: %w", err)
	}
	logT(trackID, "KML", "parsed %d markers", len(markers))

	bbox, trackID, err := processAndStoreMarkers(markers, trackID, db, dbType)
	if err != nil {
		return bbox, trackID, err
	}
	logT(trackID, "KML", "✔ done")
	return bbox, trackID, nil
}

// processKMZFile handles *.kmz (ZIP archive with KML inside).
func processKMZFile(
	file multipart.File,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {

	logT(trackID, "KMZ", "▶ start")

	data, err := io.ReadAll(file)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("read KMZ: %w", err)
	}

	zipR, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("open KMZ as ZIP: %w", err)
	}

	// accumulate bbox of *all* KML entries inside KMZ
	global := database.Bounds{MinLat: 90, MinLon: 180, MaxLat: -90, MaxLon: -180}

	for _, zf := range zipR.File {
		if filepath.Ext(zf.Name) != ".kml" {
			continue
		}

		kmlF, err := zf.Open()
		if err != nil {
			return global, trackID, fmt.Errorf("open %s: %w", zf.Name, err)
		}
		kmlMarkers, err := parseKML(trackID, kmlF)
		_ = kmlF.Close()
		if err != nil {
			return global, trackID, fmt.Errorf("parse %s: %w", zf.Name, err)
		}
		logT(trackID, "KMZ", "parsed %d markers from %q", len(kmlMarkers), zf.Name)

		bbox, trackID, err := processAndStoreMarkers(kmlMarkers, trackID, db, dbType)
		if err != nil {
			return global, trackID, err
		}

		// expand global bbox
		if bbox.MinLat < global.MinLat {
			global.MinLat = bbox.MinLat
		}
		if bbox.MaxLat > global.MaxLat {
			global.MaxLat = bbox.MaxLat
		}
		if bbox.MinLon < global.MinLon {
			global.MinLon = bbox.MinLon
		}
		if bbox.MaxLon > global.MaxLon {
			global.MaxLon = bbox.MaxLon
		}

		logT(trackID, "KMZ", "✔ done")
		return global, trackID, nil

	}

	logT(trackID, "KMZ", "✔ done")
	return global, trackID, nil
}

// -----------------------------------------------------------------------------
// processRCTRKFile — принимает *.rctrk (Radiacode) в JSON- или текстовом виде.
// Поддерживает оба признака единиц: "sv" (новый Android) и "isSievert" (старый iOS).
// Если ни одного флага нет — считаем, что числа уже в µSv/h и конвертацию НЕ делаем.
// -----------------------------------------------------------------------------
func processRCTRKFile(
	file multipart.File,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {

	logT(trackID, "RCTRK", "▶ start")

	raw, err := io.ReadAll(file)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("read RCTRK: %w", err)
	}

	// ---------- JSON ----------------------------------------------------------
	// По умолчанию оба флага TRUE ⇒ «единицы уже µSv/h».
	data := database.Data{
		IsSievert:       true,
		IsSievertLegacy: true,
	}

	if err := json.Unmarshal(raw, &data); err == nil && len(data.Markers) > 0 {
		logT(trackID, "RCTRK", "JSON detected, %d markers", len(data.Markers))

		// Выясняем, были ли в файле хоть какие-то флаги.
		// Для этого дешево парсим ключи верхнего уровня.
		var keys map[string]json.RawMessage
		_ = json.Unmarshal(raw, &keys) // ошибок игнорируем — структура уже распарсена

		_, hasSV := keys["sv"]
		_, hasOld := keys["isSievert"]

		flagPresent := hasSV || hasOld

		needConvert := flagPresent && (!data.IsSievert || !data.IsSievertLegacy)
		if needConvert {
			logT(trackID, "RCTRK", "µR/h detected → converting to µSv/h")
			data.Markers = convertRhToSv(data.Markers)
		}

		return processAndStoreMarkers(data.Markers, trackID, db, dbType)
	}

	// ---------- plain-text fallback ------------------------------------------
	markers, err := parseTextRCTRK(trackID, raw)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("parse text RCTRK: %w", err)
	}
	logT(trackID, "RCTRK", "parsed %d markers (text)", len(markers))

	return processAndStoreMarkers(markers, trackID, db, dbType)
}

// processN42File handles ANSI N42.42 XML files with gamma spectrum data.
func processN42File(
	file multipart.File,
	filename string,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {
	logT(trackID, "N42", "▶ start")

	raw, err := io.ReadAll(file)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("read N42: %w", err)
	}

	// Parse N42 file to extract spectra and markers
	spectra, markers, err := spectrum.ParseN42(raw)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("parse N42: %w", err)
	}

	logT(trackID, "N42", "parsed %d spectra, %d markers", len(spectra), len(markers))

	if len(markers) == 0 {
		return database.Bounds{}, trackID, fmt.Errorf("no markers found in N42 file")
	}

	// Set trackID for all markers
	for i := range markers {
		if markers[i].TrackID == "" {
			markers[i].TrackID = trackID
		}
	}

	// Store markers and get their IDs
	bounds, storedTrackID, err := processAndStoreMarkers(markers, trackID, db, dbType)
	if err != nil {
		return bounds, storedTrackID, fmt.Errorf("store markers: %w", err)
	}

	// Get the marker IDs from database to link spectra
	ctx := context.Background()
	for i := range spectra {
		// Find the corresponding marker by matching coordinates and timestamp
		// For simplicity, we'll assume the order is preserved
		if i < len(markers) {
			// Query for ALL marker IDs with this location and timestamp (across all zoom levels)
			query := "SELECT id FROM markers WHERE lat = ? AND lon = ? AND date = ?"
			args := []interface{}{markers[i].Lat, markers[i].Lon, markers[i].Date}

			if dbType == "pgx" {
				query = "SELECT id FROM markers WHERE lat = $1 AND lon = $2 AND date = $3"
			}

			rows, err := db.DB.QueryContext(ctx, query, args...)
			if err != nil {
				logT(trackID, "N42", "warning: could not query markers for spectrum %d: %v", i, err)
				continue
			}

			var markerIDs []int64
			for rows.Next() {
				var mid int64
				if err := rows.Scan(&mid); err != nil {
					logT(trackID, "N42", "warning: failed to scan marker ID: %v", err)
					continue
				}
				markerIDs = append(markerIDs, mid)
			}
			rows.Close()

			if len(markerIDs) == 0 {
				logT(trackID, "N42", "warning: no markers found for spectrum %d", i)
				continue
			}

			// Insert spectrum using the first marker ID (typically zoom=0)
			spectra[i].MarkerID = markerIDs[0]
			spectra[i].RawData = raw // Store original N42 file
			spectra[i].Filename = filename

			spectrumID, err := db.InsertSpectrum(ctx, spectra[i])
			if err != nil {
				logT(trackID, "N42", "warning: failed to insert spectrum %d: %v", i, err)
				continue
			}

			// Update has_spectrum flag for ALL markers at this location/time (all zoom levels)
			updateQuery := "UPDATE markers SET has_spectrum = ? WHERE lat = ? AND lon = ? AND date = ?"
			updateArgs := []interface{}{true, markers[i].Lat, markers[i].Lon, markers[i].Date}
			if dbType == "pgx" {
				updateQuery = "UPDATE markers SET has_spectrum = $1 WHERE lat = $2 AND lon = $3 AND date = $4"
			}

			result, err := db.DB.ExecContext(ctx, updateQuery, updateArgs...)
			if err != nil {
				logT(trackID, "N42", "warning: failed to update has_spectrum flags: %v", err)
			} else {
				if count, _ := result.RowsAffected(); count > 0 {
					logT(trackID, "N42", "inserted spectrum %d for marker %d (updated %d zoom levels)", spectrumID, markerIDs[0], count)
				}
			}
		}
	}

	logT(trackID, "N42", "✓ complete")
	return bounds, storedTrackID, nil
}

// processSPEFile handles IAEA SPE format files with gamma spectrum data.
func processSPEFile(
	file multipart.File,
	filename string,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {
	logT(trackID, "SPE", "▶ start")

	raw, err := io.ReadAll(file)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("read SPE: %w", err)
	}

	// Parse SPE file to extract spectra and markers
	spectra, markers, err := spectrum.ParseSPE(raw)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("parse SPE: %w", err)
	}

	logT(trackID, "SPE", "parsed %d spectra, %d markers", len(spectra), len(markers))

	if len(markers) == 0 {
		return database.Bounds{}, trackID, fmt.Errorf("no markers found in SPE file")
	}

	// Set trackID for all markers
	for i := range markers {
		if markers[i].TrackID == "" {
			markers[i].TrackID = trackID
		}
	}

	// Store markers and get their IDs
	bounds, storedTrackID, err := processAndStoreMarkers(markers, trackID, db, dbType)
	if err != nil {
		return bounds, storedTrackID, fmt.Errorf("store markers: %w", err)
	}

	// Get the marker IDs from database to link spectra
	ctx := context.Background()
	for i := range spectra {
		// Find the corresponding marker by matching coordinates and timestamp
		// For simplicity, we'll assume the order is preserved
		if i < len(markers) {
			// Query for ALL marker IDs with this location and timestamp (across all zoom levels)
			query := "SELECT id FROM markers WHERE lat = ? AND lon = ? AND date = ?"
			args := []interface{}{markers[i].Lat, markers[i].Lon, markers[i].Date}

			if dbType == "pgx" {
				query = "SELECT id FROM markers WHERE lat = $1 AND lon = $2 AND date = $3"
			}

			rows, err := db.DB.QueryContext(ctx, query, args...)
			if err != nil {
				logT(trackID, "SPE", "warning: could not query markers for spectrum %d: %v", i, err)
				continue
			}

			var markerIDs []int64
			for rows.Next() {
				var mid int64
				if err := rows.Scan(&mid); err != nil {
					logT(trackID, "SPE", "warning: failed to scan marker ID: %v", err)
					continue
				}
				markerIDs = append(markerIDs, mid)
			}
			rows.Close()

			if len(markerIDs) == 0 {
				logT(trackID, "SPE", "warning: no markers found for spectrum %d", i)
				continue
			}

			// Insert spectrum using the first marker ID (typically zoom=0)
			spectra[i].MarkerID = markerIDs[0]
			spectra[i].RawData = raw // Store original SPE file
			spectra[i].Filename = filename

			spectrumID, err := db.InsertSpectrum(ctx, spectra[i])
			if err != nil {
				logT(trackID, "SPE", "warning: failed to insert spectrum %d: %v", i, err)
				continue
			}

			// Update has_spectrum flag for ALL markers at this location/time (all zoom levels)
			updateQuery := "UPDATE markers SET has_spectrum = ? WHERE lat = ? AND lon = ? AND date = ?"
			updateArgs := []interface{}{true, markers[i].Lat, markers[i].Lon, markers[i].Date}
			if dbType == "pgx" {
				updateQuery = "UPDATE markers SET has_spectrum = $1 WHERE lat = $2 AND lon = $3 AND date = $4"
			}

			result, err := db.DB.ExecContext(ctx, updateQuery, updateArgs...)
			if err != nil {
				logT(trackID, "SPE", "warning: failed to update has_spectrum flags: %v", err)
			} else {
				if count, _ := result.RowsAffected(); count > 0 {
					logT(trackID, "SPE", "inserted spectrum %d for marker %d (updated %d zoom levels)", spectrumID, markerIDs[0], count)
				}
			}
		}
	}

	logT(trackID, "SPE", "✓ complete")
	return bounds, storedTrackID, nil
}

// processRCXMLFile handles RadiaCode proprietary XML spectrum files.
func processRCXMLFile(
	file multipart.File,
	filename string,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {
	logT(trackID, "RCXML", "▶ start")

	raw, err := io.ReadAll(file)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("read RCXML: %w", err)
	}

	// Parse RadiaCode XML file to extract spectra and markers
	spectra, markers, err := spectrum.ParseRCXML(raw)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("parse RCXML: %w", err)
	}

	logT(trackID, "RCXML", "parsed %d spectra, %d markers", len(spectra), len(markers))

	if len(markers) == 0 {
		return database.Bounds{}, trackID, fmt.Errorf("no markers found in RCXML file")
	}

	// Set trackID for all markers
	for i := range markers {
		if markers[i].TrackID == "" {
			markers[i].TrackID = trackID
		}
	}

	// Store markers and get their IDs
	bounds, storedTrackID, err := processAndStoreMarkers(markers, trackID, db, dbType)
	if err != nil {
		return bounds, storedTrackID, fmt.Errorf("store markers: %w", err)
	}

	// Get the marker IDs from database to link spectra
	ctx := context.Background()
	for i := range spectra {
		if i < len(markers) {
			// Query for ALL marker IDs with this location and timestamp (across all zoom levels)
			query := "SELECT id FROM markers WHERE lat = ? AND lon = ? AND date = ?"
			args := []interface{}{markers[i].Lat, markers[i].Lon, markers[i].Date}

			if dbType == "pgx" {
				query = "SELECT id FROM markers WHERE lat = $1 AND lon = $2 AND date = $3"
			}

			rows, err := db.DB.QueryContext(ctx, query, args...)
			if err != nil {
				logT(trackID, "RCXML", "warning: could not query markers for spectrum %d: %v", i, err)
				continue
			}

			var markerIDs []int64
			for rows.Next() {
				var mid int64
				if err := rows.Scan(&mid); err != nil {
					logT(trackID, "RCXML", "warning: failed to scan marker ID: %v", err)
					continue
				}
				markerIDs = append(markerIDs, mid)
			}
			rows.Close()

			if len(markerIDs) == 0 {
				logT(trackID, "RCXML", "warning: no markers found for spectrum %d", i)
				continue
			}

			// Insert spectrum using the first marker ID (typically zoom=0)
			spectra[i].MarkerID = markerIDs[0]
			spectra[i].RawData = raw // Store original RCXML file
			spectra[i].Filename = filename

			spectrumID, err := db.InsertSpectrum(ctx, spectra[i])
			if err != nil {
				logT(trackID, "RCXML", "warning: failed to insert spectrum %d: %v", i, err)
				continue
			}

			// Update has_spectrum flag for ALL markers at this location/time (all zoom levels)
			updateQuery := "UPDATE markers SET has_spectrum = ? WHERE lat = ? AND lon = ? AND date = ?"
			updateArgs := []interface{}{true, markers[i].Lat, markers[i].Lon, markers[i].Date}
			if dbType == "pgx" {
				updateQuery = "UPDATE markers SET has_spectrum = $1 WHERE lat = $2 AND lon = $3 AND date = $4"
			}

			result, err := db.DB.ExecContext(ctx, updateQuery, updateArgs...)
			if err != nil {
				logT(trackID, "RCXML", "warning: failed to update has_spectrum flags: %v", err)
			} else {
				if count, _ := result.RowsAffected(); count > 0 {
					logT(trackID, "RCXML", "inserted spectrum %d for marker %d (updated %d zoom levels)", spectrumID, markerIDs[0], count)
				}
			}
		}
	}

	logT(trackID, "RCXML", "✓ complete")
	return bounds, storedTrackID, nil
}

// processAtomFastFile handles Atom Fast JSON export (*.json).
func processAtomFastFile(
	file multipart.File,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {

	logT(trackID, "AtomFast", "▶ start")

	data, err := io.ReadAll(file)
	if err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("read AtomFast JSON: %w", err)
	}

	return processAtomFastData(data, trackID, db, dbType)
}

func processAtomFastData(
	data []byte,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, error) {
	var records []struct {
		D   float64 `json:"d"`
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
		T   int64   `json:"t"`
	}
	if err := json.Unmarshal(data, &records); err != nil {
		return database.Bounds{}, trackID, fmt.Errorf("parse AtomFast JSON: %w", err)
	}
	logT(trackID, "AtomFast", "parsed %d markers", len(records))

	markers := make([]database.Marker, 0, len(records))
	for _, r := range records {
		markers = append(markers, database.Marker{
			DoseRate:  r.D,
			CountRate: r.D, // AtomFast stores cps in same field
			Lat:       r.Lat,
			Lon:       r.Lng,
			Date:      r.T / 1000, // ms → s
		})
	}

	return processAndStoreMarkers(markers, trackID, db, dbType)
}

// processTrackExportFile ingests a single JSON track generated by our exporter.
// The helper performs an optimistic duplicate probe before touching the
// database so re-uploaded exports simply reuse the existing track.
func processTrackExportFile(
	ctx context.Context,
	file multipart.File,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, bool, error) {
	return processTrackExportReader(ctx, file, trackID, db, dbType)
}

// archiveProgress keeps the streaming import loop and the logger connected via
// a channel so we avoid mutexes. Each update notes how many entries we have
// consumed and the most recent filename so operators can track forward motion
// on gigantic archives without staring at a static log line.
type archiveProgress struct {
	entries  int
	filename string
}

// archiveEntryResult carries the outcome of a single archive item import through
// a channel so the tar reader loop can enforce timeouts without blocking on a
// stuck DB call. Using a struct keeps the select statement tidy while remaining
// explicit about the data we expect back.
type archiveEntryResult struct {
	bounds   database.Bounds
	track    string
	inserted bool
	err      error
}

// observeContext lets long-running import stages bail out quickly when the caller
// cancels, keeping goroutines from running past their deadlines. Returning the
// context error keeps the caller aware that nothing beyond this point was
// committed, which in turn prevents duplicate progress accounting.
func observeContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

// processTrackExportArchive handles the weekly tgz bundle produced by the API.
// We iterate entries sequentially because tar readers are streaming, yet still
// lean on channels inside the parser so each export stays memory-light.
func processTrackExportArchive(
	ctx context.Context,
	file multipart.File,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, bool, error) {
	return processTrackExportArchiveReader(ctx, file, trackID, db, dbType, nil)
}

// processTrackExportArchiveReader is the shared implementation for multipart
// uploads and remote downloads. The optional progress channel lets callers
// stream status to logs without blocking the parsing loop.
func processTrackExportArchiveReader(
	ctx context.Context,
	r io.Reader,
	trackID string,
	db *database.Database,
	dbType string,
	updates chan<- archiveProgress,
) (database.Bounds, string, bool, error) {
	logT(trackID, "Export-TGZ", "▶ start")

	gz, err := gzip.NewReader(r)
	if err != nil {
		return database.Bounds{}, trackID, false, fmt.Errorf("open tgz: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)

	combined := database.Bounds{}
	haveBounds := false
	importedAny := false
	primaryTrack := trackID
	entryIndex := 0

	sendProgress := func(name string) {
		if updates == nil {
			return
		}
		select {
		case updates <- archiveProgress{entries: entryIndex, filename: name}:
		default:
		}
	}

	for {
		select {
		case <-ctx.Done():
			return combined, primaryTrack, importedAny, ctx.Err()
		default:
		}

		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return combined, primaryTrack, importedAny, fmt.Errorf("read tgz entry: %w", err)
		}
		if hdr.FileInfo().IsDir() {
			continue
		}
		name := strings.ToLower(strings.TrimSpace(hdr.Name))
		if !strings.HasSuffix(name, ".cim") && !strings.HasSuffix(name, ".json") {
			continue
		}

		logT(trackID, "Export-TGZ", "processing %s", hdr.Name)
		sendProgress(hdr.Name)

		limited := io.LimitReader(tr, hdr.Size)
		payload, readErr := io.ReadAll(limited)
		if readErr != nil {
			logT(trackID, "Export-TGZ", "skip %s: read failure: %v", hdr.Name, readErr)
			entryIndex++
			sendProgress(hdr.Name)
			continue
		}

		// Avoid fixed per-entry timeouts so gigantic tracks can finish peacefully.
		// A cancellation from the caller still flows through entryCtx, keeping the
		// goroutine responsive without discarding legitimate long-running imports.
		entryCtx, cancel := context.WithCancel(ctx)
		results := make(chan archiveEntryResult, 1)
		go func(data []byte) {
			defer close(results)
			entryBounds, entryTrack, inserted, err := processTrackExportReader(entryCtx, bytes.NewReader(data), GenerateSerialNumber(), db, dbType)
			results <- archiveEntryResult{bounds: entryBounds, track: entryTrack, inserted: inserted, err: err}
		}(payload)

		select {
		case <-ctx.Done():
			cancel()
			// Wait briefly so the worker can observe the cancellation before we move on.
			select {
			case <-results:
			case <-time.After(15 * time.Second):
				logT(trackID, "Export-TGZ", "waited for cancelled %s worker to exit", hdr.Name)
			}
			return combined, primaryTrack, importedAny, ctx.Err()
		case result := <-results:
			cancel()
			if result.err != nil {
				logT(trackID, "Export-TGZ", "skip %s: %v", hdr.Name, result.err)
				entryIndex++
				sendProgress(hdr.Name)
				continue
			}
			if result.track != "" && primaryTrack == trackID {
				primaryTrack = result.track
			}
			combined, haveBounds = mergeBounds(combined, result.bounds, haveBounds)
			if result.inserted {
				importedAny = true
			}
			entryIndex++
			sendProgress(hdr.Name)
		}
	}

	if entryIndex == 0 {
		return database.Bounds{}, trackID, false, fmt.Errorf("tgz archive contained no track export files")
	}
	logT(trackID, "Export-TGZ", "✔ done (entries=%d imported=%v)", entryIndex, importedAny)
	return combined, primaryTrack, importedAny, nil
}

// processTrackExportReader centralises the decoding and duplicate detection shared by
// single-file and archive imports. We keep the logic small so the upload handler
// simply forwards the context and fallback TrackID.
func processTrackExportReader(
	ctx context.Context,
	r io.Reader,
	fallbackTrackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, bool, error) {
	// Early exit keeps the goroutine responsive to caller cancellation so timeouts
	// never leak work into the next archive entry.
	if err := observeContext(ctx); err != nil {
		return database.Bounds{}, fallbackTrackID, false, err
	}

	payload, err := mapimport.Parse(r)
	if err != nil {
		return database.Bounds{}, fallbackTrackID, false, err
	}

	// Respect overrides in the payload but keep a fallback for archive-level defaults.
	parsedTrackID := strings.TrimSpace(payload.TrackID)
	trackID := fallbackTrackID
	if parsedTrackID != "" {
		trackID = parsedTrackID
	}

	// Ignore live-only exports because they belong to the realtime cache.
	// Skipping them keeps weekly archives from stalling on transient snapshots
	// while the map still renders live points directly from the realtime table.
	if strings.HasPrefix(trackID, "live:") {
		logT(trackID, "Export", "skip live track export payload")
		return database.Bounds{}, trackID, false, nil
	}

	markers, bounds := payload.ToDatabaseMarkers(trackID)
	if len(markers) == 0 {
		return bounds, trackID, false, fmt.Errorf("track export import: no usable markers")
	}

	logT(trackID, "Export", "parsed %d markers", len(markers))

	if err := observeContext(ctx); err != nil {
		return bounds, trackID, false, err
	}

	probe := pickIdentityProbe(markers, 128)
	threshold := min(len(probe), 10)
	if len(probe) > 0 {
		if threshold == 0 {
			threshold = len(probe)
		}
		existing, detectErr := db.DetectExistingTrackID(probe, threshold, dbType)
		if detectErr != nil {
			return bounds, trackID, false, fmt.Errorf("detect duplicate: %w", detectErr)
		}
		if existing != "" {
			logT(existing, "CIM", "duplicate payload matches existing track; skipping import")
			return bounds, existing, false, nil
		}
	}

	storedBounds, finalTrackID, err := processAndStoreMarkersWithContext(ctx, markers, trackID, db, dbType)
	if err != nil {
		return storedBounds, finalTrackID, false, err
	}
	return storedBounds, finalTrackID, true, nil
}

// processTrackExportPayload lets callers attempt to import a raw JSON payload
// as an export while allowing graceful fallback to other JSON parsers when the
// payload does not match our schema. Using a byte slice keeps retries cheap for
// callers that need multiple attempts.
func processTrackExportPayload(
	ctx context.Context,
	data []byte,
	trackID string,
	db *database.Database,
	dbType string,
) (database.Bounds, string, bool, error) {
	reader := bytes.NewReader(data)
	return processTrackExportReader(ctx, reader, trackID, db, dbType)
}
