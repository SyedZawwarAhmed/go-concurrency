package ratelimit

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func makeRows(n int) []Row {
	rows := make([]Row, n)
	for i := range rows {
		rows[i] = Row{ID: i, Data: fmt.Sprintf("row-%d", i)}
	}
	return rows
}

func TestBatchWriteRespectsLimit(t *testing.T) {
	const maxConcurrent = 3
	rows := makeRows(20)

	var active int64 // currently-running writers
	var peak int64   // high-water mark of active
	var total int64  // total completed writes

	write := func(r Row) error {
		cur := atomic.AddInt64(&active, 1)

		// Track the running maximum with a CAS loop.
		for {
			p := atomic.LoadInt64(&peak)
			if cur <= p || atomic.CompareAndSwapInt64(&peak, p, cur) {
				break
			}
		}
		if cur > maxConcurrent {
			t.Errorf("observed %d concurrent writers, limit is %d", cur, maxConcurrent)
		}

		time.Sleep(20 * time.Millisecond) // hold the slot so overlap is observable
		atomic.AddInt64(&total, 1)
		atomic.AddInt64(&active, -1)
		return nil
	}

	if err := BatchWrite(rows, maxConcurrent, write); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := atomic.LoadInt64(&peak); got > maxConcurrent {
		t.Errorf("peak concurrency = %d, must never exceed %d", got, maxConcurrent)
	}
	if got := atomic.LoadInt64(&peak); got != maxConcurrent {
		t.Errorf("peak concurrency = %d, want exactly %d (should fully use the budget)", got, maxConcurrent)
	}
	if got := atomic.LoadInt64(&total); got != 20 {
		t.Errorf("processed %d rows, want 20", got)
	}
}

func TestBatchWriteAllProcessedOnce(t *testing.T) {
	rows := makeRows(20)
	var counts [20]int64
	write := func(r Row) error {
		atomic.AddInt64(&counts[r.ID], 1)
		return nil
	}
	if err := BatchWrite(rows, 4, write); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for id, c := range counts {
		if c != 1 {
			t.Errorf("row %d processed %d times, want exactly 1", id, c)
		}
	}
}

func TestBatchWritePropagatesError(t *testing.T) {
	rows := makeRows(10)
	boom := errors.New("write failed")
	write := func(r Row) error {
		if r.ID == 7 {
			return boom
		}
		return nil
	}
	if err := BatchWrite(rows, 3, write); err == nil {
		t.Fatal("expected a non-nil error when a write fails, got nil")
	}
}
