package etl

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/SyedZawwarAhmed/go-concurrency/sandbox"
)

func TestRunLoadsAllRecords(t *testing.T) {
	dir := t.TempDir()
	const nFiles, linesPer = 10, 50
	paths, err := sandbox.SeedFiles(dir, nFiles, linesPer)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	const workers = 8
	db := sandbox.NewDB(sandbox.DBOptions{Latency: 2 * time.Millisecond, MaxConns: workers})

	n, err := Run(context.Background(), db, paths, workers)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	want := nFiles * linesPer
	if n != want {
		t.Errorf("wrote %d records, want %d", n, want)
	}
	if got := db.QueryCount(); got != int64(want) {
		t.Errorf("DB served %d writes, want %d", got, want)
	}
}

func TestRunStopsOnCancel(t *testing.T) {
	dir := t.TempDir()
	paths, err := sandbox.SeedFiles(dir, 20, 100) // 2000 records
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	db := sandbox.NewDB(sandbox.DBOptions{Latency: 5 * time.Millisecond, MaxConns: 4})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(25 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	n, err := Run(ctx, db, paths, 4)
	elapsed := time.Since(start)

	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
	if n >= 2000 {
		t.Errorf("wrote all %d records despite cancellation", n)
	}
	if elapsed > 500*time.Millisecond {
		t.Errorf("Run took %v; should stop promptly after cancel", elapsed)
	}
}

func TestRunReportsFileError(t *testing.T) {
	dir := t.TempDir()
	paths, err := sandbox.SeedFiles(dir, 3, 10)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	paths = append(paths, filepath.Join(dir, "missing.txt"))
	db := sandbox.NewDB(sandbox.DBOptions{MaxConns: 4})

	if _, err := Run(context.Background(), db, paths, 4); err == nil {
		t.Error("expected an error for a missing file, got nil")
	}
}
