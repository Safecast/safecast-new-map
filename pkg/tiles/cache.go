package tiles

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru/v2"
)

const (
	maxMemoryTiles = 256
	maxDiskTiles   = 20000
)

// Cache stores rendered tiles in memory and in a directory on disk.
// The directory is not part of the source database. Entries stay until
// InvalidateBounds removes them or the disk cap evicts the oldest files.
type Cache struct {
	dir string
	mem *lru.Cache[string, []byte]
	mu  sync.Mutex
}

// OpenCache creates dir if needed and returns a cache rooted there.
func OpenCache(dir string) (*Cache, error) {
	if dir == "" {
		return nil, fmt.Errorf("tiles: cache directory is empty")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("tiles: create cache dir: %w", err)
	}
	mem, err := lru.New[string, []byte](maxMemoryTiles)
	if err != nil {
		return nil, fmt.Errorf("tiles: memory cache: %w", err)
	}
	return &Cache{dir: dir, mem: mem}, nil
}

// Get returns a copy of the cached bytes.
func (c *Cache) Get(key string) ([]byte, bool) {
	if c == nil {
		return nil, false
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if body, ok := c.mem.Get(key); ok {
		return append([]byte(nil), body...), true
	}
	body, err := os.ReadFile(c.path(key))
	if err != nil {
		return nil, false
	}
	c.mem.Add(key, append([]byte(nil), body...))
	return body, true
}

// Put stores body under key.
func (c *Cache) Put(key string, body []byte) error {
	if c == nil {
		return fmt.Errorf("tiles: cache is nil")
	}
	stored := append([]byte(nil), body...)
	path := c.path(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("tiles: create tile dir: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, stored, 0o644); err != nil {
		return fmt.Errorf("tiles: write tile: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("tiles: rename tile: %w", err)
	}
	c.mem.Add(key, stored)
	c.evictOldestLocked()
	return nil
}

// InvalidateBounds drops every cached variant whose tile meets b.
// Variants differ by day window and speed, so the whole tile directory goes.
func (c *Cache) InvalidateBounds(b Bounds) error {
	// An inverted box means the caller had no points. A single-point upload
	// has equal edges and still has to drop the tile under that point.
	if c == nil || b.MinLat > b.MaxLat || b.MinLon > 180 || b.MaxLon < -180 {
		return nil
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range c.mem.Keys() {
		id, ok := tileIDFromKey(key)
		if !ok {
			continue
		}
		tb, ok := TileBounds(id)
		if ok && Overlaps(tb, b) {
			c.mem.Remove(key)
		}
	}
	entries, err := os.ReadDir(c.dir)
	if err != nil {
		return fmt.Errorf("tiles: read cache: %w", err)
	}
	for _, zent := range entries {
		if !zent.IsDir() {
			continue
		}
		z, err := strconv.Atoi(zent.Name())
		if err != nil {
			continue
		}
		xs, err := os.ReadDir(filepath.Join(c.dir, zent.Name()))
		if err != nil {
			return fmt.Errorf("tiles: read zoom dir: %w", err)
		}
		for _, xent := range xs {
			if !xent.IsDir() {
				continue
			}
			x, err := strconv.Atoi(xent.Name())
			if err != nil {
				continue
			}
			ys, err := os.ReadDir(filepath.Join(c.dir, zent.Name(), xent.Name()))
			if err != nil {
				return fmt.Errorf("tiles: read x dir: %w", err)
			}
			for _, yent := range ys {
				if !yent.IsDir() {
					continue
				}
				y, err := strconv.Atoi(yent.Name())
				if err != nil {
					continue
				}
				id := TileID{Z: z, X: x, Y: y}
				tb, ok := TileBounds(id)
				if !ok || !Overlaps(tb, b) {
					continue
				}
				if err := os.RemoveAll(filepath.Join(c.dir, zent.Name(), xent.Name(), yent.Name())); err != nil {
					return fmt.Errorf("tiles: remove tile: %w", err)
				}
			}
		}
	}
	return nil
}

func (c *Cache) path(key string) string {
	return filepath.Join(c.dir, filepath.FromSlash(key))
}

func (c *Cache) evictOldestLocked() {
	type aged struct {
		path string
		mod  time.Time
	}
	var files []aged
	_ = filepath.WalkDir(c.dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || strings.HasSuffix(path, ".tmp") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		files = append(files, aged{path: path, mod: info.ModTime()})
		return nil
	})
	if len(files) <= maxDiskTiles {
		return
	}
	sort.Slice(files, func(i, j int) bool { return files[i].mod.Before(files[j].mod) })
	keep := maxDiskTiles * 9 / 10
	for i := 0; i < len(files)-keep; i++ {
		_ = os.Remove(files[i].path)
		rel, err := filepath.Rel(c.dir, files[i].path)
		if err == nil {
			c.mem.Remove(filepath.ToSlash(rel))
		}
	}
}

func tileIDFromKey(key string) (TileID, bool) {
	parts := strings.Split(key, "/")
	if len(parts) < 3 {
		return TileID{}, false
	}
	z, err1 := strconv.Atoi(parts[0])
	x, err2 := strconv.Atoi(parts[1])
	y, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil {
		return TileID{}, false
	}
	return TileID{Z: z, X: x, Y: y}, true
}
