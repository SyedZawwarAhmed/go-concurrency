package waitgroups

import (
	"runtime"
	"testing"
	"time"
)

func TestMassageStrings(t *testing.T) {
	tests := []struct {
		name   string
		inputs []string
		want   []string
	}{
		{"basic", []string{" hello ", "World", "go"}, []string{"HELLO", "WORLD", "GO"}},
		{"single", []string{"x"}, []string{"X"}},
		{"empty", []string{}, []string{}},
		{"nil", nil, []string{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := MassageStrings(tc.inputs)
			if len(got) != len(tc.want) {
				t.Fatalf("len = %d, want %d (got %v)", len(got), len(tc.want), got)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("index %d = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestMassageStringsRunsConcurrently(t *testing.T) {
	inputs := make([]string, 10)
	for i := range inputs {
		inputs[i] = "item"
	}

	start := time.Now()
	got := MassageStrings(inputs)
	elapsed := time.Since(start)

	if len(got) != 10 {
		t.Fatalf("got %d results, want 10", len(got))
	}
	// Serial processing would take ~500ms (10 * 50ms). Concurrent should be ~50ms.
	if elapsed >= 150*time.Millisecond {
		t.Errorf("took %v; expected ~50ms (concurrent), not ~500ms (serial)", elapsed)
	}
	// Guard against "cheating" by skipping the simulated work entirely.
	if elapsed < processDelay {
		t.Errorf("took %v; expected at least one %v delay — is the work actually happening?", elapsed, processDelay)
	}
}

func TestMassageStringsNoGoroutineLeak(t *testing.T) {
	// Let any startup goroutines settle before measuring the baseline.
	time.Sleep(50 * time.Millisecond)
	before := runtime.NumGoroutine()

	_ = MassageStrings([]string{"a", "b", "c", "d", "e"})

	// Give the spawned goroutines time to exit.
	time.Sleep(100 * time.Millisecond)
	after := runtime.NumGoroutine()

	if after > before {
		t.Errorf("goroutine leak: before=%d after=%d (every goroutine should have exited)", before, after)
	}
}
