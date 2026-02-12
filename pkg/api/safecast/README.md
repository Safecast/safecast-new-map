# Safecast API (pkg/api/safecast)

This package implements the Safecast API (api.safecast.org replacement) in Go. It follows the migration plan in [docs/SAFECAST_API_MIGRATION_PLAN.md](../../docs/SAFECAST_API_MIGRATION_PLAN.md).

## Architecture

- **Single implementation**: One set of core logic handles all requests. No duplicate handlers for v1 vs v2.
- **Adapter layer**: Unversioned (`/measurements`, etc.) and v1 (`/api/v1/measurements`) routes adapt request/response to the Rails contract, then call the shared core.
- **v2** (`/api/v2/measurements`): Calls the same core directly. No Rails adaptation.

## Route registration

- **Adapter**: `/measurements`, `/bgeigie_imports`, `/users`, `/devices`, `/count`, `/radiation_index`, `/ingest`, `/device_stories`, `/airnote/`, and `/api/v1/*`
- **v2**: `/api/v2/*`

## Endpoints

| Endpoint | Status |
|----------|--------|
| GET /api/v1, /api/v2 (root) | ✓ |
| GET /measurements, GET /bgeigie_imports | ✓ |
| GET /measurements/:id | ✓ |
| POST /measurements | ✓ |
| GET /count | ✓ |
| GET /bgeigie_imports/:id | ✓ |
| POST /bgeigie_imports | Stub (501) |
| Users, devices, radiation_index, ingest, device_stories | Stub |

## Auth

- `api_key` query param or `X-API-Key` header
- Maps to `users.api_key` via `pkg/auth.GetUserByAPIKey`
