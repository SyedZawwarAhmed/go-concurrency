package httpfetch

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/SyedZawwarAhmed/go-concurrency/sandbox"
)

func urls(base string, n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = fmt.Sprintf("%s/item/%d", base, i)
	}
	return out
}

func TestFetchAllBoundsConcurrency(t *testing.T) {
	srv := sandbox.NewServer(20 * time.Millisecond)
	defer srv.Close()

	const maxConcurrent = 5
	us := urls(srv.URL, 50)

	results := FetchAll(context.Background(), srv.Client(), us, maxConcurrent)

	if len(results) != len(us) {
		t.Fatalf("got %d results, want %d", len(results), len(us))
	}
	for _, r := range results {
		if r.Err != nil {
			t.Errorf("%s: unexpected error %v", r.URL, r.Err)
			continue
		}
		if r.Status != 200 {
			t.Errorf("%s: status %d, want 200", r.URL, r.Status)
		}
	}
	if peak := srv.MaxConcurrent(); peak > maxConcurrent {
		t.Errorf("server saw %d concurrent requests, want <= %d", peak, maxConcurrent)
	}
}

func TestFetchAllRespectsContextDeadline(t *testing.T) {
	srv := sandbox.NewServer(100 * time.Millisecond)
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()

	us := urls(srv.URL, 10)
	start := time.Now()
	results := FetchAll(ctx, srv.Client(), us, 10)
	elapsed := time.Since(start)

	if elapsed > 90*time.Millisecond {
		t.Errorf("FetchAll took %v; should have bailed out at the ~40ms deadline", elapsed)
	}
	if len(results) != len(us) {
		t.Fatalf("got %d results, want %d", len(results), len(us))
	}
	for _, r := range results {
		if r.Err == nil {
			t.Errorf("%s: err = nil, want a deadline error (server is slower than ctx)", r.URL)
		}
	}
}
