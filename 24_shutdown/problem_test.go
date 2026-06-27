package shutdown

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestProcessesAllSubmittedBeforeShutdown(t *testing.T) {
	var counter int64
	s := NewServer(4, func(int) {
		atomic.AddInt64(&counter, 1)
	})

	for i := 0; i < 100; i++ {
		if !s.Submit(i) {
			t.Fatalf("Submit(%d) returned false before shutdown", i)
		}
	}

	processed, err := s.Shutdown(context.Background())
	if err != nil {
		t.Fatalf("Shutdown error = %v, want nil", err)
	}
	if processed != 100 {
		t.Fatalf("processed = %d, want 100", processed)
	}
	if got := atomic.LoadInt64(&counter); got != 100 {
		t.Fatalf("counter = %d, want 100", got)
	}
}

func TestSubmitRejectedAfterShutdown(t *testing.T) {
	var counter int64
	s := NewServer(4, func(int) {
		time.Sleep(1 * time.Millisecond)
		atomic.AddInt64(&counter, 1)
	})

	for i := 0; i < 10; i++ {
		if !s.Submit(i) {
			t.Fatalf("Submit(%d) returned false before shutdown", i)
		}
	}

	if _, err := s.Shutdown(context.Background()); err != nil {
		t.Fatalf("Shutdown error = %v, want nil", err)
	}

	if s.Submit(999) {
		t.Fatal("Submit after shutdown returned true, want false")
	}
}

func TestShutdownRespectsDeadline(t *testing.T) {
	s := NewServer(2, func(int) {
		time.Sleep(100 * time.Millisecond)
	})

	for i := 0; i < 20; i++ {
		s.Submit(i)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	processed, err := s.Shutdown(ctx)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Shutdown error = %v, want context.DeadlineExceeded", err)
	}
	if processed >= 20 {
		t.Fatalf("processed = %d, want < 20 (could not finish in 50ms)", processed)
	}
	if elapsed >= 500*time.Millisecond {
		t.Fatalf("Shutdown took %v, want it to return promptly on deadline", elapsed)
	}
}
