package retrypool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetriesUntilSuccess(t *testing.T) {
	items := []int{0, 1, 2, 3, 4}
	attempts := make([]atomic.Int64, len(items))

	fn := func(ctx context.Context, item int) (int, error) {
		n := attempts[item].Add(1)
		if n < 3 {
			return 0, errors.New("transient")
		}
		return item * 10, nil
	}

	outcomes := Process(context.Background(), items, 3, 3, fn)

	if len(outcomes) != len(items) {
		t.Fatalf("got %d outcomes, want %d", len(outcomes), len(items))
	}
	for i, o := range outcomes {
		if o.Item != items[i] {
			t.Errorf("outcome %d: Item=%d, want %d (out of order)", i, o.Item, items[i])
		}
		if o.Err != nil {
			t.Errorf("outcome %d: Err=%v, want nil", i, o.Err)
		}
		if o.Value != items[i]*10 {
			t.Errorf("outcome %d: Value=%d, want %d", i, o.Value, items[i]*10)
		}
		if o.Attempts != 3 {
			t.Errorf("outcome %d: Attempts=%d, want 3", i, o.Attempts)
		}
	}
}

var errBoom = errors.New("boom")

func TestRecordsErrorAfterMaxAttempts(t *testing.T) {
	items := []int{0, 1, 2}

	fn := func(ctx context.Context, item int) (int, error) {
		return 0, errBoom
	}

	outcomes := Process(context.Background(), items, 2, 4, fn)

	if len(outcomes) != len(items) {
		t.Fatalf("got %d outcomes, want %d", len(outcomes), len(items))
	}
	for i, o := range outcomes {
		if !errors.Is(o.Err, errBoom) {
			t.Errorf("outcome %d: Err=%v, want errBoom", i, o.Err)
		}
		if o.Attempts != 4 {
			t.Errorf("outcome %d: Attempts=%d, want 4", i, o.Attempts)
		}
		if o.Value != 0 {
			t.Errorf("outcome %d: Value=%d, want 0", i, o.Value)
		}
	}
}

func TestRespectsContext(t *testing.T) {
	items := make([]int, 100)
	for i := range items {
		items[i] = i
	}

	fn := func(ctx context.Context, item int) (int, error) {
		select {
		case <-time.After(20 * time.Millisecond):
			return item, nil
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()

	done := make(chan []Outcome, 1)
	go func() {
		done <- Process(ctx, items, 4, 3, fn)
	}()

	var outcomes []Outcome
	select {
	case outcomes = <-done:
	case <-time.After(500 * time.Millisecond):
		t.Fatal("Process did not return promptly after cancellation")
	}

	if len(outcomes) != 100 {
		t.Fatalf("got %d outcomes, want 100 (one per item)", len(outcomes))
	}

	allNil := true
	for _, o := range outcomes {
		if o.Err != nil {
			allNil = false
			break
		}
	}
	if allNil {
		t.Error("expected some outcomes to be cancelled, but all succeeded")
	}
}
