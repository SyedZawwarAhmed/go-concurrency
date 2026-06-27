package boundedqueue

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBackpressureBlocksThenUnblocks(t *testing.T) {
	q := NewBoundedQueue(2)
	ctx := context.Background()

	if err := q.Put(ctx, 1); err != nil {
		t.Fatalf("Put(1) = %v, want nil (room available)", err)
	}
	if err := q.Put(ctx, 2); err != nil {
		t.Fatalf("Put(2) = %v, want nil (room available)", err)
	}

	// The queue is now full (capacity 2). A third Put must block (backpressure).
	thirdDone := make(chan error, 1)
	go func() {
		thirdDone <- q.Put(ctx, 3)
	}()

	// Give the third Put a chance to run; it should still be blocked.
	time.Sleep(30 * time.Millisecond)
	select {
	case err := <-thirdDone:
		t.Fatalf("third Put returned (err=%v) while queue is full, want it blocked (backpressure)", err)
	default:
	}

	// Draining one item frees a slot, so the blocked Put should proceed.
	if v, ok := q.Get(); !ok || v != 1 {
		t.Fatalf("Get() = (%d, %v), want (1, true)", v, ok)
	}

	select {
	case err := <-thirdDone:
		if err != nil {
			t.Fatalf("third Put = %v, want nil once a slot freed up", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("third Put did not unblock within 2s after a slot freed up")
	}
}

func TestPutCancelledWhileBlocked(t *testing.T) {
	q := NewBoundedQueue(1)

	if err := q.Put(context.Background(), 1); err != nil {
		t.Fatalf("Put(1) = %v, want nil (room available)", err)
	}

	// The queue is full; this Put blocks until the context's deadline fires.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- q.Put(ctx, 2)
	}()

	select {
	case err := <-done:
		if err == nil {
			t.Fatal("Put on a full queue with a cancelled context returned nil, want a context error")
		}
		if !errors.Is(err, context.Canceled) && !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("Put returned %v, want context.Canceled or context.DeadlineExceeded", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("blocked Put did not return after its context was done within 2s")
	}
}

func TestDrainsAfterClose(t *testing.T) {
	q := NewBoundedQueue(5)
	ctx := context.Background()

	for _, v := range []int{10, 20, 30} {
		if err := q.Put(ctx, v); err != nil {
			t.Fatalf("Put(%d) = %v, want nil", v, err)
		}
	}

	q.Close()

	for _, want := range []int{10, 20, 30} {
		v, ok := q.Get()
		if !ok || v != want {
			t.Fatalf("Get() = (%d, %v), want (%d, true)", v, ok, want)
		}
	}

	if v, ok := q.Get(); ok {
		t.Fatalf("Get() after drain = (%d, %v), want (_, false) on closed+empty queue", v, ok)
	}
}
