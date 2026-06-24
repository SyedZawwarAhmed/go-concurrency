package safecache

import (
	"sync"
	"testing"
)

func TestSafeCounterConcurrentInc(t *testing.T) {
	c := NewSafeCounter()
	const n = 1000

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			c.Inc("key")
		}()
	}
	wg.Wait()

	if got := c.Value("key"); got != n {
		t.Errorf("Value(key) = %d, want %d", got, n)
	}
}

func TestSafeCounterMultipleKeys(t *testing.T) {
	c := NewSafeCounter()
	tests := []struct {
		key  string
		incr int
	}{
		{"a", 100},
		{"b", 250},
		{"c", 50},
	}

	var wg sync.WaitGroup
	for _, tc := range tests {
		for i := 0; i < tc.incr; i++ {
			wg.Add(1)
			go func(k string) {
				defer wg.Done()
				c.Inc(k)
			}(tc.key)
		}
	}
	wg.Wait()

	for _, tc := range tests {
		if got := c.Value(tc.key); got != tc.incr {
			t.Errorf("Value(%q) = %d, want %d", tc.key, got, tc.incr)
		}
	}
}

// TestSafeCounterConcurrentReadWrite calls Value while Inc is still running.
// This exercises the concurrent map read+write path that the other tests miss
// (they only read after wg.Wait()). With the lock removed from Value, this test
// trips the race detector and may panic with "concurrent map read and map write".
func TestSafeCounterConcurrentReadWrite(t *testing.T) {
	c := NewSafeCounter()
	const n = 1000

	var wg sync.WaitGroup

	// Writers.
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			c.Inc("key")
		}()
	}

	// Readers running concurrently with the writers.
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_ = c.Value("key")
		}()
	}

	wg.Wait()

	if got := c.Value("key"); got != n {
		t.Errorf("Value(key) = %d, want %d", got, n)
	}
}

func TestSafeCounterUnknownKey(t *testing.T) {
	c := NewSafeCounter()
	if got := c.Value("never-incremented"); got != 0 {
		t.Errorf("Value(unknown) = %d, want 0", got)
	}
}
