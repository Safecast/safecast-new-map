// Package web provides HTTP handlers for the Safecast map application.
// Handlers are organized on a Server that holds shared dependencies (database,
// content, config) so main can wire routes without globals.
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
		w.Header().Add("Vary", "Accept-Encoding")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gw := &gzipResponseWriter{ResponseWriter: w, Writer: gz}
		next(gw, r)
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
