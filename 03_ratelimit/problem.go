// Package ratelimit — Problem 03: Rate limiting with a buffered-channel semaphore.
//
// CONCEPT: Bounding the number of concurrently-active goroutines using a
// buffered channel as a counting semaphore.
//
// SCENARIO: A batch DB writer. Given a slice of rows, write them all
// concurrently — but never allow more than maxConcurrent writes to be in flight
// at the same instant.
//
// The write function is INJECTED by the caller (this is how the test measures
// concurrency without seeing your internals). Your job is purely to bound the
// concurrency and collect any error.
//
// REQUIREMENTS:
//   - Process every row, each in a goroutine, calling write(row).
//   - At no instant may more than maxConcurrent calls to write be running.
//   - Use a buffered channel of capacity maxConcurrent as a semaphore: acquire a
//     slot before calling write, release it after.
//   - Wait for all rows with a sync.WaitGroup before returning.
//   - Return the FIRST non-nil error from any write (or nil if all succeed).
//   - No data races — verify with: go test -race -v ./03_ratelimit/
package ratelimit

// Row is a single unit of data to be written.
type Row struct {
	ID   int
	Data string
}

// BatchWrite writes all rows concurrently with at most maxConcurrent in flight.
//
// TODO: Implement the semaphore pattern.
//   - Make a buffered channel `sem` of capacity maxConcurrent.
//   - For each row: acquire a slot (sem <- struct{}{}), then launch a goroutine
//     that defers releasing the slot (<-sem), calls write(row), and records any error.
//   - Use a sync.WaitGroup to wait for every goroutine.
//   - Capture the first error safely (guard the shared error variable with a sync.Mutex).
//   - Return that error (or nil).
func BatchWrite(rows []Row, maxConcurrent int, write func(Row) error) error {
	panic("TODO: implement BatchWrite")
}
