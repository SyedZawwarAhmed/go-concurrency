// Package etl — Problem 14: A cancellable, pool-bounded ETL pipeline.
//
// CONCEPT: The capstone composition — a real pipeline that reads files, turns
// their lines into records, and writes each record to a database whose
// connection pool is capped. Every stage is cancellation-aware, so the whole
// thing tears down promptly on error or ctx cancellation, and the writers are
// bounded so the DB is never overrun.
//
// SCENARIO: Run reads each path (real os.ReadFile), emits one record per
// non-empty line, and writes every record to db with exactly `dbWorkers`
// concurrent writers. sandbox.DB caps simultaneous calls — exceed the cap and
// Put fails with ErrTooManyConns, so dbWorkers must not outrun the pool.
//
// REQUIREMENTS:
//   - Reading and writing happen concurrently (a reader feeding a pool of
//     writers over a channel), not in two serial passes.
//   - At most dbWorkers DB writes in flight at once.
//   - Stop promptly on the first error (a bad file or a DB error) OR when ctx is
//     cancelled, cancelling the rest. No goroutine leaks on any exit path.
//   - Return the number of records successfully written and the first error
//     (nil on full success; ctx.Err() if the caller cancelled).
//   - No data races — go test -race -v ./14_etl/
package etl

import (
	"context"

	"github.com/SyedZawwarAhmed/go-concurrency/sandbox"
)

// Run executes the ETL pipeline and returns the number of records written and
// the first error encountered.
//
// HINT: derive a cancellable context; one reader goroutine produces records onto
// a channel; dbWorkers consumers write them; capture the first error once and
// cancel; count writes atomically.
func Run(ctx context.Context, db *sandbox.DB, paths []string, dbWorkers int) (int, error) {
	panic("TODO: implement Run")
}
