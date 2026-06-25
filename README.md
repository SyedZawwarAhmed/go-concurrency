# Go Concurrency — A Test-Driven Course

A self-paced, **test-driven** course to master Go concurrency. Each problem ships with a
description, a skeleton you fill in, and a robust test suite. **No solutions are included** —
your job is to make the failing tests pass (red → green).

The problems are ordered as an incremental ladder: start at `01` and climb. Each package is
self-contained, so you can also jump around.

> 💡 **Want to peek at my solutions?** This `main` branch is kept solution-free so you can
> practice. My worked-out answers live on the [`solutions`](../../tree/solutions) branch —
> switch to it only if you want to compare after giving a problem an honest try.

## The 10 problems

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
