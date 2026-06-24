package errgroup

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunAllSuccess(t *testing.T) {
	var done int64
	var tasks []func(context.Context) error
	for i := 0; i < 5; i++ {
		tasks = append(tasks, func(ctx context.Context) error {
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt64(&done, 1)
			return nil
		})
	}

	if err := RunAll(context.Background(), tasks); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if done != 5 {
		t.Errorf("completed %d tasks, want 5", done)
	}
}

func TestRunAllReturnsErrorAndCancels(t *testing.T) {
	boom := errors.New("task failed")
	var canceledObserved int64

	tasks := []func(context.Context) error{
		// Fails quickly.
		func(ctx context.Context) error {
			time.Sleep(20 * time.Millisecond)
			return boom
		},
		// Long-running but cancellation-aware: must observe ctx.Done().
		func(ctx context.Context) error {
			select {
			case <-time.After(2 * time.Second):
				return nil
			case <-ctx.Done():
				atomic.AddInt64(&canceledObserved, 1)
				return ctx.Err()
			}
		},
	}

	start := time.Now()
	err := RunAll(context.Background(), tasks)
	elapsed := time.Since(start)

	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want %v", err, boom)
	}
	if elapsed >= time.Second {
		t.Errorf("RunAll took %v; the long task should have been cancelled, not run to completion", elapsed)
	}
	if atomic.LoadInt64(&canceledObserved) == 0 {
		t.Error("expected the second task to observe context cancellation")
	}
}

func TestRunAllEmpty(t *testing.T) {
	if err := RunAll(context.Background(), nil); err != nil {
		t.Errorf("RunAll(nil) = %v, want nil", err)
	}
}
