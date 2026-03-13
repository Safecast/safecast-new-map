// Chat question logging to DuckDB analytics
// Captures user questions from the web-chat and map widget with request metadata

package main

import (
	"log"
	"net"
	"net/http"
	"strings"
)

// logChatQuestion logs a user's chat question to DuckDB asynchronously.
func logChatQuestion(r *http.Request, question, source, model, sessionID string, historyLen int) {
	if !duckDBAvailable() {
		return
	}

	if len(question) > 5000 {
		question = question[:5000]
	}

	ip := getClientIP(r)
	ua := r.Header.Get("User-Agent")
	isMobile, osName, browser := parseUserAgent(ua)
	country := r.Header.Get("CloudFront-Viewer-Country")
	acceptLang := r.Header.Get("Accept-Language")
	if len(acceptLang) > 200 {
		acceptLang = acceptLang[:200]
	}
	referer := r.Header.Get("Referer")
	isCloudFront := r.Header.Get("CloudFront-Viewer-Country") != "" ||
		r.Header.Get("CloudFront-Forwarded-Proto") != "" ||
		r.Header.Get("X-Amz-Cf-Id") != ""

	go func() {
		_, err := duckDB.Exec(`
			INSERT INTO chat_questions (
				question, source, ip_address, user_agent, is_mobile,
				os, browser, country, accept_language, referer,
				session_id, history_length, model, cloudfront
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			question, source, ip, ua, isMobile,
			osName, browser, country, acceptLang, referer,
			sessionID, historyLen, model, isCloudFront,
		)
		if err != nil {
			log.Printf("chat_questions insert error: %v", err)
		}
	}()
}

// getClientIP extracts the client IP from the request, respecting proxy headers.
func getClientIP(r *http.Request) string {
	// X-Forwarded-For may contain multiple IPs: client, proxy1, proxy2
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		ip := strings.TrimSpace(parts[0])
		if ip != "" {
			return ip
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// parseUserAgent extracts mobile/desktop, OS, and browser from User-Agent string.
func parseUserAgent(ua string) (isMobile bool, osName, browser string) {
	lower := strings.ToLower(ua)

	// Mobile detection
	isMobile = strings.Contains(lower, "mobile") ||
		strings.Contains(lower, "android") && !strings.Contains(lower, "tablet") ||
		strings.Contains(lower, "iphone") ||
		strings.Contains(lower, "ipod")

	// OS detection
	switch {
	case strings.Contains(lower, "iphone") || strings.Contains(lower, "ipad") || strings.Contains(lower, "ipod"):
		osName = "iOS"
	case strings.Contains(lower, "android"):
		osName = "Android"
	case strings.Contains(lower, "windows"):
		osName = "Windows"
	case strings.Contains(lower, "macintosh") || strings.Contains(lower, "mac os"):
		osName = "macOS"
	case strings.Contains(lower, "linux"):
		osName = "Linux"
	case strings.Contains(lower, "cros"):
		osName = "ChromeOS"
	default:
		osName = "Unknown"
	}

	// Browser detection (order matters — check specific before generic)
	switch {
	case strings.Contains(lower, "edg/") || strings.Contains(lower, "edge/"):
		browser = "Edge"
	case strings.Contains(lower, "opr/") || strings.Contains(lower, "opera"):
		browser = "Opera"
	case strings.Contains(lower, "firefox/"):
		browser = "Firefox"
	case strings.Contains(lower, "chrome/") && !strings.Contains(lower, "chromium"):
		browser = "Chrome"
	case strings.Contains(lower, "safari/") && !strings.Contains(lower, "chrome"):
		browser = "Safari"
	case strings.Contains(lower, "chromium"):
		browser = "Chromium"
	default:
		browser = "Unknown"
	}

	return
}
