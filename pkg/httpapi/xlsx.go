// xlsx.go — minimal XLSX generator using archive/zip + encoding/xml.
// Produces a valid .xlsx file with a single worksheet containing track marker data.
package httpapi

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"safecast-new-map/pkg/database"
	"safecast-new-map/pkg/trackjson"
)

// xlsxPackage holds the data needed to build a minimal .xlsx workbook.
type xlsxPackage struct {
	SharedStrings []string
	Markers       []database.Marker
}

// Write produces a valid .xlsx file to w.
func (p *xlsxPackage) Write(w io.Writer) error {
	zw := zip.NewWriter(w)

	// 1. [Content_Types].xml
	ctXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
<Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
<Override PartName="/xl/sharedStrings.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sharedStrings+xml"/>
</Types>`
	if err := addFile(zw, "[Content_Types].xml", ctXML); err != nil {
		return err
	}

	// 2. _rels/.rels
	relsXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`
	if err := addFile(zw, "_rels/.rels", relsXML); err != nil {
		return err
	}

	// 3. xl/workbook.xml
	wbXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
<sheets><sheet name="Track" sheetId="1" r:id="rId1"/></sheets>
</workbook>`
	if err := addFile(zw, "xl/workbook.xml", wbXML); err != nil {
		return err
	}

	// 4. xl/_rels/workbook.xml.rels
	wbRelsXML := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/sharedStrings" Target="sharedStrings.xml"/>
</Relationships>`
	if err := addFile(zw, "xl/_rels/workbook.xml.rels", wbRelsXML); err != nil {
		return err
	}

	// 5. xl/sharedStrings.xml
	if err := p.writeSharedStrings(zw); err != nil {
		return err
	}

	// 6. xl/worksheets/sheet1.xml
	if err := p.writeSheet(zw); err != nil {
		return err
	}

	return zw.Close()
}

func (p *xlsxPackage) writeSharedStrings(zw *zip.Writer) error {
	fw, err := zw.Create("xl/sharedStrings.xml")
	if err != nil {
		return err
	}
	_, _ = fw.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"))
	_, _ = fmt.Fprintf(fw, `<sst xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" count="%d" uniqueCount="%d">`, len(p.SharedStrings), len(p.SharedStrings))
	for _, s := range p.SharedStrings {
		_, _ = fmt.Fprintf(fw, `<si><t xml:space="preserve">%s</t></si>`, escapeXML(s))
	}
	_, _ = fw.Write([]byte(`</sst>`))
	return nil
}

func (p *xlsxPackage) writeSheet(zw *zip.Writer) error {
	fw, err := zw.Create("xl/worksheets/sheet1.xml")
	if err != nil {
		return err
	}
	_, _ = fw.Write([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"))
	_, _ = fw.Write([]byte(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`))

	// Header row (row 1): shared-string references for column names
	_, _ = fmt.Fprintf(fw, `<row r="1">`)
	for colIdx, header := range p.SharedStrings {
		// Find the index in shared strings
		idx := 0
		for i, s := range p.SharedStrings {
			if s == header {
				idx = i
				break
			}
		}
		cell := fmt.Sprintf(`<c r="%s1" t="s"><v>%d</v></c>`, colLetter(colIdx), idx)
		_, _ = fw.Write([]byte(cell))
	}
	_, _ = fw.Write([]byte(`</row>`))

	// Data rows
	for rowIdx, marker := range p.Markers {
		rowNum := rowIdx + 2 // row 1 is header
		payload, _ := trackjson.MakeMarkerPayload(marker)

		// Build cell values: column order matches SharedStrings
		// id, time_utc, lat, lon, altitude_m, dose_rate_usvh, dose_rate_uroenth, count_rate_cps, speed_ms, speed_kmh, temperature_c, humidity_pct, detector
		vals := []string{
			fmt.Sprintf("%d", payload.ID),    // id (number)
			payload.TimeUTC,                  // time_utc (string)
			fmt.Sprintf("%.6f", payload.Lat), // lat (number)
			fmt.Sprintf("%.6f", payload.Lon), // lon (number)
		}

		// altitude (optional)
		if marker.AltitudeValid {
			vals = append(vals, fmt.Sprintf("%.1f", marker.Altitude))
		} else {
			vals = append(vals, "")
		}

		vals = append(vals,
			fmt.Sprintf("%.4f", payload.DoseRateMicroSvH),
			fmt.Sprintf("%.4f", payload.DoseRateMicroRoentgenH),
			fmt.Sprintf("%.2f", payload.CountRateCPS),
			fmt.Sprintf("%.2f", payload.SpeedMS),
			fmt.Sprintf("%.2f", payload.SpeedKMH),
		)

		// temperature (optional)
		if marker.TemperatureValid {
			vals = append(vals, fmt.Sprintf("%.1f", marker.Temperature))
		} else {
			vals = append(vals, "")
		}

		// humidity (optional)
		if marker.HumidityValid {
			vals = append(vals, fmt.Sprintf("%.1f", marker.Humidity))
		} else {
			vals = append(vals, "")
		}

		vals = append(vals, payload.DetectorName)

		_, _ = fmt.Fprintf(fw, `<row r="%d">`, rowNum)
		for colIdx, val := range vals {
			// Decide cell type: string or number
			if val == "" {
				// Empty cell
				_, _ = fmt.Fprintf(fw, `<c r="%s%d"/>`, colLetter(colIdx), rowNum)
			} else if isNumericCell(colIdx) {
				_, _ = fmt.Fprintf(fw, `<c r="%s%d"><v>%s</v></c>`, colLetter(colIdx), rowNum, escapeXML(val))
			} else {
				// Use shared string
				idx := findSharedString(p.SharedStrings, val)
				if idx >= 0 {
					_, _ = fmt.Fprintf(fw, `<c r="%s%d" t="s"><v>%d</v></c>`, colLetter(colIdx), rowNum, idx)
				} else {
					// Fallback: inline string
					_, _ = fmt.Fprintf(fw, `<c r="%s%d" t="inlineStr"><is><t>%s</t></is></c>`, colLetter(colIdx), rowNum, escapeXML(val))
				}
			}
		}
		_, _ = fw.Write([]byte(`</row>`))
	}

	_, _ = fw.Write([]byte(`</sheetData></worksheet>`))
	return nil
}

func isNumericCell(colIdx int) bool {
	// Numeric columns: id(0), lat(2), lon(3), altitude(4), dose_usvh(5), dose_uroenth(6), cps(7), speed_ms(8), speed_kmh(9), temp(10), humidity(11)
	return colIdx == 0 || (colIdx >= 2 && colIdx <= 11)
}

func findSharedString(headers []string, val string) int {
	for i, s := range headers {
		if s == val {
			return i
		}
	}
	return -1
}

func colLetter(idx int) string {
	// 0→A, 1→B, ..., 25→Z, 26→AA, ...
	result := ""
	for idx >= 0 {
		result = string(rune('A'+idx%26)) + result
		idx = idx/26 - 1
	}
	return result
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func addFile(zw *zip.Writer, name, content string) error {
	fw, err := zw.Create(name)
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(content))
	return err
}

// Suppress unused import warning for xml package (used in struct definitions above)
var _ = xml.Name{}
