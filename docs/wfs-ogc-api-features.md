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

### 1. Add the connection

1. Layer → Add Layer → **Add WFS / OGC API - Features Layer…**
   (or Data Source Manager → **WFS / OGC API - Features**).
2. **New** connection:
   - **Name:** e.g. `Safecast`
   - **URL:** `https://simplemap.safecast.org/ogc`
   - **Version:** must be **`OGC API - Features`** (not WFS 1.0/1.1/2.0 — those
     expect classic XML GetCapabilities, which this endpoint does not serve).
   - Authentication: **No Authentication**.
3. **OK** → **Connect**. The `sensors — Safecast fixed sensors (live)`
   collection appears.
4. Select it → **Add** (**once** — clicking Add repeatedly stacks duplicate
   layers) → **Close**.

Most live sensors are in Japan; pan there if the canvas looks empty. The layer
re-queries by bounding box as you pan/zoom.

### 2. Style by radiation level

Each point carries the current `value` + `unit`.

1. Right-click the layer → **Properties → Symbology**.
2. Top dropdown: **Graduated**. **Value:** `value`.
3. Pick a **Color ramp** (e.g. white→red), set **Classes** to ~7.
4. Click **Classify** (required — the class list is empty until you do),
   then **Apply** → **OK**.

**Mode:** *Equal Count (Quantile)* highlights relative highs; *Natural Breaks
(Jenks)* or *Equal Interval* better reflects absolute levels.

> **Unit caveat.** The `value` field mixes units across sensor types (`unit` is
> often `lnd_7318c_cps`, but varies by device). A single graduated ramp
> therefore compares slightly different quantities — fine for a quick overview,
> but do not read the colours as calibrated µSv/h across all devices. Filter or
> facet by `unit` (or `type`) for like-for-like comparison.

### Troubleshooting

- **"Download of landing page failed: Missing information in response"** — the
  connection Version is set to a classic WFS version. Set it to
  **OGC API - Features**. (QGIS also requires the landing page to expose a
  `service-desc` link to the OpenAPI doc at `/ogc/api`; the server provides
  this.)
- **Empty canvas** — pan to Japan, or check the layer's bbox filter isn't
  excluding everything.
- **Multiple identical `sensors` layers** — Add was clicked more than once;
  right-click the extras → **Remove Layer**.

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

## Deployment note (nginx / CloudFront)

**No dedicated nginx block is required.** Both site configs
(`origin-simplemap.safecast.org` and `simplemap.safecast.org`) already have a
`location /` catch-all that proxies everything to unified-server on port 8765,
which covers `/ogc`. (Only paths with their own specific `location` rule — the
`/api/*`, `/mcp*`, `/assistant/`, `/docs/` blocks — bypass the catch-all.)

Verify after deploy:

```bash
curl -s https://simplemap.safecast.org/ogc
curl -s "https://simplemap.safecast.org/ogc/collections/sensors/items?bbox=139,35,140,36&limit=1"
```

**CloudFront host caveat.** `simplemap.safecast.org` is fronted by CloudFront,
which rewrites the `Host` header to the origin (`origin-simplemap.safecast.org`)
before it reaches nginx. The server builds its self/next/items links from that
host, so discovery links in responses resolve to `origin-simplemap.safecast.org`
rather than `simplemap.safecast.org`. This works for QGIS because the origin host
is publicly reachable with a valid cert — link-following just bypasses the CDN
cache. To make links use the public host, honor `X-Forwarded-Host` in
`ogcBaseURL()` (and have nginx/CloudFront forward the original host).

## Example

```bash
# Fixed sensors around Tokyo, first 2
curl "https://simplemap.safecast.org/ogc/collections/sensors/items?bbox=139,35,140,36&limit=2"
```
