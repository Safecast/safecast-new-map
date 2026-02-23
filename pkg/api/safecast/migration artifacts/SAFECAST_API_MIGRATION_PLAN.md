# Safecast API Migration Plan

**API contract (paths, parameters, response shapes) is maintained in pkg/api/safecast as swag comments and generated Swagger; this document is for migration and strategy context only.**

Migrate all API functionality from the Rails app in OLD_API_CODE/ to pkg/api/safecast/ in Go, following the [Go Proverbs](https://go-proverbs.github.io/). The Safecast API replaces api.safecast.org; the existing pkg/api (map/track browsing) remains untouched.

---

## Versioning and Adapter Architecture

**Single implementation:** One set of core logic handles all requests. No duplicate handlers for v1 vs v2.

**Adapter layer (unversioned + v1):** Receives **unversioned** and **v1** routes. Adapts request (Rails-style → internal) and response (internal → Rails format), then calls the shared core implementation. Legacy and v1 clients continue to work without change.

**v2 (direct):** Receives **v2** routes. Bypasses the adapter; calls the **same** core implementation directly. v2 can use a different request/response format at the HTTP boundary. The adapter and v2 differ only in the HTTP translation layer, not in business logic.

**Flow:**
```
Unversioned: GET /measurements.json     ─┐
                                        ├─→ Adapter (adapt in/out) → Core implementation → Response
v1:         GET /api/v1/measurements   ─┘

v2:         GET /api/v2/measurements   ───→ Direct → Core implementation (same) → Response (no adapter)
```

**Implementation:** The adapter (`pkg/api/safecast/adapter.go`) registers unversioned and v1 routes, adapts request/response, and forwards to the core handlers. The v2 routes (`/api/v2/*`) call the same core handlers directly. Only the HTTP layer (param parsing, response serialization) differs; the core logic is shared.

**Guardrails (required):**

- Adapter and v2 handlers must both call the same core service methods; do not duplicate database/business logic by version.
- Keep Rails-compat transformations in adapter translators only.
- Keep v2 request/response shaping in v2 translators only.
- Core service types remain version-neutral.
- Add parity tests that assert equivalent adapter/v2 requests produce equivalent core outcomes.

---

## API Contract Definitions

Each endpoint must match the following contracts so the new API is backward-compatible with existing clients (Safecast fetcher, bGeigie apps, scripts).

### Authentication

- **api_key**: Query param or header. Rails uses `config.token_authentication_key = :api_key` (Devise token). Lookup `users.authentication_token`.
- **Required for**: POST /measurements, POST /bgeigie_imports, POST /devices, PATCH/PUT/DELETE on owned resources, GET /users/me.
- **Optional for**: GET /measurements, GET /bgeigie_imports (index: all if moderator, else own only), GET /users, GET /devices.

### Path / Locale

- Routes may be under `/:locale/` with `locale` in `{cs-CZ, en-US, ja, pt}`. Also `/api/*path` redirects to `/%{path}.%{format}`.
- Support both: `/measurements.json` and `/en-US/measurements.json` (or equivalent).

---

### 1. GET / (root JSON)

**Paths**: `GET /` or `GET /:locale/` (adapter) → `GET /api/v1` (versioned).  
**Format**: `Accept: application/json` or `.json` format.

**Response 200**:
```json
{
  "name": "Safecast API",
  "uri": "https://api.safecast.org/",
  "subresource_uris": [
    "https://api.safecast.org/users.json",
    "https://api.safecast.org/measurements.json",
    "https://api.safecast.org/bgeigie_imports.json",
    "https://api.safecast.org/devices.json"
  ]
}
```

**Source**: OLD_API_CODE/app/views/home/show.json.jbuilder

---

### 2. GET /measurements.json

**Paths**: `GET /measurements`, `GET /measurements.json` (adapter) → `GET /api/v1/measurements` (versioned).  
**Auth**: None (public).  
**Query params**:

| Param | Type | Behaviour |
|-------|------|-----------|
| `latitude` | float | With `longitude` + `distance`: geo filter (ST_DWithin) |
| `longitude` | float | |
| `distance` | int | Distance in metres for `nearby_to` |
| `captured_after` | datetime | `captured_at > value` |
| `captured_before` | datetime | `captured_at < value` |
| `user_id` | int | Filter by user |
| `device_id` | int | Filter by device |
| `measurement_import_id` | int | Filter by import |
| `original_id` | int | `original_id = value OR id = value` |
| `since` | datetime | `updated_at > value` |
| `until` | datetime | `updated_at < value` |
| `unit` | string | Filter by unit |
| `order` | string | e.g. `created_at asc`, `created_at desc` (NULLS LAST) |
| `page` | int | Kaminari page (1-based) |
| `per_page` | int | Default 100 (Measurement.per_page) |

**Response 200**: JSON array of Measurement objects (no pagination metadata in body; Kaminari uses `without_count`).

**Measurement JSON** (from serializable_hash):
```json
{
  "id": 1,
  "value": 123,
  "height": null,
  "user_id": 1,
  "unit": "cpm",
  "device_id": 1,
  "location_name": "Some place",
  "original_id": 1,
  "captured_at": "2025/01/15 12:00:00 +0000",
  "devicetype_id": null,
  "sensor_id": null,
  "channel_id": null,
  "station_id": null,
  "measurement_import_id": 1,
  "latitude": 1.1,
  "longitude": 2.2
}
```

**Source**: OLD_API_CODE/app/controllers/measurements_controller.rb, OLD_API_CODE/app/models/measurement.rb

---

### 3. GET /measurements/:id.json

**Path**: `GET /measurements/:id.json`  
**Auth**: None.  
**Response 200**: Single Measurement object (same schema as above).  
**Response 404**: Record not found.  
**Note**: Returns the *original* measurement (by id); `original_id` links revisions. Use `original_id` query for revision history.

**Source**: OLD_API_CODE/app/controllers/measurements_controller.rb, OLD_API_CODE/spec/integration/api/measurements_spec.rb

---

### 4. POST /measurements.json

**Path**: `POST /measurements.json`  
**Auth**: Required (`api_key`).  
**Body**: `{ "measurement": { "value": 123, "unit": "cpm", "latitude": 1.1, "longitude": 2.2, ... } }`  
**Wrapped params**: Rails `wrap_parameters` with `measurement` key.

**Permitted fields**: `value`, `unit`, `location`, `location_name`, `device_id`, `height`, `surface`, `radiation`, `latitude`, `longitude`, `captured_at`, `devicetype_id`, `sensor_id`, `channel_id`, `station_id`

**Response 200**: Created Measurement object.  
**Response 422**: Validation errors `{ "errors": { "value": ["can't be blank"], ... } }`  
**Response 401**: Unauthorized without api_key.

**Post-create**: `original_id` set to `id` in same transaction.

**Source**: OLD_API_CODE/app/controllers/measurements_controller.rb, OLD_API_CODE/app/models/concerns/swagger_blocks/models/measurement_input.rb, OLD_API_CODE/spec/integration/api/measurements_spec.rb

---

### 4b. PUT /measurements/:id.json, PATCH /measurements/:id.json

**Path**: `PUT /measurements/:id.json` or `PATCH /measurements/:id.json`  
**Auth**: Required (`api_key`).  
**Body**: Same as POST (partial updates allowed).  
**Behaviour**: Non-destructive revise — creates new measurement, sets original's `expired_at`, `replaced_by`; new one gets `original_id` from original.  
**Response 202**: New measurement (the revision).  
**Response 401/404**: Unauthorized or not found.

**Source**: OLD_API_CODE/app/controllers/measurements_controller.rb, OLD_API_CODE/app/models/measurement.rb revise

---

### 5. GET /count and GET /measurements/count.json

**Path**: `GET /count` or `GET /measurements/count.json`  
**Auth**: None.  
**Response 200**: `{ "count": N }`  
**Count logic**: Uses PostgreSQL estimation: `(reltuples/NULLIF(relpages,0)) * (pg_relation_size(...) / block_size)`. On small datasets may return 0.

**Source**: OLD_API_CODE/app/controllers/measurements_controller.rb, OLD_API_CODE/spec/integration/api/measurements_spec.rb

---

### 6. GET /bgeigie_imports.json

**Path**: `GET /bgeigie_imports.json` (or `Accept: application/json`)  
**Auth**: If not moderator: only own imports. If moderator: all.  
**Query params**:

| Param | Behaviour |
|-------|-----------|
| `q` | filter_by_text_fields on name, source, description, cities, credits (lowercase LIKE) |
| `by_status` | status filter |
| `by_user_id` | user_id |
| `by_user_name` | user.by_name |
| `by_rejected` | rejected |
| `rejected_by` | rejected_by LIKE |
| `uploaded_after` | created_at > |
| `uploaded_before` | created_at < |
| `status` | `approved` \| `rejected` \| `not_moderated` |
| `subtype` | subtype (comma-separated: None, Drive, Surface, Cosmic) |
| `order` | column + asc/desc, e.g. `created_at asc`, `created_at desc` |
| `page` | Pagination |
| `per_page` | Default Kaminari |

**Response 200**: JSON array of BgeigieImport objects (paginated).  
**BgeigieImport shape**: `id`, `user_id`, `approved`, `created_at`, `updated_at`, `measurements_count`, `md5sum`, `name`, `status`, `source: { url }` (S3 URL when available).

**Source**: OLD_API_CODE/app/controllers/bgeigie_imports_controller.rb, OLD_API_CODE/db/structure.sql, pkg/safecast-fetcher/client.go

---

### 7. GET /bgeigie_imports/:id.json

**Path**: `GET /bgeigie_imports/:id.json`  
**Auth**: Index scope (own or moderator).  
**Response 200**: Single BgeigieImport.  
**Response 404**: Not found.

**Status values**: `unprocessed`, `processed`, `submitted`, `approved`, `done`.

**Source**: OLD_API_CODE/app/controllers/bgeigie_imports_controller.rb, OLD_API_CODE/spec/integration/api/bgeigie_imports_spec.rb

---

### 8. POST /bgeigie_imports.json

**Path**: `POST /bgeigie_imports`  
**Auth**: Required (`api_key`).  
**Body**: Multipart form; `bgeigie_import[source]` = file upload.  
**Headers**: `Accept: application/json` for JSON response.

**Response 200**: Created BgeigieImport: `id`, `status` (e.g. `unprocessed`), `md5sum` (e.g. `aad36f9743753b490745c9656cd8fc79`).  
**Processing**: Async job. After processing, `status` → `done`, `measurements_count` populated.

**bGeigie format**: `$BNRDD`, `$BMRDD`, `$BGRDD`, `$BNXRDD`, `$PNTDD`, `$CZRDD`; CPM validity A/V; GPS A/V; valid date.

**Source**: OLD_API_CODE/app/controllers/bgeigie_imports_controller.rb, OLD_API_CODE/app/models/bgeigie_import.rb is_sane?, OLD_API_CODE/spec/integration/api/bgeigie_imports_spec.rb

---

### 8b. PATCH /bgeigie_imports/:id/submit

**Path**: `PATCH /bgeigie_imports/:id/submit` or `PUT /bgeigie_imports/:id/submit`  
**Auth**: Required (owner).  
**Behaviour**: If `would_auto_approve`, approve immediately; else set status to `submitted`, send notification.  
**Response**: Redirect (HTML) or 200 (JSON).

**Source**: OLD_API_CODE/app/controllers/bgeigie_imports_controller.rb, OLD_API_CODE/doc/upload_demo/api_demo.py

---

### 9. GET /users.json

**Path**: `GET /users.json`  
**Auth**: None (public).  
**Query params**: `name` (filter by name, LIKE).  
**Response 200**: Array of users. `where.not(confirmed_at: nil)` (only confirmed).

**User JSON**: `id`, `name`, `measurements_count`, `authentication_token` (only if current_user == user).

**Source**: OLD_API_CODE/app/controllers/users_controller.rb, OLD_API_CODE/app/views/users/_user.json.jbuilder

---

### 10. GET /users/:id.json, GET /users/me.json

**Path**: `GET /users/:id.json`, `GET /users/me.json`  
**Auth**: `me` requires auth.  
**Response 200**: Single user.  
**Response 401**: Unauthorized for /me without api_key.

**Source**: OLD_API_CODE/app/controllers/users_controller.rb

---

### 11. POST /users.json

**Path**: `POST /users.json`  
**Auth**: None (registration).  
**Body**: `{ "user": { "email": "...", "name": "...", "password": "..." } }`  
**Response 200**: Created user with `id`, `email`, `name`, `authentication_token` (for immediate API use).

**Source**: OLD_API_CODE/spec/integration/api/users_spec.rb

---

### 12. GET /devices.json

**Path**: `GET /devices.json`  
**Auth**: None.  
**Query params**: `manufacturer`, `model`, `sensor` (LIKE %value%).  
**Response 200**: Array of `{ "id", "manufacturer", "model", "sensor" }`.

**Source**: OLD_API_CODE/app/controllers/devices_controller.rb, OLD_API_CODE/spec/integration/api/devices_spec.rb

---

### 13. GET /devices/:id.json

**Path**: `GET /devices/:id.json`  
**Response 200**: Single device.

**Source**: OLD_API_CODE/app/controllers/devices_controller.rb

---

### 14. POST /devices.json

**Path**: `POST /devices.json`  
**Auth**: Required (`api_key`).  
**Body**: `{ "device": { "manufacturer": "...", "model": "...", "sensor": "..." } }`  
**Response 200**: Device (get_or_create: returns existing if same manufacturer+model+sensor).  
**Response 422**: `{ "errors": { "manufacturer": [...], "model": [...], "sensor": [...] } }` if empty.

**Source**: OLD_API_CODE/app/controllers/devices_controller.rb, OLD_API_CODE/app/models/device.rb get_or_create, OLD_API_CODE/spec/integration/api/devices_spec.rb

---

### 15. GET /radiation_index

**Path**: `GET /radiation_index` or `GET /:locale/radiation_index`  
**Query params**: `index` (optional): `:minimum` → index 1, `:maximum` → index 2, else 0.  
**Data**: Reads `public/system/g20.csv`; selects column by index; sorts by value descending; nils as -1.  
**Response**: HTML or JSON (depends on format). Returns sorted array of `{ country_code: value }`.

**Source**: OLD_API_CODE/app/controllers/radiation_index_controller.rb

---

### 16. GET /ingest, GET /ingest.csv

**Path**: `GET /ingest`, `GET /ingest.csv`  
**Query params**: `area`, `field`, `uploaded_after`, `uploaded_before`.  
**Data**: Elasticsearch. `area` in DEVICE_GROUPS (central_japan, fukushima, washington, etc.); `field` in document; date range.  
**Response 200**: HTML table or CSV. CSV: `when_captured`, `value`, `device`.

**Source**: OLD_API_CODE/app/controllers/ingest_controller.rb  
**Dependency**: Elasticsearch.

---

### 17. GET /device_stories, GET /device_stories/:id, GET /airnote/:device_urn

**Path**: `GET /device_stories`, `GET /device_stories/:id`, `GET /airnote/:device_urn`  
**Query params**: `search` (lowercase LIKE on device_urn, custodian_name, last_seen); `page`, `per_page`; `order`.  
**Response**: Device story objects. `device_urn` identifies AirNote devices.

**Source**: OLD_API_CODE/app/controllers/device_stories_controller.rb, OLD_API_CODE/config/routes.rb

---

### 18. CORS

**Headers** (from OLD_API_CODE/app/controllers/application_controller.rb):

- `Access-Control-Allow-Origin`: request origin or `safecast.org` if no user
- `Access-Control-Allow-Methods`: POST, GET, OPTIONS
- `Access-Control-Allow-Headers`: *, X-Requested-With
- `Access-Control-Max-Age`: 100000

---

## Implementation Plan

### Phase 1: Foundation

1. Create `pkg/api/safecast/` package (handler.go, dto.go)
2. **Core implementation**: Single set of handlers (measurements, bgeigie, users, devices, etc.) used by both adapter and v2
3. **Adapter layer**: `adapter.go` receives unversioned + v1 routes, adapts request/response to Rails contract, calls core
4. **v2 routes**: Register `/api/v2/*` that call the **same** core directly (no adapter)
5. Wire adapter and v2 routes in safecast-new-map.go
6. Database extensions for markers→measurements and uploads→bgeigie_imports mapping

### Phase 2: Core Endpoints

Measurements GET/POST/PUT, count, bgeigie_imports GET/POST/submit (single core, adapter + v2 both use it)

### Phase 3: Auxiliary

Users, devices, radiation_index, ingest (stub), device_stories

### Phase 4: Polish

Root JSON, locale handling, tests, documentation

### Phase 5: Rollout and Verification

- Deploy behind feature flag.
- Enable adapter endpoints first (legacy compatibility).
- Run sampled parity checks between adapter and v2.
- Enable v2 for selected clients, then full rollout.
- Keep rollback path documented (disable flag, route traffic back to stable path).

---

## File Structure

**pkg/api/safecast/** (new package):

```
handler.go       # Handler struct; Register(mux) mounts adapter + v2 routes
adapter.go       # Adapter: unversioned + v1 → adapt in/out → core
dto.go           # Rails (v1) + v2 response structs
measurements.go  # Core logic (shared by adapter + v2)
bgeigie.go       # Core logic (shared)
users.go         # Core logic (shared)
devices.go       # Core logic (shared)
radiation.go     # Core logic (shared)
```

**Route registration**:
- **Adapter** (unversioned + v1): `/`, `/measurements`, `/bgeigie_imports`, `/users`, `/devices`, `/count`, `/radiation_index`, `/ingest`, `/device_stories`, `/airnote/`, and `/api/v1/*` — adapts to/from Rails format, calls core
- **v2** (direct): `/api/v2/*` — calls the same core directly, no adapter (v2 request/response format)

---

## Definition of Done

- All endpoint contracts in this document are implemented and tested.
- Adapter and v2 both use the same core service logic for each endpoint.
- Contract tests validate status codes, auth rules, and payload/error shapes.
- Parity tests validate equivalent semantic results across adapter and v2.
- Production telemetry shows stable error rate/latency after rollout.
