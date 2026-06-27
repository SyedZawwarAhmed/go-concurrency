// Package shutdown — Problem 24: Graceful shutdown — stop accepting, drain, deadline.
//
// CONCEPT: Graceful shutdown means three things in order: stop accepting new
// work, drain the work already in flight or queued, and bound the whole thing
// by a deadline. A worker pool over one channel makes this natural — close the
// channel to signal "no more work" and let workers finish the range.
//
// SCENARIO: A job server runs a fixed pool of workers pulling jobs off a queue.
// Callers Submit jobs from many goroutines. When the server shuts down it must
// finish what it already accepted, but reject anything new — and give up if the
// caller's context deadline passes before the queue drains.
//
// REQUIREMENTS:
//   - NewServer starts exactly `workers` goroutines, each calling handle per job.
//   - Submit is safe for concurrent use; it returns true once the job is queued.
//   - Once Shutdown has begun, Submit returns false and accepts no new work.
//   - Submit must NEVER send on a closed channel (the classic panic): once
//     shutdown begins it returns false instead of sending. Gate the send and the
//     close on the same `closing` flag under the same mutex so they are mutually
//     exclusive.
//   - Shutdown stops accepting, waits for in-flight and queued jobs to drain, and
//     returns the count processed — unless ctx's deadline passes first, in which
//     case it returns the count so far and ctx.Err(). Safe to call once.
//   - No data races — go test -race -v ./24_shutdown/
package shutdown

import "context"

// Server processes submitted jobs concurrently with a fixed worker pool.
// Submit hands a job to the pool and returns true; once Shutdown has begun it
// returns false (no new work accepted). Shutdown stops accepting, waits for
// in-flight and already-queued jobs to drain, and returns how many jobs were
// processed — unless ctx's deadline passes first, in which case it returns the
// count so far and ctx.Err().
type Server struct{ /* student chooses fields */ }

// NewServer starts `workers` goroutines that each call handle for jobs.
//
// HINT: range a `jobs chan int` in each worker; a closer goroutine wg.Wait()s
// then close(done) so Shutdown can select on it.
func NewServer(workers int, handle func(int)) *Server {
	panic("TODO: implement NewServer")
}

// Submit enqueues job; returns false if the server is shutting down.
//
// HINT: under mu, check `closing` then send on jobs (both gated by mu so a send
// can never race the close).
func (s *Server) Submit(job int) bool {
	panic("TODO: implement Submit")
}

// Shutdown performs graceful shutdown bounded by ctx. Safe to call once.
//
// HINT: set closing+close(jobs) once under mu, then select on done vs ctx.Done().
func (s *Server) Shutdown(ctx context.Context) (processed int, err error) {
	panic("TODO: implement Shutdown")
}
