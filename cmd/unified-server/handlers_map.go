package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"

	"safecast-new-map/pkg/database"
)

// =====================
// WEB
// =====================
// =====================
// WEB  — главная карта
// =====================
// marshalTemplateJS encodes the provided value into JSON and tags it as safe
// JavaScript for html/template. We return template.JS so scripts can embed the
// literal without tripping the context analyser, following "A little copying is
// better than a little dependency" by keeping the helper tiny and local.
func marshalTemplateJS(value interface{}) (template.JS, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return template.JS(""), err
	}
	return template.JS(payload), nil
}

// translationsForLang returns a filtered copy of the translations map containing
// only the requested language and "en" (fallback). This keeps the embedded JSON
// small (~30KB instead of ~850KB for all 30 languages) so the AI widget stays
// well within Claude's context limit.
// parseDebugAllowlist converts the comma-separated flag payload into a lookup
// map. Returning a new map keeps the zero value useful and avoids hidden shared
// state that could surprise future callers.
func parseDebugAllowlist(raw string) map[string]struct{} {
	allow := make(map[string]struct{})
	for _, part := range strings.Split(raw, ",") {
		ip := strings.TrimSpace(part)
		if ip == "" {
			continue
		}
		allow[ip] = struct{}{}
	}
	return allow
}

// requestClientIP mirrors the API rate-limiter helper so template handlers can
// decide whether to surface diagnostics. We respect X-Forwarded-For first so the
// overlay still works when the service sits behind a proxy.
func requestClientIP(r *http.Request) string {
	forwarded := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		candidate := strings.TrimSpace(parts[0])
		if candidate != "" {
			return candidate
		}
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && strings.TrimSpace(host) != "" {
		return host
	}
	if trimmed := strings.TrimSpace(r.RemoteAddr); trimmed != "" {
		return trimmed
	}
	return ""
}

// debugEnabledForRequest checks whether the caller IP is in the allowlist.
// Keeping the lookup in one spot makes it simple to extend later with CIDR
// matching or runtime toggles without touching handlers.
func debugEnabledForRequest(r *http.Request) bool {
	if len(debugIPAllowlist) == 0 || r == nil {
		return false
	}
	ip := requestClientIP(r)
	if ip == "" {
		return false
	}
	_, ok := debugIPAllowlist[ip]
	return ok
}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent CloudFront from caching pages that show different content based on login status
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	lang := getPreferredLanguage(r)

	// Готовим шаблон
	tmpl := template.Must(template.New("map.html").Funcs(template.FuncMap{
		"translate": func(key string) string {
			if val, ok := translations[lang][key]; ok {
				return val
			}
			return translations["en"][key]
		},
	}).ParseFS(content, "public_html/map.html"))

	if CompileVersion == "dev" {
		CompileVersion = "latest"
	}

	translationsJSON, err := marshalTemplateJS(translationsForLang(translations, lang))
	if err != nil {
		log.Printf("map handler: marshal translations failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	markers := doseData.Markers
	if markers == nil {
		markers = []database.Marker{}
	}
	markersJSON, err := marshalTemplateJS(markers)
	if err != nil {
		log.Printf("map handler: marshal markers failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Данные для шаблона
	data := struct {
		Version           string
		Translations      map[string]map[string]string
		Lang              string
		DefaultLat        float64
		DefaultLon        float64
		DefaultZoom       int
		DefaultLayer      string
		AutoLocateDefault bool
		RealtimeAvailable bool
		SupportEmail      string
		TranslationsJSON  template.JS
		MarkersJSON       template.JS
		DebugEnabled      bool
	}{
		Version:           CompileVersion,
		Translations:      translations,
		Lang:              lang,
		DefaultLat:        *defaultLat,
		DefaultLon:        *defaultLon,
		DefaultZoom:       *defaultZoom,
		DefaultLayer:      *defaultLayer,
		AutoLocateDefault: *autoLocateDefault,
		RealtimeAvailable: *safecastRealtimeEnabled,
		SupportEmail:      strings.TrimSpace(*supportEmail),
		TranslationsJSON:  translationsJSON,
		MarkersJSON:       markersJSON,
		DebugEnabled:      debugEnabledForRequest(r),
	}

	// Рендерим в буфер, чтобы не дублировать WriteHeader
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := buf.WriteTo(w); err != nil {
		if isClientDisconnect(err) {
			log.Printf("client disconnected while writing response")
		} else {
			log.Printf("Error writing response: %v", err)
		}
	}
}

// homeHandler serves the home page with location search interface.
// The home page shows an empty map with a centered modal prompting for location entry.
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Prevent CloudFront from caching pages that show different content based on login status
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate, private")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	lang := getPreferredLanguage(r)

	// Prepare template
	tmpl := template.Must(template.New("home.html").Funcs(template.FuncMap{
		"translate": func(key string) string {
			if val, ok := translations[lang][key]; ok {
				return val
			}
			return translations["en"][key]
		},
	}).ParseFS(content, "public_html/home.html"))

	// Template data
	data := struct {
		Version      string
		Translations map[string]map[string]string
		Lang         string
	}{
		Version:      CompileVersion,
		Translations: translations,
		Lang:         lang,
	}

	// Render to buffer
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("Error executing home template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := buf.WriteTo(w); err != nil {
		if isClientDisconnect(err) {
			log.Printf("client disconnected while writing home response")
		} else {
			log.Printf("Error writing home response: %v", err)
		}
	}
}

