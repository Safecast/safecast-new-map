package tiles

import "testing"

func TestSafecastBreakpoints(t *testing.T) {
	for _, stop := range SafecastStops {
		r, g, b := ColorForDose(stop.Dose)
		if r != stop.R || g != stop.G || b != stop.B {
			t.Errorf("dose %g = rgb(%d,%d,%d), want rgb(%d,%d,%d)", stop.Dose, r, g, b, stop.R, stop.G, stop.B)
		}
	}
	// Below the first stop and above the last stop clamp, matching the map page.
	r, g, b := ColorForDose(0)
	if r != SafecastStops[0].R || g != SafecastStops[0].G || b != SafecastStops[0].B {
		t.Fatalf("low dose rgb(%d,%d,%d)", r, g, b)
	}
	last := SafecastStops[len(SafecastStops)-1]
	r, g, b = ColorForDose(100)
	if r != last.R || g != last.G || b != last.B {
		t.Fatalf("high dose rgb(%d,%d,%d)", r, g, b)
	}
	// Midway between the first two stops is the rounded blend the page uses.
	lo, hi := SafecastStops[0], SafecastStops[1]
	mid := (lo.Dose + hi.Dose) / 2
	r, g, b = ColorForDose(mid)
	wantR := lerp(lo.R, hi.R, 0.5)
	wantG := lerp(lo.G, hi.G, 0.5)
	wantB := lerp(lo.B, hi.B, 0.5)
	if r != wantR || g != wantG || b != wantB {
		t.Fatalf("mid rgb(%d,%d,%d), want rgb(%d,%d,%d)", r, g, b, wantR, wantG, wantB)
	}
}
