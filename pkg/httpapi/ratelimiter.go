// ratelimiter.go implements per-IP rate limiting with a dedicated goroutine per IP.
package httpapi

import (
	"context"
	"time"
)

// RequestKind distinguishes lightweight metadata calls from heavy streaming responses.
// Heavy requests get a cooldown to prevent a single client from flooding the server.
type RequestKind int

const (
	RequestGeneral RequestKind = iota // Inexpensive metadata; still queued per IP.
	RequestHeavy                       // Streams large payloads; cooldown applied after each.
)

// RateLimiter coordinates per-IP request sequencing. Each IP has its own worker;
// heavy requests are throttled after a burst so normal browsing stays responsive.
type RateLimiter struct {
	heavyCooldown time.Duration
	heavyBurst    int
	noticeFloor   time.Duration
	requests      chan keyedRequest
	now           func() time.Time
}

// completionNotice is sent back when a handler finishes so the worker can update heavy history.
type completionNotice struct {
	kind        RequestKind
	completedAt time.Time
}

// keyedRequest carries the client IP and the actual request.
type keyedRequest struct {
	ip  string
	req ipRequest
}

type ipRequest struct {
	ctx      context.Context
	kind     RequestKind
	arrived  time.Time
	response chan acquireResponse
}

type acquireResponse struct {
	release      chan struct{}
	wait         bool
	waitDuration time.Duration
	err          error
}

// Permit represents an acquired slot. The handler must call Release when done so the next request can proceed.
type Permit struct {
	release      chan struct{}
	WaitNotice   bool
	WaitDuration time.Duration
}

// Release signals the limiter that the request is done.
func (p *Permit) Release() {
	if p == nil || p.release == nil {
		return
	}
	close(p.release)
	p.release = nil
}

// NewRateLimiter constructs a limiter. heavyCooldown is the minimum time between heavy requests per IP.
func NewRateLimiter(heavyCooldown time.Duration) *RateLimiter {
	limiter := &RateLimiter{
		heavyCooldown: heavyCooldown,
		heavyBurst:    8,
		noticeFloor:   50 * time.Millisecond,
		requests:      make(chan keyedRequest),
		now:            time.Now,
	}
	go limiter.loop()
	return limiter
}

// Acquire reserves a slot for the given IP and request kind. Returns nil, nil if the limiter is nil.
func (l *RateLimiter) Acquire(ctx context.Context, ip string, kind RequestKind) (*Permit, error) {
	if l == nil {
		return nil, nil
	}
	respCh := make(chan acquireResponse, 1)
	req := ipRequest{
		ctx:      ctx,
		kind:     kind,
		arrived:  l.now(),
		response: respCh,
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case l.requests <- keyedRequest{ip: ip, req: req}:
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case resp := <-respCh:
		if resp.err != nil {
			return nil, resp.err
		}
		return &Permit{
			release:      resp.release,
			WaitNotice:   resp.wait,
			WaitDuration: resp.waitDuration,
		}, nil
	}
}

func (l *RateLimiter) loop() {
	workers := make(map[string]chan ipRequest)
	for keyed := range l.requests {
		ch, ok := workers[keyed.ip]
		if !ok {
			ch = make(chan ipRequest)
			workers[keyed.ip] = ch
			go l.runIPWorker(keyed.ip, ch)
		}
		select {
		case ch <- keyed.req:
		case <-keyed.req.ctx.Done():
			keyed.req.response <- acquireResponse{err: keyed.req.ctx.Err()}
		}
	}
}

func (l *RateLimiter) runIPWorker(ip string, requests <-chan ipRequest) {
	var heavyHistory []time.Time
	activeHeavy := 0
	doneCh := make(chan completionNotice)
	for {
		select {
		case req, ok := <-requests:
			if !ok {
				return
			}
			select {
			case <-req.ctx.Done():
				req.response <- acquireResponse{err: req.ctx.Err()}
				continue
			default:
			}
			now := l.now()
			queueWait := now.Sub(req.arrived)
			if queueWait < 0 {
				queueWait = 0
			}
			totalWait := queueWait
			if req.kind == RequestHeavy && l.heavyCooldown > 0 && l.heavyBurst > 0 {
				cutoff := now.Add(-l.heavyCooldown)
				dst := heavyHistory[:0]
				for _, ts := range heavyHistory {
					if ts.After(cutoff) {
						dst = append(dst, ts)
					}
				}
				heavyHistory = dst
				heavySeen := activeHeavy + len(heavyHistory)
				if heavySeen >= l.heavyBurst {
					readyAt := heavyHistory[0].Add(l.heavyCooldown)
					now = l.now()
					if now.Before(readyAt) {
						cooldownWait := readyAt.Sub(now)
						timer := time.NewTimer(cooldownWait)
						select {
						case <-req.ctx.Done():
							if !timer.Stop() {
								<-timer.C
							}
							req.response <- acquireResponse{err: req.ctx.Err()}
							continue
						case <-timer.C:
							totalWait += cooldownWait
						}
					}
				}
			}
			release := make(chan struct{})
			effectiveWait := totalWait
			if effectiveWait < l.noticeFloor {
				effectiveWait = 0
			}
			resp := acquireResponse{
				release:      release,
				wait:         effectiveWait > 0,
				waitDuration: totalWait,
			}
			select {
			case <-req.ctx.Done():
				req.response <- acquireResponse{err: req.ctx.Err()}
				continue
			case req.response <- resp:
			}
			if req.kind == RequestHeavy {
				activeHeavy++
			}
			go l.waitForRelease(req, release, doneCh)
		case done := <-doneCh:
			if done.kind == RequestHeavy {
				if activeHeavy > 0 {
					activeHeavy--
				}
				heavyHistory = append(heavyHistory, done.completedAt)
				cutoff := l.now().Add(-l.heavyCooldown)
				dst := heavyHistory[:0]
				for _, ts := range heavyHistory {
					if ts.After(cutoff) {
						dst = append(dst, ts)
					}
				}
				heavyHistory = dst
			}
		}
	}
}

func (l *RateLimiter) waitForRelease(req ipRequest, release <-chan struct{}, done chan<- completionNotice) {
	select {
	case <-release:
	case <-req.ctx.Done():
		<-release
	}
	done <- completionNotice{kind: req.kind, completedAt: l.now()}
}
