package tiles

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type fakeProc struct {
	name    string
	version string
	delay   time.Duration
	calls   atomic.Int32
	current atomic.Int32
	max     atomic.Int32
}

func (f *fakeProc) Name() string        { return f.name }
func (f *fakeProc) ContentType() string { return "image/png" }
func (f *fakeProc) Version() string     { return f.version }
func (f *fakeProc) Render(ctx context.Context, req Request) (Result, error) {
	f.calls.Add(1)
	now := f.current.Add(1)
	for {
		old := f.max.Load()
		if now <= old || f.max.CompareAndSwap(old, now) {
			break
		}
	}
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			f.current.Add(-1)
			return Result{}, ctx.Err()
		}
	}
	f.current.Add(-1)
	return Result{Body: []byte{byte(req.Tile.X), byte(req.Tile.Y)}, ContentType: "image/png"}, nil
}

func TestSingleflightAndConcurrencyLimit(t *testing.T) {
	cache, err := OpenCache(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	slow := &fakeProc{name: "aggregate", version: "1", delay: 40 * time.Millisecond}
	s := &Service{
		cache: cache,
		procs: map[string]Processor{slow.Name(): slow},
		slots: make(chan struct{}, maxConcurrentRenders),
	}
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := s.Render(context.Background(), "aggregate", Request{Tile: TileID{Z: 1, X: 1, Y: 1}})
			if err != nil {
				t.Errorf("render: %v", err)
			}
		}()
	}
	wg.Wait()
	if slow.calls.Load() != 1 {
		t.Fatalf("cold tile rendered %d times, want 1", slow.calls.Load())
	}

	other := &fakeProc{name: "other", version: "1", delay: 30 * time.Millisecond}
	s.RegisterProcessor(other)
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, err := s.Render(context.Background(), "other", Request{Tile: TileID{Z: 2, X: i, Y: 0}})
			if err != nil {
				t.Errorf("render: %v", err)
			}
		}(i)
	}
	wg.Wait()
	if other.max.Load() > maxConcurrentRenders {
		t.Fatalf("concurrent renders %d, limit %d", other.max.Load(), maxConcurrentRenders)
	}
	if other.calls.Load() != 6 {
		t.Fatalf("distinct tiles rendered %d times, want 6", other.calls.Load())
	}
}

func TestInvalidateLeavesNeighbor(t *testing.T) {
	cache, err := OpenCache(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	keep := TileID{Z: 4, X: 1, Y: 1}
	drop := TileID{Z: 4, X: 3, Y: 3}
	keepKey := cacheKey("aggregate", "1", Request{Tile: keep, Speed: "all"})
	dropKey := cacheKey("aggregate", "1", Request{Tile: drop, Speed: "all"})
	if err := cache.Put(keepKey, []byte("keep")); err != nil {
		t.Fatal(err)
	}
	if err := cache.Put(dropKey, []byte("drop")); err != nil {
		t.Fatal(err)
	}
	bounds, ok := TileBounds(drop)
	if !ok {
		t.Fatal("bounds")
	}
	if err := cache.InvalidateBounds(bounds); err != nil {
		t.Fatal(err)
	}
	if _, ok := cache.Get(dropKey); ok {
		t.Fatal("invalidated tile is still cached")
	}
	body, ok := cache.Get(keepKey)
	if !ok || string(body) != "keep" {
		t.Fatalf("neighbor cache = %q, ok=%v", body, ok)
	}
}
