// Package batcher — Problem 20: Time- and size-triggered batching with time.Ticker.
//
// CONCEPT: A batcher coalesces a stream of individual items into batches,
// flushing on whichever trigger fires first: the buffer filling to maxBatch, or
// a periodic time.Ticker tick. One goroutine owns the buffer, so flush never
// runs concurrently with itself — Add just hands items over a channel.
//
// SCENARIO: Items arrive from many goroutines (Add) and must be written out in
// batches — for throughput you want full batches, but for latency you cap how
// long an item waits before being flushed.
//
// REQUIREMENTS:
//   - flush is called with a non-empty slice, always from one goroutine.
//   - A batch flushes when it reaches maxBatch items OR every flushInterval.
//   - Add is safe for concurrent use and must not block indefinitely.
//   - Stop flushes any remaining buffered items, stops the goroutine and ticker,
//     and returns only after the final flush completes; safe to call once.
//   - No data races — go test -race -v ./20_batcher/
package batcher

import "time"

// Batcher buffers items and flushes them to `flush` in batches: whenever the
// buffer reaches maxBatch items, OR every flushInterval, whichever comes
// first. flush is always called from a single goroutine (never concurrently).
// Add is safe for concurrent use. Stop flushes any remaining items and shuts
// the internal goroutine down.
type Batcher struct{ /* student chooses fields */ }

// NewBatcher starts the batcher. flush is called with a non-empty slice; the
// slice is owned by flush (the batcher will not mutate it afterward).
//
// HINT: start a goroutine looping over select on an `in` channel, ticker.C, and
// a `done` channel; keep a buf []int and flush a copy on the size/time/stop triggers.
func NewBatcher(maxBatch int, flushInterval time.Duration, flush func([]int)) *Batcher {
	panic("TODO: implement NewBatcher")
}

// Add enqueues v for the next batch (must not block indefinitely).
func (b *Batcher) Add(v int) {
	panic("TODO: implement Add")
}

// Stop flushes remaining buffered items, stops the goroutine and ticker, and
// returns after the final flush completes. Safe to call once.
//
// HINT: close `done` via sync.Once, then wait on a `stopped` channel the goroutine closes.
func (b *Batcher) Stop() {
	panic("TODO: implement Stop")
}
