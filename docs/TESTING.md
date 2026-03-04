# Testing Guide

## Unit & Integration Tests

Tests use an in-memory SQLite database — no running server needed.

### Run all tests

```bash
go test ./pkg/... -timeout 120s
```

### Run a specific package

```bash
go test ./pkg/api/...   -v    # API endpoints
go test ./pkg/auth/...  -v    # Auth endpoints (slow: ~3s due to bcrypt)
go test ./pkg/web/...   -v    # Web/spectrum/QR handlers
```

### Run a single test

```bash
go test ./pkg/api/...  -run TestHandleShorten -v
go test ./pkg/auth/... -run TestLoginHandler  -v
```

### Expected output

```
ok  safecast-new-map/pkg/api              0.1s
ok  safecast-new-map/pkg/auth             2.9s   ← slow due to bcrypt cost 12
ok  safecast-new-map/pkg/countryresolver  0.7s
ok  safecast-new-map/pkg/safecast-realtime 0.0s
ok  safecast-new-map/pkg/web              0.5s
```

---

## Smoke Tests (end-to-end against a running server)

Requires `curl`. Tests all HTTP endpoints by making real requests.

### Prerequisites

Start the server locally (see `local-server-config.sh` for credentials):

```bash
./safecast-new-map -db-type pgx -db-conn "postgres://..." -allow-registration -admin-password test123
```

### Run smoke tests

```bash
# Local server (default)
./test/smoke_test.sh

# Custom URL and admin password
./test/smoke_test.sh http://localhost:8765 test123

# Production (read-only — registration creates a test user each run)
./test/smoke_test.sh https://simplemap.safecast.org your_admin_pw

# Also test the MCP server REST endpoints
MCP_BASE=http://localhost:3333 ./test/smoke_test.sh http://localhost:8765 test123
```

### What is tested

| Category | Endpoints |
|---|---|
| Pages | `GET /`, `/home` |
| Utilities | `/qrpng`, `/licenses/mit`, `/licenses/cc0`, `/api/docs` |
| API data | `/api`, `/api/tracks`, `/api/countries`, `/api/latest` |
| Shorten | `POST /api/shorten` (preview + bad requests) |
| Spectrum | `/api/markers/spectra`, `/api/track-info/{id}` |
| Markers | `/get_markers` |
| Auth | register, login, logout, profile, forgot-password |
| Admin | `/api/admin/uploads`, `/api/admin/tracks`, `/api/admin/users` |
| MCP REST | `/api/radiation`, `/api/stats` (if `MCP_BASE` set) |

### Notes

- Each smoke test run registers a new user with a timestamped email to avoid conflicts.
- Admin endpoints are tested with `?password=ADMIN_PW` query auth.
- The script exits with code `1` if any test fails.

---

## Coverage Report

```bash
go test ./pkg/... -coverprofile=coverage.out -timeout 120s
go tool cover -html=coverage.out -o coverage.html
open coverage.html
```
