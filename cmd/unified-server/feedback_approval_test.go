package main

import "testing"

func TestFirstVote(t *testing.T) {
	if !firstVote(9001, "1.2.3.4") {
		t.Fatal("first vote should count")
	}
	if firstVote(9001, "1.2.3.4") {
		t.Fatal("repeat vote should not count")
	}
	if !firstVote(9001, "5.6.7.8") || !firstVote(9002, "1.2.3.4") {
		t.Fatal("different client or answer should count")
	}
}

func TestAutoApproveScoreDefault(t *testing.T) {
	t.Setenv("CHAT_AUTO_APPROVE_SCORE", "")
	if autoApproveScore() != 3 {
		t.Fatal("default should be 3")
	}
	t.Setenv("CHAT_AUTO_APPROVE_SCORE", "0")
	if autoApproveScore() != 0 {
		t.Fatal("0 disables")
	}
}
