// Package hotconfig — Problem 18: Lock-free hot-swap with atomic.Pointer.
//
// CONCEPT: Some state is read constantly on a hot path but changes rarely, and
// when it changes it is replaced wholesale rather than mutated in place. An
// atomic.Pointer lets every reader load the current value with no lock and no
// contention, while a writer swaps in a brand-new value atomically.
//
// SCENARIO: A Store holds the current *Config. Readers call Load on a hot path;
// a writer occasionally calls Store with a fresh, fully-built Config. The config
// is immutable once stored, so readers never see a half-updated value.
//
// REQUIREMENTS:
//   - Load takes no lock and never returns nil after construction.
//   - Store atomically replaces the current config; in-flight readers see either
//     the old or the new config, never a torn one.
//   - Safe under many concurrent readers racing a writer.
//   - No data races — go test -race -v ./18_atomicvalue/
package hotconfig

import (
	"sync/atomic"
	"time"
)

// Config is read on a hot path by many goroutines and occasionally replaced
// wholesale. Treat it as immutable once stored.
type Config struct {
	MaxConns int
	Timeout  time.Duration
}

// Store holds the current *Config and allows lock-free reads.
type Store struct {
	// TODO: hold the current config, e.g. an atomic.Pointer[Config].
	_ atomic.Pointer[Config]
}

// NewStore seeds the store with an initial config (non-nil).
//
// HINT: allocate a Store and atomically store the initial pointer into it.
func NewStore(initial *Config) *Store {
	panic("TODO: implement NewStore")
}

// Load returns the current config without locking.
func (s *Store) Load() *Config {
	panic("TODO: implement Load")
}

// Store atomically replaces the current config.
func (s *Store) Store(c *Config) {
	panic("TODO: implement Store")
}
