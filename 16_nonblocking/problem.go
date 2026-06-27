// Package nonblocking — Problem 16: Non-blocking channel ops with select + default.
//
// CONCEPT: A select with a default clause never blocks: if no other case is
// immediately ready, default runs. This turns a channel send/receive into a
// "try" operation — try-send drops the item when there's no room, try-receive
// returns nothing when there's nothing buffered.
//
// SCENARIO: The Queue below wraps a buffered channel as a fixed-size FIFO. Hot
// producers push values that may be dropped under load, while a consumer drains
// whatever has accumulated so far — neither side ever waits on the other.
//
// REQUIREMENTS:
//   - TryPush adds v and returns true when there's room; when full it returns
//     false immediately and drops v (never blocks).
//   - Drain removes and returns all currently-buffered items in FIFO order
//     without blocking; an empty queue yields an empty (possibly nil) slice.
//   - Safe for concurrent producers and a draining consumer.
//   - No data races — go test -race -v ./16_nonblocking/
package nonblocking

// Queue is a fixed-capacity, non-blocking FIFO of ints.
type Queue struct {
	// TODO: choose your fields (a buffered chan int is all you need).
}

// NewQueue returns a Queue that buffers at most capacity items.
//
// HINT: a buffered channel of the given capacity is the whole queue.
func NewQueue(capacity int) *Queue {
	panic("TODO: implement NewQueue")
}

// TryPush adds v if there is room and returns true; if the buffer is full it
// returns false immediately WITHOUT blocking (the item is dropped).
//
// HINT: select { case q.ch <- v: ...; default: ... }
func (q *Queue) TryPush(v int) bool {
	panic("TODO: implement TryPush")
}

// Drain removes and returns all currently-buffered items (in FIFO order)
// without blocking; returns an empty (possibly nil) slice if none.
//
// HINT: loop a select with a default that returns once the channel is drained.
func (q *Queue) Drain() []int {
	panic("TODO: implement Drain")
}
