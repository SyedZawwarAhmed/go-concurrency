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

func TestSafeCounterUnknownKey(t *testing.T) {
	c := NewSafeCounter()
	if got := c.Value("never-incremented"); got != 0 {
		t.Errorf("Value(unknown) = %d, want 0", got)
	}
}
