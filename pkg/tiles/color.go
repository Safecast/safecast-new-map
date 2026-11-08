package tiles

import "math"

// Stop is one breakpoint on the Safecast dose-rate ramp.
// The values match schemeSafecast in public_html/map.html. A pixel at a
// breakpoint uses that color exactly, so the overlay and the circle markers
// agree where the legend is labeled.
type Stop struct {
	Dose    float64
	R, G, B uint8
}

// SafecastStops is the scientific ramp. Doses at or above the last stop use
// that stop's color, matching getGradientColor in the map page.
var SafecastStops = []Stop{
	{0.03, 0, 0, 139},
	{0.05, 0, 71, 171},
	{0.08, 65, 105, 225},
	{0.12, 135, 206, 250},
	{0.16, 0, 191, 255},
	{0.23, 0, 255, 255},
	{0.31, 64, 224, 208},
	{0.43, 147, 112, 219},
	{0.60, 186, 85, 211},
	{0.87, 255, 0, 255},
	{1.31, 255, 20, 147},
	{2.13, 255, 0, 0},
	{3.99, 255, 69, 0},
	{10.09, 255, 165, 0},
	{65.54, 255, 255, 0},
}

// ColorForDose returns the fully opaque ramp color for a dose in µSv/h.
// Callers paint unsampled pixels as transparent instead of fading this color.
func ColorForDose(dose float64) (r, g, b uint8) {
	if dose < 0 {
		dose = 0
	}
	stops := SafecastStops
	if dose <= stops[0].Dose {
		return stops[0].R, stops[0].G, stops[0].B
	}
	last := stops[len(stops)-1]
	if dose >= last.Dose {
		return last.R, last.G, last.B
	}
	for i := 0; i < len(stops)-1; i++ {
		lo, hi := stops[i], stops[i+1]
		if dose >= lo.Dose && dose <= hi.Dose {
			f := (dose - lo.Dose) / (hi.Dose - lo.Dose)
			return lerp(lo.R, hi.R, f), lerp(lo.G, hi.G, f), lerp(lo.B, hi.B, f)
		}
	}
	return last.R, last.G, last.B
}

func lerp(a, b uint8, f float64) uint8 {
	return uint8(math.Round(float64(a) + (float64(b)-float64(a))*f))
}
