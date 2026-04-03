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
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "safecast-new-map/cmd/unified-server/docs"
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

	// Favicon endpoints
	mux.HandleFunc("/mcp-api/favicon.ico", serveFavicon)
	mux.HandleFunc("/mcp-api/favicon-16x16.png", serveFavicon16)
	mux.HandleFunc("/mcp-api/favicon-32x32.png", serveFavicon32)

	// Swagger UI — themed to match simplemap admin pages
	mux.HandleFunc("/mcp-api/swagger-theme.css", serveSwaggerTheme)
	mux.Handle("/mcp-api/", httpSwagger.Handler(
		httpSwagger.URL("/mcp-api/doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"onComplete": `function() {
				// Change page title
				document.title = 'Safecast MCP Docs';

				// Remove Swagger logo completely from DOM
				const swaggerLogo = document.querySelector('.topbar-wrapper .link');
				if (swaggerLogo) {
					swaggerLogo.remove();
				}
				// Also remove any img tags in topbar as backup
				const logoImgs = document.querySelectorAll('.topbar-wrapper img');
				logoImgs.forEach(img => img.remove());

				// Inject Safecast favicon
				const link16 = document.createElement('link');
				link16.rel = 'icon';
				link16.type = 'image/png';
				link16.sizes = '16x16';
				link16.href = '/mcp-api/favicon-16x16.png';
				document.head.appendChild(link16);

				const link32 = document.createElement('link');
				link32.rel = 'icon';
				link32.type = 'image/png';
				link32.sizes = '32x32';
				link32.href = '/mcp-api/favicon-32x32.png';
				document.head.appendChild(link32);

				const linkICO = document.createElement('link');
				linkICO.rel = 'shortcut icon';
				linkICO.href = '/mcp-api/favicon.ico';
				document.head.appendChild(linkICO);

				// Inject custom CSS
				const style = document.createElement('link');
				style.rel = 'stylesheet';
				style.href = '/mcp-api/swagger-theme.css';
				document.head.appendChild(style);

				// Create dark mode toggle button
				const btn = document.createElement('button');
				btn.id = 'dark-mode-toggle';
				btn.textContent = '🌙 Dark Mode';

				// Check localStorage for saved preference
				const isDark = localStorage.getItem('darkMode') === 'true';
				if (isDark) {
					document.body.classList.add('dark-mode');
					btn.textContent = '☀️ Light Mode';
				}

				btn.onclick = function() {
					document.body.classList.toggle('dark-mode');
					const nowDark = document.body.classList.contains('dark-mode');
					btn.textContent = nowDark ? '☀️ Light Mode' : '🌙 Dark Mode';
					localStorage.setItem('darkMode', nowDark);
				};

				document.body.appendChild(btn);

				// Add small route switcher so users can jump between the two Swagger UIs.
				const existing = document.getElementById('safecast-doc-nav');
				if (existing) existing.remove();
				const nav = document.createElement('div');
				nav.id = 'safecast-doc-nav';
				nav.style.cssText = 'position:fixed;top:8px;right:12px;z-index:9999;background:#111827;color:#fff;padding:8px 12px;border-radius:8px;font:12px/1.3 -apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif;box-shadow:0 2px 8px rgba(0,0,0,.25)';
				nav.innerHTML = 'This is the API documentation for the <strong>MCP API</strong>. &nbsp;|&nbsp; <a href="/map-api/" style="color:#93c5fd;text-decoration:none">Open Map API docs</a>';
				document.body.appendChild(nav);
			}`,
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

// serveSwaggerTheme serves the MCP API Swagger theme CSS (teal accent).
func serveSwaggerTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(mcpSwaggerThemeCSS))
}

// serveAPIDocsPage serves the combined Map API + MCP API documentation page with tabs.
func serveAPIDocsPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write([]byte(apiDocsPageHTML))
}

const apiDocsPageHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Safecast API Documentation</title>
<link rel="stylesheet" href="/map-api/swagger-ui.css">
<style>
/* ── Admin-style CSS variables (matches admin pages exactly) ── */
:root {
  --bg-primary:   #f5f5f5;
  --bg-card:      #fff;
  --text-primary: #333;
  --text-secondary:#666;
  --text-muted:   #999;
  --border-color: #ddd;
  --link-color:   #0066cc;
  --shadow:       0 1px 3px rgba(0,0,0,0.1);
  --hover-bg:     #f9f9f9;
  --th-bg:        #424242;
}
@media (prefers-color-scheme: dark) {
  :root {
    --bg-primary:   #1a1a1a;
    --bg-card:      #2b2b2b;
    --text-primary: #eee;
    --text-secondary:#aaa;
    --text-muted:   #777;
    --border-color: #444;
    --link-color:   #90caf9;
    --shadow:       0 1px 3px rgba(255,255,255,0.07);
    --hover-bg:     #333;
    --th-bg:        #616161;
    color-scheme: dark;
  }
}
:root[data-theme='light'] {
  --bg-primary:   #f5f5f5; --bg-card: #fff; --text-primary: #333;
  --text-secondary:#666; --text-muted: #999; --border-color: #ddd;
  --link-color:   #0066cc; --shadow: 0 1px 3px rgba(0,0,0,0.1);
  --hover-bg:     #f9f9f9; --th-bg: #424242; color-scheme: light;
}
:root[data-theme='dark'] {
  --bg-primary:   #1a1a1a; --bg-card: #2b2b2b; --text-primary: #eee;
  --text-secondary:#aaa; --text-muted: #777; --border-color: #444;
  --link-color:   #90caf9; --shadow: 0 1px 3px rgba(255,255,255,0.07);
  --hover-bg:     #333; --th-bg: #616161; color-scheme: dark;
}

/* ── Layout ── */
*, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans", Arial, sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  min-height: 100vh;
  transition: background 0.2s, color 0.2s;
}

/* ── Top nav bar (matches admin pages) ── */
.top-nav {
  background: #1a3a5c;
  color: #fff;
  padding: 0 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  height: 52px;
  box-shadow: 0 2px 6px rgba(0,0,0,0.3);
  position: sticky;
  top: 0;
  z-index: 100;
}
.top-nav .logo {
  font-weight: 700;
  font-size: 17px;
  letter-spacing: 0.02em;
  color: #fff;
  text-decoration: none;
  white-space: nowrap;
}
.top-nav .logo span { color: #7ecbff; }
.top-nav .back-link {
  color: #afd4f5;
  text-decoration: none;
  font-size: 13px;
  white-space: nowrap;
}
.top-nav .back-link:hover { color: #fff; }
.top-nav .nav-title {
  flex: 1;
  font-size: 15px;
  font-weight: 600;
  color: #d0e8ff;
  text-align: center;
}
#theme-toggle {
  background: rgba(255,255,255,0.12);
  color: #fff;
  border: 1px solid rgba(255,255,255,0.25);
  border-radius: 6px;
  padding: 6px 14px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: background 0.2s;
}
#theme-toggle:hover { background: rgba(255,255,255,0.22); }

/* ── Page content ── */
.page-content {
  max-width: 1400px;
  margin: 0 auto;
  padding: 24px 20px 48px;
}

/* ── Description card ── */
.desc-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 20px 24px;
  margin-bottom: 20px;
  box-shadow: var(--shadow);
}
.desc-card h1 {
  font-size: 20px;
  font-weight: 700;
  color: var(--text-primary);
  margin-bottom: 6px;
}
.desc-card p {
  font-size: 14px;
  color: var(--text-secondary);
  line-height: 1.6;
}
.desc-card a { color: var(--link-color); }

/* ── API tabs (matches admin-tabs exactly) ── */
.api-tabs {
  display: flex;
  gap: 2px;
  margin-bottom: 16px;
  background: var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  width: fit-content;
}
.api-tab {
  padding: 10px 28px;
  border: none;
  background: var(--bg-card);
  color: var(--text-secondary);
  font-family: inherit;
  font-size: 0.95em;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}
.api-tab:hover { background: var(--hover-bg); color: var(--text-primary); }
.api-tab.active { background: #2196F3; color: #fff; }

/* ── Swagger containers ── */
.swagger-panel { display: none; }
.swagger-panel.active { display: block; }
.swagger-wrap {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: var(--shadow);
}

/* ── Swagger UI overrides — shared for both panels ── */
.swagger-wrap .swagger-ui { background: var(--bg-card) !important; }
.swagger-wrap .swagger-ui .info .title,
.swagger-wrap .swagger-ui .info h1,
.swagger-wrap .swagger-ui .info h2 { color: var(--text-primary) !important; }
.swagger-wrap .swagger-ui .info .base-url,
.swagger-wrap .swagger-ui .info p { color: var(--text-secondary) !important; }
.swagger-wrap .swagger-ui .opblock-tag { color: var(--text-primary) !important; border-bottom: 1px solid var(--border-color) !important; }
.swagger-wrap .swagger-ui .opblock { border-radius: 6px !important; margin-bottom: 6px !important; border: 1px solid var(--border-color) !important; background: var(--bg-card) !important; box-shadow: var(--shadow) !important; }
.swagger-wrap .swagger-ui .opblock.opblock-get { border-color: #0066cc !important; background: #f0f6ff !important; }
.swagger-wrap .swagger-ui .opblock.opblock-get .opblock-summary-method { background: #0066cc !important; border-radius: 4px !important; }
.swagger-wrap .swagger-ui .opblock.opblock-post { border-color: #4caf50 !important; background: #f0fff4 !important; }
.swagger-wrap .swagger-ui .opblock.opblock-post .opblock-summary-method { background: #4caf50 !important; border-radius: 4px !important; }
.swagger-wrap .swagger-ui table thead tr th { background: var(--th-bg) !important; color: #fff !important; font-weight: 600 !important; }
.swagger-wrap .swagger-ui table tbody tr:hover { background: var(--hover-bg) !important; }
.swagger-wrap .swagger-ui .responses-inner { background: var(--bg-card) !important; border-radius: 6px !important; }
.swagger-wrap .swagger-ui input[type=text],
.swagger-wrap .swagger-ui textarea,
.swagger-wrap .swagger-ui select { border-radius: 6px !important; border: 1px solid var(--border-color) !important; background: var(--bg-card) !important; color: var(--text-primary) !important; }
.swagger-wrap .swagger-ui .btn.execute { background: #0066cc !important; border-color: #0066cc !important; border-radius: 6px !important; color: #fff !important; }
.swagger-wrap .swagger-ui .btn.authorize { border-radius: 6px !important; color: #0066cc !important; border-color: #0066cc !important; }
/* Code blocks — ensure readable text in light mode */
.swagger-wrap .swagger-ui .microlight,
.swagger-wrap .swagger-ui pre.microlight { background: #f6f8fa !important; color: #24292e !important; border-radius: 6px !important; border: 1px solid var(--border-color) !important; font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace !important; font-size: 0.88rem !important; }
.swagger-wrap .swagger-ui .microlight span { color: inherit !important; }
/* Hide swagger topbar (we have our own nav) */
.swagger-wrap .swagger-ui .topbar { display: none !important; }
/* Hide doc.json info link */
.swagger-wrap .swagger-ui .info .link,
.swagger-wrap .swagger-ui .info a[href*="doc.json"] { display: none !important; }
/* Scheme selector */
.swagger-wrap .swagger-ui .scheme-container { background: var(--bg-card) !important; box-shadow: none !important; border-bottom: 1px solid var(--border-color) !important; padding: 10px 20px !important; }
.swagger-wrap .swagger-ui .wrapper { padding: 16px 20px !important; }

/* ── Dark mode overrides for swagger ── */
[data-theme='dark'] .swagger-wrap .swagger-ui,
[data-theme='dark'] .swagger-wrap .swagger-ui .wrapper { background: #1e1e1e !important; color: #eee !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .opblock { background: #2b2b2b !important; border-color: #444 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .opblock.opblock-get { background: #1a2940 !important; border-color: #90caf9 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .opblock.opblock-get .opblock-summary-method { background: #1565c0 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .opblock.opblock-post { background: #1a2e1a !important; border-color: #81c784 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .opblock.opblock-post .opblock-summary-method { background: #388e3c !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .opblock-summary-description,
[data-theme='dark'] .swagger-wrap .swagger-ui .opblock-summary-path,
[data-theme='dark'] .swagger-wrap .swagger-ui .opblock-tag,
[data-theme='dark'] .swagger-wrap .swagger-ui p,
[data-theme='dark'] .swagger-wrap .swagger-ui td,
[data-theme='dark'] .swagger-wrap .swagger-ui label { color: #ddd !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui a { color: #90caf9 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .scheme-container { background: #1e1e1e !important; border-bottom-color: #444 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui select { background: #2b2b2b !important; color: #eee !important; border-color: #555 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui input[type=text],
[data-theme='dark'] .swagger-wrap .swagger-ui textarea { background: #2b2b2b !important; color: #eee !important; border-color: #555 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .microlight,
[data-theme='dark'] .swagger-wrap .swagger-ui pre.microlight { background: #0d1117 !important; color: #c9d1d9 !important; border-color: #444 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .responses-inner { background: #1e1e1e !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .model-box,
[data-theme='dark'] .swagger-wrap .swagger-ui section.models { background: #2b2b2b !important; border-color: #444 !important; }
[data-theme='dark'] .swagger-wrap { border-color: #444 !important; }
[data-theme='dark'] .desc-card { border-color: #444 !important; }
[data-theme='dark'] .swagger-wrap .swagger-ui .opblock-body { background: #222 !important; }
@media (prefers-color-scheme: dark) {
  :root:not([data-theme='light']) .swagger-wrap .swagger-ui { background: #1e1e1e !important; color: #eee !important; }
  :root:not([data-theme='light']) .swagger-wrap .swagger-ui .opblock { background: #2b2b2b !important; border-color: #444 !important; }
  :root:not([data-theme='light']) .swagger-wrap .swagger-ui .opblock.opblock-get { background: #1a2940 !important; border-color: #90caf9 !important; }
  :root:not([data-theme='light']) .swagger-wrap .swagger-ui .microlight,
  :root:not([data-theme='light']) .swagger-wrap .swagger-ui pre.microlight { background: #0d1117 !important; color: #c9d1d9 !important; border-color: #444 !important; }
}
</style>
</head>
<body>

<nav class="top-nav">
  <a href="/" class="logo">Safe<span>cast</span></a>
  <span style="color:rgba(255,255,255,0.3);font-size:18px;">|</span>
  <a href="/" class="back-link">← Back to Map</a>
  <span class="nav-title">API Documentation</span>
  <button id="theme-toggle" onclick="toggleTheme()">🌙 Dark Mode</button>
</nav>

<div class="page-content">
  <div class="desc-card">
    <h1>Safecast API Documentation</h1>
    <p>Safecast has collected over 200 million radiation measurements from citizen scientists worldwide.
       All data is <a href="https://creativecommons.org/publicdomain/zero/1.0/" target="_blank">CC0 licensed</a>
       and freely accessible — no account or API key required.</p>
  </div>

  <div class="api-tabs">
    <button class="api-tab active" id="tab-btn-map" onclick="switchTab('map')">Map API</button>
    <button class="api-tab" id="tab-btn-mcp" onclick="switchTab('mcp')">MCP API</button>
  </div>

  <div id="panel-map" class="swagger-panel active">
    <div class="swagger-wrap"><div id="swagger-map"></div></div>
  </div>
  <div id="panel-mcp" class="swagger-panel">
    <div class="swagger-wrap"><div id="swagger-mcp"></div></div>
  </div>
</div>

<script src="/map-api/swagger-ui-bundle.js"></script>
<script src="/map-api/swagger-ui-standalone-preset.js"></script>
<script>
(function() {
  // ── Theme ──
  var saved = localStorage.getItem('safecastDocTheme');
  var prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  var theme = saved || (prefersDark ? 'dark' : 'light');
  document.documentElement.setAttribute('data-theme', theme);

  function applyThemeLabel() {
    var btn = document.getElementById('theme-toggle');
    var isDark = document.documentElement.getAttribute('data-theme') === 'dark';
    btn.textContent = isDark ? '\u2600\uFE0F Light Mode' : '\uD83C\uDF19 Dark Mode';
  }
  applyThemeLabel();

  window.toggleTheme = function() {
    var isDark = document.documentElement.getAttribute('data-theme') === 'dark';
    var next = isDark ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', next);
    localStorage.setItem('safecastDocTheme', next);
    applyThemeLabel();
  };

  // ── Tab switching ──
  var mcpInitialized = false;
  window.switchTab = function(tab) {
    document.getElementById('panel-map').classList.toggle('active', tab === 'map');
    document.getElementById('panel-mcp').classList.toggle('active', tab === 'mcp');
    document.getElementById('tab-btn-map').classList.toggle('active', tab === 'map');
    document.getElementById('tab-btn-mcp').classList.toggle('active', tab === 'mcp');
    localStorage.setItem('safecastDocTab', tab);
    if (tab === 'mcp' && !mcpInitialized) {
      mcpInitialized = true;
      initMCP();
    }
  };

  // ── Initialize Map API swagger ──
  SwaggerUIBundle({
    url: '/map-api/doc.json',
    domNode: document.getElementById('swagger-map'),
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    layout: 'BaseLayout',
    deepLinking: false,
    displayRequestDuration: true,
    defaultModelsExpandDepth: -1,
  });

  // ── Initialize MCP API swagger (lazy — on first tab click) ──
  function initMCP() {
    SwaggerUIBundle({
      url: '/mcp-api/doc.json',
      domNode: document.getElementById('swagger-mcp'),
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
      layout: 'BaseLayout',
      deepLinking: false,
      displayRequestDuration: true,
      defaultModelsExpandDepth: -1,
    });
  }

  // ── Restore last active tab ──
  var lastTab = localStorage.getItem('safecastDocTab');
  if (lastTab === 'mcp') switchTab('mcp');
})();
</script>
</body>
</html>`

// serveMapSwaggerTheme serves the Map API Swagger theme CSS (blue/navy accent).
func serveMapSwaggerTheme(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(mapSwaggerThemeCSS))
}

// mapSwaggerThemeCSS is the Swagger UI theme for the Map REST API (blue/navy palette).
const mapSwaggerThemeCSS = `
/* ── Safecast Map API Swagger theme — blue/navy palette ── */

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
  background: #0066cc;
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
  background: #0055aa;
}
body.dark-mode #dark-mode-toggle {
  background: #1565c0;
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

/* Top bar — deep navy to signal Map API */
.swagger-ui .topbar {
  background: #1a3a5c !important;
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

/* GET method colour — matches --link-color (#0066cc) */
.swagger-ui .opblock.opblock-get {
  border-color: #0066cc !important;
  background: #f0f6ff !important;
}
.swagger-ui .opblock.opblock-get .opblock-summary-method {
  background: #0066cc !important;
  border-radius: 4px !important;
}

/* Accent links */
.swagger-ui a,
.swagger-ui .opblock-summary-path,
.swagger-ui .info a {
  color: #0066cc !important;
}
.swagger-ui a:hover {
  text-decoration: underline !important;
}

/* Buttons */
.swagger-ui .btn.execute {
  background: #0066cc !important;
  border-color: #0066cc !important;
  border-radius: 8px !important;
  color: #fff !important;
}
.swagger-ui .btn.execute:hover {
  background: #0055aa !important;
}
.swagger-ui .btn.cancel {
  border-radius: 8px !important;
}
.swagger-ui .btn.authorize {
  border-radius: 8px !important;
  color: #0066cc !important;
  border-color: #0066cc !important;
}

/* Section headers */
.swagger-ui .opblock-tag {
  border-bottom: 1px solid #ddd !important;
  color: #333 !important;
  font-weight: 600 !important;
}

/* Parameter tables */
.swagger-ui table thead tr th {
  background: #1a3a5c !important;
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
  color: #24292e !important;
  border-radius: 8px !important;
  border: 1px solid #ddd !important;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace !important;
  font-size: 0.9rem !important;
}
.swagger-ui .microlight span {
  color: inherit !important;
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
  background: #2a4a6c !important;
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
  background: #1a2940 !important;
  border-color: #90caf9 !important;
}
body.dark-mode .swagger-ui .opblock.opblock-get .opblock-summary-method {
  background: #1565c0 !important;
}

body.dark-mode .swagger-ui a,
body.dark-mode .swagger-ui .opblock-summary-path,
body.dark-mode .swagger-ui .info a {
  color: #90caf9 !important;
}

body.dark-mode .swagger-ui .opblock-tag {
  border-bottom-color: #444 !important;
  color: #eee !important;
}

body.dark-mode .swagger-ui table thead tr th {
  background: #2a4a6c !important;
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
  background: #1565c0 !important;
  border-color: #1565c0 !important;
}
body.dark-mode .swagger-ui .btn.authorize {
  color: #90caf9 !important;
  border-color: #90caf9 !important;
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

// mcpSwaggerThemeCSS is the Swagger UI theme for the MCP API (teal palette).
const mcpSwaggerThemeCSS = `
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
/* Hide the info link that shows doc.json URL */
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

/* Operation blocks */
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
  color: #24292e !important;
  border-radius: 8px !important;
  border: 1px solid #ddd !important;
  font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace !important;
  font-size: 0.9rem !important;
}
.swagger-ui .microlight span {
  color: inherit !important;
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

/* ── Dark mode ── */
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
