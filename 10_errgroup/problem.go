// Package errgroup — Problem 10: First-error-cancels-the-rest.
//
// CONCEPT: Run a set of tasks concurrently. If any task fails, cancel the others
// (via context) and return the FIRST error. If all succeed, return nil. This is
// the pattern that golang.org/x/sync/errgroup provides — here you build it by
// hand from context + sync.Once + sync.WaitGroup, no external dependencies.
//
// SCENARIO: RunAll launches every task in its own goroutine, passing each a
// DERIVED context that is cancelled as soon as the first task returns an error.
// Well-behaved tasks watch ctx.Done() and bail out early.
//
// REQUIREMENTS:
//   - Derive a cancellable child context from ctx (context.WithCancel).
//   - Run each task(childCtx) in its own goroutine.
//   - Capture the FIRST non-nil error exactly once (sync.Once) and cancel the
//     child context when it happens.
//   - Wait for all goroutines (sync.WaitGroup) before returning.
//   - Always release the context (call cancel) to avoid a context leak.
//   - Return the captured error (nil if every task succeeded).
//   - No data races — verify with: go test -race -v ./10_errgroup/
package errgroup

import "context"

// RunAll runs all tasks concurrently, cancelling the rest on the first error.
//
// TODO:
//   - childCtx, cancel := context.WithCancel(ctx); defer cancel().
//   - Declare a sync.WaitGroup, a sync.Once, and a variable to hold the first error.
//   - For each task, launch a goroutine that calls task(childCtx); on a non-nil
//     error, use the Once to store it once and call cancel() to stop the others.
//   - wg.Wait(), then return the stored error.
func RunAll(ctx context.Context, tasks []func(context.Context) error) error {
	panic("TODO: implement RunAll")
}
