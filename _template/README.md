# Problem template

A copy-pasteable starting point for a **new problem**. Go tooling ignores directories
whose names start with `_`, so nothing in here is compiled by `go build ./...` or run by
`go test ./...`.

## How to use

1. Copy this directory to the next free number and a topic name:

   ```bash
   cp -r _template NN_topic        # e.g. cp -r _template 11_synccond
   ```

2. In **both** files, change `package template` to your real package name
   (the topic, not the number — e.g. `package synccond`).

3. Fill in `problem.go`:
   - The package doc comment: `CONCEPT`, `SCENARIO`, `REQUIREMENTS`, run command.
   - The skeleton function(s) with a `TODO:` hint comment and a `panic("TODO: ...")` body.
   - Any helpers that are *given* to the learner, fully implemented.

4. Fill in `problem_test.go` with a thorough, race-aware suite (happy path, edge cases,
   and the concurrency property — timing / leak / watchdog as appropriate).

5. Add a row to the README problem table, and provide the reference solution on the
   `solutions` branch.

See [CONTRIBUTING.md](../CONTRIBUTING.md#adding-a-new-problem) for the full workflow and
conventions.
