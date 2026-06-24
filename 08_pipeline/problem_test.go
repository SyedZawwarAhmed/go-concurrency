package pipeline

import (
	"testing"
	"time"
)

// drain reads ch into a slice, with a watchdog so a stage that never closes its
// output fails the test instead of hanging forever.
func drain(t *testing.T, ch <-chan int) []int {
	t.Helper()
	var got []int
	done := make(chan struct{})
	go func() {
		defer close(done)
		for v := range ch {
			got = append(got, v)
		}
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("a pipeline stage never closed its output (drain timed out)")
	}
	return got
}

func equalOrdered(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestGenerate(t *testing.T) {
	got := drain(t, Generate(1, 2, 3))
	if !equalOrdered(got, []int{1, 2, 3}) {
		t.Errorf("Generate = %v, want [1 2 3]", got)
	}
}

func TestSquare(t *testing.T) {
	got := drain(t, Square(Generate(1, 2, 3, 4)))
	if !equalOrdered(got, []int{1, 4, 9, 16}) {
		t.Errorf("Square = %v, want [1 4 9 16]", got)
	}
}

func TestFilter(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	got := drain(t, Filter(Generate(1, 2, 3, 4, 5, 6), isEven))
	if !equalOrdered(got, []int{2, 4, 6}) {
		t.Errorf("Filter = %v, want [2 4 6]", got)
	}
}

func TestFullPipeline(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	got := drain(t, Filter(Square(Generate(1, 2, 3, 4, 5, 6)), isEven))
	want := []int{4, 16, 36}
	if !equalOrdered(got, want) {
		t.Errorf("pipeline = %v, want %v", got, want)
	}
}
