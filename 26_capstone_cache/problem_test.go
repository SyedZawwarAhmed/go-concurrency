package ttlcache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/SyedZawwarAhmed/go-concurrency/sandbox"
)

func TestLoadsOncePerKeyUnderStampede(t *testing.T) {
	db := sandbox.NewDB(sandbox.DBOptions{
		Latency: 50 * time.Millisecond,
		Seed:    map[string]string{"u1": "alice"},
	})
	c := New(db, 5*time.Second, 1*time.Second)
	defer c.Close()

	const n = 50
	got := make([]string, n)
	var errCount atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			v, err := c.Get(context.Background(), "u1")
			if err != nil {
				errCount.Add(1)
				return
			}
			got[i] = v
		}(i)
	}
	wg.Wait()

	if errCount.Load() != 0 {
		t.Fatalf("got %d errors, want 0", errCount.Load())
	}
	for i, v := range got {
		if v != "alice" {
			t.Errorf("got[%d] = %q, want alice", i, v)
		}
	}
	if q := db.QueryCount(); q != 1 {
		t.Fatalf("DB served %d queries, want 1 (singleflight should collapse the stampede)", q)
	}

	// Subsequent reads come from the cache; no new queries.
	for i := 0; i < 5; i++ {
		v, err := c.Get(context.Background(), "u1")
		if err != nil || v != "alice" {
			t.Fatalf("cached Get = %q, %v; want alice, nil", v, err)
		}
	}
	if q := db.QueryCount(); q != 1 {
		t.Fatalf("DB served %d queries after cached reads, want 1", q)
	}
}

func TestExpiredEntryReloads(t *testing.T) {
	db := sandbox.NewDB(sandbox.DBOptions{
		Latency: 5 * time.Millisecond,
		Seed:    map[string]string{"k": "v"},
	})
	c := New(db, 40*time.Millisecond, 10*time.Millisecond)
	defer c.Close()

	v, err := c.Get(context.Background(), "k")
	if err != nil || v != "v" {
		t.Fatalf("first Get = %q, %v; want v, nil", v, err)
	}
	if q := db.QueryCount(); q != 1 {
		t.Fatalf("after first Get QueryCount = %d, want 1", q)
	}

	// Entry expires (ttl 40ms) and the janitor (every 10ms) likely evicts it.
	time.Sleep(120 * time.Millisecond)

	v, err = c.Get(context.Background(), "k")
	if err != nil || v != "v" {
		t.Fatalf("reload Get = %q, %v; want v, nil", v, err)
	}
	if q := db.QueryCount(); q != 2 {
		t.Fatalf("after reload QueryCount = %d, want 2 (expired entry should reload)", q)
	}
}

func TestMissReturnsErrorNotCached(t *testing.T) {
	db := sandbox.NewDB(sandbox.DBOptions{
		Seed: map[string]string{},
	})
	c := New(db, 1*time.Second, 1*time.Second)
	defer c.Close()

	if _, err := c.Get(context.Background(), "absent"); !errors.Is(err, sandbox.ErrNotFound) {
		t.Fatalf("first Get error = %v, want ErrNotFound", err)
	}
	if _, err := c.Get(context.Background(), "absent"); !errors.Is(err, sandbox.ErrNotFound) {
		t.Fatalf("second Get error = %v, want ErrNotFound", err)
	}
	if q := db.QueryCount(); q != 2 {
		t.Fatalf("QueryCount = %d, want 2 (errors must not be cached)", q)
	}
}

func TestConcurrentDifferentKeysRaceClean(t *testing.T) {
	seed := map[string]string{
		"a": "A", "b": "B", "c": "C", "d": "D", "e": "E",
	}
	keys := []string{"a", "b", "c", "d", "e"}
	want := map[string]string{"a": "A", "b": "B", "c": "C", "d": "D", "e": "E"}

	db := sandbox.NewDB(sandbox.DBOptions{Seed: seed})
	c := New(db, 1*time.Second, 200*time.Millisecond)
	defer c.Close()

	const n = 100
	var errCount atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := keys[i%5]
			v, err := c.Get(context.Background(), k)
			if err != nil {
				errCount.Add(1)
				return
			}
			if v != want[k] {
				t.Errorf("Get(%q) = %q, want %q", k, v, want[k])
			}
		}(i)
	}
	wg.Wait()

	if errCount.Load() != 0 {
		t.Fatalf("got %d errors, want 0", errCount.Load())
	}
}
