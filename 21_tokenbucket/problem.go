// Package tokenbucket — Problem 21: A token-bucket rate limiter with time.Ticker.
//
// CONCEPT: A semaphore caps CONCURRENCY (how many run at once); a token bucket
// caps RATE over TIME (how many start per interval). Tokens accrue on a
// time.Ticker up to a `burst` ceiling, so callers may spend a saved-up burst
// quickly but average out to one token per refill interval.
//
// SCENARIO: An outbound API client must respect a "N requests per interval"
// quota. Each request asks Allow(); a true means a token was spent and the call
// may proceed, a false means the quota is exhausted right now and the caller
// should skip or retry later. Allow never blocks.
//
// REQUIREMENTS:
//   - The bucket starts full with `burst` tokens; Allow consumes one and reports
//     whether it did, never blocking.
//   - A background goroutine refills one token every `refill` interval, never
//     exceeding `burst` (extra ticks are dropped).
//   - Stop halts the refill goroutine; it is called exactly once per Limiter.
//   - No data races — go test -race -v ./21_tokenbucket/
package tokenbucket

import "time"

// Limiter is a token bucket: it starts full with `burst` tokens and refills
// one token every `refill` interval, never exceeding burst. Allow consumes a
// token if one is available and reports whether it did — it never blocks.
type Limiter struct{ /* student chooses fields, e.g. chan struct{} + *time.Ticker + stop chan */ }

// NewLimiter starts the refill loop. The bucket starts full (burst tokens).
//
// HINT: a buffered channel of cap `burst`, pre-filled with `burst` tokens, is
// the bucket; a goroutine on time.Ticker adds one token per tick.
func NewLimiter(refill time.Duration, burst int) *Limiter {
	panic("TODO: implement NewLimiter")
}

// Allow consumes one token and returns true, or returns false immediately if
// the bucket is empty. Never blocks.
//
// HINT: select with a default case turns a channel receive non-blocking.
func (l *Limiter) Allow() bool {
	panic("TODO: implement Allow")
}

// Stop halts the background refill goroutine. Call exactly once.
//
// HINT: close a stop channel the refill goroutine selects on, then ticker.Stop.
func (l *Limiter) Stop() {
	panic("TODO: implement Stop")
}
