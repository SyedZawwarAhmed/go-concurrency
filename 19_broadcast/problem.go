// Package broadcast — Problem 19: One-to-many signaling by closing a channel.
//
// CONCEPT: Closing a channel wakes EVERY receiver blocked on it at once — a
// receive on a closed channel returns immediately, forever. This is the
// idiomatic Go broadcast: one Open frees all current and future waiters with no
// per-waiter bookkeeping. sync.Cond (Broadcast) is the alternative when you need
// repeatable signaling, but a one-shot gate is simplest as a closed channel.
//
// SCENARIO: A startup gate. Many goroutines call Wait and park until some
// coordinator calls Open exactly once (the "go" signal), after which everyone
// proceeds. Latecomers that Wait after Open must not block.
//
// REQUIREMENTS:
//   - Wait blocks until Open is called; after Open, every Wait returns at once.
//   - Wait calls made AFTER Open return immediately.
//   - Open is idempotent: calling it more than once is a safe no-op (never panics).
//   - Open is safe to call concurrently from multiple goroutines.
//   - No data races — go test -race -v ./19_broadcast/
package broadcast

// Gate blocks every caller of Wait until Open is called, after which all
// current AND future Wait calls return immediately. Open is idempotent and
// safe to call from multiple goroutines. The zero value is NOT ready; use
// NewGate.
type Gate struct{ /* student chooses fields, e.g. chan struct{} + sync.Once */ }

// NewGate returns a Gate that is closed (waiters block) until Open is called.
func NewGate() *Gate {
	panic("TODO: implement NewGate")
}

// Wait blocks until the gate is opened.
//
// HINT: a receive on a closed channel returns immediately, forever.
func (g *Gate) Wait() {
	panic("TODO: implement Wait")
}

// Open releases all waiters; calling it more than once is a safe no-op.
//
// HINT: sync.Once guards the close so a second Open can't double-close (panic).
func (g *Gate) Open() {
	panic("TODO: implement Open")
}
