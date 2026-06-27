// Package sandbox provides realistic, hermetic dependencies for the concurrency
// problems: a fake database with latency and a connection-pool cap, a real HTTP
// test server, and helpers that write real files to disk. The point is to make
// the problems feel like systems work — file reads, DB calls, HTTP round-trips —
// without any external services.
package sandbox

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Errors returned by DB.
var (
	ErrNotFound     = errors.New("sandbox: record not found")
	ErrTooManyConns = errors.New("sandbox: too many concurrent connections")
)

// DB is a fake database that imitates a real one:
//   - every call costs `latency` (interruptible by context),
//   - the number of simultaneous in-flight calls is capped like a connection
//     pool — exceed it and the call fails with ErrTooManyConns,
//   - it records how many queries it actually served, so a test can assert that
//     a cache or singleflight really did cut down on load.
//
// DB is safe for concurrent use.
type DB struct {
	latency  time.Duration
	maxConns int32

	inflight atomic.Int32
	queries  atomic.Int64

	mu   sync.RWMutex
	data map[string]string
}

// DBOptions configures a DB. The zero value yields a DB with no latency and no
// connection cap.
type DBOptions struct {
	Latency  time.Duration     // per-call delay
	MaxConns int               // max simultaneous calls; 0 means unlimited
	Seed     map[string]string // initial contents
}

// NewDB returns a ready-to-use DB.
func NewDB(opts DBOptions) *DB {
	data := make(map[string]string, len(opts.Seed))
	for k, v := range opts.Seed {
		data[k] = v
	}
	return &DB{
		latency:  opts.Latency,
		maxConns: int32(opts.MaxConns),
		data:     data,
	}
}

func (db *DB) enter() error {
	n := db.inflight.Add(1)
	if db.maxConns > 0 && n > db.maxConns {
		db.inflight.Add(-1)
		return ErrTooManyConns
	}
	db.queries.Add(1)
	return nil
}

func (db *DB) leave() { db.inflight.Add(-1) }

// Get returns the value for key, or ErrNotFound. It respects ctx.
func (db *DB) Get(ctx context.Context, key string) (string, error) {
	if err := db.enter(); err != nil {
		return "", err
	}
	defer db.leave()
	if err := sleep(ctx, db.latency); err != nil {
		return "", err
	}
	db.mu.RLock()
	v, ok := db.data[key]
	db.mu.RUnlock()
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// Put writes key=val. It respects ctx.
func (db *DB) Put(ctx context.Context, key, val string) error {
	if err := db.enter(); err != nil {
		return err
	}
	defer db.leave()
	if err := sleep(ctx, db.latency); err != nil {
		return err
	}
	db.mu.Lock()
	db.data[key] = val
	db.mu.Unlock()
	return nil
}

// QueryCount reports how many calls the DB has served (successful enters).
func (db *DB) QueryCount() int64 { return db.queries.Load() }

// sleep waits for d or until ctx is done, whichever comes first.
func sleep(ctx context.Context, d time.Duration) error {
	if d <= 0 {
		return ctx.Err()
	}
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
