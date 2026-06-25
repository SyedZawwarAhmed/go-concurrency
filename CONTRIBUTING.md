# Contributing to Go Concurrency

Thanks for wanting to make this course better! 🐹 Contributions of every size are
welcome — a new problem, a sharper test case, a clearer problem statement, a typo fix,
or a bug report all add real value.

This guide explains **how the repo is laid out**, the **conventions every problem
follows**, and the **exact steps** for the most common contributions.

---

## The most important thing to understand: the two-branch model

This repo is a **test-driven course**. Learners read a problem, fill in a skeleton, and
make failing tests pass (red → green). Because of that, the two long-lived branches play
very different roles:

| Branch      | Contains                                                       | Do tests pass?                         |
| ----------- | -------------------------------------------------------------- | -------------------------------------- |
| `main`      | Problem statements, **skeletons** (`panic("TODO: ...")`), tests | ❌ **No — by design.** Skeletons panic. |
| `solutions` | The same problems with worked-out implementations              | ✅ **Yes — fully green.**               |

> ⚠️ **Never commit a working solution to `main`.** The whole point of `main` is that a
> learner can clone it and practice. The function body of every `problem.go` on `main`
> must remain `panic("TODO: implement ...")`.

When you add or change a problem you will usually touch **both** branches:

1. The **skeleton + tests** go on `main` (via your PR).
2. The **reference solution** goes on `solutions` (so we can prove the tests are
   passable). See [Adding a new problem](#adding-a-new-problem) for the workflow.

---

## Repo layout & conventions

Each problem lives in a numbered, self-contained package:

```
NN_topic/
├── problem.go        # problem statement (doc comment) + types + skeleton func(s)
└── problem_test.go   # the test suite that defines "done"
```

Conventions that keep the course consistent — please match them:

### `problem.go`

- **Package name** matches the topic, not the number: `package waitgroups`, not
  `package _01`.
- A **package doc comment** at the top states, in this order:
  - `CONCEPT:` — the one Go concurrency idea being taught.
  - `SCENARIO:` — a concrete, realistic framing of the task.
  - `REQUIREMENTS:` — a bulleted contract the implementation must satisfy (ordering,
    no leaks, no races, edge cases like nil/empty, closing channels exactly once, etc.).
  - The run command, e.g. `Run the tests:  go test -race -v ./NN_topic/`.
- The skeleton function has a **`TODO:` doc comment** sketching the intended approach as
  hints (not a full solution), and a body of exactly `panic("TODO: implement <Name>")`.
- Any helper that is *given* to the learner (e.g. `massage` in `01`) is fully
  implemented, with a comment noting it's already done.
- Constants the tests depend on (delays, etc.) carry a comment like
  `// Do not change it — the timing test relies on this value.`

### `problem_test.go`

- Same package as the code (white-box tests).
- **Always race-friendly** — every problem is run with `-race`. Use `sync/atomic` or
  proper synchronization in the tests themselves.
- Cover, at minimum: the happy path, **edge cases** (nil/empty input, zero channels,
  single element), and the **concurrency property** the problem teaches:
  - **Timing tests** to prove work happens concurrently, not serially (compare elapsed
    time against a serial baseline).
  - **Goroutine-leak checks** where relevant (`runtime.NumGoroutine()` before/after).
  - A **watchdog** (`select` + `time.After`) so a deadlocked implementation fails after
    ~2s with a clear message instead of hanging the whole suite forever.
- Failure messages are **specific and actionable** — say what was expected, what was
  got, and ideally *why it likely happened* (e.g. "the long task should have been
  cancelled, not run to completion").

> 📁 A copy-pasteable starting point lives in [`_template/`](_template/). Go tooling
> ignores directories that start with `_`, so the template never affects `go build` or
> `go test ./...`.

---

## Common contributions

### Improving or adding test cases

The highest-leverage contribution. Better tests catch more subtle concurrency bugs.

1. Branch off `main`.
2. Edit the relevant `NN_topic/problem_test.go` — add table cases, a leak check, a
   tighter timing bound, a watchdog, etc.
3. **Verify your tests pass against a real solution.** Check out the `solutions` branch,
   copy your new test file over, and run `go test -race -v ./NN_topic/`. Good tests must
   be **green against a correct implementation** and **red against the `main` skeleton**.
4. Open a PR to `main` describing what bug class the new test catches.

### Improving a problem statement or docs

Typos, clearer wording, better hints, fixing a misleading requirement. Branch off `main`,
edit the doc comments / README, open a PR. No solution branch changes needed.

### Adding a new problem

1. Pick the **next free number** (or propose where it fits in the ladder — open a
   [New problem proposal](../../issues/new/choose) issue first if you'd like feedback on
   the concept before building it).
2. Copy `_template/` to `NN_topic/`, rename the package, and write the problem statement,
   skeleton, and a thorough test suite following the conventions above.
3. On a branch off **`solutions`**, add the *working implementation* so the tests can be
   proven green. Include this in your PR (or link a companion PR / paste the solution in
   the PR description) — maintainers will land it on `solutions`.
4. Add a row to the problem table in `README.md`.
5. Open your PR to `main` with the skeleton + tests. In the description, show the test
   output **passing against your solution** and **failing (panicking) against the
   skeleton**.

A good new problem teaches **one** distinct concept that isn't already covered, fits the
incremental difficulty ladder, and has tests that assert the *concurrency property*, not
just the return value.

---

## Before you open a PR

Run the local checks — these mirror CI:

```bash
# Formatting (CI fails on unformatted code)
gofmt -l .

# Static analysis
go vet ./...

# Everything compiles — skeletons must build (they panic, but they must compile)
go build ./...
```

On the `solutions` branch, also run the full suite with the race detector:

```bash
go test -race ./...
```

- Keep PRs **focused** — one problem or one coherent improvement per PR.
- Match the surrounding style; run `gofmt`.
- Use clear commit messages (e.g. `feat: add problem 11 (sync.Cond)`,
  `test: tighten 03 race detection`, `fix: typo in 06 requirements`).

---

## Reporting bugs

Found a flaky test, a race the suite misses, or an incorrect problem statement? Open a
[bug report](../../issues/new/choose). Include the problem number, the command you ran
(with `-race`), and the full output.

---

## Code of Conduct

This project follows a [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you agree
to uphold it.

## License

By contributing, you agree that your contributions will be licensed under the
[MIT License](LICENSE) that covers this project.
