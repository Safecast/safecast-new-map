package main

import "testing"

func TestSummarizeCorrection(t *testing.T) {
	got := summarizeCorrection("The level is 0.5 uSv/h near the plant.", "The level is 0.05 uSv/h near the plant.")
	if got != `replaced "0.5" with "0.05"` {
		t.Fatalf("got %q", got)
	}
	if summarizeCorrection("a b", "a b") != "" {
		t.Fatal("identical must give empty note")
	}
}
