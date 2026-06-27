// Package httpfetch — Problem 13: Bounded concurrent HTTP with context.
//
// CONCEPT: Fan out real network requests, but cap how many run at once (a
// buffered-channel semaphore) and thread a context through every request so a
// cancellation or deadline tears the whole batch down promptly.
//
// SCENARIO: A batch URL fetcher hitting a real HTTP server. Fetch every URL
// concurrently, at most maxConcurrent in flight, using the given client. Each
// request must be bound to ctx (use http.NewRequestWithContext).
//
// REQUIREMENTS:
//   - One Result per input URL (order does not matter).
//   - Never more than maxConcurrent requests in flight at the same instant.
//   - A per-URL failure (transport error, cancelled context) goes in Result.Err
//     for that URL — it does NOT abort the others and is not returned separately.
//   - Respect ctx: when it's done, in-flight and pending requests stop promptly.
//   - No data races — go test -race -v ./13_httpfetch/
package httpfetch

import (
	"context"
	"net/http"
)

// Result is the outcome of fetching one URL.
type Result struct {
	URL    string
	Status int
	Body   string
	Err    error
}

// FetchAll GETs every url concurrently with at most maxConcurrent in flight,
// returning one Result per url.
//
// HINT: a buffered channel of size maxConcurrent bounds concurrency; have each
// goroutine write to its own slot so no mutex is needed.
func FetchAll(ctx context.Context, client *http.Client, urls []string, maxConcurrent int) []Result {
	panic("TODO: implement FetchAll")
}
