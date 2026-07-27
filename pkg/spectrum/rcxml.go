package spectrum

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"safecast-new-map/pkg/database"
)

// RadiaCode XML format structures (ResultDataFile)
// This is the proprietary XML format exported by RadiaCode devices and apps.

// RCXMLResultDataFile is the root element of RadiaCode XML spectrum files.
type RCXMLResultDataFile struct {
	XMLName        xml.Name          `xml:"ResultDataFile"`
	FormatVersion  string            `xml:"FormatVersion"`
	ResultDataList RCXMLResultDataList `xml:"ResultDataList"`
}

// RCXMLResultDataList contains one or more spectrum results.
type RCXMLResultDataList struct {
	ResultData []RCXMLResultData `xml:"ResultData"`
}

// RCXMLResultData represents a single spectrum measurement.
type RCXMLResultData struct {
	DeviceConfigReference RCXMLDeviceConfig   `xml:"DeviceConfigReference"`
	SampleInfo            RCXMLSampleInfo     `xml:"SampleInfo"`
	BackgroundSpectrumFile string             `xml:"BackgroundSpectrumFile"`
	StartTime             string              `xml:"StartTime"`
	EndTime               string              `xml:"EndTime"`
	EnergySpectrum        RCXMLEnergySpectrum `xml:"EnergySpectrum"`
	Visible               bool                `xml:"Visible"`
}

// RCXMLDeviceConfig contains device identification.
type RCXMLDeviceConfig struct {
	Name string `xml:"Name"`
}

// RCXMLSampleInfo contains sample metadata.
type RCXMLSampleInfo struct {
	Name string `xml:"Name"`
	Note string `xml:"Note"`
}

// RCXMLEnergySpectrum contains the actual spectrum data.
type RCXMLEnergySpectrum struct {
	NumberOfChannels  int                    `xml:"NumberOfChannels"`
	ChannelPitch      int                    `xml:"ChannelPitch"`
	SpectrumName      string                 `xml:"SpectrumName"`
	Comment           string                 `xml:"Comment"`
	SerialNumber      string                 `xml:"SerialNumber"`
	EnergyCalibration RCXMLEnergyCalibration `xml:"EnergyCalibration"`
	MeasurementTime   float64                `xml:"MeasurementTime"`
	LiveTime          float64                `xml:"LiveTime"`
	Spectrum          RCXMLSpectrum          `xml:"Spectrum"`
}

// RCXMLEnergyCalibration contains polynomial calibration coefficients.
type RCXMLEnergyCalibration struct {
	PolynomialOrder int                   `xml:"PolynomialOrder"`
	Coefficients    RCXMLCoefficients     `xml:"Coefficients"`
}

// RCXMLCoefficients contains the calibration coefficient values.
type RCXMLCoefficients struct {
	Coefficient []string `xml:"Coefficient"`
}

// RCXMLSpectrum contains the channel data points.
type RCXMLSpectrum struct {
	DataPoint []string `xml:"DataPoint"`
}

// ParseRCXML parses RadiaCode proprietary XML format and extracts spectrum data.
// Returns a slice of Spectrum objects with their associated marker data.
func ParseRCXML(data []byte) ([]database.Spectrum, []database.Marker, error) {
	var doc RCXMLResultDataFile
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, nil, fmt.Errorf("parse RadiaCode XML: %w", err)
	}

	var spectra []database.Spectrum
	var markers []database.Marker

	fmt.Printf("RCXML Debug: Found %d result data entries\n", len(doc.ResultDataList.ResultData))

	for i, rd := range doc.ResultDataList.ResultData {
		spectrum, marker, err := convertRCXMLSpectrum(&rd, data)
		if err != nil {
			fmt.Printf("RCXML Debug: Skipping entry[%d]: %v\n", i, err)
			continue
		}
		if spectrum != nil {
			spectra = append(spectra, *spectrum)
		}
		if marker != nil {
			markers = append(markers, *marker)
		}
	}

	if len(spectra) == 0 {
		return nil, nil, fmt.Errorf("no valid spectra found in RadiaCode XML file (found %d entries)",
			len(doc.ResultDataList.ResultData))
	}

	return spectra, markers, nil
}

// convertRCXMLSpectrum converts a RadiaCode XML spectrum entry to database format.
func convertRCXMLSpectrum(rd *RCXMLResultData, rawData []byte) (*database.Spectrum, *database.Marker, error) {
	es := &rd.EnergySpectrum

	// Parse channel data from DataPoint elements
	channels, err := parseRCXMLChannelData(es.Spectrum.DataPoint)
	if err != nil {
		return nil, nil, fmt.Errorf("parse channel data: %w", err)
	}

	if len(channels) == 0 {
		return nil, nil, fmt.Errorf("empty channel data")
	}

	// Parse energy calibration
	calibration, err := parseRCXMLCalibration(&es.EnergyCalibration)
	if err != nil {
		// Use default calibration if parsing fails
		fmt.Printf("RCXML Debug: Using default calibration: %v\n", err)
		calibration = &database.EnergyCalibration{
			A: 0,
			B: 3000.0 / float64(len(channels)),
			C: 0,
		}
	}

	// Parse measurement time
	liveTime := es.LiveTime
	realTime := es.MeasurementTime
	if realTime == 0 {
		realTime = liveTime
	}
	// RadiaCode files often omit LiveTime and only report MeasurementTime (real time).
	// Fall back to it for rate calculations so markers aren't filtered as zero-dose.
	effectiveTime := liveTime
	if effectiveTime == 0 {
		effectiveTime = realTime
	}

	// Parse timestamp from StartTime
	timestamp := time.Now().Unix()
	if rd.StartTime != "" {
		if t, err := parseRCXMLTime(rd.StartTime); err == nil {
			timestamp = t.Unix()
		}
	}

	// Extract device model
	deviceModel := rd.DeviceConfigReference.Name
	if deviceModel == "" {
		deviceModel = "RadiaCode"
	}
	// Append serial number if available
	if es.SerialNumber != "" {
		deviceModel = fmt.Sprintf("%s (%s)", deviceModel, es.SerialNumber)
	}

	// Calculate energy range from calibration
	energyMin := ChannelToEnergy(0, calibration)
	energyMax := ChannelToEnergy(len(channels)-1, calibration)
	if energyMin > energyMax {
		energyMin, energyMax = energyMax, energyMin
	}
	if energyMin < 0 {
		energyMin = 0
	}

	spectrum := &database.Spectrum{
		Channels:     channels,
		ChannelCount: len(channels),
		EnergyMinKeV: energyMin,
		EnergyMaxKeV: energyMax,
		LiveTimeSec:  liveTime,
		RealTimeSec:  realTime,
		DeviceModel:  deviceModel,
		Calibration:  calibration,
		SourceFormat: "rcxml",
		RawData:      rawData,
		CreatedAt:    time.Now().Unix(),
	}

	// Calculate dose rate and count rate
	doseRate := 0.0
	if effectiveTime > 0 {
		doseRate = CalculateDoseRate(channels, effectiveTime, calibration)
	}

	countRate := 0.0
	if effectiveTime > 0 {
		countRate = float64(sumChannels(channels)) / effectiveTime
	}

	// Create marker (no GPS coordinates in this format - user can add manually)
	marker := &database.Marker{
		Lat:         0,
		Lon:         0,
		DoseRate:    doseRate,
		CountRate:   countRate,
		Date:        timestamp,
		HasSpectrum: true,
		Detector:    rd.DeviceConfigReference.Name,
	}

	fmt.Printf("RCXML Debug: Parsed spectrum with %d channels, live=%.1fs, real=%.1fs, device=%s\n",
		len(channels), liveTime, realTime, deviceModel)

	return spectrum, marker, nil
}

// parseRCXMLChannelData parses channel counts from DataPoint elements.
func parseRCXMLChannelData(dataPoints []string) ([]int, error) {
	if len(dataPoints) == 0 {
		return nil, fmt.Errorf("no data points")
	}

	channels := make([]int, 0, len(dataPoints))
	for i, dp := range dataPoints {
		dp = strings.TrimSpace(dp)
		if dp == "" {
			channels = append(channels, 0)
			continue
		}
		count, err := strconv.Atoi(dp)
		if err != nil {
			return nil, fmt.Errorf("parse data point %d (%q): %w", i, dp, err)
		}
		channels = append(channels, count)
	}

	return channels, nil
}

// parseRCXMLCalibration parses energy calibration coefficients.
// RadiaCode uses polynomial: Energy = A + B*ch + C*ch^2
func parseRCXMLCalibration(cal *RCXMLEnergyCalibration) (*database.EnergyCalibration, error) {
	if cal == nil || len(cal.Coefficients.Coefficient) < 2 {
		return nil, fmt.Errorf("insufficient calibration coefficients")
	}

	coeffs := cal.Coefficients.Coefficient

	// Parse coefficient A (offset)
	a, err := strconv.ParseFloat(strings.TrimSpace(coeffs[0]), 64)
	if err != nil {
		return nil, fmt.Errorf("parse coefficient A: %w", err)
	}

	// Parse coefficient B (linear)
	b, err := strconv.ParseFloat(strings.TrimSpace(coeffs[1]), 64)
	if err != nil {
		return nil, fmt.Errorf("parse coefficient B: %w", err)
	}

	// Parse coefficient C (quadratic) if available
	c := 0.0
	if len(coeffs) >= 3 {
		c, _ = strconv.ParseFloat(strings.TrimSpace(coeffs[2]), 64)
	}

	return &database.EnergyCalibration{A: a, B: b, C: c}, nil
}

// parseRCXMLTime parses RadiaCode XML timestamp format.
// Format: "2026-01-06T16:28:52"
func parseRCXMLTime(timeStr string) (time.Time, error) {
	timeStr = strings.TrimSpace(timeStr)
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("empty time string")
	}

	// Try common formats
	formats := []string{
		"2006-01-02T15:04:05",       // Standard ISO without timezone
		"2006-01-02T15:04:05Z07:00", // ISO with timezone
		time.RFC3339,
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time: %s", timeStr)
}

// IsRCXMLFormat checks if the data appears to be RadiaCode XML format.
// This helps with format detection when file extension is ambiguous.
func IsRCXMLFormat(data []byte) bool {
	// Quick check for ResultDataFile root element
	return strings.Contains(string(data[:min(500, len(data))]), "<ResultDataFile")
}
