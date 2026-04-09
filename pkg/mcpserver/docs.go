package mcpserver

import (
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

type DocsConfig struct {
	Mux              *http.ServeMux
	MapDocsURL       string
	IncludeAssistant bool
	FaviconICO       http.HandlerFunc
	Favicon16        http.HandlerFunc
	Favicon32        http.HandlerFunc
	ThemeCSS         http.HandlerFunc
}

// RegisterMCPDocs wires /mcp-api/* Swagger routes.
func RegisterMCPDocs(cfg DocsConfig) {
	if cfg.Mux == nil {
		return
	}
	if cfg.FaviconICO != nil {
		cfg.Mux.HandleFunc("/mcp-api/favicon.ico", cfg.FaviconICO)
	}
	if cfg.Favicon16 != nil {
		cfg.Mux.HandleFunc("/mcp-api/favicon-16x16.png", cfg.Favicon16)
	}
	if cfg.Favicon32 != nil {
		cfg.Mux.HandleFunc("/mcp-api/favicon-32x32.png", cfg.Favicon32)
	}
	if cfg.ThemeCSS != nil {
		cfg.Mux.HandleFunc("/mcp-api/swagger-theme.css", cfg.ThemeCSS)
	}
	cfg.Mux.Handle("/mcp-api/", httpSwagger.Handler(
		httpSwagger.URL("/mcp-api/doc.json"),
		httpSwagger.UIConfig(map[string]string{
			"onComplete": buildMCPDocsOnComplete(cfg.MapDocsURL, cfg.IncludeAssistant),
		}),
	))
}

func buildMCPDocsOnComplete(mapDocsURL string, includeAssistant bool) string {
	assistantBlock := ""
	assistantButton := ""
	if includeAssistant {
		assistantBlock = `' + '<div style="background:#f0fdfa;border:1px solid #99f6e4;border-radius:8px;padding:12px;">' +
								'<strong style="color:#0f3d38;">AI Web Chat</strong>' +
								'<p style="margin:4px 0 0;font-size:13px;color:#555;">Human-friendly interface at <code>/assistant/</code> for querying data with natural language.</p>' +
							'</div>'`
		assistantButton = `' + '<a href="/assistant/" style="display:inline-block;padding:7px 16px;background:#0f3d38;color:#fff;border-radius:6px;font:600 13px/1.4 sans-serif;text-decoration:none;">Open AI Chat</a>'`
	}

	return fmt.Sprintf(`function() {
				document.title = 'Safecast MCP API Docs';
				const mapDocsURL = %q;

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

				const style = document.createElement('link');
				style.rel = 'stylesheet'; style.href = '/mcp-api/swagger-theme.css';
				document.head.appendChild(style);

				const swaggerLogo = document.querySelector('.topbar-wrapper .link');
				if (swaggerLogo) swaggerLogo.remove();
				document.querySelectorAll('.topbar-wrapper img').forEach(img => img.remove());

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
						'applications can query radiation data programmatically.' +
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
							'</div>'%s +
						'</div>' +
						'<p style="margin:12px 0 0;font-size:13px;color:#777;">No authentication required. All data is CC0 licensed.</p>' +
					'</details>' +
					'<div style="display:flex;gap:10px;flex-wrap:wrap;align-items:center;">'%s +
						'<a href="' + mapDocsURL + '" style="display:inline-block;padding:7px 16px;background:#1a3a5c;color:#fff;border-radius:6px;font:600 13px/1.4 sans-serif;text-decoration:none;">\u2190 Switch to Map API</a>' +
					'</div>' +
				'</div>';
				const swaggerUIEl = document.getElementById('swagger-ui');
				if (swaggerUIEl) {
					document.body.insertBefore(preamble, swaggerUIEl);
				} else {
					document.body.prepend(preamble);
				}
			}`, mapDocsURL, assistantBlock, assistantButton)
}
