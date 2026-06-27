package readcache

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/SyedZawwarAhmed/go-concurrency/sandbox"
)

func TestCacheLoadsEachKeyOnce(t *testing.T) {
	seed := map[string]string{"a": "1", "b": "2", "c": "3"}
	db := sandbox.NewDB(sandbox.DBOptions{Latency: 20 * time.Millisecond, Seed: seed})
	c := New(db)

	keys := []string{"a", "b", "c"}
	var wg sync.WaitGroup
	for i := 0; i < 150; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := keys[i%len(keys)]
			v, err := c.Get(context.Background(), k)
			if err != nil {
				t.Errorf("Get(%q): %v", k, err)
				return
			}
			if v != seed[k] {
				t.Errorf("Get(%q) = %q, want %q", k, v, seed[k])
			}
		}(i)
	}
	wg.Wait()

	if got := db.QueryCount(); got != int64(len(keys)) {
		t.Errorf("DB served %d queries, want %d (each key should load exactly once)", got, len(keys))
	}
}

func TestCacheDoesNotCacheErrors(t *testing.T) {
	db := sandbox.NewDB(sandbox.DBOptions{Seed: map[string]string{}})
	c := New(db)

	if _, err := c.Get(context.Background(), "missing"); err == nil {
		t.Fatal("Get(missing) = nil error, want ErrNotFound")
	}
	// A second miss must hit the DB again — errors are not cached.
	if _, err := c.Get(context.Background(), "missing"); err == nil {
		t.Fatal("second Get(missing) = nil error, want ErrNotFound")
	}
	if got := db.QueryCount(); got != 2 {
		t.Errorf("DB served %d queries, want 2 (errors must not be cached)", got)
	}
}
