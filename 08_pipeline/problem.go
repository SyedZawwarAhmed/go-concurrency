// Package pipeline — Problem 08: Channel pipelines (chained stages).
//
// CONCEPT: A pipeline is a series of stages connected by channels. Each stage is
// a goroutine that receives from an inbound channel, does work, sends on an
// outbound channel, and CLOSES its outbound channel when its input is exhausted.
// Because each stage closes downstream, the whole pipeline tears down cleanly.
//
// SCENARIO: Build three composable stages so that
//
//	Filter(ctx, Square(ctx, Generate(ctx, 1, 2, 3, 4, 5, 6)), isEven)
//
// yields the even squares of 1..6, in order: 4, 16, 36.
//
// REQUIREMENTS (for EACH stage):
//   - Create an output channel.
//   - Launch a goroutine that produces values (Generate ranges its slice; Square
//     and Filter range their input channel), sends results out, and CLOSES the
//     output when done.
//   - Return the output channel immediately (don't block the caller).
//   - Order is preserved end to end.
//   - Every send must also watch ctx.Done() so the stage unblocks and exits if
//     the consumer abandons the pipeline (cancellation) — no goroutine leaks.
//   - No data races — verify with: go test -race -v ./08_pipeline/
package pipeline

import "context"

// Generate emits each of nums on the returned channel, in order, then closes it.
// It stops early if ctx is cancelled.
func Generate(ctx context.Context, nums ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, num := range nums {
			select {
			case ch <- num:
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

// Square reads ints from in and emits their squares, preserving order.
// It stops early if ctx is cancelled.
func Square(ctx context.Context, in <-chan int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for v := range in {
			select {
			case ch <- v * v:
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}

// Filter emits only the values from in for which pred returns true, in order.
// It stops early if ctx is cancelled.
func Filter(ctx context.Context, in <-chan int, pred func(int) bool) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for v := range in {
			if !pred(v) {
				continue
			}
			select {
			case ch <- v:
			case <-ctx.Done():
				return
			}
		}
	}()
	return ch
}
