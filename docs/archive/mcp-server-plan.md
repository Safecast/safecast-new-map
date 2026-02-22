# MCP Server Plan for Safecast Map

**Status:** ✅ IMPLEMENTED
**Date Archived:** 2026-02-22
**Implementation:** [safecast-map-MCP repository](https://github.com/Safecast/safecast-map-MCP)
**Live Endpoint:** https://simplemap.safecast.org/mcp-http
**Documentation:** See the [safecast-map-MCP README](https://github.com/Safecast/safecast-map-MCP#readme)

---

## Overview

An MCP (Model Context Protocol) server that allows AI assistants (Claude, etc.) to query the Safecast radiation map data through structured tools. Users chatting with an AI can ask questions like "What's the radiation level near Fukushima?" and the AI will call MCP tools to get real answers from the database.

**No authentication required** - all tools are read-only queries against public data.

---

## Architecture

### Recommended: Standalone MCP Server

A separate Go binary that connects to the same PostgreSQL database and exposes MCP tools over `stdio` or `SSE` transport. Keeps concerns separated and can be deployed independently.

```
┌──────────────┐     ┌──────────────┐     ┌────────────┐
│  AI Client   │────▶│  MCP Server  │────▶│ PostgreSQL │
│  (Claude)    │◀────│  (Go binary) │◀────│  Database  │
└──────────────┘     └──────────────┘     └────────────┘
                      Uses pkg/database/
```

### Alternative: Embedded in Existing Server

Add MCP endpoints alongside the existing HTTP API. Simpler deployment but couples the MCP protocol to the map server.

---

## Tools to Expose

### 1. `query_radiation` - Query radiation levels near a location

- **Description:** Find radiation measurements near a geographic location
- **Inputs:**
  - `lat` (float, required) - Latitude
  - `lon` (float, required) - Longitude
  - `radius_m` (int, optional, default 1500, range 25-50000) - Search radius in meters
  - `limit` (int, optional, default 25, max 200) - Maximum results
- **Returns:** Nearest measurements with dose rates (µSv/h), count rates (CPS), dates, detector info, coordinates
- **Maps to:** Existing `StreamLatestMarkersNear()` database method

### 2. `get_track` - Get measurements from a specific track

- **Description:** Retrieve all radiation measurements recorded during a specific track/journey
- **Inputs:**
  - `track_id` (string, required) - Track identifier
  - `from` (int, optional) - Start marker ID for pagination
  - `to` (int, optional) - End marker ID for pagination
- **Returns:** All markers in the track with dose rates, coordinates, timestamps, detector info, speed, altitude
- **Maps to:** Existing track query + `StreamMarkersInRange()`

### 3. `list_tracks` - Browse available tracks

- **Description:** List available radiation measurement tracks, optionally filtered by date
- **Inputs:**
  - `year` (int, optional) - Filter by year
  - `month` (int, optional) - Filter by month (requires year)
  - `limit` (int, optional, default 50) - Maximum results
- **Returns:** Track summaries with track ID, marker count, date ranges, API URLs
- **Maps to:** `StreamTrackSummaries()` and date-range variants

### 4. `get_spectrum` - Get gamma spectroscopy data

- **Description:** Retrieve gamma spectrum data for a specific measurement point
- **Inputs:**
  - `marker_id` (int, required) - Marker identifier
- **Returns:** Channel data array, energy range (keV), calibration coefficients, live/real time, device model, source format
- **Maps to:** `GetSpectrum()`

### 5. `search_area` - Search for measurements in a bounding box

- **Description:** Find radiation measurements within a geographic bounding box
- **Inputs:**
  - `min_lat` (float, required) - Southern boundary
  - `max_lat` (float, required) - Northern boundary
  - `min_lon` (float, required) - Western boundary
  - `max_lon` (float, required) - Eastern boundary
  - `limit` (int, optional, default 100, max 200) - Maximum results
- **Returns:** Markers within the bounding box with full metadata
- **Maps to:** Existing `get_markers` tile query logic

### 6. `radiation_info` - General radiation reference information

- **Description:** Get reference information about radiation units, safety levels, and terminology
- **Inputs:**
  - `topic` (string, required) - One of: `units`, `dose_rates`, `safety_levels`, `detectors`, `background_levels`, `isotopes`
- **Returns:** Static reference content about radiation science
- **Implementation:** Hardcoded reference content (no database query needed)

### 7. `device_history` - Realtime device measurement history

- **Description:** Get historical measurements from a specific monitoring device
- **Inputs:**
  - `device_id` (string, required) - Device identifier
  - `days` (int, optional, default 30) - Number of days of history
- **Returns:** Historical measurements with values, units, coordinates, timestamps
- **Maps to:** `realtime_measurements` table queries

---

## Code Structure

### Standalone Binary (Recommended)

```
cmd/mcp-server/
  main.go              # Entry point, flag parsing, DB connection
  tools.go             # Tool definitions and handlers
  reference.go         # Static radiation reference data
```

### Embedded in Existing Server

```
pkg/mcp/
  server.go            # MCP server setup
  tools.go             # Tool handlers
  reference.go         # Static reference data
```

---

## Implementation Details

### MCP SDK

Use [github.com/mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) - the most established Go MCP SDK. Handles JSON-RPC protocol, tool registration, and transport (stdio/SSE).

### Database Connection

Reuse the existing `pkg/database` package directly. The database layer already has clean interfaces (`DatabaseOperations`) and streaming query methods. The MCP server imports this package and connects to the same PostgreSQL instance.

### Rate Limiting

MCP requests should go through the existing workload-aware rate limiter, either as `RequestGeneral` or a new `RequestMCP` lane to prevent abuse.

### Streaming Results

The database layer uses channel-based streaming (`<-chan T, <-chan error`). MCP tool handlers collect results from channels and return them as JSON responses.

### Result Size Limits

MCP tool responses must be bounded. Cap results (e.g., max 200 markers per query) to avoid overwhelming the AI context window.

### Transport Options

| Transport | Use Case |
|-----------|----------|
| **stdio** | Local use - user runs the binary alongside Claude Desktop/CLI |
| **SSE (HTTP)** | Remote/public use - hosted server that any AI client can connect to |

For public access without login, **SSE transport** is recommended - host it alongside the map server.

---

## Dependencies

```
github.com/mark3labs/mcp-go  # MCP protocol SDK
```

All other dependencies (database drivers, streaming, spatial queries) already exist in the codebase.

---

## Data Models (Response Shapes)

### Marker (from `query_radiation`, `get_track`, `search_area`)

```json
{
  "id": 12345,
  "doseRate": 0.15,
  "countRate": 450.5,
  "date": 1609459200,
  "lat": 44.08832,
  "lon": 42.97577,
  "altitude": 150.25,
  "speed": 15.5,
  "trackID": "2021-01-01-track",
  "detector": "RadiaCode-102",
  "radiation": "gamma",
  "hasSpectrum": true
}
```

### Track Summary (from `list_tracks`)

```json
{
  "trackID": "2021-01-01-track",
  "firstID": 100,
  "lastID": 1500,
  "markerCount": 1401,
  "index": 42,
  "apiURL": "/api/track/2021-01-01-track.json"
}
```

### Spectrum (from `get_spectrum`)

```json
{
  "id": 999,
  "markerID": 12345,
  "channels": [10, 15, 20, "..."],
  "channelCount": 1024,
  "energyMinKeV": 0,
  "energyMaxKeV": 3000,
  "liveTimeSec": 60.5,
  "realTimeSec": 65.0,
  "deviceModel": "RadiaCode-102",
  "calibration": { "a": 0.0, "b": 2.93, "c": 0.0 },
  "sourceFormat": "rctrk"
}
```

### Realtime Measurement (from `device_history`)

```json
{
  "id": 555,
  "deviceID": "SC-98765",
  "value": 0.12,
  "unit": "µSv/h",
  "lat": 44.08832,
  "lon": 42.97577,
  "measuredAt": 1609459200,
  "deviceName": "Mobile Unit B",
  "tube": "SBM-20",
  "country": "JP"
}
```

---

## Example User Interactions

Once deployed, users talking to an AI assistant could ask:

| User Question | MCP Tool Called |
|---------------|-----------------|
| "What are the radiation levels near Chernobyl?" | `query_radiation(lat=51.39, lon=30.10, radius_m=5000)` |
| "Show me tracks recorded in Japan in 2024" | `list_tracks(year=2024)` |
| "What does 0.5 µSv/h mean? Is that dangerous?" | `radiation_info(topic="safety_levels")` |
| "Get the gamma spectrum for marker 12345" | `get_spectrum(marker_id=12345)` |
| "What radiation data is available in the Tokyo area?" | `search_area(min_lat=35.5, max_lat=35.8, min_lon=139.5, max_lon=139.9)` |
| "Show me the last 7 days from device SC-12345" | `device_history(device_id="SC-12345", days=7)` |

---

## Estimated Scope

| Component | Estimate |
|-----------|----------|
| MCP protocol layer | Handled by `mcp-go` SDK |
| Tool definitions | ~7 tools, each ~50-100 lines |
| Database glue | Thin wrappers around existing `pkg/database` methods |
| Reference data | Static content about radiation science |
| Main/config | Flag parsing, DB connection setup |
| **Total new code** | **~800-1200 lines of Go** |

---

## Startup Command (Example)

```bash
./safecast-mcp-server \
  -db-type pgx \
  -db-conn "postgres://postgres:PASSWORD@127.0.0.1:5432/safecast?sslmode=disable" \
  -transport sse \
  -port 8766
```

Or for stdio (local Claude Desktop use):

```bash
./safecast-mcp-server \
  -db-type pgx \
  -db-conn "postgres://postgres:PASSWORD@127.0.0.1:5432/safecast?sslmode=disable" \
  -transport stdio
```
