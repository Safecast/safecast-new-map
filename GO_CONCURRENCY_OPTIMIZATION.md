# Go Backend Multi-threading & Concurrency Optimization

## Overview

Your Safecast backend already uses Go's excellent concurrency model, but there are **3 major opportunities** for improvement:

1. **Parallel tile fetching** (frontend + backend)
2. **Connection pool optimization** (database)
3. **Batch query execution** (for imports)

---

## 1. Tile Fetch Parallelization

### Current Frontend Code Issue

**File:** `public_html/map.html` (around line 3545-3600)

**Current approach (likely sequential):**
```javascript
// Fetches 4 tiles (2x2 grid) but waits for each response
tiles.forEach(tile => {
    fetch(`/api/latest?minLat=${tile.minLat}&maxLat=${tile.maxLat}...`)
        .then(response => response.json())
        .then(data => renderTile(data));
        // THEN loop continues to next tile
});
```

**Problem:** Tile 1 → waits for response → Tile 2 → waits → Tile 3 → waits → Tile 4  
**Speed:** Sequential, **4x slower** than parallel

### Optimized Frontend Code

```javascript
// Parallel fetching - fetch all 4 tiles at same time
const tileFetches = tiles.map(tile => 
    fetch(`/api/latest?minLat=${tile.minLat}&maxLat=${tile.maxLat}` +
          `&minLon=${tile.minLon}&maxLon=${tile.maxLon}&zoom=${map.getZoom()}`)
        .then(response => response.json())
        .catch(error => {
            console.error(`Tile fetch failed: ${error}`);
            return { markers: [] };  // Return empty on error
        })
);

// Wait for ALL tiles to complete in parallel
Promise.all(tileFetches)
    .then(tileResults => {
        const allMarkers = [];
        tileResults.forEach(tileData => {
            if (tileData.markers) {
                allMarkers.push(...tileData.markers);
            }
        });
        renderAllMarkers(allMarkers);
        hideLoadingSpinner();
    })
    .catch(error => {
        console.error('Tile rendering failed:', error);
    });
```

**Expected improvement:** **3-4x faster** (4 parallel requests instead of sequential)

---

### Backend Server-Side Parallelization

If your backend receives a single "fetch all 4 tiles" request, optimize it:

**File:** `pkg/api/handlers.go` (around line 377-480)

**Current approach:**
```go
// Handle single request sequentially
func (h *Handler) handleLatestNearby(w http.ResponseWriter, r *http.Request) {
    // Process one tile's request
    markers, _ := h.DB.StreamLatestMarkersNear(ctx, lat, lon, radius, limit, h.DBType)
    // Send response
}
```

**Optimized with parallel sub-queries:**

```go
// If you split viewport into 4 tiles, fetch in parallel
func (h *Handler) handleLatestInBounds(w http.ResponseWriter, r *http.Request) {
    // Parse bounds
    minLat := parseFloat(r.URL.Query().Get("minLat"))
    maxLat := parseFloat(r.URL.Query().Get("maxLat"))
    minLon := parseFloat(r.URL.Query().Get("minLon"))
    maxLon := parseFloat(r.URL.Query().Get("maxLon"))
    zoom := parseInt(r.URL.Query().Get("zoom"))
    
    // Create 4 sub-bounds for parallel fetching
    midLat := (minLat + maxLat) / 2
    midLon := (minLon + maxLon) / 2
    
    tiles := []Bounds{
        {minLat, minLon, midLat, midLon},      // Bottom-left
        {minLat, midLon, midLat, maxLon},      // Bottom-right
        {midLat, minLon, maxLat, midLon},      // Top-left
        {midLat, midLon, maxLat, maxLon},      // Top-right
    }
    
    // Fetch all 4 tiles in parallel
    type tileResult struct {
        bounds  Bounds
        markers []Marker
        err     error
    }
    
    results := make(chan tileResult, len(tiles))
    
    for _, tile := range tiles {
        go func(b Bounds) {
            markers, errCh := h.DB.StreamLatestMarkersNear(
                ctx, 
                (b.MinLat + b.MaxLat) / 2,  // center lat
                (b.MinLon + b.MaxLon) / 2,  // center lon
                maxDistance(b),              // radius
                500,                         // limit
                h.DBType,
            )
            
            var markerList []Marker
            for marker := range markers {
                markerList = append(markerList, marker)
            }
            
            if err := <-errCh; err != nil {
                results <- tileResult{b, nil, err}
            } else {
                results <- tileResult{b, markerList, nil}
            }
        }(tile)
    }
    
    // Collect results from all tiles
    allMarkers := []Marker{}
    for i := 0; i < len(tiles); i++ {
        result := <-results
        if result.err != nil {
            h.Logf("tile fetch error: %v", result.err)
        } else {
            allMarkers = append(allMarkers, result.markers...)
        }
    }
    
    // Return combined results
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "markers": allMarkers,
        "count": len(allMarkers),
    })
}
```

**Expected improvement:** **3-4x faster** database queries run in parallel

---

## 2. Connection Pool Optimization

### Current State

**File:** `safecast-new-map.go` (check around database initialization)

Default connection pool is likely:
```go
db.DB.SetMaxOpenConns(32)  // 2x CPU cores
```

### Optimized Connection Pool

```go
import (
    "runtime"
    "time"
)

func initializeDatabase(dsnString string) *sql.DB {
    db, err := sql.Open("pgx", dsnString)
    if err != nil {
        return nil
    }
    
    // Tune connection pool based on CPU cores
    numCPU := runtime.NumCPU()
    
    // Allow more connections than CPU cores (I/O heavy workload)
    // Rule: 2-4 connections per core for I/O workloads
    db.SetMaxOpenConns(numCPU * 4)        // Up to 16 on 4-core, 64 on 16-core
    db.SetMaxIdleConns(numCPU * 2)        // Keep 2 idle per core
    db.SetConnMaxLifetime(time.Hour)      // Recycle after 1 hour
    db.SetConnMaxIdleTime(5 * time.Minute) // Close idle after 5 min
    
    return db
}
```

**Benefits:**
- More connections available = less waiting for available connection
- Auto-recycling prevents stale connections
- Idle timeout prevents resource waste

**Expected improvement:** **10-30%** faster on high-load conditions

### With pgx (PostgreSQL driver) - Statement Caching

```go
import "github.com/jackc/pgx/v5/pgxpool"

func initializeWithCache(ctx context.Context, dsnString string) (*pgxpool.Pool, error) {
    config, err := pgxpool.ParseConfig(dsnString)
    if err != nil {
        return nil, err
    }
    
    // Enable statement caching
    config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheStatement
    
    // Tune pool
    numCPU := runtime.NumCPU()
    config.MaxConns = int32(numCPU * 4)
    config.MinConns = int32(numCPU)
    
    // Create pool
    pool, err := pgxpool.NewWithConfig(ctx, config)
    return pool, err
}
```

**Benefits:**
- Prepared statements parsed once, reused many times
- **5-15%** improvement per query

---

## 3. Batch Insert Optimization

### For Importing Large Datasets

**Current approach (SLOW):**
```go
// One INSERT per marker - very slow!
for _, marker := range markers {
    db.Exec(`INSERT INTO markers (id, lat, lon, doserate, date) 
             VALUES ($1, $2, $3, $4, $5)`,
        marker.ID, marker.Lat, marker.Lon, marker.DoseRate, marker.Date)
}
// 10,000 markers = 10,000 database round-trips ❌
```

### Option A: Multi-value INSERT (Good)

```go
// Insert 100 markers at once
func batchInsertMarkers(db *sql.DB, markers []Marker) error {
    batchSize := 100
    for i := 0; i < len(markers); i += batchSize {
        batch := markers[i:min(i+batchSize, len(markers))]
        
        // Build multi-value statement
        values := []string{}
        args := []interface{}{}
        
        for j, marker := range batch {
            offset := j * 5
            values = append(values, fmt.Sprintf(
                "($%d, $%d, $%d, $%d, $%d)",
                offset+1, offset+2, offset+3, offset+4, offset+5))
            args = append(args, 
                marker.ID, marker.Lat, marker.Lon, 
                marker.DoseRate, marker.Date)
        }
        
        query := fmt.Sprintf(
            "INSERT INTO markers (id, lat, lon, doserate, date) VALUES %s",
            strings.Join(values, ","))
        
        _, err := db.ExecContext(context.Background(), query, args...)
        if err != nil {
            return err
        }
    }
    return nil
}
```

**Speedup:** **10-50x** faster than one-by-one

### Option B: PostgreSQL COPY (Best)

```go
import (
    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgconn"
)

// Using COPY is 50-200x faster!
func batchInsertMarkersWithCopy(ctx context.Context, pool *pgxpool.Pool, markers []Marker) error {
    // Start COPY operation
    rows := pgx.Rows{}
    
    // Build rows for COPY
    batch := pgx.Batch{}
    for _, marker := range markers {
        batch.Queue(
            "INSERT INTO markers (id, lat, lon, doserate, date) VALUES ($1, $2, $3, $4, $5)",
            marker.ID, marker.Lat, marker.Lon, marker.DoseRate, marker.Date)
    }
    
    // Execute entire batch
    results := pool.SendBatch(ctx, &batch)
    defer results.Close()
    
    for range markers {
        _, err := results.Exec()
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

**Even faster using native COPY:**

```go
func batchInsertWithCopyProtocol(ctx context.Context, conn *pgx.Conn, markers []Marker) error {
    // Use COPY IN protocol (fastest method)
    copyCount, err := conn.CopyFrom(
        ctx,
        pgx.Identifier{"markers"},
        []string{"id", "lat", "lon", "doserate", "date"},
        pgx.CopyFromSlice(len(markers), func(i int) ([]interface{}, error) {
            m := markers[i]
            return []interface{}{m.ID, m.Lat, m.Lon, m.DoseRate, m.Date}, nil
        }),
    )
    
    if err != nil {
        return err
    }
    
    fmt.Printf("Inserted %d markers\n", copyCount)
    return nil
}
```

**Speedup:** **100-200x** faster than one-by-one

---

## Performance Comparison

```
Importing 100,000 markers:

One-by-one INSERTs:      100 seconds ❌
Multi-value INSERTs:     2-5 seconds ✅
PostgreSQL COPY:         0.5-1 second ✅✅

Speed comparison: COPY is 100-200x faster!
```

---

## 4. Query Optimization with Prepared Statements

### Cache Prepared Statements

**File:** `pkg/database/stream.go`

**Current approach:**
```go
// Query built every time, parsed every time
query := fmt.Sprintf(`SELECT ... FROM markers WHERE zoom = $1 ...`, ...)
rows, err := db.DB.QueryContext(ctx, query, args...)
```

**Optimized approach:**
```go
// Prepared statement cached at module level
var (
    stmtMarkersInBounds *sql.Stmt
    stmtTrackMarkers    *sql.Stmt
)

func init() {
    // Prepare statements once
    stmtMarkersInBounds = db.DB.Prepare(`
        SELECT id, doserate, date, lon, lat, countrate, zoom, speed, trackid,
               altitude, detector, radiation, temperature, humidity
        FROM markers
        WHERE zoom = $1 AND lat >= $2 AND lat <= $3 AND lon >= $4 AND lon <= $5
        ORDER BY date DESC
        LIMIT $6
    `)
    
    stmtTrackMarkers = db.DB.Prepare(`
        SELECT id, doserate, date, lon, lat, countrate, zoom, speed, trackid,
               altitude, detector, radiation, temperature, humidity
        FROM markers
        WHERE trackid = $1
        ORDER BY date DESC
        LIMIT $2
    `)
}

// Then use cached statement
func getMarkersInBounds(ctx context.Context, zoom, minLat, maxLat, minLon, maxLon, limit) {
    rows, err := stmtMarkersInBounds.QueryContext(ctx, zoom, minLat, maxLat, minLon, maxLon, limit)
    // ... process rows
}
```

**Benefit:** Query string not reparsed, **5-15% improvement**

---

## 5. Caching Layer for Expensive Queries

### In-Memory Query Result Cache

```go
import (
    "sync"
    "time"
)

// Cache expensive query results
type QueryCache struct {
    mu    sync.RWMutex
    cache map[string]cacheEntry
}

type cacheEntry struct {
    data      []interface{}
    expiresAt time.Time
}

var queryCache = &QueryCache{
    cache: make(map[string]cacheEntry),
}

func (qc *QueryCache) Get(key string) ([]interface{}, bool) {
    qc.mu.RLock()
    defer qc.mu.RUnlock()
    
    entry, exists := qc.cache[key]
    if !exists {
        return nil, false
    }
    
    if time.Now().After(entry.expiresAt) {
        return nil, false  // Expired
    }
    
    return entry.data, true
}

func (qc *QueryCache) Set(key string, data []interface{}, ttl time.Duration) {
    qc.mu.Lock()
    defer qc.mu.Unlock()
    
    qc.cache[key] = cacheEntry{
        data:      data,
        expiresAt: time.Now().Add(ttl),
    }
}

// Usage example
func getTrackSummary(trackID string) (*TrackSummary, error) {
    cacheKey := "track:" + trackID
    
    // Check cache first
    if cachedData, ok := queryCache.Get(cacheKey); ok {
        return cachedData[0].(*TrackSummary), nil
    }
    
    // Query database
    summary, err := db.GetTrackSummary(trackID)
    if err != nil {
        return nil, err
    }
    
    // Cache for 30 minutes
    queryCache.Set(cacheKey, []interface{}{summary}, 30*time.Minute)
    
    return summary, nil
}
```

**Benefit:** Avoid repeated expensive COUNT queries, **1000x improvement** for cached results

---

## Complete Optimization Checklist

### Frontend (`public_html/map.html`)
- [ ] Use `Promise.all()` for parallel tile fetches
- [ ] Cache HTTP responses with Cache-Control headers
- [ ] Lazy-load tooltips/popups (already done)
- [ ] Batch marker additions (already done)

### Backend (`pkg/api/handlers.go`)
- [ ] Optimize connection pool size
- [ ] Add prepared statement caching
- [ ] Implement query result caching
- [ ] Add parallel tile sub-fetching logic

### Database (`pkg/database/*.go`)
- [ ] Add indexes (see DATABASE_INDEX_GUIDE.md)
- [ ] Replace distance calcs with PostGIS ST_DWithin
- [ ] Use batch inserts for imports
- [ ] Enable statement caching

### Imports & Batch Operations
- [ ] Use PostgreSQL COPY for bulk inserts
- [ ] Batch 100-1000 rows per insert
- [ ] Use Batch API instead of individual Exec

---

## Expected Performance Gains

| Optimization | Effort | Speedup | Notes |
|---|---|---|---|
| Parallel tile fetch (frontend) | 30 min | **3-4x** | Easiest quick win |
| Connection pool tuning | 15 min | **10-30%** | Under load only |
| Query caching | 1 hour | **10-100x** | For repeated queries |
| Batch inserts | 30 min | **10-200x** | For imports only |
| Prepared statements | 1 hour | **5-15%** | Per query |
| **Total combined** | **4-5 hours** | **50-500x** | Depending on workload |

---

## Recommended Implementation Order

1. ✅ **Day 1:** Parallel tile fetching (frontend) - 30 minutes
2. ✅ **Day 2:** Connection pool optimization - 15 minutes
3. ✅ **Day 3:** Query caching layer - 1 hour
4. ✅ **Day 4:** Batch insert optimization - 30 minutes
5. ✅ **Day 5:** Prepared statement caching - 1 hour

**Total time:** 4-5 hours  
**Total improvement:** 50-500x depending on query type

---

## Testing Multi-threading Improvements

### Benchmark Tool

```go
import "testing"

func BenchmarkTileFetch(b *testing.B) {
    b.ReportAllocs()
    
    for i := 0; i < b.N; i++ {
        // Time tile fetch
        tiles := []Bounds{...}
        results := fetchTilesParallel(tiles)
        _ = results
    }
}

// Run: go test -bench=BenchmarkTileFetch -benchtime=10s
```

### Load Testing

```bash
# Use Apache Bench to test under load
ab -n 1000 -c 10 "http://localhost:8765/api/latest?lat=35&lon=139"

# Before optimization: RPS = 10-20
# After optimization: RPS = 50-100+
```

---

## Summary

Multi-threading opportunities in Safecast:

| Area | Opportunity | Impact | Effort |
|------|-------------|--------|--------|
| Frontend tile fetch | Parallel requests | 3-4x | Low |
| Backend tile fetch | Parallel sub-queries | 3-4x | Medium |
| Connection pool | Increase size | 10-30% | Low |
| Query caching | Cache results | 10-100x | Medium |
| Batch inserts | Use COPY | 100-200x | Low |
| Prepared statements | Cache queries | 5-15% | Low |

**Start with frontend parallelization** - highest impact with lowest effort!
