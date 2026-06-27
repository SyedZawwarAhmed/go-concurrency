package hotconfig

import (
	"sync"
	"testing"
	"time"
)

func TestLoadReturnsStored(t *testing.T) {
	s := NewStore(&Config{MaxConns: 10, Timeout: time.Second})
	if got := s.Load(); got.MaxConns != 10 {
		t.Fatalf("Load().MaxConns = %d, want 10", got.MaxConns)
	}

	s.Store(&Config{MaxConns: 20, Timeout: 2 * time.Second})
	if got := s.Load(); got.MaxConns != 20 {
		t.Fatalf("after Store, Load().MaxConns = %d, want 20", got.MaxConns)
	}
}

func TestConcurrentReadsDuringSwaps(t *testing.T) {
	values := []int{10, 20, 30}
	valid := map[int]bool{10: true, 20: true, 30: true}

	s := NewStore(&Config{MaxConns: values[0], Timeout: time.Second})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			s.Store(&Config{MaxConns: values[i%len(values)], Timeout: time.Second})
		}
	}()

	const readers = 50
	wg.Add(readers)
	for r := 0; r < readers; r++ {
		go func() {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				c := s.Load()
				if c == nil {
					t.Errorf("Load returned nil")
					return
				}
				if !valid[c.MaxConns] {
					t.Errorf("Load returned invalid MaxConns %d", c.MaxConns)
					return
				}
			}
		}()
	}

	wg.Wait()
}

func TestLoadNeverNil(t *testing.T) {
	s := NewStore(&Config{MaxConns: 10, Timeout: time.Second})
	if s.Load() == nil {
		t.Fatalf("Load() returned nil after construction")
	}

	for i := 0; i < 5; i++ {
		s.Store(&Config{MaxConns: 10 + i, Timeout: time.Second})
		if s.Load() == nil {
			t.Fatalf("Load() returned nil after Store #%d", i)
		}
	}
}
