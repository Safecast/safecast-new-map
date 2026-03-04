// Package web provides HTTP handlers for the Safecast map application.
//
// Handlers are organized on a Server that holds shared dependencies (database,
// content, config) so main can wire routes without globals.
//
// Route overview (see each handlers_*.go file for details):
//
//   - /api/docs          — API usage documentation (HTML)
//   - /licenses/{mit,cc0} — License text
//   - /api/geoip         — Client lat/lon by IP (when AutoLocateDefault)
//   - /s/{code}          — Short URL redirect
//   - /api/spectrum/{id} — Gamma spectrum JSON or download (json/csv/n42/spe)
//   - /api/track-info/{id} — Track metadata (username, detector, date)
//   - /api/markers/spectra — Markers with spectra in a bounding box
//   - /api/update-coordinates — POST to set lat/lon for a track
//   - /api/tracks/bounds — Combined bounds for multiple tracks
//   - /qrpng             — QR code PNG for a URL
package web

import (
	"compress/gzip"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"safecast-new-map/pkg/database"
)

// Config holds application settings that handlers need. Main fills it from flags.
type Config struct {
	Domain                    string
	Port                      int
	DefaultLat                 float64
	DefaultLon                 float64
	DefaultZoom                int
	DefaultLayer               string
	AutoLocateDefault          bool
	SupportEmail               string
	CompileVersion             string
	DBType                     string
	AdminPassword              string
	APIDocsArchiveEnabled      bool
	APIDocsArchiveRoute        string
	APIDocsArchiveFrequency    string
	DebugIPAllowlist           map[string]struct{}
}

// Server holds dependencies for HTTP handlers. Create with NewServer and
// register routes with Register.
type Server struct {
	DB      *database.Database
	Content fs.FS
	Config  Config
	Logf    func(string, ...any)
}

// NewServer builds a Server with the given dependencies. Logf may be nil;
// if nil, no logging is done from handlers.
func NewServer(db *database.Database, content fs.FS, cfg Config, logf func(string, ...any)) *Server {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	return &Server{
		DB:      db,
		Content: content,
		Config:  cfg,
		Logf:    logf,
	}
}

// Register installs all web routes onto mux. Call this after configuring Server.
func (s *Server) Register(mux *http.ServeMux) {
	// Static/docs
	mux.HandleFunc("/api/docs", s.apiDocs)
	mux.HandleFunc("/licenses/", s.license)
	mux.HandleFunc("/api/geoip", s.gzipWrap(s.geoIP))
	mux.HandleFunc("/s/", s.shortRedirect)

	// Spectrum & track API
	mux.HandleFunc("/api/spectrum/", s.spectrum)
	mux.HandleFunc("/api/track-info/", s.trackInfo)
	mux.HandleFunc("/api/markers/spectra", s.markersWithSpectra)
	mux.HandleFunc("/api/update-coordinates", s.updateCoordinates)
	mux.HandleFunc("/api/tracks/bounds", s.apiTracksBounds)

	// QR
	mux.HandleFunc("/qrpng", s.qrPng)
}

// gzipWrap wraps a handler to apply gzip compression when client supports it.
func (s *Server) gzipWrap(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		grw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}
		next(grw, r)
	}
}

type gzipResponseWriter struct {
	http.ResponseWriter
	io.Writer
	written bool
}

func (g *gzipResponseWriter) Write(b []byte) (int, error) {
	if !g.written {
		g.ResponseWriter.WriteHeader(http.StatusOK)
		g.written = true
	}
	return g.Writer.Write(b)
}
