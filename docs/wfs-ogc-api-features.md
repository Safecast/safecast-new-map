# WFS / OGC API - Features endpoint

Safecast radiation data is served to GIS clients (QGIS, ArcGIS, etc.) via an
**OGC API - Features** endpoint — the current OGC standard that supersedes the
classic XML-based WFS 2.0.0. QGIS consumes it through its built-in
*WFS / OGC API - Features* data source; no plugin needed.

OGC API - Features was chosen over classic WFS 2.0.0 because it is GeoJSON
throughout (no GetCapabilities/DescribeFeatureType XML to hand-maintain) and it
fixes coordinate/`bbox` order to lon,lat everywhere, avoiding the classic-WFS
EPSG:4326 axis-order trap that causes empty or misplaced results in QGIS.

- **Implementation:** [`cmd/unified-server/ogcapi.go`](../cmd/unified-server/ogcapi.go)
- **Route wiring:** `cmd/unified-server/main.go` (registered on `http.DefaultServeMux`)
- **Served by:** unified-server, HTTP port 8765

## Base URL

```
https://simplemap.safecast.org/ogc
```

## Using it in QGIS

1. Layer → Add Layer → **Add WFS / OGC API - Features Layer…**
2. **New** connection, URL: `https://simplemap.safecast.org/ogc`
3. Connect → QGIS discovers the `sensors` collection → Add.

The layer refreshes by bounding box as you pan/zoom. Each point carries the
current `value` + `unit`, so you can style/graduate by radiation level.

## Endpoints

| Method | Path | Returns |
|--------|------|---------|
| GET | `/ogc` | Landing page (JSON links) |
| GET | `/ogc/conformance` | Conformance declaration |
| GET | `/ogc/collections` | List of feature collections |
| GET | `/ogc/collections/{id}` | Collection metadata |
| GET | `/ogc/collections/{id}/items` | GeoJSON `FeatureCollection` |

### `items` query parameters

| Param | Default | Notes |
|-------|---------|-------|
| `bbox` | world | `minLon,minLat,maxLon,maxLat` (CRS84 / lon,lat order) |
| `limit` | 1000 | max 10000 |
| `offset` | 0 | paging; responses include a `next` link while more rows remain |

Item responses include `numberMatched`, `numberReturned`, and `links`
(`self`, `collection`, and `next` when paged).

## Collections

### `sensors` (v1)

Most recent reading from each **live fixed sensor** (bGeigieZen, Pointcast,
Solarcast, nGeigie, Notehub, …). Feature properties:

- `device_id`, `device_name`, `type`
- `value`, `unit` — the latest radiation reading
- `last_reading_at`

Backed by the realtime measurements table (latest row per `device_id` within
the bbox). Small dataset (~hundreds); served as a single page.

### `measurements` (planned v2)

Historical bGeigie drive data (the dataset behind the old te512 raster map,
millions of points). Requires real bbox paging via `offset`/`next`; add a case
to `ogcCollections()` / `ogcItems()` in `ogcapi.go` when the paged query lands.

## Deployment note (nginx)

`/ogc` and `/ogc/` must be routed to unified-server on **port 8765**, like the
other `/api/` locations. Add to both nginx site configs
(`origin-simplemap.safecast.org` and `simplemap.safecast.org`):

```nginx
location /ogc { proxy_pass http://127.0.0.1:8765; }
location /ogc/ { proxy_pass http://127.0.0.1:8765; }
```

Then reload nginx and verify: `curl -s https://simplemap.safecast.org/ogc`.

## Example

```bash
# Fixed sensors around Tokyo, first 2
curl "https://simplemap.safecast.org/ogc/collections/sensors/items?bbox=139,35,140,36&limit=2"
```
