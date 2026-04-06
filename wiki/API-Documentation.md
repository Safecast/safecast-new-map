# API Documentation

Complete reference for the Safecast New Map REST API.

[[Home|← Back to Home]]

---

## API Overview

The platform provides a comprehensive REST API for accessing radiation data, managing users, and administering the system.

**Base URL:** `http://localhost:8765/api`

**API Documentation (Swagger UI):**
- Map API: [http://localhost:8765/map-api/](http://localhost:8765/map-api/)
- MCP API: [http://localhost:8765/mcp-api/](http://localhost:8765/mcp-api/)
- Combined: [http://localhost:8765/docs/](http://localhost:8765/docs/)

**Authentication:**
- Most endpoints are public
- Protected endpoints require authentication via:
  - Session cookie (web login)
  - API key in header: `X-API-Key: your-key`
  - API key in query parameter: `?api_key=your-key`

---

## Radiation Data Endpoints

### Query Radiation Measurements

```
GET /api/radiation?lat={lat}&lon={lon}&radius={meters}
```

**Parameters:**
- `lat` (required) - Latitude coordinate
- `lon` (required) - Longitude coordinate  
- `radius` (optional, default 1000) - Search radius in meters

**Response:**
```json
[
  {
    "id": 12345,
    "doseRate": 0.12,
    "lat": 35.6762,
    "lon": 139.6503,
    "date": "2024-01-15T10:30:00Z",
    "countRate": 120,
    "zoom": 12,
    "speed": 5.2
  }
]
```

### Query Area Data

```
GET /api/area?minLat={minLat}&minLon={minLon}&maxLat={maxLat}&maxLon={maxLon}
```

**Parameters:**
- `minLat`, `minLon`, `maxLat`, `maxLon` (required) - Bounding box coordinates

**Use Case:** Load all measurements within map viewport

### Track Management

#### List Tracks

```
GET /api/tracks?page={page}&limit={limit}
```

**Parameters:**
- `page` (optional, default 1) - Page number
- `limit` (optional, default 50) - Results per page

**Response:**
```json
{
  "tracks": [
    {
      "trackId": "abc123",
      "count": 1250,
      "startDate": "2024-01-15T10:00:00Z",
      "endDate": "2024-01-15T12:30:00Z"
    }
  ],
  "total": 5000,
  "page": 1,
  "limit": 50
}
```

#### Get Track Data

```
GET /api/track/{trackId}
```

Returns all measurements for a specific track.

#### Get Device History

```
GET /api/device/{deviceId}/history
```

Returns measurement history for a specific device.

---

## Real-Time Sensors

### List Sensors

```
GET /api/sensors
```

Returns all active real-time sensors.

**Response:**
```json
[
  {
    "id": "sensor-123",
    "name": "Tokyo Pointcast",
    "lat": 35.6762,
    "lon": 139.6503,
    "doseRate": 0.08,
    "lastUpdate": "2024-01-15T10:30:00Z"
  }
]
```

### Get Current Sensor Reading

```
GET /api/sensor/{sensorId}/current
```

Returns the latest reading for a specific sensor.

### Get Sensor History

```
GET /api/sensor/{sensorId}/history?hours={hours}
```

**Parameters:**
- `hours` (optional, default 24) - Hours of history to return

---

## Spectroscopy Endpoints

### List Spectra

```
GET /api/spectra?page={page}&limit={limit}
```

Returns available gamma spectra.

### Get Spectrum Data

```
GET /api/spectrum/{markerId}
```

Returns spectrum data for a specific marker.

**Response:**
```json
{
  "markerId": 12345,
  "channels": [100, 120, 95, ...],
  "channelCount": 1024,
  "energyMinKev": 0,
  "energyMaxKev": 3000,
  "liveTimeSec": 300,
  "calibration": {
    "a": 0.5,
    "b": 0.3,
    "c": 0.0001
  },
  "deviceModel": "RadiaCode-102"
}
```

### Download Spectrum

```
GET /api/spectrum/{markerId}/download?format={format}
```

**Formats:**
- `json` - JSON format
- `csv` - CSV format
- `n42` - ANSI N42.42 format
- `spe` - Maestro format

---

## Reference & Statistics

### Get Statistics

```
GET /api/stats
```

Returns aggregate database statistics.

**Response:**
```json
{
  "totalMarkers": 518400000,
  "totalTracks": 45000,
  "totalUsers": 1200,
  "dateRange": {
    "start": "2012-03-11T00:00:00Z",
    "end": "2025-04-05T00:00:00Z"
  }
}
```

### Get Extreme Readings

```
GET /api/extreme
```

Returns highest radiation readings in the database.

### Get Radiation Info

```
GET /api/info/{topic}
```

Returns reference information about radiation topics.

**Topics:**
- `background` - Natural background radiation
- `units` - Radiation units and conversion
- `safety` - Safety levels and guidelines

---

## Authentication Endpoints

### Register User

```
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword",
  "username": "johndoe"
}
```

### Login

```
POST /api/auth/login
Content-Type: application/json

// Password login
{
  "email": "user@example.com",
  "password": "your-password"
}

// API key login
{
  "email": "user@example.com",
  "api_key": "your-20-char-api-key"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "user": {
    "id": 42,
    "email": "user@example.com",
    "username": "johndoe",
    "api_key": "abcdefghijklmnopqrst"
  }
}
```

### Logout

```
POST /api/auth/logout
```

### Forgot Password

```
POST /api/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

### Reset Password

```
POST /api/auth/reset-password
Content-Type: application/json

{
  "token": "reset-token-from-email",
  "password": "new-password"
}
```

### Verify Email

```
GET /api/auth/verify-email?token={token}
```

### Get User Profile

```
GET /api/user/profile
Authorization: Bearer {session-cookie or api-key}
```

### Change Password

```
POST /api/user/change-password
Content-Type: application/json

{
  "currentPassword": "old-password",
  "newPassword": "new-password"
}
```

### Get User Uploads

```
GET /api/user/uploads?page={page}&limit={limit}
```

Returns authenticated user's upload history.

---

## Admin Endpoints

**Authentication:** All admin endpoints require admin password or admin user session.

### User Management

#### List Users

```
GET /api/admin/users?password={admin-password}
```

#### Create User

```
POST /api/admin/users/create?password={admin-password}
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password",
  "username": "newuser",
  "is_admin": false
}
```

#### Update User

```
PUT /api/admin/users/{userId}?password={admin-password}
Content-Type: application/json

{
  "username": "newname",
  "is_admin": true
}
```

#### Delete User

```
DELETE /api/admin/users/{userId}?password={admin-password}
```

#### Reset User Password

```
POST /api/admin/users/{userId}/reset-password?password={admin-password}
```

#### Regenerate API Key

```
POST /api/admin/users/{userId}/regenerate-api-key?password={admin-password}
```

**Response:**
```json
{
  "message": "API key regenerated successfully",
  "api_key": "new-20-char-key",
  "user_id": 42
}
```

### Upload Management

#### List Uploads

```
GET /api/admin/uploads?password={admin-password}&page={page}
```

#### Update Upload Metadata

```
POST /api/admin/uploads/update?password={admin-password}
Content-Type: application/json

{
  "uploadId": 123,
  "recordingDate": "2024-01-15T10:00:00Z",
  "detector": "bGeigie",
  "username": "john"
}
```

### Track Management

#### Delete Track

```
POST /api/admin/delete?password={admin-password}
Content-Type: application/json

{
  "trackId": "abc123"
}
```

#### Delete Multiple Tracks

```
POST /api/admin/delete-multiple?password={admin-password}
Content-Type: application/json

{
  "trackIds": ["abc123", "def456", "ghi789"]
}
```

#### Update Track Metadata

```
POST /api/admin/tracks/update?password={admin-password}
Content-Type: application/json

{
  "trackId": "abc123",
  "name": "Tokyo Survey",
  "notes": "Urban area measurement"
}
```

### Import from Safecast

#### Import by Safecast ID

```
POST /api/admin/import-from-safecast?password={admin-password}
Content-Type: application/json

{
  "safecastId": 12345
}
```

#### Import by ID Range

```
POST /api/admin/import-by-id?password={admin-password}
Content-Type: application/json

{
  "startId": 100000,
  "endId": 110000
}
```

#### Backfill by Country

```
POST /api/admin/backfill-countries?password={admin-password}
Content-Type: application/json

{
  "countries": ["Japan", "United States"]
}
```

### Cache Management

#### Clear Cache

```
POST /api/admin/cache?password={admin-password}
```

Clears LRU cache and forces data refresh.

---

## Data Export Endpoints

### JSON Archive

```
GET /api/json/weekly.tgz
```

Download compressed JSON archive.

**Frequency options (configured via `-json-archive-frequency`):**
- `daily.tgz`
- `weekly.tgz` (default)
- `monthly.tgz`
- `yearly.tgz`

---

## Short Links

### Create Short Link

```
POST /api/shorten
Content-Type: application/json

{
  "url": "https://simplemap.safecast.org/?lat=35.6762&lon=139.6503&zoom=12"
}
```

**Response:**
```json
{
  "code": "aB3xY9zK",
  "url": "https://simplemap.safecast.org/s/aB3xY9zK"
}
```

### Resolve Short Link

```
GET /s/{code}
```

Redirects to target URL.

---

## Rate Limiting

API endpoints are rate-limited by default:

- **Public endpoints:** 100 requests per minute per IP
- **Authentication endpoints:** 10 requests per minute per IP
- **Admin endpoints:** 50 requests per minute per IP

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1617234567
```

---

## Error Responses

### Common Errors

**400 Bad Request:**
```json
{
  "error": "Missing required parameter: lat"
}
```

**401 Unauthorized:**
```json
{
  "error": "Authentication required"
}
```

**403 Forbidden:**
```json
{
  "error": "Admin access required"
}
```

**404 Not Found:**
```json
{
  "error": "Track not found: abc123"
}
```

**429 Too Many Requests:**
```json
{
  "error": "Rate limit exceeded. Try again in 30 seconds."
}
```

**500 Internal Server Error:**
```json
{
  "error": "Internal server error"
}
```

---

## GPT-Optimized Endpoints

Compact endpoints designed for Custom GPT Actions (registered in `rest_gpt.go`):

These endpoints return minimal, optimized responses for AI consumption:
- Reduced payload size
- Essential fields only
- Faster response times

---

## API Versioning

The current API is version 1 (v1). Version information is included in responses:

```
X-API-Version: 1
```

---

## Authentication Methods

### Session-Based (Web)

1. Login via `/api/auth/login`
2. Session cookie is set automatically
3. Include cookie in subsequent requests

### API Key (Programmatic)

**Header method (recommended):**
```bash
curl -H "X-API-Key: your-api-key" http://localhost:8765/api/endpoint
```

**Query parameter method:**
```bash
curl "http://localhost:8765/api/endpoint?api_key=your-api-key"
```

---

## Regenerating API Documentation

After API changes, regenerate Swagger documentation:

```bash
# Map API docs
cd cmd/unified-server && swag init \
  -g doc.go \
  -o docs/api \
  --parseDependency \
  --parseInternal \
  --parseDependencyLevel 2 \
  --instanceName unifiedapi

# MCP server docs
cd cmd/mcp-server && swag init -g rest.go
```

---

## Best Practices

### Performance

1. **Use bounding box queries** for map data instead of radius queries
2. **Paginate results** for large datasets
3. **Cache responses** when possible
4. **Use track streaming** for efficient large dataset handling

### Security

1. **Use HTTPS** in production
2. **Rotate API keys** periodically
3. **Enable authentication** for uploads
4. **Monitor logs** for suspicious activity

### Error Handling

1. **Check HTTP status codes**
2. **Parse error messages** for debugging
3. **Implement retry logic** with exponential backoff
4. **Handle rate limits** gracefully

---

## See Also

- [User Authentication](User-Authentication) - Authentication setup and API keys
- [MCP Server & AI Integration](MCP-Server-AI-Integration) - AI-powered API access
- [Admin Panel](Admin-Panel) - Web-based administration
- [Swagger UI](http://localhost:8765/map-api/) - Interactive API documentation
