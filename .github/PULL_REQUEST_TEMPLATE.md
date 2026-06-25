<!-- Thanks for contributing! Please fill this out so we can review quickly. -->

## What does this PR do?

<!-- One or two sentences. -->

## Type of change

- [ ] 🆕 New problem
- [ ] ✅ New / improved test case(s)
- [ ] 📝 Problem statement or docs improvement
- [ ] 🐛 Bug fix (flaky test, missed race, incorrect statement)
- [ ] 🧹 Other (chore, CI, tooling)

## Checklist

- [ ] I read [CONTRIBUTING.md](../CONTRIBUTING.md).
- [ ] This PR targets the **right branch** — skeletons + tests on `main`, working
      solutions on `solutions`. **No working solution is committed to `main`.**
- [ ] `gofmt -l .` reports nothing, `go vet ./...` is clean, and `go build ./...` succeeds.
- [ ] Tests pass with the race detector against a **real solution**
      (`go test -race ./...` on the `solutions` branch).
- [ ] The `main` skeleton still **panics with `TODO`** (tests are red until implemented).
- [ ] If adding a problem: README table updated and the new concept isn't already covered.

## Verification output

<!--
Paste evidence. For new problems / test changes, show BOTH:
  1. Tests passing against your solution (solutions branch)
  2. Tests failing/panicking against the main skeleton
-->

```
$ go test -race -v ./NN_topic/
...
```
