// Package fanin — Problem 06: Fan-in (merging many channels into one).
//
// CONCEPT: Combine the output of N input channels into a single output channel
// that a single consumer can range over. This is the "fan-in" half of the
// fan-out / fan-in pattern.
//
// SCENARIO: Several producers each emit ints on their own channel. Merge them
// into one channel.
//
// REQUIREMENTS:
//   - Start one goroutine per input channel that forwards everything it receives
//     to a shared output channel.
//   - Close the output channel exactly once, AFTER every input has drained
//     (use a sync.WaitGroup + a dedicated closer goroutine).
//   - Merging zero channels returns a channel that is already closed, so ranging
//     over it terminates immediately (no hang).
//   - No data races — verify with: go test -race -v ./06_fanin/
package fanin

// Merge fans in all input channels into a single output channel.
//
// TODO:
//   - Make the output channel: out := make(chan int).
//   - Size a sync.WaitGroup to len(channels). For each input, launch a goroutine
//     that ranges the input, forwards each value to out, then calls wg.Done().
//   - Launch one more goroutine that calls wg.Wait() then close(out).
//   - Return out.
func Merge(channels ...<-chan int) <-chan int {
	panic("TODO: implement Merge")
}
