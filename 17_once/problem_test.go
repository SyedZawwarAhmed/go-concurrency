package lazyinit

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestInitRunsExactlyOnce(t *testing.T) {
	var counter int64
	lazy := NewLazy(func() int {
		atomic.AddInt64(&counter, 1)
		return 42
	})

	const n = 200
	results := make([]int, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			results[i] = lazy.Get()
		}(i)
	}
	wg.Wait()

	if got := atomic.LoadInt64(&counter); got != 1 {
		t.Fatalf("init ran %d times, want 1", got)
	}
	for i, v := range results {
		if v != 42 {
			t.Fatalf("results[%d] = %d, want 42", i, v)
		}
	}
}

func TestReturnsInitValue(t *testing.T) {
	var counter int64
	lazy := NewLazy(func() string {
		atomic.AddInt64(&counter, 1)
		return "ready"
	})

	if got := lazy.Get(); got != "ready" {
		t.Fatalf("first Get = %q, want %q", got, "ready")
	}
	if got := lazy.Get(); got != "ready" {
		t.Fatalf("second Get = %q, want %q", got, "ready")
	}
	if got := atomic.LoadInt64(&counter); got != 1 {
		t.Fatalf("init ran %d times, want 1", got)
	}
}

func TestGetIsConcurrencySafe(t *testing.T) {
	var counter int64
	lazy := NewLazy(func() []int {
		atomic.AddInt64(&counter, 1)
		return []int{1, 2, 3}
	})

	const n = 100
	results := make([][]int, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			results[i] = lazy.Get()
		}(i)
	}
	wg.Wait()

	if got := atomic.LoadInt64(&counter); got != 1 {
		t.Fatalf("init ran %d times, want 1", got)
	}

	want := &results[0][0]
	for i, s := range results {
		if len(s) != 3 {
			t.Fatalf("results[%d] len = %d, want 3", i, len(s))
		}
		if &s[0] != want {
			t.Fatalf("results[%d] is a different underlying slice", i)
		}
	}
}
