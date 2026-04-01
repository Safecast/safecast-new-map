// Package main provides the Safecast MCP server with an optional REST API layer.
//
// @title           Safecast Map API
// @version         1.0
// @description     REST access to the Safecast radiation monitoring dataset via simplemap.safecast.org — 200M+ measurements from citizen scientists worldwide. All data is CC0-licensed and read-only. Powered by PostgreSQL+PostGIS.
// @contact.name    Safecast
// @contact.url     https://safecast.org
// @license.name    CC0 1.0 Universal
// @license.url     https://creativecommons.org/publicdomain/zero/1.0/
// @host            simplemap.safecast.org
// @BasePath        /api
// @schemes         https http
//
// @tag.name        historical
// @tag.description Historical bGeigie mobile radiation measurements
// @tag.name        realtime
// @tag.description Real-time fixed sensor readings (Pointcast, Solarcast, bGeigieZen)
// @tag.name        spectroscopy
// @tag.description Gamma spectroscopy records
// @tag.name        reference
// @tag.description Aggregate statistics and reference information
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "safecast-new-map/cmd/mcp-server/docs"
)

//go:embed static/favicon.ico
var faviconICO []byte

//go:embed static/favicon-16x16.png
var favicon16 []byte

//go:embed static/favicon-32x32.png
var favicon32 []byte

// RESTHandler wires all REST API routes onto a mux.
type RESTHandler struct{}

// Register attaches all /api/* routes and the /mcp-api/ Swagger UI to mux.
func (h *RESTHandler) Register(mux *http.ServeMux) {
	// Historical data
	mux.HandleFunc("/api/radiation", h.handleRadiation)
	mux.HandleFunc("/api/area", h.handleArea)
	mux.HandleFunc("/api/tracks", h.handleTracks)
	mux.HandleFunc("/api/track/", h.handleTrack)   // /api/track/{id}
	mux.HandleFunc("/api/device/", h.handleDevice) // /api/device/{id}/history

	// Real-time sensors
	mux.HandleFunc("/api/sensors", h.handleSensors)
	mux.HandleFunc("/api/sensor/", h.handleSensor) // /api/sensor/{id}/current or /history

	// Spectroscopy
	mux.HandleFunc("/api/spectra", h.handleSpectra)
	mux.HandleFunc("/api/spectrum/", h.handleSpectrum) // /api/spectrum/{marker_id}

	// Reference / stats
	mux.HandleFunc("/api/stats", h.handleStats)
	mux.HandleFunc("/api/extreme", handleRESTExtremeReadings)
	mux.HandleFunc("/api/info/", h.handleInfo) // /api/info/{topic}

	// GPT-optimised compact endpoints (for Custom GPT Actions)
	h.RegisterGPT(mux)

	mapBaseForDocs := strings.TrimSpace(os.Getenv("MAP_BASE_URL"))
	if mapBaseForDocs == "" {
		mapBaseForDocs = "http://localhost:8765"
	}
	mapDocsURL := strings.TrimRight(mapBaseForDocs, "/") + "/map-api/"
	mcpAPINavScript := fmt.Sprintf(`function() {
				document.title = 'Safecast MCP API Docs';
				const mapDocsURL = %q;

				// ── Favicons ──
				const link16 = document.createElement('link');
				link16.rel = 'icon'; link16.type = 'image/png'; link16.sizes = '16x16';
				link16.href = '/mcp-api/favicon-16x16.png';
				document.head.appendChild(link16);
				const link32 = document.createElement('link');
				link32.rel = 'icon'; link32.type = 'image/png'; link32.sizes = '32x32';
				link32.href = '/mcp-api/favicon-32x32.png';
				document.head.appendChild(link32);
				const linkICO = document.createElement('link');
				linkICO.rel = 'shortcut icon'; linkICO.href = '/mcp-api/favicon.ico';
				document.head.appendChild(linkICO);

				// ── Theme CSS ──
				const style = document.createElement('link');
				style.rel = 'stylesheet'; style.href = '/mcp-api/swagger-theme.css';
				document.head.appendChild(style);

				// ── Remove Swagger logo ──
				const swaggerLogo = document.querySelector('.topbar-wrapper .link');
				if (swaggerLogo) swaggerLogo.remove();
				document.querySelectorAll('.topbar-wrapper img').forEach(img => img.remove());

				// ── Inject "Switch to Map API" button into topbar ──
				const topbar = document.querySelector('.swagger-ui .topbar');
				if (topbar) {
					const existingBtn = document.getElementById('safecast-switch-btn');
					if (existingBtn) existingBtn.remove();
					const switchBtn = document.createElement('a');
					switchBtn.id = 'safecast-switch-btn';
					switchBtn.href = mapDocsURL;
					switchBtn.textContent = '\u2190 Switch to Map API';
					switchBtn.style.cssText = 'display:inline-block;margin-left:auto;margin-right:16px;padding:6px 14px;background:#085e58;color:#fff;border-radius:6px;font:600 13px/1.4 -apple-system,BlinkMacSystemFont,sans-serif;text-decoration:none;border:1px solid rgba(255,255,255,0.25);white-space:nowrap;';
					switchBtn.onmouseover = function() { this.style.background = '#064a45'; };
					switchBtn.onmouseout  = function() { this.style.background = '#085e58'; };
					const wrapper = topbar.querySelector('.topbar-wrapper');
					if (wrapper) {
						wrapper.style.display = 'flex';
						wrapper.style.alignItems = 'center';
						wrapper.style.width = '100%%';
						wrapper.appendChild(switchBtn);
					} else {
						topbar.appendChild(switchBtn);
					}
				}

				// ── Dark mode toggle ──
				const btn = document.createElement('button');
				btn.id = 'dark-mode-toggle';
				btn.textContent = '\u{1F319} Dark Mode';
				const isDark = localStorage.getItem('safecastMCPDarkMode') === 'true';
				if (isDark) { document.body.classList.add('dark-mode'); btn.textContent = '\u2600\uFE0F Light Mode'; }
				btn.onclick = function() {
					document.body.classList.toggle('dark-mode');
					const nowDark = document.body.classList.contains('dark-mode');
					btn.textContent = nowDark ? '\u2600\uFE0F Light Mode' : '\u{1F319} Dark Mode';
					localStorage.setItem('safecastMCPDarkMode', nowDark);
				};
				document.body.appendChild(btn);

				// ── Preamble section (inserted before #swagger-ui) ──
				const existing = document.getElementById('safecast-preamble');
				if (existing) existing.remove();
				const preamble = document.createElement('div');
				preamble.id = 'safecast-preamble';
				preamble.style.cssText = 'font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;background:#fff;border-bottom:3px solid #0d9488;padding:28px 32px 24px;line-height:1.6;';
				preamble.innerHTML = '<div style="max-width:900px;margin:0 auto;">' +
					'<h2 style="margin:0 0 6px;font-size:22px;color:#0f3d38;">Safecast MCP API</h2>' +
					'<p style="margin:0 0 16px;font-size:15px;color:#555;">' +
						'This API is designed for AI assistants and automated tools. It exposes the Safecast radiation dataset ' +
						'as a set of callable functions (MCP tools) and standard REST endpoints, so that language models and ' +
						'applications can query radiation data programmatically. If you want to ask questions about the data ' +
						'in plain English instead, try the <a href="/assistant/" style="color:#0d9488;">AI web chat</a>.' +
					'</p>' +
					'<details style="margin-bottom:16px;">' +
						'<summary style="cursor:pointer;font-weight:600;font-size:14px;color:#0f3d38;user-select:none;">For developers &mdash; transports &amp; tools overview</summary>' +
						'<div style="margin-top:10px;display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:10px;">' +
							'<div style="background:#f0fdfa;border:1px solid #99f6e4;border-radius:8px;padding:12px;">' +
								'<strong style="color:#0f3d38;">MCP Streamable HTTP</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Primary MCP transport at <code>/mcp-http</code>. Compatible with Claude, GPT, and most MCP clients.</p>' +
							'</div>' +
							'<div style="background:#f0fdfa;border:1px solid #99f6e4;border-radius:8px;padding:12px;">' +
								'<strong style="color:#0f3d38;">MCP SSE</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Server-Sent Events transport at <code>/mcp/sse</code> for streaming-capable clients.</p>' +
							'</div>' +
							'<div style="background:#f0fdfa;border:1px solid #99f6e4;border-radius:8px;padding:12px;">' +
								'<strong style="color:#0f3d38;">REST Endpoints</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">All MCP tools are also accessible as plain HTTP GET/POST calls under <code>/api/</code>.</p>' +
							'</div>' +
							'<div style="background:#f0fdfa;border:1px solid #99f6e4;border-radius:8px;padding:12px;">' +
								'<strong style="color:#0f3d38;">AI Web Chat</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Human-friendly interface at <code>/assistant/</code> for querying data with natural language.</p>' +
							'</div>' +
						'</div>' +
						'<p style="margin:12px 0 0;font-size:13px;color:#777;">No authentication required. All data is CC0 licensed.</p>' +
					'</details>' +
					'<div style="display:flex;gap:10px;flex-wrap:wrap;align-items:center;">' +
						'<a href="/assistant/" style="display:inline-block;padding:7px 16px;background:#0f3d38;color:#fff;border-radius:6px;font:600 13px/1.4 sans-serif;text-decoration:none;">Open AI Chat</a>' +
						'<a href="' + mapDocsURL + '" style="display:inline-block;padding:7px 16px;background:#1a3a5c;color:#fff;border-radius:6px;font:600 13px/1.4 sans-serif;text-decoration:none;">\u2190 Switch to Map API</a>' +
					'</div>' +
				'</div>';
				const swaggerUIEl = document.getElementById('swagger-ui');
				if (swaggerUIEl) {
					document.body.insertBefore(preamble, swaggerUIEl);
				} else {
					document.body.prepend(preamble);
				}
			}`, mapDocsURL)

	// Favicon endpoints
	mux.HandleFunc("/mcp-api/favicon.ico", serveFavicon)
	mux.HandleFunc("/mcp-api/favicon-16x16.png", serveFavicon16)
	mux.HandleFunc("/mcp-api/favicon-32x32.png", serveFavicon32)

	// Swagger UI — teal theme for MCP API
	mux.HandleFunc("/mcp-api/swagger-theme.css", serveSwaggerTheme)
	mux.Handle("/mcp-api/", httpSwagger.Handler(
		httpSwagger.URL("/mcp-api/doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"onComplete": mcpAPINavScript,
		}),
	))
}

// writeJSON writes v as a JSON response with the given HTTP status code.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(status)
	_ = jsonEncode(w, v)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// jsonEncode writes v as JSON to w.
func jsonEncode(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

// serveMCPResult pipes an MCP tool result directly to an HTTP response.
// The tool functions already produce indented JSON, so we write the text content
// straight through. Tool errors become HTTP 400 responses.
func serveMCPResult(w http.ResponseWriter, result *mcp.CallToolResult, err error) {
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if result == nil || len(result.Content) == 0 {
		writeError(w, http.StatusInternalServerError, "empty result")
		return
	}
	// Extract the text payload. Content is an interface; use AsTextContent to unwrap it.
	text := ""
	for _, c := range result.Content {
		if tc, ok := mcp.AsTextContent(c); ok && tc.Text != "" {
			text = tc.Text
			break
		}
	}
	if result.IsError {
		writeError(w, http.StatusBadRequest, text)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, text)
}

// serveFavicon serves the Safecast favicon.ico
func serveFavicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(faviconICO)
}

// serveFavicon16 serves the Safecast 16x16 favicon
func serveFavicon16(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(favicon16)
}

// serveFavicon32 serves the Safecast 32x32 favicon
func serveFavicon32(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(favicon32)
}

// serveSwaggerTheme serves the MCP API Swagger theme CSS (teal palette).
func serveSwaggerTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerThemeCSS))
}

// swaggerThemeCSS is the Swagger UI theme for the MCP API (teal palette).
const swaggerThemeCSS = `
/* ── Safecast MCP API Swagger theme — teal palette ── */

/* Hide Swagger logo and collapse the space */
.swagger-ui .topbar-wrapper {
  padding-left: 20px !important;
}
.swagger-ui .topbar-wrapper img,
.swagger-ui .topbar-wrapper a,
.swagger-ui .topbar-wrapper .link,
.swagger-ui .topbar-wrapper svg,
.swagger-ui .topbar-wrapper .svg-assets,
.swagger-ui svg-assets,
.svg-assets {
  display: none !important;
  visibility: hidden !important;
  opacity: 0 !important;
  width: 0 !important;
  height: 0 !important;
  position: absolute !important;
  left: -9999px !important;
  max-height: 0 !important;
}
/* Hide the info link that shows /mcp-api/doc.json */
.swagger-ui .info .link,
.swagger-ui .info a[href*="doc.json"] {
  display: none !important;
  visibility: hidden !important;
}

/* Dark mode toggle button */
#dark-mode-toggle {
  position: fixed;
  top: 12px;
  right: 20px;
  z-index: 10001;
  padding: 8px 16px;
  background: #0d9488;
  color: white;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans", Arial, sans-serif;
  font-size: 14px;
  font-weight: 600;
  box-shadow: 0 4px 8px rgba(0,0,0,0.3);
  transition: background 0.2s;
}
/* Hide any buttons that might show behind dark mode toggle */
.swagger-ui .topbar .download-url-button {
  display: none !important;
}
#dark-mode-toggle:hover {
  background: #0b7a70;
}
body.dark-mode #dark-mode-toggle {
  background: #0a6b62;
}

/* Font stack and base background (light mode) */
body,
.swagger-ui,
.swagger-ui .wrapper {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans", Arial, sans-serif !important;
  background: #f5f5f5 !important;
  color: #333 !important;
  transition: background 0.3s, color 0.3s;
}

/* Top bar — deep teal-dark to signal MCP API */
.swagger-ui .topbar {
  background: #0f3d38 !important;
  padding: 8px 0 !important;
}
.swagger-ui .topbar .download-url-wrapper input[type=text] {
  border-radius: 8px !important;
}
.swagger-ui .topbar-wrapper a span {
  color: #fff !important;
  font-weight: 600 !important;
}

/* Info block */
.swagger-ui .info .title,
.swagger-ui .info h1,
.swagger-ui .info h2,
.swagger-ui .info h3 {
  color: #333 !important;
}

/* Operation blocks — card style matching --shadow and --btn-border-radius */
.swagger-ui .opblock {
  border-radius: 8px !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1) !important;
  margin-bottom: 8px !important;
  border: 1px solid #ddd !important;
  background: #fff !important;
}
.swagger-ui .opblock .opblock-summary {
  border-radius: 8px !important;
}
.swagger-ui .opblock.is-open {
  border-radius: 8px !important;
}

/* GET method colour — teal */
.swagger-ui .opblock.opblock-get {
  border-color: #0d9488 !important;
  background: #f0fdfa !important;
}
.swagger-ui .opblock.opblock-get .opblock-summary-method {
  background: #0d9488 !important;
  border-radius: 4px !important;
}

/* Accent links */
.swagger-ui a,
.swagger-ui .opblock-summary-path,
.swagger-ui .info a {
  color: #0d9488 !important;
}
.swagger-ui a:hover {
  text-decoration: underline !important;
}

/* Buttons */
.swagger-ui .btn.execute {
  background: #0d9488 !important;
  border-color: #0d9488 !important;
  border-radius: 8px !important;
  color: #fff !important;
}
.swagger-ui .btn.execute:hover {
  background: #0b7a70 !important;
}
.swagger-ui .btn.cancel {
  border-radius: 8px !important;
}
.swagger-ui .btn.authorize {
  border-radius: 8px !important;
  color: #0d9488 !important;
  border-color: #0d9488 !important;
}

/* Section headers */
.swagger-ui .opblock-tag {
  border-bottom: 1px solid #ddd !important;
  color: #333 !important;
  font-weight: 600 !important;
}

/* Parameter tables */
.swagger-ui table thead tr th {
  background: #0f3d38 !important;
  color: #fff !important;
  font-weight: 600 !important;
}
.swagger-ui table tbody tr:hover {
  background: #f9f9f9 !important;
}

/* Response blocks */
.swagger-ui .responses-inner {
  background: #fff !important;
  border-radius: 8px !important;
}

/* Code blocks */
.swagger-ui .microlight,
.swagger-ui pre.microlight {
  background: #f6f8fa !important;
  border-radius: 8px !important;
  border: 1px solid #ddd !important;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace !important;
  font-size: 0.9rem !important;
}

/* Input fields */
.swagger-ui input[type=text],
.swagger-ui textarea,
.swagger-ui select {
  border-radius: 8px !important;
  border: 1px solid #ddd !important;
  background: #fff !important;
  color: #333 !important;
}

/* ── Dark mode — toggle via .dark-mode class ── */
body.dark-mode,
body.dark-mode .swagger-ui,
body.dark-mode .swagger-ui .wrapper {
  background: #1a1a1a !important;
  color: #eee !important;
}

body.dark-mode .swagger-ui .topbar {
  background: #1a5048 !important;
}

body.dark-mode .swagger-ui .info .title,
body.dark-mode .swagger-ui .info h1,
body.dark-mode .swagger-ui .info h2,
body.dark-mode .swagger-ui .info h3 {
  color: #eee !important;
}

body.dark-mode .swagger-ui .opblock {
  background: #2b2b2b !important;
  border-color: #444 !important;
  box-shadow: 0 1px 3px rgba(255, 255, 255, 0.07) !important;
}

body.dark-mode .swagger-ui .opblock.opblock-get {
  background: #0f2924 !important;
  border-color: #5eead4 !important;
}
body.dark-mode .swagger-ui .opblock.opblock-get .opblock-summary-method {
  background: #0a6b62 !important;
}

body.dark-mode .swagger-ui a,
body.dark-mode .swagger-ui .opblock-summary-path,
body.dark-mode .swagger-ui .info a {
  color: #5eead4 !important;
}

body.dark-mode .swagger-ui .opblock-tag {
  border-bottom-color: #444 !important;
  color: #eee !important;
}

body.dark-mode .swagger-ui table thead tr th {
  background: #1a5048 !important;
}
body.dark-mode .swagger-ui table tbody tr:hover {
  background: #333 !important;
}
body.dark-mode .swagger-ui table tbody tr td {
  color: #eee !important;
  border-bottom-color: #444 !important;
}

body.dark-mode .swagger-ui .responses-inner {
  background: #2b2b2b !important;
}

body.dark-mode .swagger-ui .microlight,
body.dark-mode .swagger-ui pre.microlight {
  background: #161b22 !important;
  border-color: #444 !important;
  color: #e6edf3 !important;
}

body.dark-mode .swagger-ui input[type=text],
body.dark-mode .swagger-ui textarea,
body.dark-mode .swagger-ui select {
  background: #2b2b2b !important;
  border-color: #444 !important;
  color: #eee !important;
}

body.dark-mode .swagger-ui .btn.execute {
  background: #0a6b62 !important;
  border-color: #0a6b62 !important;
}
body.dark-mode .swagger-ui .btn.authorize {
  color: #5eead4 !important;
  border-color: #5eead4 !important;
}

body.dark-mode .swagger-ui .opblock-summary-description,
body.dark-mode .swagger-ui .parameter__name,
body.dark-mode .swagger-ui .parameter__type,
body.dark-mode .swagger-ui label {
  color: #eee !important;
}

body.dark-mode .swagger-ui section.models {
  background: #2b2b2b !important;
  border-color: #444 !important;
}

body.dark-mode .swagger-ui .scheme-container {
  background: #2b2b2b !important;
  border-color: #444 !important;
  box-shadow: none !important;
}

body.dark-mode .swagger-ui .schemes > label {
  color: #eee !important;
}

body.dark-mode .swagger-ui select {
  background: #2b2b2b !important;
  color: #eee !important;
  border-color: #444 !important;
}
`
