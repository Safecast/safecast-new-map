package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// requestClientIP returns the best-effort client IP (X-Forwarded-For or RemoteAddr).
func requestClientIP(r *http.Request) string {
	forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		if c := strings.TrimSpace(parts[0]); c != "" {
			return c
		}
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && strings.TrimSpace(host) != "" {
		return host
	}
	if t := strings.TrimSpace(r.RemoteAddr); t != "" {
		return t
	}
	return ""
}

// geoIPLookup returns approximate lat/lon for an IP via ipapi.co.
func geoIPLookup(ctx context.Context, ip string) (float64, float64, error) {
	if strings.TrimSpace(ip) == "" {
		return 0, 0, fmt.Errorf("missing ip")
	}
	endpoint := fmt.Sprintf("https://ipapi.co/%s/json/", url.PathEscape(ip))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("geoip status %d", resp.StatusCode)
	}
	var payload struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return 0, 0, err
	}
	return payload.Latitude, payload.Longitude, nil
}

// geoIP returns JSON with lat/lon for the client IP when AutoLocateDefault is on.
func (s *Server) geoIP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if !s.Config.AutoLocateDefault {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ip := requestClientIP(r)
	if ip == "" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	lat, lon, err := geoIPLookup(ctx, ip)
	if err != nil {
		log.Printf("geoip lookup failed for %s: %v", ip, err)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if err := json.NewEncoder(w).Encode(map[string]float64{"lat": lat, "lon": lon}); err != nil {
		log.Printf("geoip response encode failed: %v", err)
	}
}
