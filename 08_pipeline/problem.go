// Package pipeline — Problem 08: Channel pipelines (chained stages).
//
// CONCEPT: A pipeline is a series of stages connected by channels. Each stage is
// a goroutine that receives from an inbound channel, does work, sends on an
// outbound channel, and CLOSES its outbound channel when its input is exhausted.
// Because each stage closes downstream, the whole pipeline tears down cleanly.
//
// SCENARIO: Build three composable stages so that
//
//	Filter(Square(Generate(1, 2, 3, 4, 5, 6)), isEven)
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
//   - No data races — verify with: go test -race -v ./08_pipeline/
package pipeline

// Generate emits each of nums on the returned channel, in order, then closes it.
//
// TODO: make a chan int; in a goroutine send each num then close the channel; return it.
func Generate(nums ...int) <-chan int {
	panic("TODO: implement Generate")
}

// Square reads ints from in and emits their squares, preserving order.
//
// TODO: make a chan int; in a goroutine range over in, send v*v, then close; return it.
func Square(in <-chan int) <-chan int {
	panic("TODO: implement Square")
}

// Filter emits only the values from in for which pred returns true, in order.
//
// TODO: make a chan int; in a goroutine range over in, send v when pred(v) is true,
// then close; return it.
func Filter(in <-chan int, pred func(int) bool) <-chan int {
	panic("TODO: implement Filter")
}
