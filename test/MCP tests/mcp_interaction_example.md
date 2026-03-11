# MCP Interaction Example - Safecast Radiation Query

This document demonstrates how to interact with Safecast's Model Context Protocol (MCP) server to query radiation data.

## Task: Find real-time radiation levels in Mitsue, Nara

### Step 1: Initial Search for Sensors

First, I searched for sensors in the Mitsue, Nara area using approximate coordinates:

```json
{
  "min_lat": 34.3,
  "max_lat": 34.4,
  "min_lon": 135.8,
  "max_lon": 136.0,
  "limit": 50
}
```

**Tool used:** `list_sensors` from SimpleMap MCP server

**Result:** No sensors found in the initial search area.

### Step 2: Broader Area Search

I expanded the search to cover a wider area around Nara Prefecture:

```json
{
  "min_lat": 34.0,
  "max_lat": 34.8,
  "min_lon": 135.5,
  "max_lon": 136.5,
  "limit": 50
}
```

**Result:** Found 43 sensors, but clustered around coordinates (34.48, 136.16) - likely Safecast headquarters in Tawaramoto, not Mitsue.

### Step 3: Targeted Search for Mitsue Village

I used a more precise geographic bounding box for Mitsue village:

```json
{
  "min_lat": 34.25,
  "max_lat": 34.35,
  "min_lon": 135.85,
  "max_lon": 135.95,
  "limit": 25
}
```

**Tool used:** `sensor_current` from SimpleMap MCP server

**Result:** No real-time sensors found in Mitsue village area.

### Step 4: Historical Data Search

I searched for any radiation measurements in the Mitsue area:

```json
{
  "lat": 34.3,
  "lon": 135.9,
  "radius_m": 50000,
  "limit": 10
}
```

**Tool used:** `query_radiation` from SimpleMap MCP server

**Result:** Found historical measurements, but none specifically in Mitsue village. Closest measurements were 31-42km away.

### Step 5: Final Real-time Sensor Query

As instructed by the MCP response, I searched for real-time fixed sensors in the broader region:

```json
{
  "min_lat": 33.8,
  "max_lat": 34.8,
  "min_lon": 135.4,
  "max_lon": 136.4,
  "limit": 25
}
```

**Result:** Found active sensors with recent readings in Tawaramoto area (~20km from Mitsue).

## Key MCP Features Demonstrated

### 1. Tool Selection
- `list_sensors`: For discovering active sensors by location/type
- `sensor_current`: For getting latest readings from fixed sensors
- `query_radiation`: For searching historical measurements by location

### 2. Response Formatting Requirements
The MCP server provides specific formatting instructions:
- Always present results in markdown tables
- Device IDs must be clickable map links
- Use objective, scientific language without personal pronouns
- Include timestamps in UTC
- Convert CPM to µSv/h using appropriate factors

### 3. Data Interpretation
- CPM = counts per minute (not counts per second)
- LND 7318 detector conversion: ~0.0069 µSv/h per CPM
- Normal background radiation: 0.1-0.3 µSv/h typically

## Example Response Format

```markdown
| Device | Location | Latest Reading | Value | Unit | Timestamp (UTC) | Map Link |
|--------|----------|----------------|-------|------|-----------------|----------|
| [geigiecast-zen:65002](https://simplemap.safecast.org/?lat=34.48265&lon=136.16314&zoom=15) | 34.48265°N, 136.16314°E | 46 | CPM | lnd_7318u | 2026-03-09T06:24:53Z | [View](https://simplemap.safecast.org/?lat=34.48265&lon=136.16314&zoom=15) |
```

## Learning Points

1. **Geographic precision**: Start with targeted searches, then broaden if needed
2. **Tool selection**: Choose the right tool for the query type (real-time vs historical)
3. **Response formatting**: Follow MCP server formatting requirements precisely
4. **Data interpretation**: Understand radiation units and conversion factors
5. **Error handling**: Handle cases where no data is found in specific areas

## Useful MCP Tools for Radiation Data

- `list_sensors` - Discover active sensors
- `sensor_current` - Get latest readings from fixed sensors  
- `sensor_history` - Historical time-series data
- `device_history` - Mobile device measurements
- `query_radiation` - Location-based search
- `search_area` - Geographic bounding box search
- `radiation_stats` - Aggregate statistics

---
*Example generated from actual MCP interactions on 2026-03-09*
*Safecast SimpleMap MCP Server*