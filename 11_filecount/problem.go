// Package filecount — Problem 11: Lock-free aggregation with sync/atomic.
//
// CONCEPT: When many goroutines only ever ADD to a few shared counters, a mutex
// is overkill. sync/atomic gives race-free read-modify-write without blocking.
//
// SCENARIO: A log-directory analyzer. Given the paths of real files on disk,
// read them concurrently with a fixed pool of `workers` goroutines and report
// the aggregate file / line / byte counts. (A line is a '\n'.)
//
// REQUIREMENTS:
//   - Read files concurrently using exactly `workers` goroutines (a worker pool,
//     not one goroutine per file).
//   - Accumulate the totals across goroutines WITHOUT a mutex — use sync/atomic.
//   - Return the FIRST read error encountered (e.g. a missing file); files that
//     did read successfully may still be counted.
//   - Empty input returns a zero Counts and nil, and must not hang.
//   - No data races — go test -race -v ./11_filecount/
package filecount

// Counts is the aggregate result across all files.
type Counts struct {
	Files int64
	Lines int64
	Bytes int64
}

// CountFiles reads every path using a pool of `workers` goroutines and returns
// the aggregate counts plus the first read error (nil if all succeeded).
//
// HINT: feed paths through a channel to the pool; atomically add into Counts.
func CountFiles(paths []string, workers int) (Counts, error) {
	panic("TODO: implement CountFiles")
}
