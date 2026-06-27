// Package retrypool — Problem 25: A resilient worker pool with retries and per-item errors.
//
// CONCEPT: Real work fails transiently. A robust pool retries a failing item a
// bounded number of times, records the outcome per item, and stops promptly when
// the caller cancels — all while writing each result without locks or races.
//
// SCENARIO: Process a batch of items with a fixed pool of `workers` goroutines.
// Each item is handed to fn; on error it is retried (up to maxAttempts total
// calls) before its error is recorded. A cancelled context halts work quickly.
//
// REQUIREMENTS:
//   - Run fn over every item using exactly `workers` goroutines (a worker pool).
//   - Retry a failing item immediately (no backoff) up to maxAttempts total calls.
//   - Record one Outcome per item, in the SAME ORDER as items.
//   - Honor ctx cancellation promptly; cancelled items carry a context error.
//   - No data races — go test -race -v ./25_retrypool/
package retrypool

import "context"

// Outcome is the per-item result.
type Outcome struct {
	Item     int   // the input value
	Value    int   // fn's successful result (zero if it never succeeded)
	Err      error // nil on success, else the last error after exhausting retries
	Attempts int   // how many times fn was called for this item (1..maxAttempts)
}

// Process runs fn over every item using a pool of `workers` goroutines. If fn
// returns an error, the item is retried (immediately, no backoff needed) up to
// maxAttempts total calls before its error is recorded. ctx cancellation stops
// work promptly. Returns one Outcome per item, in the SAME ORDER as items.
//
// HINT: feed indices through a channel; each index is owned by one worker, so
// writing results[i] needs no lock; check ctx.Err() before each attempt.
func Process(ctx context.Context, items []int, workers, maxAttempts int, fn func(context.Context, int) (int, error)) []Outcome {
	panic("TODO: implement Process")
}
