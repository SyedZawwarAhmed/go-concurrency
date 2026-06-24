// Package safecache — Problem 05: Protecting shared state with sync.Mutex.
//
// CONCEPT: When many goroutines mutate shared state, you get a data race — and
// for a Go map, concurrent writes are a hard runtime panic. A sync.Mutex
// serializes access so only one goroutine touches the state at a time.
//
// SCENARIO: A concurrency-safe counter keyed by string. Many goroutines call
// Inc(key) simultaneously; Value(key) reports the current count.
//
// REQUIREMENTS:
//   - SafeCounter must be safe for concurrent use by multiple goroutines.
//   - Inc and Value must lock/unlock around every access to the underlying map.
//   - Must pass the race detector: go test -race -v ./05_mutex/
package safecache

import "sync"

// SafeCounter is a concurrency-safe, string-keyed counter.
//
// The mutex field is provided for you — use it to guard every access to counts.
type SafeCounter struct {
	mu     sync.Mutex
	counts map[string]int
}

// NewSafeCounter returns a ready-to-use SafeCounter.
func NewSafeCounter() *SafeCounter {
	return &SafeCounter{counts: make(map[string]int)}
}

// Inc increments the counter for key by one.
//
// TODO: lock c.mu, increment c.counts[key], unlock (defer is idiomatic).
func (c *SafeCounter) Inc(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[key]++
}

// Value returns the current count for key.
//
// TODO: lock c.mu, read c.counts[key], unlock, and return the value.
func (c *SafeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.counts[key]
}
