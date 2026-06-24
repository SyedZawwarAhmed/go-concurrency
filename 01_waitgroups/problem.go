// Package waitgroups — Problem 01: Basic WaitGroups.
//
// CONCEPT: Launching goroutines and waiting for them all to finish using a
// sync.WaitGroup.
//
// SCENARIO: A parallel string massager. Given a slice of strings, process each
// string concurrently — each unit of work simulates a 50ms delay — and return
// the processed results IN THE SAME ORDER as the input.
//
// The transformation itself (massage) is already written for you: it trims
// surrounding whitespace and upper-cases the string. Your ONLY job is the
// concurrency: run all the work in parallel and wait for completion.
//
// REQUIREMENTS:
//   - Each input is processed in its own goroutine.
//   - Results are returned in input order (inputs[i] -> output[i]).
//   - Use a sync.WaitGroup to wait for every goroutine before returning.
//   - No data races (have each goroutine write to its own slice slot).
//   - nil / empty input returns a length-0 slice and must not hang.
//
// Run the tests:  go test -race -v ./01_waitgroups/
package waitgroups

import (
	"strings"
	"time"
)

// processDelay is the simulated per-item work duration. Do not change it — the
// timing test relies on this value.
const processDelay = 50 * time.Millisecond

// massage is the (already implemented) transformation applied to each string.
// It sleeps for processDelay to simulate real work, then normalizes the string.
func massage(s string) string {
	time.Sleep(processDelay)
	return strings.ToUpper(strings.TrimSpace(s))
}

// MassageStrings processes every input concurrently and returns the results in
// the same order as inputs.
//
// TODO: Implement this using goroutines + sync.WaitGroup.
//   - Allocate a result slice of len(inputs).
//   - For each index i, launch a goroutine that sets results[i] = massage(inputs[i]).
//   - Use a sync.WaitGroup to wait until all goroutines complete.
//   - Return the results.
func MassageStrings(inputs []string) []string {
	panic("TODO: implement MassageStrings")
}
