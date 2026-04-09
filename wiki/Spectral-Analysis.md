# Spectral Analysis

Guide to gamma spectroscopy, spectrum analysis, and isotope identification.

[[Home|← Back to Home]]

---

## Overview

The platform supports gamma spectrum analysis for identifying radioactive isotopes:

- **Spectrum upload** - Import spectrum files in multiple formats
- **Spectrum storage** - Store spectrum data with calibration
- **Spectrum visualization** - Interactive spectrum viewer
- **Isotope identification** - Automatic peak detection and isotope matching
- **Spectrum export** - Download in various formats

---

## Supported Spectrum Formats

### Maestro (.spe)

Standard spectrum format from ORTEC Maestro software.

**Characteristics:**
- Text-based format
- Contains channel data
- Includes calibration coefficients
- Widely supported

### ANSI N42.42 (.n42)

Industry standard for radiation data exchange.

**Characteristics:**
- XML-based format
- Rich metadata
- Standardized structure
- Interoperable

### RadiaCode (.rctrk, .rcxml, .xml)

RadiaCode device formats with embedded spectra.

**Characteristics:**
- JSON or XML format
- Includes GPS data
- Real-time capable
- Device-specific calibration

---

## Spectrum Database Schema

### Spectra Table

```sql
CREATE TABLE spectra (
  id BIGSERIAL PRIMARY KEY,
  marker_id BIGINT REFERENCES markers(id),
  channels JSON,              -- Channel data array
  channel_count INTEGER,      -- Number of channels (typically 1024 or 2048)
  energy_min_kev DOUBLE PRECISION,  -- Minimum energy (keV)
  energy_max_kev DOUBLE PRECISION,  -- Maximum energy (keV)
  live_time_sec DOUBLE PRECISION,   -- Live time (seconds)
  real_time_sec DOUBLE PRECISION,   -- Real time (seconds)
  device_model TEXT,          -- Device model name
  calibration JSON,           -- Calibration coefficients {A, B, C}
  source_format TEXT,         -- Original format (spe, n42, rctrk)
  filename TEXT,              -- Original filename
  raw_data BYTEA,             -- Original file data
  created_at TIMESTAMPTZ
);
```

### Calibration Coefficients

Quadratic energy calibration:

```
Energy (keV) = A + B × channel + C × channel²
```

**Calibration JSON:**
```json
{
  "a": 0.5,
  "b": 0.3,
  "c": 0.0001
}
```

---

## Upload Spectrum Data

### Via Web Interface

1. Log in to your account
2. Upload spectrum file (.spe, .n42, .rctrk, etc.)
3. System auto-detects format
4. Spectrum extracted and stored
5. Linked to measurement marker

### Via API

```bash
curl -X POST http://localhost:8765/api/upload \
  -H "X-API-Key: your-api-key" \
  -F "file=@spectrum.spe" \
  -F "source=manual"
```

### Automatic Extraction

During upload:
- Spectrum data extracted automatically
- Calibration coefficients parsed
- Channel data stored as JSON
- Linked to parent measurement

---

## Spectrum Viewer

### Access Spectrum

1. Navigate to map
2. Click marker with spectrum data
3. Click "View Spectrum" button
4. Interactive spectrum viewer opens

### Viewer Features

**Interactive Chart:**
- Canvas 2D rendering
- Smooth rendering for large datasets
- Efficient memory usage

**Channel Display:**
- Raw channel counts on Y-axis
- Channel number on X-axis
- Hover for exact values

**Energy Conversion:**
- Automatic channel-to-energy conversion
- Energy scale on secondary X-axis
- Calibration coefficients displayed

### Zoom and Pan

**Zoom:**
- Drag to select region
- Double-click to zoom in
- Right-click to zoom out
- Mouse wheel for continuous zoom

**Pan:**
- Click and drag to pan
- Arrow keys for fine control

### Isotope Peak Overlays

**Available Isotopes:**
- Cs-137 (661.7 keV)
- Co-60 (1173.2, 1332.5 keV)
- K-40 (1460.8 keV)
- I-131 (364.5 keV)
- U-238 series peaks
- Th-232 series peaks
- And more...

**Toggle Peaks:**
- Click isotope list button
- Select isotopes to display
- Peaks overlaid on spectrum
- Labels show energy and isotope

### Peak Detection

**Automatic Detection:**
- System identifies significant peaks
- Matches to known isotopes
- Displays confidence level

**Manual Peak Selection:**
- Click on spectrum to mark peak
- Enter energy manually
- Search isotope database

---

## Spectrum Export

### Export Formats

**JSON:**
```json
{
  "markerId": 12345,
  "channels": [100, 120, 95, ...],
  "channelCount": 1024,
  "energyMinKev": 0,
  "energyMaxKev": 3000,
  "calibration": {
    "a": 0.5,
    "b": 0.3,
    "c": 0.0001
  },
  "deviceModel": "RadiaCode-102"
}
```

**CSV:**
```csv
channel,count,energy_kev
0,100,0.5
1,120,0.8
2,95,1.1
...
```

**ANSI N42.42 (.n42):**
- Standard XML format
- Full metadata
- Interoperable

**Maestro (.spe):**
- ORTEC Maestro format
- Legacy compatibility
- Wide software support

### Export via API

```bash
# JSON format
GET /api/spectrum/{markerId}/download?format=json

# CSV format
GET /api/spectrum/{markerId}/download?format=csv

# N42 format
GET /api/spectrum/{markerId}/download?format=n42

# SPE format
GET /api/spectrum/{markerId}/download?format=spe
```

### Download from Viewer

1. Open spectrum viewer
2. Click "Download" button
3. Select format
4. File downloads automatically

---

## Isotope Identification

### Known Isotope Peaks

The system maintains a database of common isotopes:

| Isotope | Half-life | Energy (keV) | Source |
|---------|-----------|--------------|--------|
| **Cs-137** | 30.17 years | 661.7 | Fission product |
| **Co-60** | 5.27 years | 1173.2, 1332.5 | Activation product |
| **K-40** | 1.25 billion years | 1460.8 | Natural |
| **I-131** | 8.02 days | 364.5 | Fission product |
| **U-238** | 4.47 billion years | Various | Natural decay series |
| **Th-232** | 14.05 billion years | Various | Natural decay series |

### Peak Matching Algorithm

1. **Detect peaks** - Find local maxima above threshold
2. **Calculate energy** - Convert channel to keV using calibration
3. **Search database** - Match energy to known isotopes
4. **Calculate confidence** - Based on peak shape and intensity
5. **Display results** - Show matched isotopes with confidence

### Manual Identification

For unknown peaks:
1. Mark peak position
2. View energy value
3. Search isotope database
4. Consider context (location, history)
5. Consult expert if needed

---

## Spectrum Analysis Tools

### Via MCP Server

Use AI-powered spectrum analysis:

```bash
# Connect Claude to Safecast MCP
claude mcp add --transport http safecast https://simplemap.safecast.org/mcp-http

# Ask Claude to analyze spectrum
claude "Analyze the spectrum at marker 12345 and identify isotopes"
```

### MCP Spectrum Tools

**`get_spectrum`:**
```json
{
  "markerId": 12345
}
```

**`list_spectra`:**
```json
{
  "page": 1,
  "limit": 20
}
```

### Analysis Workflow

1. **Retrieve spectrum** via MCP tool
2. **Detect peaks** automatically
3. **Match isotopes** from database
4. **Generate report** with findings
5. **Provide recommendations**

---

## Spectrum Migration

### Add Spectrum Support

If your database doesn't have spectrum tables:

**PostgreSQL:**
```bash
psql -d your_database -f migrations/add_spectrum_support.sql
```

**SQLite:**
```bash
sqlite3 data.db < migrations/add_spectrum_support_sqlite.sql
```

**DuckDB:**
```bash
duckdb data.duckdb < migrations/add_spectrum_support_duckdb.sql
```

See [SPECTRAL_MIGRATION_GUIDE.md](/SPECTRAL_MIGRATION_GUIDE.md) for detailed instructions.

### Migration Script Contents

The migration adds:
- `spectra` table
- Indexes on `marker_id`
- Calibration support
- Raw data storage

---

## Spectrum Architecture

### Data Flow

```
Upload File
    ↓
Parse Format (.spe, .n42, .rctrk)
    ↓
Extract Spectrum Data
    ↓
Store in Database
    ↓
Link to Marker
    ↓
Display in Viewer
```

### Components

**Parsers:**
- `pkg/spectrum/` - Spectrum parsing library
- Format-specific parsers
- Calibration extraction
- Channel data processing

**Storage:**
- JSON for channel data
- BYTEA for raw files
- Efficient compression
- Indexed for fast retrieval

**Visualization:**
- Canvas 2D rendering
- Efficient large dataset handling
- Interactive controls
- Isotope overlays

See [spectral-data-flow.mmd](/docs/spectral-data-flow.mmd) for detailed architecture diagram.

---

## Best Practices

### Data Quality

1. **Use proper calibration** for your detector
2. **Record live time** accurately
3. **Include device model** for compatibility
4. **Store raw files** for future reference

### Analysis

1. **Check calibration** before analysis
2. **Consider background subtraction**
3. **Account for detector efficiency**
4. **Verify peak identification** manually

### Performance

1. **Use PostgreSQL** for large spectrum datasets
2. **Index by marker_id** for fast retrieval
3. **Compress channel data** (JSON compression)
4. **Cache frequently accessed spectra**

---

## Troubleshooting

### Spectrum Not Displaying

**Check spectrum exists:**
```bash
curl http://localhost:8765/api/spectrum/12345
```

**Verify calibration:**
```sql
SELECT calibration FROM spectra WHERE marker_id = 12345;
```

### Upload Fails

**Check file format:**
```bash
# Verify SPE file
head -20 spectrum.spe

# Check N42 structure
xmllint --format spectrum.n42
```

**Check logs:**
```bash
grep "spectrum" /var/log/safecast.log
```

### Incorrect Energy Scale

**Verify calibration coefficients:**
```sql
SELECT calibration FROM spectra WHERE id = 1;
```

**Recalibrate if needed:**
```sql
UPDATE spectra
SET calibration = '{"a": 0.5, "b": 0.3, "c": 0.0001}'
WHERE marker_id = 12345;
```

---

## See Also

- [Map Features](Map-Features) - Spectrum viewer integration
- [Data Import & Export](Data-Import-Export) - Spectrum file formats
- [API Documentation](API-Documentation) - Spectrum API endpoints
- [MCP Server & AI Integration](MCP-Server-AI-Integration) - AI spectrum analysis
- [SPECTRAL_MIGRATION_GUIDE.md](/SPECTRAL_MIGRATION_GUIDE.md) - Migration guide
- [spectral-data-flow.mmd](/docs/spectral-data-flow.mmd) - Architecture diagram
