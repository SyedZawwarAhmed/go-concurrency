package filecount

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/SyedZawwarAhmed/go-concurrency/sandbox"
)

func TestCountFiles(t *testing.T) {
	dir := t.TempDir()
	const nFiles, linesPer = 20, 100
	paths, err := sandbox.SeedFiles(dir, nFiles, linesPer)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	var wantBytes int64
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatalf("read %s: %v", p, err)
		}
		wantBytes += int64(len(data))
	}

	got, err := CountFiles(paths, 8)
	if err != nil {
		t.Fatalf("CountFiles: %v", err)
	}
	if got.Files != nFiles {
		t.Errorf("Files = %d, want %d", got.Files, nFiles)
	}
	if got.Lines != nFiles*linesPer {
		t.Errorf("Lines = %d, want %d", got.Lines, nFiles*linesPer)
	}
	if got.Bytes != wantBytes {
		t.Errorf("Bytes = %d, want %d", got.Bytes, wantBytes)
	}
}

func TestCountFilesEmpty(t *testing.T) {
	got, err := CountFiles(nil, 4)
	if err != nil {
		t.Errorf("err = %v, want nil", err)
	}
	if got != (Counts{}) {
		t.Errorf("Counts = %+v, want zero", got)
	}
}

func TestCountFilesReportsError(t *testing.T) {
	dir := t.TempDir()
	paths, err := sandbox.SeedFiles(dir, 3, 10)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}
	paths = append(paths, filepath.Join(dir, "does_not_exist.txt"))

	if _, err := CountFiles(paths, 4); err == nil {
		t.Error("expected an error for a missing file, got nil")
	}
}
