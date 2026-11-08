package tiles

import "sync"

// flight collapses concurrent renders of the same cache key into one call.
type flight struct {
	mu sync.Mutex
	m  map[string]*flightCall
}

type flightCall struct {
	wg  sync.WaitGroup
	val Result
	err error
}

func (f *flight) Do(key string, fn func() (Result, error)) (Result, error) {
	f.mu.Lock()
	if f.m == nil {
		f.m = make(map[string]*flightCall)
	}
	if c, ok := f.m[key]; ok {
		f.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	c := &flightCall{}
	c.wg.Add(1)
	f.m[key] = c
	f.mu.Unlock()

	defer func() {
		c.wg.Done()
		f.mu.Lock()
		delete(f.m, key)
		f.mu.Unlock()
	}()
	c.val, c.err = fn()
	return c.val, c.err
}
