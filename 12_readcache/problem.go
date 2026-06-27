// Package readcache — Problem 12: A read-through cache with sync.RWMutex.
//
// CONCEPT: sync.RWMutex lets unlimited readers proceed in parallel while writers
// get exclusive access — the right tool for state read far more than written.
// Paired with double-checked locking, it ensures a cache miss hits the backing
// store only once per key, even under a stampede of concurrent Gets.
//
// SCENARIO: A cache in front of a slow database (sandbox.DB adds latency to every
// call and counts how many it serves). Get returns the cached value if present;
// on a miss it loads from the DB, stores the result, and returns it.
//
// REQUIREMENTS:
//   - Safe for concurrent use. Reads take the read lock; only a miss takes the
//     write lock.
//   - A given key must be loaded from the DB at most once, no matter how many
//     goroutines request it simultaneously (re-check after acquiring the write
//     lock).
//   - A DB error is returned and NOT cached.
//   - No data races — go test -race -v ./12_readcache/
package readcache

import (
	"context"

	"github.com/SyedZawwarAhmed/go-concurrency/sandbox"
)

// Cache is a concurrency-safe read-through cache backed by db.
type Cache struct {
	// TODO: add the fields you need (the db, a sync.RWMutex, a map).
}

// New returns a ready-to-use Cache backed by db.
func New(db *sandbox.DB) *Cache {
	panic("TODO: implement New")
}

// Get returns the value for key, loading and caching it from the DB on a miss.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	panic("TODO: implement Get")
}
