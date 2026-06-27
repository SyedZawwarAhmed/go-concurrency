// Package ttlcache — Problem 26 (capstone): A concurrent TTL cache with singleflight + background eviction.
//
// CONCEPT: This capstone integrates four earlier patterns into one component:
// the RWMutex read-through cache (problem 12), singleflight call collapsing
// (problem 15), a time.Ticker janitor that sweeps on an interval (problem 20),
// and a clean Stop/lifecycle that shuts a background goroutine down (problem 24).
// Together they make a cache that is fast to read, loads each key exactly once
// under a stampede, expires stale data, and shuts down without leaking goroutines.
//
// SCENARIO: A read-through cache in front of a slow database (sandbox.DB adds
// latency and counts the calls it serves). Many goroutines read the same hot
// keys; entries go stale after a TTL and must be reloaded; a janitor reclaims
// expired entries in the background; Close tears the janitor down.
//
// REQUIREMENTS:
//   - Get returns a fresh cached value without touching the DB; on a miss or an
//     expired entry it loads from the DB and caches the result with the TTL.
//   - Under a stampede of concurrent Gets for the same key, the DB is hit EXACTLY
//     once (singleflight); every caller receives that one result.
//   - A DB error is returned to the caller and NOT cached.
//   - A background janitor evicts expired entries every cleanupInterval.
//   - Close stops the janitor goroutine; safe to call once.
//   - No data races — go test -race -v ./26_capstone_cache/
package ttlcache

import (
	"context"
	"time"

	"github.com/SyedZawwarAhmed/go-concurrency/sandbox"
)

// Cache is a concurrent, read-through cache in front of db with per-entry TTL.
// A background janitor evicts expired entries every cleanupInterval. On a miss
// (or expired entry) Get loads the value from db EXACTLY ONCE even under a
// stampede of concurrent callers (singleflight), caches it with the TTL, and
// returns it. Close stops the janitor goroutine.
type Cache struct {
	// TODO: add the fields you need (the db + ttl, a sync.RWMutex guarding a
	// map[string]entry, an inline singleflight (mutex + map[string]*call), and a
	// stop channel + sync.Once for the janitor's lifecycle).
}

// New starts the cache (and its janitor). ttl is the freshness window per
// entry; cleanupInterval is how often the janitor sweeps for expired entries.
//
// HINT: build the Cache, then `go` a janitor loop that selects on a stop channel
// and a time.NewTicker(cleanupInterval).C, deleting entries whose TTL has passed.
func New(db *sandbox.DB, ttl time.Duration, cleanupInterval time.Duration) *Cache {
	panic("TODO: implement New")
}

// Get returns the value for key: from cache if fresh, otherwise loaded from db
// (once per concurrent miss) and cached. DB errors are returned, not cached.
//
// HINT: RLock and return a fresh entry; else singleflight on key — register a
// *call, double-check the cache, db.Get on a real miss, store entry{val, now+ttl}
// on success (never on error), signal waiters, and forget the call.
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	panic("TODO: implement Get")
}

// Close stops the background janitor. Safe to call once.
//
// HINT: close the stop channel inside a sync.Once so repeat calls don't panic.
func (c *Cache) Close() {
	panic("TODO: implement Close")
}
