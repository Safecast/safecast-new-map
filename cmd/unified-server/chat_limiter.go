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
//
// Client identity: the server sits behind CloudFront and nginx, so the TCP
// peer is the proxy, not the visitor. CloudFront sets CloudFront-Viewer-Address
// to the visitor's real address, but the origin is also reachable directly and
// a direct caller can send that header themselves. If CHAT_ORIGIN_SECRET is
// set, the header is trusted only when the request also carries
// X-Origin-Secret with that value (configure CloudFront to add it as a custom
// origin header). Requests without it are keyed by their TCP peer, so direct
// callers share one bucket instead of choosing their own.
package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// chatMaxTrackedClients bounds memory used by per-client windows. When full,
// new clients share one overflow bucket until old windows expire.
const chatMaxTrackedClients = 100000

const chatOverflowKey = "overflow"

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
	maxClients     int
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
		maxClients:   chatMaxTrackedClients,
		clients:      make(map[string]*chatWindow),
		now:          time.Now,
	}
}

func (l *chatLimiter) sweep(now time.Time) {
	for k, w := range l.clients {
		if now.Sub(w.dayStart) >= 24*time.Hour {
			delete(l.clients, k)
		}
	}
	l.lastSweep = now
}

// allow records a request for key and reports whether it may proceed. When it
// may not, it returns how long until the blocking limit resets and a reason
// code ("minute", "day" or "global"). Blocked requests are not counted and do
// not create a client entry.
func (l *chatLimiter) allow(key string) (bool, time.Duration, string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()

	if now.Sub(l.globalDayStart) >= 24*time.Hour {
		l.globalDayStart = now
		l.globalCount = 0
	}
	if l.globalPerDay > 0 && l.globalCount >= l.globalPerDay {
		return false, l.globalDayStart.Add(24 * time.Hour).Sub(now), "global"
	}

	// Drop expired client windows now and then so the map can't grow forever.
	if now.Sub(l.lastSweep) > 10*time.Minute {
		l.sweep(now)
	}

	w := l.clients[key]
	if w == nil {
		if len(l.clients) >= l.maxClients {
			l.sweep(now)
		}
		if len(l.clients) >= l.maxClients {
			key = chatOverflowKey
			w = l.clients[key]
		}
	}
	if w == nil {
		w = &chatWindow{minuteStart: now, dayStart: now}
	}
	if now.Sub(w.minuteStart) >= time.Minute {
		w.minuteStart = now
		w.minuteCount = 0
	}
	if now.Sub(w.dayStart) >= 24*time.Hour {
		w.dayStart = now
		w.dayCount = 0
	}

	if l.perDay > 0 && w.dayCount >= l.perDay {
		return false, w.dayStart.Add(24 * time.Hour).Sub(now), "day"
	}
	if l.perMinute > 0 && w.minuteCount >= l.perMinute {
		return false, w.minuteStart.Add(time.Minute).Sub(now), "minute"
	}

	l.clients[key] = w
	w.minuteCount++
	w.dayCount++
	l.globalCount++
	return true, 0, ""
}

// normalizeIP returns the canonical form of an IP address (with or without a
// port), or "" if s is not an IP. This keeps keys short and stops one address
// from being counted under many spellings.
func normalizeIP(s string) string {
	s = strings.TrimSpace(s)
	if host, _, err := net.SplitHostPort(s); err == nil {
		s = host
	} else if i := strings.LastIndex(s, ":"); i > 0 && strings.Count(s, ":") > 1 && net.ParseIP(s) == nil {
		// CloudFront-Viewer-Address writes IPv6 as "2001:db8::1:46532" (no brackets).
		s = s[:i]
	}
	ip := net.ParseIP(strings.Trim(s, "[]"))
	if ip == nil {
		return ""
	}
	return ip.String()
}

// chatClientKey identifies the client for per-client limits.
func chatClientKey(r *http.Request, originSecret string) string {
	trusted := originSecret == "" ||
		subtle.ConstantTimeCompare([]byte(r.Header.Get("X-Origin-Secret")), []byte(originSecret)) == 1
	if trusted {
		if ip := normalizeIP(r.Header.Get("CloudFront-Viewer-Address")); ip != "" {
			return ip
		}
		if originSecret == "" {
			if ip := normalizeIP(getClientIP(r)); ip != "" {
				return ip
			}
		}
	}
	if ip := normalizeIP(r.RemoteAddr); ip != "" {
		return ip
	}
	return chatOverflowKey
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
func withChatLimits(l *chatLimiter, maxBodyBytes int64, originSecret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS preflight is free; let the chat handler answer it.
		if r.Method == http.MethodOptions {
			next(w, r)
			return
		}
		ok, retryAfter, reason := l.allow(chatClientKey(r, originSecret))
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
func newChatLimiterFromEnv() (*chatLimiter, int64, string) {
	l := newChatLimiter(
		getEnvIntOrDefault("CHAT_LIMIT_PER_MINUTE", 6),
		getEnvIntOrDefault("CHAT_LIMIT_PER_DAY", 60),
		getEnvIntOrDefault("CHAT_LIMIT_GLOBAL_PER_DAY", 3000),
	)
	maxBody := int64(getEnvIntOrDefault("CHAT_MAX_BODY_BYTES", 131072))
	secret := os.Getenv("CHAT_ORIGIN_SECRET")
	log.Printf("[chat] limits: %d/min and %d/day per client, %d/day total, %d byte requests (0 = off); origin secret %s",
		l.perMinute, l.perDay, l.globalPerDay, maxBody, map[bool]string{true: "set", false: "not set"}[secret != ""])
	return l, maxBody, secret
}
