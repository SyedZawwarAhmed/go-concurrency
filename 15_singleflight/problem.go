// Package singleflight — Problem 15: Collapse duplicate in-flight calls.
//
// CONCEPT: When N goroutines ask for the same expensive thing at the same moment
// (a cache stampede), you want the work done ONCE and the single result shared
// with all callers. This is the singleflight pattern, built by hand from a mutex,
// a map of in-flight calls, and a sync.WaitGroup used as a one-shot barrier.
//
// SCENARIO: A guard in front of an expensive backend (sandbox.DB counts calls).
// Many goroutines call Do with the same key concurrently; fn must run once and
// every caller gets that result. Different keys proceed independently.
//
// REQUIREMENTS:
//   - For a given key, concurrent Do calls share ONE execution of fn; late
//     callers wait for the in-flight one and receive its value and error.
//   - Distinct keys do not block each other.
//   - Once a call completes, it is forgotten — a later Do(key, ...) runs fn
//     again (no permanent caching; that's problem 12's job).
//   - The zero value of Group is ready to use.
//   - No data races — go test -race -v ./15_singleflight/
package singleflight

// Group collapses concurrent calls sharing the same key into one execution.
// It is generic over the call's result type and safe for concurrent use.
type Group[T any] struct {
	// TODO: add the fields you need (a sync.Mutex and a map of in-flight calls).
}

// Do executes fn (once) for key, returning its result to every concurrent caller
// that shared the same in-flight call.
//
// HINT: under the lock, look for an in-flight call for key; if none, register one
// (a struct holding a WaitGroup plus result fields), release the lock, run fn,
// then signal the WaitGroup and remove the entry. Waiters block on the WaitGroup.
func (g *Group[T]) Do(key string, fn func() (T, error)) (T, error) {
	panic("TODO: implement Do")
}
