package pipeline

import (
	"context"
	"runtime"
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
	got := drain(t, Generate(context.Background(), 1, 2, 3))
	if !equalOrdered(got, []int{1, 2, 3}) {
		t.Errorf("Generate = %v, want [1 2 3]", got)
	}
}

func TestSquare(t *testing.T) {
	ctx := context.Background()
	got := drain(t, Square(ctx, Generate(ctx, 1, 2, 3, 4)))
	if !equalOrdered(got, []int{1, 4, 9, 16}) {
		t.Errorf("Square = %v, want [1 4 9 16]", got)
	}
}

func TestFilter(t *testing.T) {
	ctx := context.Background()
	isEven := func(n int) bool { return n%2 == 0 }
	got := drain(t, Filter(ctx, Generate(ctx, 1, 2, 3, 4, 5, 6), isEven))
	if !equalOrdered(got, []int{2, 4, 6}) {
		t.Errorf("Filter = %v, want [2 4 6]", got)
	}
}

func TestFullPipeline(t *testing.T) {
	ctx := context.Background()
	isEven := func(n int) bool { return n%2 == 0 }
	got := drain(t, Filter(ctx, Square(ctx, Generate(ctx, 1, 2, 3, 4, 5, 6)), isEven))
	want := []int{4, 16, 36}
	if !equalOrdered(got, want) {
		t.Errorf("pipeline = %v, want %v", got, want)
	}
}

// TestPipelineNoLeakOnAbandon proves the gap: if a consumer stops reading early
// (the common real-world case — an error, a `break`, a cancelled request), the
// pipeline must still tear down. With the current implementation each stage
// blocks forever on `ch <- v`, leaking one goroutine per stage.
//
// We start a long pipeline, read exactly ONE value, then walk away. After a
// settle period the goroutine count should return to baseline.
func TestPipelineNoLeakOnAbandon(t *testing.T) {
	isPositive := func(n int) bool { return n > 0 }

	// Let goroutines from prior tests settle, then take a baseline.
	time.Sleep(50 * time.Millisecond)
	runtime.GC()
	before := runtime.NumGoroutine()

	// A pipeline with plenty of work still waiting to flow through.
	nums := make([]int, 1000)
	for i := range nums {
		nums[i] = i + 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	out := Filter(ctx, Square(ctx, Generate(ctx, nums...)), isPositive)

	// Consume a single value, then abandon the pipeline by cancelling.
	<-out
	cancel()

	// Give the stages time to observe cancellation and exit.
	time.Sleep(100 * time.Millisecond)
	runtime.GC()
	after := runtime.NumGoroutine()

	if leaked := after - before; leaked > 0 {
		t.Fatalf("pipeline leaked %d goroutine(s) after the consumer abandoned it "+
			"(before=%d, after=%d); stages need to stop on cancellation",
			leaked, before, after)
	}
}
