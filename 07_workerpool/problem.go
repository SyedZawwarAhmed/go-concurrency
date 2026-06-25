// Package workerpool — Problem 07: A fixed-size worker pool.
//
// CONCEPT: A bounded set of worker goroutines that pull jobs off a shared jobs
// channel and push results onto a results channel. Unlike Problem 01 (one
// goroutine per item), here a FIXED number of workers handle an arbitrary number
// of jobs.
//
// SCENARIO: Apply fn to every job using exactly `workers` goroutines. Results
// may come back in any order.
//
// REQUIREMENTS:
//   - Feed all jobs onto a jobs channel, then close it.
//   - Spawn exactly `workers` goroutines; each ranges the jobs channel, applies
//     fn, and sends the result on a results channel.
//   - Close the results channel once all workers finish (WaitGroup + closer).
//   - Collect and return all results (order does not matter).
//   - Each job must be processed exactly once.
//   - No data races — verify with: go test -race -v ./07_workerpool/
package workerpool

import (
	"sync"
)

// Process runs fn over all jobs using a pool of `workers` goroutines and returns
// every result (in any order).
//
// TODO:
//   - jobsCh := make(chan int, len(jobs)); push every job, then close(jobsCh).
//     (Buffering to len(jobs) lets you enqueue without a feeder goroutine.)
//   - resultsCh := make(chan int, len(jobs)); a sync.WaitGroup sized to workers.
//   - Spawn `workers` goroutines: `for j := range jobsCh { resultsCh <- fn(j) }`,
//     then wg.Done().
//   - A closer goroutine: wg.Wait() then close(resultsCh).
//   - range resultsCh into a slice and return it.
func Process(jobs []int, workers int, fn func(int) int) []int {
	var wg sync.WaitGroup

	jobsCh := make(chan int, len(jobs))

	for _, job := range jobs {
		jobsCh <- job
	}
	close(jobsCh)

	resultsCh := make(chan int, len(jobs))

	wg.Add(workers)
	for range workers {
		go func() {
			for j := range jobsCh {
				resultsCh <- fn(j)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	result := make([]int, 0, len(jobs))
	for res := range resultsCh {
		result = append(result, res)
	}

	return result
}
