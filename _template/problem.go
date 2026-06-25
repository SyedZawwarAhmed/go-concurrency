// Package template — Problem NN: <Title>.
//
// CONCEPT: <The single Go concurrency idea this problem teaches.>
//
// SCENARIO: <A concrete, realistic framing of the task the learner implements.>
//
// REQUIREMENTS:
//   - <The observable contract: ordering guarantees, return shape, etc.>
//   - <Resource hygiene: close channels exactly once, no goroutine leaks.>
//   - <Edge cases: nil/empty input, zero items — and that they must not hang.>
//   - No data races — verify with: go test -race -v ./NN_topic/
package template

// processDelay is a sample simulated-work duration. Keep any timing constant the
// tests depend on here, with a note that it must not change.
//
// const processDelay = 50 * time.Millisecond

// Solve is the function the learner implements.
//
// TODO: Sketch the intended approach as HINTS (not a full solution):
//   - <step 1>
//   - <step 2>
//   - <step 3>
func Solve(input []int) []int {
	panic("TODO: implement Solve")
}
