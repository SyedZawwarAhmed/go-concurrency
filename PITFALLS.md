# Concurrency Pitfalls — the bugs this course trains you to avoid

A field guide to the mistakes that bite real Go programs. Each one names the bug, shows it, and
gives the fix. The race detector (`go test -race`) catches some of these; the rest only bite in
production. Knowing all of them is most of what "comfortable with concurrent Go" means.

---

## 1. Sending on a closed channel → panic

A `close`d channel panics on send (and panics again if you close it twice). The sender must never
be the one racing a close.

```go
// BAD: a worker closes the channel it sends on; other workers then panic.
go func() { out <- v; close(out) }()

// GOOD: the senders signal done via WaitGroup; ONE closer closes after Wait.
go func() { wg.Wait(); close(out) }()
```

Rule of thumb: **the producer side closes, exactly once, and only after all sends are done.**
Receivers never close. (Problems 02, 06, 07, 24.)

## 2. Goroutine leaks — a blocked send/recv with no one on the other end

A goroutine parked forever on `ch <- v` or `<-ch` never gets collected. The classic case: a
pipeline stage whose consumer walked away.

```go
// BAD: if the consumer stops reading, this blocks forever.
for v := range in { out <- transform(v) }

// GOOD: also watch a cancellation signal so the send can give up.
for v := range in {
    select {
    case out <- transform(v):
    case <-ctx.Done():
        return
    }
}
```

Any send/recv whose lifetime you don't fully control needs an escape hatch (`ctx.Done()` or a
`done` channel). (Problems 08, 13, 14, 22.)

## 3. The loop-variable trap (pre-Go 1.22)

Before Go 1.22, the loop variable was shared across iterations, so goroutines all saw the final
value. Go 1.22+ gives each iteration its own copy, so this is fixed — **but only if you're on
1.22+.** Know which world you're in.

```go
// BAD on Go <1.22: every goroutine likely prints the same i.
for i := 0; i < n; i++ { go func() { use(i) }() }

// SAFE everywhere: shadow it explicitly.
for i := 0; i < n; i++ { i := i; go func() { use(i) }() }
```

This course targets Go 1.25, so the implicit per-iteration copy is in effect — but write code that
survives a downgrade. (Problem 01.)

## 4. Data race on a map or counter

Concurrent writes to a built-in `map` are a hard runtime panic, not just a race. Shared counters
need synchronization too.

```go
// BAD: concurrent c[key]++ → race / panic.
// GOOD: guard with a sync.Mutex (05), or use sync/atomic for plain counters (11),
//       or atomic.Pointer to swap whole values lock-free (18).
```

Reach for the cheapest correct tool: `atomic` for a counter, `Mutex` for a small critical
section, `RWMutex` when reads vastly outnumber writes (12).

## 5. Holding a lock across slow I/O

A mutex held during a network/DB/disk call serializes every other goroutine behind that call.

```go
// BAD: every reader blocks for the whole DB round-trip.
mu.Lock(); v := db.Get(key); cache[key] = v; mu.Unlock()

// GOOD: load outside the lock; or use singleflight so concurrent misses share
//       ONE load instead of each holding the lock in turn. (Problems 15, 26.)
```

## 6. `WaitGroup` misuse

`Add` must happen **before** the goroutine starts (not inside it), or `Wait` can race past it.
And a `WaitGroup` must not be copied after first use.

```go
// BAD: Add inside the goroutine — Wait may return before this runs.
go func() { wg.Add(1); defer wg.Done(); work() }()

// GOOD: Add before launching.
wg.Add(1)
go func() { defer wg.Done(); work() }()
// (Go 1.25's wg.Go(func(){...}) does the Add/Done for you — prefer it.)
```

## 7. Forgetting to release a `context`

`context.WithCancel`/`WithTimeout` return a `cancel` you must call — always — or you leak the
context (and its timer). `defer cancel()` right after creating it.

```go
ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
defer cancel() // even on the success path
```

(Problems 09, 10, 13, 14.)

## 8. Deadlock: unbuffered channel with no concurrent peer

A send on an unbuffered channel blocks until someone receives. Do both in the same goroutine and
you deadlock instantly.

```go
ch := make(chan int)
ch <- 1      // BAD: blocks forever; no receiver is running yet.
fmt.Println(<-ch)

// GOOD: receiver runs concurrently, or buffer the channel if that fits.
go func() { ch <- 1 }()
fmt.Println(<-ch)
```

## 9. `select` with no `default` when you needed non-blocking

If every case would block and there's no `default`, `select` blocks. Add `default` for try-send /
try-receive; **omit** it when you actually want to wait.

```go
select {
case ch <- v:     // sent
default:          // full — drop or handle, don't block
}
```

(Problems 16, 21, 23.)

## 10. Unbounded fan-out

`for _, x := range hugeSlice { go work(x) }` can spawn millions of goroutines and exhaust memory or
a downstream connection pool. Bound it: a worker pool (07, 25) or a semaphore (03, 13). The
`sandbox.DB` in this course will fail with `ErrTooManyConns` if you don't (14).

---

### How to catch these

- **Always** `go test -race`. It finds #3, #4, and many #2/#5 races. It does **not** find deadlocks,
  leaks that never get scheduled, or logic bugs — those need watchdogs and `QueryCount`-style
  assertions (which the tests here use).
- Run timing/concurrency tests with `-count=3` or higher; a bug that shows up 1-in-3 runs is still
  a bug.
- When a test hangs, you have a deadlock or a leak (#1, #2, #8). When it panics, suspect #1 or #4.
