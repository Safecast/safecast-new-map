package tiles

import (
	"context"

	"safecast-new-map/pkg/database"
)

const maxConcurrentRenders = 4

// Service caches processor output and limits how many cold tiles generate at once.
type Service struct {
	cache  *Cache
	procs  map[string]Processor
	agg    *Aggregate
	flight flight
	slots  chan struct{}
}

// New builds a service that serves the aggregate processor from db.
// cacheDir is created if it does not exist. It is not a database file.
func New(db *database.Database, dbType, cacheDir string) (*Service, error) {
	cache, err := OpenCache(cacheDir)
	if err != nil {
		return nil, err
	}
	agg := &Aggregate{DB: db, DBType: dbType}
	s := &Service{
		cache: cache,
		procs: map[string]Processor{agg.Name(): agg},
		agg:   agg,
		slots: make(chan struct{}, maxConcurrentRenders),
	}
	return s, nil
}

// RegisterProcessor adds or replaces a processor. Tests use it to avoid a database.
func (s *Service) RegisterProcessor(p Processor) {
	if s.procs == nil {
		s.procs = map[string]Processor{}
	}
	s.procs[p.Name()] = p
}

// Render returns the cached tile or generates it once for this key.
func (s *Service) Render(ctx context.Context, name string, req Request) (Result, error) {
	p, ok := s.procs[name]
	if !ok {
		return Result{}, ErrUnknownProcessor
	}
	req = normalizeRequest(req)
	key := cacheKey(name, p.Version(), req)
	if body, ok := s.cache.Get(key); ok {
		return Result{Body: body, ContentType: p.ContentType()}, nil
	}
	return s.flight.Do(key, func() (Result, error) {
		if body, ok := s.cache.Get(key); ok {
			return Result{Body: body, ContentType: p.ContentType()}, nil
		}
		select {
		case s.slots <- struct{}{}:
		case <-ctx.Done():
			return Result{}, ctx.Err()
		}
		defer func() { <-s.slots }()
		res, err := p.Render(ctx, req)
		if err != nil {
			return Result{}, err
		}
		if res.ContentType == "" {
			res.ContentType = p.ContentType()
		}
		if err := s.cache.Put(key, res.Body); err != nil {
			return Result{}, err
		}
		return res, nil
	})
}

// InvalidateBounds drops cached overview tiles that meet b.
// Tiles that do not meet b stay on disk.
func (s *Service) InvalidateBounds(b Bounds) error {
	if s == nil || s.cache == nil {
		return nil
	}
	return s.cache.InvalidateBounds(b)
}
