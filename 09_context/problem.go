// Package cancellable — Problem 09: Cooperative cancellation with context.
//
// CONCEPT: A long-running worker should stop promptly when its context is
// cancelled or its deadline passes, returning ctx.Err(). The key tool is a
// select that watches ctx.Done() alongside the real work.
//
// SCENARIO: Worker consumes integers from the work channel, "processing" each
// (a short simulated delay). It keeps a running count of items processed and
// stops as soon as EITHER the work channel is exhausted OR the context is done.
//
// REQUIREMENTS:
//   - Loop with a select over: receiving from work, and <-ctx.Done().
//   - When work is closed/drained first: return (count, nil).
//   - When ctx is done first: return (count-so-far, ctx.Err()).
//     ctx.Err() is context.Canceled or context.DeadlineExceeded.
//   - No data races — verify with: go test -race -v ./09_context/
package cancellable

import (
	"context"
	"time"
)

// workDelay is the simulated per-item processing time.
const workDelay = 20 * time.Millisecond

// Worker processes items from work until the channel drains or ctx is done.
//
// TODO:
//   - Loop forever.
//   - Each iteration, select between two cases:
//       * ctx.Done() fired       -> return the running count and ctx.Err().
//       * a value arrives from work -> if the channel is closed (ok == false),
//         return (count, nil); otherwise simulate processing (sleep workDelay)
//         and increment the count.
func Worker(ctx context.Context, work <-chan int) (processed int, err error) {
	panic("TODO: implement Worker")
}
