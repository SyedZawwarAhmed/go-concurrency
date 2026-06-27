// Package boundedqueue — Problem 22: Backpressure with a cancellable bounded queue.
//
// CONCEPT: Backpressure is what keeps a fast producer from outrunning a slow
// consumer. A bounded buffer enforces it: once the buffer is full, Put has
// nowhere to write and must WAIT until the consumer frees a slot — the queue's
// fullness pushes back on the producer. But a blocked Put can't wait forever:
// it must stay cancellable, so a done context unblocks it with an error.
//
// SCENARIO: A bounded FIFO connecting producers to consumers. Producers Put
// (blocking under load); a consumer Gets in order. When upstream is finished it
// Closes the queue; Get keeps returning buffered items until the queue is both
// closed AND empty, then reports ok=false.
//
// REQUIREMENTS:
//   - Put blocks while the queue is full and unblocks when a slot frees up.
//   - A blocked Put returns ctx.Err() if ctx is done before room appears; else nil.
//   - Get returns items in FIFO order; ok is false only once closed AND drained.
//   - Close is called once to signal no more Puts.
//   - No data races — go test -race -v ./22_backpressure/
package boundedqueue

import "context"

// BoundedQueue is a blocking FIFO with a fixed capacity. Put blocks while the
// queue is full (backpressure) but unblocks when space frees up or ctx is
// done. Get blocks while empty and returns ok=false once the queue is closed
// AND fully drained.
type BoundedQueue struct {
	// TODO: choose your fields (a buffered chan int of the given capacity is enough).
}

// NewBoundedQueue returns a queue that buffers at most capacity items.
//
// HINT: a buffered channel of size capacity is the whole queue.
func NewBoundedQueue(capacity int) *BoundedQueue {
	panic("TODO: implement NewBoundedQueue")
}

// Put blocks until there is room or ctx is done; returns ctx.Err() if it gave
// up, else nil.
//
// HINT: select { case q.ch <- v: return nil; case <-ctx.Done(): return ctx.Err() }
func (q *BoundedQueue) Put(ctx context.Context, v int) error {
	panic("TODO: implement Put")
}

// Get returns the next item; ok is false once the queue is closed and empty.
//
// HINT: v, ok := <-q.ch — a closed, drained channel yields ok=false.
func (q *BoundedQueue) Get() (v int, ok bool) {
	panic("TODO: implement Get")
}

// Close signals that no more items will be Put. Call once.
//
// HINT: close the channel so Get drains remaining items then reports ok=false.
func (q *BoundedQueue) Close() {
	panic("TODO: implement Close")
}
