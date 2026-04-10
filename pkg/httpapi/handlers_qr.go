// handlers_qr.go — GET /qrpng returns a PNG QR code for a URL (query "u", or Referer, or current page). Includes Safecast logo.
package httpapi

import (
	"image/color"
	"net/http"

	"safecast-new-map/pkg/qrlogoext"
)

// qrPng renders a QR code PNG.
//
// @Summary     Generate QR PNG
// @Description Returns a QR image for query parameter `u`, referer, or current URL.
// @Tags        web
// @Produce     image/png
// @Param       u query string false "Target URL to encode"
// @Success     200 {file} file "QR PNG"
// @Failure     500 {string} string "QR generation failed"
// @Router      /qrpng [get]
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
		writeJSONError(w, http.StatusInternalServerError, "QR encode: "+err.Error())
	}
}
