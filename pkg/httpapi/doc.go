// Package httpapi provides the unified HTTP and API transport layer for the Safecast map application.
//
// It consolidates all API routes, web handlers, and supporting infrastructure in one place:
//
//   - Core API: overview, latest nearby markers, tracks list, track data by ID/index/year/month,
//     countries, short URLs, and optional weekly JSON archive download.
//   - Web handlers: API docs, licenses, geoip, short redirects, spectrum data, track metadata,
//     markers with spectra, update-coordinates (admin), track bounds, and QR code generation.
//   - Auth and admin: registration, login, user profile, uploads, and admin-only routes are
//     registered here from the shared Register entry point.
//
// Usage: build a RegisterConfig with WebServer (from NewWebServer), APIHandler (from NewHandler),
// AuthManager, and optional admin handlers, then call Register(mux, cfg) to attach all routes.
package httpapi
