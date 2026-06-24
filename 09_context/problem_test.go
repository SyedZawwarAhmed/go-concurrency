package cancellable

import (
	"context"
	"errors"
	"testing"
	"time"
)

// feed returns a buffered channel pre-loaded with 0..n-1, then closed, so the
// producer never blocks regardless of how fast the worker drains it.
func feed(n int) <-chan int {
	ch := make(chan int, n)
	for i := 0; i < n; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

func TestWorkerCompletes(t *testing.T) {
	got, err := Worker(context.Background(), feed(5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 5 {
		t.Errorf("processed = %d, want 5", got)
	}
}

func TestWorkerCanceled(t *testing.T) {
	work := feed(100) // 100 * 20ms = ~2s of work; we cancel long before that.

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	got, err := Worker(ctx, work)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want context.Canceled", err)
	}
	if got >= 100 {
		t.Errorf("processed = %d; expected to stop early on cancellation", got)
	}
}

func TestWorkerDeadline(t *testing.T) {
	work := feed(100)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	got, err := Worker(ctx, work)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
	if got >= 100 {
		t.Errorf("processed = %d; expected to stop early on deadline", got)
	}
}
