# Safecast API (pkg/api/safecast)

This package implements the Safecast API (api.safecast.org replacement) in Go. It follows the migration plan in [docs/SAFECAST_API_MIGRATION_PLAN.md](../../docs/SAFECAST_API_MIGRATION_PLAN.md).

**API documentation is generated from code.** From the **repository root**, run (use `doc.go` for `-g` so swag finds it inside the parse dir):

```bash
swag fmt -g doc.go -d pkg/api/safecast
swag init -g doc.go -d pkg/api/safecast -o pkg/api/safecast/docs
```

Then see `docs/` (swagger.json, swagger.yaml) or Swagger UI if mounted in the main app. The API contract (paths, parameters, response shapes) lives only in the package’s Go files (swag comments and [doc.go](doc.go)).

## Architecture

- **Single implementation**: One set of core logic handles all requests. No duplicate handlers for v1 vs v2.
- **Adapter layer**: Unversioned (`/measurements`, etc.) and v1 (`/api/v1/measurements`) routes adapt request/response to the Rails contract, then call the shared core.
- **v2** (`/api/v2/measurements`): Calls the same core directly. No Rails adaptation.

## Auth

- `api_key` query param or `X-API-Key` header
- Maps to `users.api_key` via `pkg/auth.GetUserByAPIKey`

## Limitations

This API implementation lacks testing. It's also a 1-1 recreation of the original API, and likely has structural issues. Currently, the map uses a bridge to query the old API, rather than using the new implmentation.
