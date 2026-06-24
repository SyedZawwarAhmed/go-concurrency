package workerpool

import (
	"sort"
	"sync/atomic"
	"testing"
)

func TestProcessSquares(t *testing.T) {
	tests := []struct {
		name    string
		n       int
		workers int
	}{
		{"50 jobs / 5 workers", 50, 5},
		{"more workers than jobs", 3, 10},
		{"single worker", 20, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			jobs := make([]int, tc.n)
			want := make([]int, tc.n)
			for i := 0; i < tc.n; i++ {
				jobs[i] = i
				want[i] = i * i
			}

			got := Process(jobs, tc.workers, func(x int) int { return x * x })
			if len(got) != len(want) {
				t.Fatalf("got %d results, want %d", len(got), len(want))
			}
			sort.Ints(got)
			sort.Ints(want)
			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("result multiset mismatch: got %v, want %v", got, want)
				}
			}
		})
	}
}

func TestProcessEachJobOnce(t *testing.T) {
	const n = 50
	jobs := make([]int, n)
	for i := range jobs {
		jobs[i] = i
	}

	var calls int64
	got := Process(jobs, 5, func(x int) int {
		atomic.AddInt64(&calls, 1)
		return x
	})

	if calls != n {
		t.Errorf("fn called %d times, want exactly %d", calls, n)
	}
	if len(got) != n {
		t.Errorf("got %d results, want %d", len(got), n)
	}
}

func TestProcessEmpty(t *testing.T) {
	got := Process(nil, 4, func(x int) int { return x })
	if len(got) != 0 {
		t.Errorf("empty jobs should yield no results, got %v", got)
	}
}
