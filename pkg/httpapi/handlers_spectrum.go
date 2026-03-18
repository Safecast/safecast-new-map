// handlers_spectrum.go — GET /api/spectrum/{markerID} and .../download?format=json|csv|n42|spe.
package httpapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// spectrum returns spectrum JSON for the marker.
//
// @Summary     Get spectrum by marker ID
// @Description Returns spectroscopy data for a marker. Download paths are delegated to the download handler.
// @Tags        web
// @Produce     json
// @Param       markerID path int true "Marker ID"
// @Success     200 {object} map[string]interface{} "Spectrum payload"
// @Failure     400 {string} string "Invalid marker ID"
// @Failure     404 {string} string "Spectrum not found"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/spectrum/{markerID} [get]
func (s *Server) spectrum(w http.ResponseWriter, r *http.Request) {
	if strings.Contains(r.URL.Path, "/download") {
		s.spectrumDownload(w, r)
		return
	}
	if s.DB == nil || s.DB.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "Marker ID not provided", http.StatusBadRequest)
		return
	}
	markerID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid marker ID", http.StatusBadRequest)
		return
	}
	ctx, cancel := WithMinimumDeadline(r.Context(), 30*time.Second)
	defer cancel()
	spectrum, err := s.DB.GetSpectrum(ctx, markerID)
	if err != nil {
		log.Printf("Error fetching spectrum for marker %d: %v", markerID, err)
		http.Error(w, "Error fetching spectrum", http.StatusInternalServerError)
		return
	}
	if spectrum == nil {
		http.Error(w, "No spectrum found for this marker", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spectrum)
}

// spectrumDownload returns spectrum content as a downloadable attachment.
//
// @Summary     Download spectrum file
// @Description Downloads a spectrum as json, csv, n42, or spe.
// @Tags        web
// @Produce     application/octet-stream
// @Param       markerID path int true "Marker ID"
// @Param       format query string false "Download format: json, csv, n42, spe"
// @Success     200 {file} file "Spectrum file"
// @Failure     400 {string} string "Invalid marker or format"
// @Failure     404 {string} string "Spectrum or format data unavailable"
// @Failure     503 {string} string "Database unavailable"
// @Router      /api/spectrum/{markerID}/download [get]
func (s *Server) spectrumDownload(w http.ResponseWriter, r *http.Request) {
	if s.DB == nil || s.DB.DB == nil {
		http.Error(w, "Database not available", http.StatusServiceUnavailable)
		return
	}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		http.Error(w, "Marker ID not provided", http.StatusBadRequest)
		return
	}
	markerID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid marker ID", http.StatusBadRequest)
		return
	}
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}
	ctx, cancel := WithMinimumDeadline(r.Context(), 30*time.Second)
	defer cancel()
	spectrum, err := s.DB.GetSpectrum(ctx, markerID)
	if err != nil {
		log.Printf("Error fetching spectrum for marker %d: %v", markerID, err)
		http.Error(w, "Error fetching spectrum", http.StatusInternalServerError)
		return
	}
	if spectrum == nil {
		http.Error(w, "No spectrum found for this marker", http.StatusNotFound)
		return
	}
	switch strings.ToLower(format) {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=spectrum_%d.json", markerID))
		json.NewEncoder(w).Encode(spectrum)
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=spectrum_%d.csv", markerID))
		fmt.Fprintf(w, "Channel,Energy_keV,Counts\n")
		for i, count := range spectrum.Channels {
			energy := 0.0
			if spectrum.Calibration != nil {
				ch := float64(i)
				energy = spectrum.Calibration.A + spectrum.Calibration.B*ch + spectrum.Calibration.C*ch*ch
			}
			fmt.Fprintf(w, "%d,%.2f,%d\n", i, energy, count)
		}
	case "n42":
		if len(spectrum.RawData) > 0 && spectrum.SourceFormat == "n42" {
			w.Header().Set("Content-Type", "application/xml")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=spectrum_%d.n42", markerID))
			w.Write(spectrum.RawData)
		} else {
			http.Error(w, "N42 format not available for this spectrum", http.StatusNotFound)
		}
	case "spe":
		if len(spectrum.RawData) > 0 && spectrum.SourceFormat == "spe" {
			w.Header().Set("Content-Type", "text/plain")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=spectrum_%d.spe", markerID))
			w.Write(spectrum.RawData)
		} else {
			http.Error(w, "SPE format not available for this spectrum", http.StatusNotFound)
		}
	default:
		http.Error(w, fmt.Sprintf("Unsupported format: %s", format), http.StatusBadRequest)
	}
}
