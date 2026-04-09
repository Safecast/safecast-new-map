# MCP Server & AI Integration

Guide to using the Model Context Protocol (MCP) server and AI-powered features.

[[Home|← Back to Home]]

---

## Overview

The platform includes an **MCP (Model Context Protocol)** server that provides AI models like Claude with tools to query and analyze radiation data. The server integrates seamlessly with the unified binary.

### Components

| Component | Endpoint | Description |
|-----------|----------|-------------|
| Unified Server | `/` | Map UI and REST API |
| MCP Server | `/mcp-http`, `/mcp/sse` | Model Context Protocol |
| Web Chat | `/assistant/` | Claude-powered chat interface |
| Swagger Docs | `/map-api/`, `/mcp-api/` | API documentation |

---

## Quick Start

### Enable MCP Server

The MCP server is included in the unified server binary. No additional configuration required.

```bash
./safecast-new-map
```

Access points:
- **MCP HTTP:** `http://localhost:8765/mcp-http`
- **MCP SSE:** `http://localhost:8765/mcp/sse`

### Enable Web Chat

```bash
DATABASE_URL="postgres://user:pass@localhost/db" \
ANTHROPIC_API_KEY="your-anthropic-key" \
CLAUDE_MODEL="claude-haiku-4-5-20251001" \
./safecast-new-map
```

Access web chat: `http://localhost:8765/assistant/`

---

## MCP Tools

The server provides **16 tools** for AI interaction with radiation data:

### Data Query Tools

#### `query_radiation`

Find radiation measurements near coordinates.

**Parameters:**
```json
{
  "lat": 35.6762,
  "lon": 139.6503,
  "radius": 1000
}
```

**Features:**
- PostGIS spatial queries (ST_DWithin)
- Bounding box pre-filtering
- Returns measurements with dose rates, timestamps, GPS coordinates

#### `search_area`

Search radiation data within a geographic area.

**Parameters:**
```json
{
  "minLat": 35.6,
  "minLon": 139.6,
  "maxLat": 35.7,
  "maxLon": 139.7
}
```

#### `query_extreme_readings`

Find extreme (highest) radiation readings in the database.

**Parameters:**
```json
{
  "limit": 10
}
```

### Track Tools

#### `list_tracks`

List measurement tracks with pagination.

**Parameters:**
```json
{
  "page": 1,
  "limit": 50,
  "location": "Tokyo"  // Optional location filter
}
```

#### `get_track`

Get specific track data.

**Parameters:**
```json
{
  "trackId": "abc123"
}
```

#### `search_tracks_by_location`

Search tracks by geographic location.

**Parameters:**
```json
{
  "query": "Tokyo",
  "limit": 20
}
```

#### `device_history`

Get device measurement history.

**Parameters:**
```json
{
  "deviceId": "device-123",
  "limit": 100
}
```

### Spectrum Tools

#### `get_spectrum`

Get gamma spectrum data for isotope analysis.

**Parameters:**
```json
{
  "markerId": 12345
}
```

**Returns:**
- Channel data (JSON array)
- Calibration coefficients
- Energy range
- Device model

#### `list_spectra`

List all available spectra.

**Parameters:**
```json
{
  "page": 1,
  "limit": 20
}
```

### Real-Time Sensor Tools

#### `list_sensors`

List all real-time sensors.

**Returns:** Active sensors with current readings.

#### `sensor_current`

Get current reading for a specific sensor.

**Parameters:**
```json
{
  "sensorId": "sensor-123"
}
```

#### `sensor_history`

Get historical data for a sensor.

**Parameters:**
```json
{
  "sensorId": "sensor-123",
  "hours": 24
}
```

### Reference Tools

#### `radiation_info`

Get reference radiation information.

**Parameters:**
```json
{
  "topic": "background"  // background, units, safety
}
```

#### `db_info`

Get database information and statistics.

**Returns:**
- Total markers
- Date range
- Database type
- Performance metrics

### Analytics Tools (DuckLake)

#### `query_analytics`

Query analytics via DuckLake (requires DuckLake configuration).

**Parameters:**
```json
{
  "query": "SELECT COUNT(*) FROM markers WHERE doseRate > 0.5"
}
```

#### `query_duckdb_logs`

Query DuckDB analytics logs.

**Parameters:**
```json
{
  "query": "SELECT * FROM query_log ORDER BY timestamp DESC LIMIT 10"
}
```

### Utility Tools

#### `ping`

Health check tool.

**Returns:** Server status and version.

#### `top_uploaders`

List top data contributors.

**Parameters:**
```json
{
  "limit": 10
}
```

---

## Connecting Claude to Safecast

### Claude Code CLI

```bash
# Add Safecast MCP server to Claude Code
claude mcp add --transport http safecast https://simplemap.safecast.org/mcp-http

# Use in Claude Code
claude "What's the radiation level in Tokyo?"
```

### Claude.ai Web Interface

1. Go to **Settings** → **Integrations**
2. Click **Add custom integration**
3. Enter URL: `https://simplemap.safecast.org/mcp-http`
4. Save and start using

### Custom GPT Actions

The platform includes GPT-optimized endpoints for Custom GPT Actions:
- Compact responses
- Minimal payload size
- Essential fields only

Register your instance in the Custom GPT configuration.

---

## Web Chat Interface

Access the AI-powered chat at: `http://localhost:8765/assistant/`

### Features

- **Natural language queries** - Ask questions like "What's the radiation level near Fukushima?"
- **Interactive maps** - Results include map visualizations
- **Spectrum analysis** - Upload and analyze gamma spectra
- **Historical data** - Query trends and patterns
- **Real-time monitoring** - Check live sensor readings

### Example Queries

```
"What's the current radiation level in Tokyo?"
"Show me the highest readings in Japan"
"Find tracks near Chernobyl"
"What isotopes are in this spectrum?"
"Compare radiation levels in 2023 vs 2024"
"List all sensors currently online"
```

### Chat Logging

All chat questions are logged to DuckLake for analytics:
- Question text
- Timestamp
- Tools used
- Response metadata

---

## Model Adapter

The platform includes a model adapter for model-specific optimizations:

### Features

- **Model-specific hints** - Tailored prompts for different models
- **Response formatting** - Optimized for each model's strengths
- **Tool selection** - Smart tool choice based on query type

### Configuration

```bash
CLAUDE_MODEL="claude-haiku-4-5-20251001"
```

Supported models:
- Claude Haiku 4.5 (recommended)
- Claude Sonnet
- Claude Opus

---

## Semantic Caching with RAG

The platform implements semantic caching for improved performance:

### How It Works

1. **Query embedding** - Convert query to vector
2. **Similarity search** - Find cached similar queries
3. **Cache hit** - Return cached response if similar enough
4. **Cache miss** - Execute query and cache result

### Configuration

Semantic cache is enabled by default when DuckLake is configured.

**Cache settings:**
- Similarity threshold: 0.95
- TTL: 24 hours
- Max cache size: 10,000 entries

### RAG Integration

See [rag-semantic-cache.md](/docs/rag-semantic-cache.md) for architecture details.

---

## DuckLake Analytics

DuckLake enables advanced analytics with shared tables:

### Setup

```bash
DATABASE_URL="postgres://user:pass@localhost/safecast" \
DUCKLAKE_PG_URL="dbname=ducklake_catalog host=localhost user=ducklake_rw" \
DUCKLAKE_DATA_PATH="/var/lib/safecast/ducklake/" \
./safecast-new-map
```

### Shared Tables

| Table | Description |
|-------|-------------|
| `chat_questions` | AI chat question log |
| `mcp_query_log` | MCP tool usage log |
| `mcp_ai_query_log` | AI query execution log |

### Query Analytics

```sql
-- Most popular queries
SELECT query, COUNT(*) as count
FROM chat_questions
GROUP BY query
ORDER BY count DESC
LIMIT 10;

-- Tool usage
SELECT tool_name, COUNT(*) as executions
FROM mcp_query_log
GROUP BY tool_name
ORDER BY executions DESC;
```

---

## AI Query Logging

All AI queries are logged for monitoring and analytics:

### Logged Information

- Query text
- Timestamp
- User session
- Tools invoked
- Response time
- Cache status (hit/miss)

### Access Logs

Via admin panel: `/admin/mcp`

Or query directly:
```sql
SELECT * FROM mcp_ai_query_log
ORDER BY timestamp DESC
LIMIT 100;
```

---

## MCP Transports

The server supports multiple MCP transports:

### Streamable HTTP

```
POST /mcp-http
Content-Type: application/json

{
  "method": "query_radiation",
  "params": {
    "lat": 35.6762,
    "lon": 139.6503
  }
}
```

### SSE (Server-Sent Events)

```
GET /mcp/sse
```

Provides real-time streaming for long-running queries.

---

## Best Practices

### Performance

1. **Use semantic caching** for repeated queries
2. **Limit result sizes** with pagination
3. **Use bounding box queries** for large areas
4. **Enable DuckLake** for analytics queries

### Security

1. **Rate limit MCP endpoints**
2. **Monitor query logs** for suspicious patterns
3. **Use authentication** for sensitive data
4. **Sanitize user inputs** in chat interface

### Integration

1. **Use HTTP transport** for web applications
2. **Enable SSE** for real-time updates
3. **Test with Claude Code** before production

---

## Troubleshooting

### MCP Server Not Responding

**Check logs:**
```bash
grep "MCP" /var/log/safecast.log
```

**Verify endpoint:**
```bash
curl http://localhost:8765/mcp-http -X POST \
  -H "Content-Type: application/json" \
  -d '{"method": "ping"}'
```

### Web Chat Not Loading

**Check environment variables:**
```bash
echo $ANTHROPIC_API_KEY
echo $CLAUDE_MODEL
```

**Verify database connection:**
```bash
echo $DATABASE_URL
```

### Slow AI Responses

**Enable caching:**
```bash
# Ensure DuckLake is configured
echo $DUCKLAKE_PG_URL
echo $DUCKLAKE_DATA_PATH
```

**Check query complexity:**
- Large area searches take longer
- Use smaller radius for faster results
- Enable result limits

---

## See Also

- [API Documentation](API-Documentation) - REST API reference
- [Database Setup](Database-Setup) - DuckLake configuration
- [Deployment](Deployment) - Production deployment with AI features
- [Development](Development) - Extending MCP tools
- [RAG Semantic Cache](/docs/rag-semantic-cache.md) - Caching architecture
