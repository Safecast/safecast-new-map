// chat_limiter.go limits how often the AI assistant (/chat) can be used.
//
// Every /chat request can call a paid model API, so without a limit anyone
// could script the endpoint and run up the bill. Three limits apply:
//
//   - per client, per minute   (CHAT_LIMIT_PER_MINUTE, default 6)
//   - per client, per day      (CHAT_LIMIT_PER_DAY, default 60)
//   - all clients, per day     (CHAT_LIMIT_GLOBAL_PER_DAY, default 3000)
//
// The global daily limit is the hard ceiling on spend: it cannot be avoided by
// changing IP address or forging headers. Set any limit to 0 to turn it off.
// Request bodies are capped at CHAT_MAX_BODY_BYTES (default 131072) so a
// single request cannot carry an oversized conversation history.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type chatWindow struct {
	minuteStart time.Time
	minuteCount int
	dayStart    time.Time
	dayCount    int
}

type chatLimiter struct {
	mu             sync.Mutex
	perMinute      int
	perDay         int
	globalPerDay   int
	clients        map[string]*chatWindow
	globalDayStart time.Time
	globalCount    int
	lastSweep      time.Time
	now            func() time.Time
}

func newChatLimiter(perMinute, perDay, globalPerDay int) *chatLimiter {
	return &chatLimiter{
		perMinute:    perMinute,
		perDay:       perDay,
		globalPerDay: globalPerDay,
		clients:      make(map[string]*chatWindow),
		now:          time.Now,
	}
}

// allow records a request for key and reports whether it may proceed. When it
// may not, it returns how long until the blocking limit resets and a reason
// code ("minute", "day" or "global").
func (l *chatLimiter) allow(key string) (bool, time.Duration, string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()

	// Drop expired client windows now and then so the map can't grow forever.
	if now.Sub(l.lastSweep) > 10*time.Minute {
		for k, w := range l.clients {
			if now.Sub(w.dayStart) >= 24*time.Hour {
				delete(l.clients, k)
			}
		}
		l.lastSweep = now
	}

	if now.Sub(l.globalDayStart) >= 24*time.Hour {
		l.globalDayStart = now
		l.globalCount = 0
	}
	w := l.clients[key]
	if w == nil {
		w = &chatWindow{minuteStart: now, dayStart: now}
		l.clients[key] = w
	}
	if now.Sub(w.minuteStart) >= time.Minute {
		w.minuteStart = now
		w.minuteCount = 0
	}
	if now.Sub(w.dayStart) >= 24*time.Hour {
		w.dayStart = now
		w.dayCount = 0
	}

	if l.globalPerDay > 0 && l.globalCount >= l.globalPerDay {
		return false, l.globalDayStart.Add(24 * time.Hour).Sub(now), "global"
	}
	if l.perDay > 0 && w.dayCount >= l.perDay {
		return false, w.dayStart.Add(24 * time.Hour).Sub(now), "day"
	}
	if l.perMinute > 0 && w.minuteCount >= l.perMinute {
		return false, w.minuteStart.Add(time.Minute).Sub(now), "minute"
	}

	w.minuteCount++
	w.dayCount++
	l.globalCount++
	return true, 0, ""
}

// chatClientKey identifies the client for per-client limits. CloudFront's
// CloudFront-Viewer-Address header is set by CloudFront itself and can't be
// forged by the client, so it is preferred when the distribution forwards it.
// Otherwise this falls back to getClientIP, whose X-Forwarded-For value can be
// set by the client; the global daily limit still applies in that case.
func chatClientKey(r *http.Request) string {
	if v := strings.TrimSpace(r.Header.Get("CloudFront-Viewer-Address")); v != "" {
		if host, _, err := net.SplitHostPort(v); err == nil {
			return host
		}
		// IPv6 values look like "2001:db8::1:46532" without brackets.
		if i := strings.LastIndex(v, ":"); i > 0 {
			return v[:i]
		}
		return v
	}
	return getClientIP(r)
}

// chatLimitMessage returns the visitor-facing message in the page language,
// falling back to English.
func chatLimitMessage(r *http.Request, reason string) string {
	key := "ai_rate_limited"
	if reason == "global" {
		key = "ai_busy_today"
	}
	lang := getPreferredLanguage(r)
	translationsMu.RLock()
	defer translationsMu.RUnlock()
	if v, ok := translations[lang][key]; ok && v != "" {
		return v
	}
	if v, ok := translations["en"][key]; ok && v != "" {
		return v
	}
	if reason == "global" {
		return "The assistant has reached its limit for today. Please try again tomorrow."
	}
	return "Too many questions in a short time. Please wait a moment and try again."
}

// withChatLimits wraps the /chat handler with the limits described above.
func withChatLimits(l *chatLimiter, maxBodyBytes int64, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS preflight is free; let the chat handler answer it.
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}
		ok, retryAfter, reason := l.allow(chatClientKey(r))
		if !ok {
			if reason == "global" {
				log.Printf("[chat] global daily limit reached (%d)", l.globalPerDay)
			}
			secs := int(math.Ceil(retryAfter.Seconds()))
			if secs < 1 {
				secs = 1
			}
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/x-ndjson")
			w.Header().Set("Cache-Control", "no-cache, no-store")
			w.Header().Set("Retry-After", strconv.Itoa(secs))
			w.WriteHeader(http.StatusTooManyRequests)
			// The page shows {"type":"error"} chunks to the visitor.
			data, _ := json.Marshal(chunk{Type: "error", Error: chatLimitMessage(r, reason)})
			fmt.Fprintf(w, "%s\n", data)
			return
		}
		if maxBodyBytes > 0 {
			r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		}
		next(w, r)
	}
}

// newChatLimiterFromEnv builds the limiter from environment variables.
func newChatLimiterFromEnv() (*chatLimiter, int64) {
	l := newChatLimiter(
		getEnvIntOrDefault("CHAT_LIMIT_PER_MINUTE", 6),
		getEnvIntOrDefault("CHAT_LIMIT_PER_DAY", 60),
		getEnvIntOrDefault("CHAT_LIMIT_GLOBAL_PER_DAY", 3000),
	)
	maxBody := int64(getEnvIntOrDefault("CHAT_MAX_BODY_BYTES", 131072))
	log.Printf("[chat] limits: %d/min and %d/day per client, %d/day total, %d byte requests (0 = off)",
		l.perMinute, l.perDay, l.globalPerDay, maxBody)
	return l, maxBody
}
