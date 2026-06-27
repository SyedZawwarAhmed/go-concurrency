package sandbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SeedFiles writes n real files into dir, each with linesPerFile non-empty
// lines, and returns their paths. Read them back with os.ReadFile — this is
// genuine disk I/O. Pair it with t.TempDir() in a test for automatic cleanup.
func SeedFiles(dir string, n, linesPerFile int) ([]string, error) {
	paths := make([]string, 0, n)
	for i := 0; i < n; i++ {
		p := filepath.Join(dir, fmt.Sprintf("file_%03d.txt", i))
		var b strings.Builder
		for l := 0; l < linesPerFile; l++ {
			fmt.Fprintf(&b, "file %d line %d alpha beta gamma\n", i, l)
		}
		if err := os.WriteFile(p, []byte(b.String()), 0o644); err != nil {
			return nil, err
		}
		paths = append(paths, p)
	}
	return paths, nil
}
