// handlers_docs.go — GET /api/docs (HTML usage page) and GET /licenses/{mit|cc0}.
package httpapi

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// apiDocs serves the API usage HTML page with placeholders (base URL, archive route, etc.) replaced from config.
func (s *Server) apiDocs(w http.ResponseWriter, r *http.Request) {
	b, err := fs.ReadFile(s.Content, "public_html/api-usage.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	scheme := "http"
	if proto := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto")); proto != "" {
		scheme = strings.ToLower(proto)
	} else if r.TLS != nil {
		scheme = "https"
	}
	host := strings.TrimSpace(r.Host)
	if host == "" {
		if strings.TrimSpace(s.Config.Domain) != "" {
			host = strings.TrimSpace(s.Config.Domain)
		} else {
			host = fmt.Sprintf("localhost:%d", s.Config.Port)
		}
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, host)
	apiRoot := strings.TrimRight(baseURL, "/") + "/api"
	page := string(b)
	page = strings.ReplaceAll(page, "__BASE_URL__", baseURL)
	page = strings.ReplaceAll(page, "__API_ROOT__", apiRoot)
	page = strings.ReplaceAll(page, "__DISPLAY_HOST__", host)
	page = strings.ReplaceAll(page, "__ARCHIVE_ENABLED__", strconv.FormatBool(s.Config.APIDocsArchiveEnabled))
	route := strings.TrimSpace(s.Config.APIDocsArchiveRoute)
	if route == "" {
		route = "/api/json/weekly.tgz"
	}
	page = strings.ReplaceAll(page, "__ARCHIVE_ROUTE__", route)
	freq := strings.TrimSpace(s.Config.APIDocsArchiveFrequency)
	if freq == "" {
		freq = "weekly"
	}
	page = strings.ReplaceAll(page, "__ARCHIVE_FREQUENCY__", freq)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page))
}

// license serves plain-text license files from the embedded FS (LICENSE for mit, LICENSE.CC0 for cc0).
func (s *Server) license(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	rawPath := strings.TrimPrefix(r.URL.Path, "/licenses/")
	if rawPath == r.URL.Path {
		http.NotFound(w, r)
		return
	}
	code := strings.ToLower(strings.Trim(rawPath, "/"))
	var file string
	switch code {
	case "mit":
		file = "LICENSE"
	case "cc0":
		file = "LICENSE.CC0"
	default:
		http.NotFound(w, r)
		return
	}
	data, err := fs.ReadFile(s.Content, file)
	if err != nil {
		s.Logf("license handler: %s read error: %v", file, err)
		http.Error(w, "unable to load license", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.ServeContent(w, r, file, time.Time{}, bytes.NewReader(data))
}
