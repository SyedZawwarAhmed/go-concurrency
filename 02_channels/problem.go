// Package channels — Problem 02: Data collection via unbuffered channels.
//
// CONCEPT: Safely gathering data from concurrent workers using an UNBUFFERED
// channel, closing the channel exactly once, and ranging over it to collect
// results — with no shared-memory data race.
//
// SCENARIO: A concurrent URL status checker. Given a slice of URLs, fetch each
// one concurrently (fetchStatus simulates a ~30ms HTTP call) and return a map of
// url -> HTTP status code.
//
// REQUIREMENTS:
//   - One goroutine per URL; each sends its Result on a shared UNBUFFERED channel.
//   - A sync.WaitGroup tracks the workers; a SEPARATE goroutine closes the channel
//     once all workers are done. Close exactly once — never from inside a worker.
//   - The caller ranges over the channel to build and return the map.
//   - Empty / nil input returns an empty (non-nil) map and MUST NOT deadlock.
//   - No data races — verify with: go test -race -v ./02_channels/
package channels

import (
	"sync"
	"time"
)

// Result is what each worker reports back over the channel.
type Result struct {
	URL    string
	Status int
}

// fetchStatus simulates an HTTP request. Already implemented — do not change.
func fetchStatus(url string) int {
	time.Sleep(30 * time.Millisecond)
	return 200
}

// CheckURLs fetches every URL concurrently and returns a map of url -> status.
//
// TODO: Implement using an UNBUFFERED channel.
//   - Make an unbuffered channel of Result.
//   - Launch one goroutine per URL that sends Result{url, fetchStatus(url)}.
//   - Track workers with a sync.WaitGroup; in a SEPARATE goroutine, wg.Wait()
//     then close(ch).
//   - range over ch to populate and return the map.
//   - Make sure the empty-input case returns an empty map without blocking.
func CheckURLs(urls []string) map[string]int {
	var wg sync.WaitGroup

	ch := make(chan Result)

	for _, url := range urls {
		wg.Go(func() {
			result := Result{
				URL:    url,
				Status: fetchStatus(url),
			}
			ch <- result

		})
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	result := make(map[string]int, len(urls))
	for res := range ch {
		result[res.URL] = res.Status
	}

	return result
}
