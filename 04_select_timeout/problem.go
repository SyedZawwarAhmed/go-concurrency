// Package selecttimeout — Problem 04: Timeouts & multiplexing with select.
//
// CONCEPT: Using a select statement to race multiple channels against a timeout
// (time.After) — the "Promise.race" pattern, plus early-out on timeout.
//
// SCENARIO: A resilient API fetcher. You're given two receive-only channels — a
// primary and a replica — each of which will (eventually) deliver one result.
// Return whichever value arrives FIRST. If neither delivers within timeout,
// return ErrTimeout.
//
// REQUIREMENTS:
//   - Use a SINGLE select statement with three cases:
//       a value from primary, a value from replica, and <-time.After(timeout).
//   - Return the first value received, with a nil error.
//   - On timeout, return "" and ErrTimeout.
//
// Run the tests:  go test -race -v ./04_select_timeout/
package selecttimeout

import (
	"errors"
	"time"
)

// ErrTimeout is returned when neither server responds within the timeout.
var ErrTimeout = errors.New("selecttimeout: both servers timed out")

// FetchResilient returns the first value from primary or replica, or ErrTimeout
// if neither arrives within timeout.
//
// TODO: Implement with a single select over primary, replica, and time.After(timeout).
func FetchResilient(primary, replica <-chan string, timeout time.Duration) (string, error) {
	panic("TODO: implement FetchResilient")
}
