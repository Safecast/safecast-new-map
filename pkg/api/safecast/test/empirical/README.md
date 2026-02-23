# Empirical API tests

These tests replay recorded HTTP responses from **api.safecast.org** and assert that the new API matches the same response format: status code, `Content-Type`, and for JSON responses the **shape** of the body (keys and value types), not exact values.

## Testdata

- **Location:** `testdata/*.json` in this package (embedded via `//go:embed`).
- **Origin:** Recordings were captured from api.safecast.org with a one-off script (since removed). They were then **sanitized**: no real user names/IDs or third-party keys remain. User names appear as `[SANITIZED] User 1`, etc. See `testdata/README.md` for details.
- **Override:** Set `SAFECAST_RECORDINGS_DIR` to a directory path to load recordings from disk instead of the embedded testdata.

## What is asserted

- Request method and URL path (and query) are replayed; the handler under test serves the request.
- Response status code is compared (with allowed overrides, e.g. 406 for missing `Accept: application/json`).
- For JSON responses, the response body is checked for the same JSON structure (object keys, array element shape). Values are not compared.
