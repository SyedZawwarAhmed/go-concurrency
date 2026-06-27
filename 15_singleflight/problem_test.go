package singleflight

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/SyedZawwarAhmed/go-concurrency/sandbox"
)

func TestDoCollapsesDuplicates(t *testing.T) {
	db := sandbox.NewDB(sandbox.DBOptions{
		Latency: 50 * time.Millisecond,
		Seed:    map[string]string{"u1": "alice"},
	})
	var g Group[string]

	const n = 50
	got := make([]string, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v, err := g.Do("u1", func() (string, error) {
				return db.Get(context.Background(), "u1")
			})
			if err != nil {
				t.Errorf("Do: %v", err)
				return
			}
			got[i] = v
		}(i)
	}
	wg.Wait()

	if q := db.QueryCount(); q != 1 {
		t.Errorf("DB served %d queries, want 1 (concurrent calls should collapse)", q)
	}
	for i, v := range got {
		if v != "alice" {
			t.Errorf("got[%d] = %q, want alice", i, v)
		}
	}
}

func TestDoDistinctKeysRunConcurrently(t *testing.T) {
	var g Group[int]
	const n = 10
	start := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			g.Do(fmt.Sprintf("key%d", i), func() (int, error) {
				time.Sleep(100 * time.Millisecond)
				return i, nil
			})
		}(i)
	}
	wg.Wait()
	if elapsed := time.Since(start); elapsed > 300*time.Millisecond {
		t.Errorf("took %v; distinct keys should run concurrently, not serialize", elapsed)
	}
}

func TestDoForgetsCompletedCalls(t *testing.T) {
	var g Group[int]
	var calls int
	fn := func() (int, error) {
		calls++
		return 42, nil
	}

	if v, _ := g.Do("k", fn); v != 42 {
		t.Fatalf("first Do = %d, want 42", v)
	}
	if v, _ := g.Do("k", fn); v != 42 {
		t.Fatalf("second Do = %d, want 42", v)
	}
	if calls != 2 {
		t.Errorf("fn ran %d times, want 2 (completed calls must not be cached)", calls)
	}
}
