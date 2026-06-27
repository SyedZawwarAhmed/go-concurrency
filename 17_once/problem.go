// Package lazyinit — Problem 17: Exactly-once lazy initialization with sync.Once.
//
// CONCEPT: Some values are expensive to build (open a connection, parse a big
// config) and may never be needed. Lazy initialization defers that work to the
// first use, and sync.Once guarantees it happens exactly once even when many
// goroutines race to be first; every later caller reuses the cached value.
//
// SCENARIO: A Lazy[T] wraps an init function. The first Get runs init and caches
// the result; concurrent and subsequent Gets all return that same value without
// re-running init.
//
// REQUIREMENTS:
//   - init runs at most once, on the first Get, across all goroutines.
//   - NewLazy stores init but does NOT call it.
//   - Every Get returns the one cached value.
//   - No data races — go test -race -v ./17_once/
package lazyinit

// Lazy builds a value of type T at most once, on the first Get, even under
// concurrent callers, then returns the cached value on every later Get.
type Lazy[T any] struct {
	// TODO: add the fields you need (e.g. a sync.Once, the cached value, and init fn).
}

// NewLazy stores the (expensive) init function; it is NOT called yet.
func NewLazy[T any](init func() T) *Lazy[T] {
	panic("TODO: implement NewLazy")
}

// Get returns the value, running init exactly once across all goroutines.
//
// HINT: use sync.Once.Do to run init once and stash its result, then return it.
func (l *Lazy[T]) Get() T {
	panic("TODO: implement Get")
}
