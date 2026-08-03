# Auto-submit bGeigie uploads to api.safecast.org

## Goal
When a user uploads a bGeigie log (`.log`/`.txt`/`$BNRDD`) to simplemap, check whether
that file already exists on api.safecast.org for that user, and if not, submit it there
too using the user's own Safecast API key.

## Scope
bGeigie logs only. Spectrum/track-export uploads have no equivalent upstream endpoint.

## Upstream API (Safecast/safecastapi, confirmed from source)

### Resolve user id from API key
```
GET https://api.safecast.org/users/me.json?api_key=<key>
```
`UsersController#me` (`before_action :authenticate_user!, only: %i(me)`) returns the
current user's record, including numeric `id`. Called once when a user saves/updates
their Safecast API key; the returned `id` is cached as `safecast_user_id`.

### Check for existing import (dedup)
```
GET https://api.safecast.org/bgeigie_imports.json?api_key=<key>&by_user_id=<id>&q=<filename>
```
`BgeigieImportsController#index` supports `has_scope :by_user_id` →
`BgeigieImport.where(user_id: user_id)`, and `has_scope :q` →
`filter_by_text_fields`, which includes `lower(source) LIKE '%query%'` (`source` is the
stored filename). Combining both scopes filters to "this user's imports matching this
filename."

### Submit a log file
```
POST https://api.safecast.org/bgeigie_imports.json?api_key=<key>
Content-Type: multipart/form-data; boundary=...
```
Two parts:
- `bgeigie_import[description]` = free text (e.g. `"Uploaded from simplemap"`)
- `bgeigie_import[source]` = raw log file bytes, `filename="<name>.log"`,
  `Content-Type: text/plain`

Response JSON on 2xx: `{"id": <n>, ...}` — new import id.

Reference implementation: bGeigieZen firmware, branch `feature/sd-log-upload`,
commit `1dc2edf`, `bgeigiezen_firmware/screens/log_viewer.cpp::upload_detail()`.
The firmware never does the dedup check (manual one-shot upload) — that part is new.

## Data model changes

`users` table (migration `add_safecast_credentials.sql`):
- `safecast_api_key` (nullable text) — user's api.safecast.org API key, pasted in `/profile`
- `safecast_user_id` (nullable text) — resolved automatically via `/users/me.json` when
  `safecast_api_key` is set/changed; not user-editable

`uploads` table (migration `add_safecast_submit_tracking.sql`):
- `safecast_import_id` (nullable text) — id returned by successful submit
- `submitted_to_safecast_at` (nullable timestamp)
- `safecast_submit_error` (nullable text) — last failure reason, if any

Both migrations follow the existing `ALTER TABLE ... ADD COLUMN IF NOT EXISTS` style
(see `migrations/add_upload_metadata.sql`, `add_user_id_column.sql`). New columns are
also added to the inline `CREATE TABLE` schema in `pkg/database/database.go` for fresh
installs (sqlite / pgx / duckdb variants).

## New package: `pkg/safecast-submit`

Mirrors `pkg/safecast-fetcher/client.go` conventions (http.Client w/ timeout,
`NewRequestWithContext`, `%w`-wrapped errors, explicit status-code checks).

```go
type Client struct { httpClient *http.Client; baseURL string }
func NewClient() *Client
func (c *Client) ResolveUserID(ctx, apiKey string) (userID string, err error)
func (c *Client) CheckExists(ctx, apiKey, userID, filename string) (bool, error)
func (c *Client) Submit(ctx, apiKey, filename string, content []byte) (importID string, err error)
```

## Hook point

In `uploadHandler`'s background goroutine (`cmd/unified-server/handlers_upload.go`),
after a bGeigie file is parsed by `processBGeigieZenFile` and `db.InsertUpload` succeeds:

1. Look up `safecast_api_key`/`safecast_user_id` for `internalUserID`; skip silently if
   either is empty (feature is opt-in per user).
2. `CheckExists` by filename; skip (log, no error) if already present upstream.
3. If not present, `Submit`; on success, `UpdateUploadSafecastStatus` sets
   `safecast_import_id` + `submitted_to_safecast_at`. On failure, sets
   `safecast_submit_error` and logs — never blocks or fails the local upload.

## Testing plan

1. Unit tests for `pkg/safecast-submit` using `httptest.NewServer` mocking
   `/users/me.json`, `/bgeigie_imports.json` (GET + POST) — verify request shape
   (query params, multipart body fields) and response parsing, including error paths
   (404, 401, 5xx).
2. Integration test in `cmd/unified-server` wiring the hook with a mocked
   `safecast-submit.Client` (interface) to verify: skip when no key, skip when
   `CheckExists` true, submit + DB update when not present, error recorded on submit
   failure — all without a local upload being blocked.
3. Manual local-server test: run via `local-server-config.sh`, upload a bGeigie log for
   a test user with a real (test-account) Safecast API key, confirm the row appears via
   admin uploads endpoint with `safecast_import_id` populated, and confirm on
   api.safecast.org that the import exists.
4. Only after local verification, deploy per the standard 4-step production process.

## Rollout order
1. Migrations (`users` + `uploads` columns)
2. `pkg/safecast-submit` client + unit tests
3. Wire into `uploadHandler` bGeigie branch (behind per-user opt-in — key must be set)
4. `/profile` UI field for the Safecast API key (triggers `ResolveUserID` on save)
5. Local end-to-end test
6. Deploy to production (manual 4-step process)
