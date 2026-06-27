# Go Concurrency — A Test-Driven Course

A self-paced, **test-driven** course to master Go concurrency. Each problem ships with a
description, a skeleton you fill in, and a robust test suite. **No solutions are included** —
your job is to make the failing tests pass (red → green).

The problems are ordered as an incremental ladder: start at `01` and climb. Each package is
self-contained, so you can also jump around.

> 💡 **Want to peek at my solutions?** This `main` branch is kept solution-free so you can
> practice. My worked-out answers live on the [`solutions`](../../tree/solutions) branch —
> switch to it only if you want to compare after giving a problem an honest try.

## The problems

**Part 1 — the primitives (toy scenarios, `time.Sleep` for "work").** Learn each tool in
isolation.

| #   | Directory           | Concept                                  | What you implement                                                        |
| --- | ------------------- | ---------------------------------------- | ------------------------------------------------------------------------- |
| 01  | `01_waitgroups`     | Goroutines + `sync.WaitGroup`            | Process a slice of strings in parallel, results in order                  |
| 02  | `02_channels`       | Unbuffered channels, close & range       | Concurrent URL status checker collecting results over a channel           |
| 03  | `03_ratelimit`      | Buffered channel as a semaphore          | Batch writer capped at N concurrent workers                               |
| 04  | `04_select_timeout` | `select` + `time.After`                  | Resilient fetcher: first of primary/replica wins, else time out           |
| 05  | `05_mutex`          | `sync.Mutex` shared-state safety         | A concurrency-safe counter (map guarded by a mutex)                       |
| 06  | `06_fanin`          | Fan-in (merging channels)                | Merge N channels into one                                                 |
| 07  | `07_workerpool`     | Fixed worker pool (jobs/results)         | N workers consuming a jobs channel, producing results                     |
| 08  | `08_pipeline`       | Channel pipelines (chained stages)       | `Generate → Square → Filter` composable stages                            |
| 09  | `09_context`        | `context.Context` cancellation/deadlines | A worker that stops on cancel or deadline                                 |
| 10  | `10_errgroup`       | First-error-cancels-the-rest             | Run tasks concurrently; first error cancels the others (errgroup by hand) |

**Part 2 — real systems work.** Same primitives, now driving **real file, DB, and HTTP I/O**
via the [`sandbox/`](sandbox/) module. Each problem introduces ~one new idea while reusing what
you already know; the realism ramps up as you climb. Do these in order.

| #   | Directory          | New concept                              | Real-world scenario                                                       |
| --- | ------------------ | ---------------------------------------- | ------------------------------------------------------------------------- |
| 11  | `11_filecount`     | Lock-free counters (`sync/atomic`)       | Pool of workers tallying real files on disk; aggregate without a mutex    |
| 12  | `12_readcache`     | `sync.RWMutex` + double-checked locking  | Read-through cache in front of a slow DB; load each key exactly once      |
| 13  | `13_httpfetch`     | Bounded concurrency + `context` deadline | Batch URL fetcher over a real HTTP server, capped & cancellable           |
| 14  | `14_etl`           | Cancellable, pool-bounded pipeline       | ETL: read files → records → DB writers; tears down cleanly on cancel      |
| 15  | `15_singleflight`  | Collapse duplicate in-flight calls       | Stampede guard: N callers, one backend call, shared result (generic)      |

### The `sandbox/` module

Part 2 problems depend on `sandbox/`, which provides realistic, hermetic dependencies so the
scenarios feel like systems work rather than `time.Sleep` toys:

- **`sandbox.DB`** — a fake database with per-call latency, a hard connection-pool cap
  (`ErrTooManyConns` if you exceed it), and a `QueryCount()` so tests can _prove_ your cache or
  singleflight actually cut backend load.
- **`sandbox.Server`** — a real `httptest` server (HTTP over loopback) that tracks peak
  concurrency, so a test can verify your rate-limiting truly bounds in-flight requests.
- **`sandbox.SeedFiles`** — writes real files to a temp dir for genuine `os.ReadFile` I/O.

You implement only the `problem.go` in each numbered package; you never edit `sandbox/`.

**Part 3 — signaling, time & lock-free state.** The patterns you reach for when channels-and-a-
mutex aren't the whole answer: non-blocking ops, lazy init, lock-free reads, broadcast, and
time-driven work. These use _injected_ dependencies (a `flush` func, a `handle` func) — the
concept is the lesson, so the tests stay deterministic.

| #   | Directory          | New concept                              | Scenario                                                                  |
| --- | ------------------ | ---------------------------------------- | ------------------------------------------------------------------------- |
| 16  | `16_nonblocking`   | `select` + `default` (try-send/recv)     | Lossy hot-path queue that never blocks producers                          |
| 17  | `17_once`          | `sync.Once` (exactly-once init)          | Generic lazy value built once under concurrent `Get`                      |
| 18  | `18_atomicvalue`   | `atomic.Pointer` (lock-free reads)       | Hot-swappable config read on the hot path, replaced wholesale             |
| 19  | `19_broadcast`     | Broadcast by closing a channel           | A gate that releases all waiting goroutines at once                       |
| 20  | `20_batcher`       | `time.Ticker` size/time batching         | Buffer that flushes every N items _or_ every interval                     |
| 21  | `21_tokenbucket`   | Token-bucket rate limiting               | Refill-over-time limiter (rate, not concurrency — contrast with `03`)     |
| 22  | `22_backpressure`  | Cancellable bounded queue                | Fast producers wait for slow consumers; a blocked `Put` is cancellable    |

**Part 4 — services & lifecycle.** Putting it together into long-lived components that start,
serve, and stop cleanly — the difference between "works in a test" and "survives production."

| #   | Directory          | New concept                              | Scenario                                                                  |
| --- | ------------------ | ---------------------------------------- | ------------------------------------------------------------------------- |
| 23  | `23_pubsub`        | Fan-out hub, safe (un)subscribe          | In-memory pub/sub; a slow subscriber can't block publish                  |
| 24  | `24_shutdown`      | Graceful shutdown (drain + deadline)     | Stop accepting, drain in-flight, honor a shutdown deadline                |
| 25  | `25_retrypool`     | Resilient pool (retries + per-item errs) | Worker pool that retries transient failures and aggregates outcomes       |

**Part 5 — capstone.** One component that integrates the whole course.

| #   | Directory             | Integrates                            | Scenario                                                               |
| --- | --------------------- | ------------------------------------- | ---------------------------------------------------------------------- |
| 26  | `26_capstone_cache`   | `12` + `15` + `20` + `24`             | Concurrent TTL cache: read-through, singleflight, janitor eviction, `Close` |

> See **[PITFALLS.md](PITFALLS.md)** for the classic concurrency bugs this course trains you to
> avoid (send-on-closed, goroutine leaks, the loop-variable trap, lock-during-I/O, and more) —
> read it once now, then again after Part 4 when it'll really land.

Each `problem.go` contains the full problem statement in its file comments, the supporting
types/stubs, and a skeleton function whose body is `panic("TODO: implement ...")`. Replace the
panic with your solution.

## How to run the tests

> Tests **fail until you implement the TODOs** — that's expected. An un-implemented skeleton
> panics with `TODO: implement ...`; that's your signal to write code.

Always run with the **race detector** on — concurrency bugs hide until `-race` finds them.

```bash
# Run a single problem (verbose, with the race detector)
go test -race -v ./01_waitgroups/

# Run everything
go test -race ./...

# Run everything, verbose
go test -race -v ./...

# Run one specific test by name
go test -race -run TestMassageStringsRunsConcurrently -v ./01_waitgroups/
```

A clean pass for a problem looks like `ok  .../01_waitgroups`. Keep iterating on a package until
it's green, then move to the next number.

### Tips

- `02`, `03`, `05`, and onward are specifically designed to catch race conditions — never skip
  `-race` on them.
- Some tests have a built-in watchdog (a `time.After` guard): if your code deadlocks, the test
  fails after ~2s with a clear message instead of hanging forever.
- Timing tests (e.g. `01`, `10`) assert that work actually happens _concurrently_, not serially.

## Contributing

Contributions are very welcome — new problems, sharper test cases, clearer problem
statements, and bug reports all add value. See **[CONTRIBUTING.md](CONTRIBUTING.md)** for
the conventions and workflow (note the two-branch model: skeletons + tests live on
`main`, working solutions on `solutions`). A copy-pasteable starting point for new
problems is in [`_template/`](_template/).

This project is released under the [MIT License](LICENSE) and follows a
[Code of Conduct](CODE_OF_CONDUCT.md).

Happy hacking. 🐹
