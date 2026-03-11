// server_web.go defines the web/API Server, its config, and route registration.
package httpapi

import (
	"compress/gzip"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"safecast-new-map/pkg/auth"
	"safecast-new-map/pkg/database"
)

// WebConfig holds application settings for web and API handlers (domain, port, map defaults,
// admin password, API docs archive settings, etc.). Filled from flags in main.
type WebConfig struct {
	Domain                  string
	Port                    int
	DefaultLat              float64
	DefaultLon              float64
	DefaultZoom             int
	DefaultLayer            string
	AutoLocateDefault       bool   // If true, /api/geoip returns client lat/lon.
	SupportEmail            string
	CompileVersion          string
	DBType                  string
	AdminPassword           string
	APIDocsArchiveEnabled   bool
	APIDocsArchiveRoute     string
	APIDocsArchiveFrequency string
	DebugIPAllowlist        map[string]struct{}
}

// Server holds dependencies for docs, spectrum, track-info, markers, bounds, short URLs, and QR handlers.
type Server struct {
	DB       *database.Database
	Services Services
	Content  fs.FS
	Config   WebConfig
	Logf     func(string, ...any)
}

// NewWebServer builds a Server with the given database, embedded content FS, config, and optional logger.
func NewWebServer(db *database.Database, content fs.FS, cfg WebConfig, logf func(string, ...any)) *Server {
	if logf == nil {
		logf = func(string, ...any) {}
	}
	return &Server{
		DB:       db,
		Services: newDatabaseServices(db),
		Content:  content,
		Config:   cfg,
		Logf:     logf,
	}
}

// Register installs all web routes onto mux: docs, licenses, geoip, short redirects,
// spectrum, track-info, markers/spectra, update-coordinates, tracks/bounds, and qrpng.
func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/docs", s.apiDocs)
	mux.HandleFunc("/licenses/", s.license)
	mux.HandleFunc("/api/geoip", s.gzipWrap(s.geoIP))
	mux.HandleFunc("/s/", s.shortRedirect)
	mux.HandleFunc("/api/spectrum/", s.spectrum)
	mux.HandleFunc("/api/track-info/", s.trackInfo)
	mux.HandleFunc("/api/markers/spectra", s.markersWithSpectra)
	mux.HandleFunc("/api/update-coordinates", auth.RequireStaticAdminBasic(s.Config.AdminPassword, s.updateCoordinates))
	mux.HandleFunc("/api/tracks/bounds", s.apiTracksBounds)
	mux.HandleFunc("/qrpng", s.qrPng)
}

// gzipWrap compresses the response with gzip when the client sends Accept-Encoding: gzip.
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

// gzipResponseWriter wraps the response writer so Write goes to the gzip writer; first Write triggers 200.
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
