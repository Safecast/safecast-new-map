package web

import (
	"image/color"
	"net/http"

	"safecast-new-map/pkg/qrlogoext"
)

// qrPng generates a QR code image for the given URL (query u or referer).
func (s *Server) qrPng(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Query().Get("u")
	if u == "" {
		if ref := r.Referer(); ref != "" {
			u = ref
		} else {
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			u = scheme + "://" + r.Host + r.URL.RequestURI()
		}
	}
	if len(u) > 4096 {
		u = u[:4096]
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Disposition", "inline; filename=\"qr.png\"")
	var logoBytes []byte
	opts := qrlogoext.Options{
		TargetPx:    1500,
		Fg:          color.RGBA{0, 0, 0, 255},
		Bg:          color.RGBA{255, 255, 255, 255},
		Logo:        color.RGBA{233, 192, 35, 255},
		LogoBoxFrac: 0.32,
		LogoPadding: 16,
	}
	if err := qrlogoext.EncodePNG(w, []byte(u), logoBytes, opts); err != nil {
		http.Error(w, "QR encode: "+err.Error(), http.StatusInternalServerError)
	}
}
