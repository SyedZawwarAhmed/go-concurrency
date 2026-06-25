package template

import (
	"testing"
	"time"
)

// Happy path + edge cases. Replace with cases real for your problem.
func TestSolve(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"basic", []int{1, 2, 3}, []int{ /* ... */ }},
		{"empty", []int{}, []int{}},
		{"nil", nil, []int{}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := Solve(tc.input)
			if len(got) != len(tc.want) {
				t.Fatalf("len = %d, want %d (got %v)", len(got), len(tc.want), got)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("index %d = %v, want %v", i, got[i], tc.want[i])
				}
			}
		})
	}
}

// Concurrency property: prove the work happens in parallel, not serially.
// Compare elapsed time against a serial baseline. (Delete if not applicable.)
func TestSolveRunsConcurrently(t *testing.T) {
	t.Skip("template: replace with a real timing assertion for your problem")
}

// Watchdog pattern: fail fast on a deadlock instead of hanging the suite.
// Wrap the call in a goroutine and race it against time.After.
func TestSolveDoesNotDeadlock(t *testing.T) {
	done := make(chan struct{})
	go func() {
		_ = Solve([]int{1, 2, 3})
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Solve did not return within 2s — likely deadlocked")
	}
}
