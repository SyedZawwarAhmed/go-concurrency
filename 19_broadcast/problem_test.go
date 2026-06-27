package broadcast

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestOpenReleasesAllWaiters(t *testing.T) {
	g := NewGate()

	const n = 50
	var released atomic.Int64
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			g.Wait()
			released.Add(1)
		}()
	}

	// Give the goroutines time to park on Wait; none should be released yet.
	time.Sleep(20 * time.Millisecond)
	if got := released.Load(); got != 0 {
		t.Fatalf("before Open: %d waiters released, want 0 (they must block)", got)
	}

	g.Open()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting for waiters after Open; %d/%d released", released.Load(), n)
	}

	if got := released.Load(); got != n {
		t.Fatalf("after Open: %d waiters released, want %d", got, n)
	}
}

func TestWaitAfterOpenReturnsImmediately(t *testing.T) {
	g := NewGate()
	g.Open()

	done := make(chan struct{})
	go func() {
		g.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait after Open did not return promptly")
	}
}

func TestOpenIsIdempotent(t *testing.T) {
	g := NewGate()

	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			g.Open() // a double close would panic
		}()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for concurrent Open calls")
	}

	waited := make(chan struct{})
	go func() {
		g.Wait()
		close(waited)
	}()
	select {
	case <-waited:
	case <-time.After(2 * time.Second):
		t.Fatal("Wait did not return after Open")
	}
}
