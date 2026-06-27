package tokenbucket

import (
	"testing"
	"time"
)

func TestBurstThenEmpty(t *testing.T) {
	l := NewLimiter(50*time.Millisecond, 3)
	defer l.Stop()

	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("Allow call %d: got false, want true (bucket starts full with 3 tokens)", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("4th Allow: got true, want false (bucket drained, no time elapsed)")
	}
}

func TestRefillAddsOneToken(t *testing.T) {
	l := NewLimiter(50*time.Millisecond, 3)
	defer l.Stop()

	// Drain the 3 initial tokens.
	for i := 0; i < 3; i++ {
		if !l.Allow() {
			t.Fatalf("drain Allow call %d: got false, want true", i+1)
		}
	}
	if l.Allow() {
		t.Fatal("after draining: Allow got true, want false (empty bucket)")
	}

	// One refill interval passes (80ms > 50ms): exactly one token added.
	time.Sleep(80 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("after one refill interval: Allow got false, want true")
	}
	if l.Allow() {
		t.Fatal("after one refill: second Allow got true, want false (only one token refilled)")
	}
}

func TestRefillCapsAtBurst(t *testing.T) {
	l := NewLimiter(10*time.Millisecond, 2)

	// Many refill intervals pass; tokens must cap at burst (2), not pile up.
	time.Sleep(100 * time.Millisecond)

	// Stop FIRST so the refiller can't add a token while we drain (avoids a race).
	l.Stop()

	available := 0
	for l.Allow() {
		available++
	}
	if available != 2 {
		t.Fatalf("drained %d tokens, want exactly 2 (capped at burst)", available)
	}
}

func TestStopHaltsRefill(t *testing.T) {
	l := NewLimiter(10*time.Millisecond, 1)

	if !l.Allow() {
		t.Fatal("first Allow: got false, want true (bucket starts full)")
	}
	l.Stop()

	// Several refill intervals would pass, but Stop halted the refiller.
	time.Sleep(50 * time.Millisecond)
	if l.Allow() {
		t.Fatal("after Stop: Allow got true, want false (no refill after Stop)")
	}
}
