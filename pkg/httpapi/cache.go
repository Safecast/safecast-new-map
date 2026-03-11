// cache.go implements an in-memory response cache with a single goroutine and channels.
package httpapi

import (
	"context"
	"errors"
	"time"
)

var (
	errCacheDisabled = errors.New("cache disabled")
	errCacheStopped  = errors.New("cache stopped")
	errNoLoader      = errors.New("no loader")
)

// cacheRequest is sent to the cache goroutine for a lookup or populate.
type cacheRequest struct {
	ctx    context.Context
	key    string
	loader func(context.Context) ([]byte, error)
	ttl    time.Duration
	reply  chan cacheResponse
}

// cacheResponse is the reply from the cache goroutine.
type cacheResponse struct {
	data []byte
	err  error
}

// cacheEntry stores cached bytes and expiry; the loop trims stale entries on access.
type cacheEntry struct {
	data    []byte
	expires time.Time
}

// ResponseCache caches expensive API responses in memory. All access is serialized
// in a single goroutine so no mutexes are needed.
type ResponseCache struct {
	ttl      time.Duration
	requests chan cacheRequest
	quit     chan struct{}
	now      func() time.Time
}

// NewResponseCache starts the caching goroutine. If ttl <= 0, returns nil (caching disabled).
func NewResponseCache(ttl time.Duration) *ResponseCache {
	if ttl <= 0 {
		return nil
	}
	cache := &ResponseCache{
		ttl:      ttl,
		requests: make(chan cacheRequest),
		quit:     make(chan struct{}),
		now:      time.Now,
	}
	go cache.loop()
	return cache
}

// Close stops the cache goroutine.
func (c *ResponseCache) Close() {
	if c == nil {
		return
	}
	select {
	case <-c.quit:
		return
	default:
	}
	close(c.quit)
}

// Get returns cached bytes for key or invokes loader to produce them. Uses the cache's default TTL.
func (c *ResponseCache) Get(ctx context.Context, key string, loader func(context.Context) ([]byte, error)) ([]byte, error) {
	return c.GetWithTTL(ctx, key, 0, loader)
}

// GetWithTTL is like Get but allows a per-request TTL override (0 means use default).
func (c *ResponseCache) GetWithTTL(ctx context.Context, key string, ttl time.Duration, loader func(context.Context) ([]byte, error)) ([]byte, error) {
	if c == nil {
		return nil, errCacheDisabled
	}
	req := cacheRequest{
		ctx:    ctx,
		key:    key,
		loader: loader,
		ttl:    ttl,
		reply:  make(chan cacheResponse, 1),
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.quit:
		return nil, errCacheStopped
	case c.requests <- req:
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.quit:
		return nil, errCacheStopped
	case resp := <-req.reply:
		if resp.err != nil {
			return nil, resp.err
		}
		if resp.data == nil {
			return nil, nil
		}
		copyBuf := make([]byte, len(resp.data))
		copy(copyBuf, resp.data)
		return copyBuf, nil
	}
}

// loop runs in a goroutine and serializes all cache access.
func (c *ResponseCache) loop() {
	store := make(map[string]cacheEntry)
	for {
		select {
		case <-c.quit:
			return
		case req := <-c.requests:
			now := c.now()
			if entry, ok := store[req.key]; ok && now.Before(entry.expires) {
				req.reply <- cacheResponse{data: entry.data}
				continue
			}
			if req.loader == nil {
				req.reply <- cacheResponse{err: errNoLoader}
				continue
			}
			data, err := req.loader(req.ctx)
			entryTTL := c.ttl
			if req.ttl > 0 {
				entryTTL = req.ttl
			}
			if err == nil && data != nil {
				if entryTTL > 0 {
					buf := make([]byte, len(data))
					copy(buf, data)
					store[req.key] = cacheEntry{data: buf, expires: now.Add(entryTTL)}
				}
			} else if err != nil {
				delete(store, req.key)
			}
			req.reply <- cacheResponse{data: data, err: err}
		}
	}
}
