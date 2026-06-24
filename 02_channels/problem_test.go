package channels

import (
	"fmt"
	"testing"
	"time"
)

func TestCheckURLs(t *testing.T) {
	tests := []struct {
		name string
		urls []string
	}{
		{"several", []string{"https://a.com", "https://b.com", "https://c.com"}},
		{"single", []string{"https://only.com"}},
		{"empty", []string{}},
		{"nil", nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Run in a goroutine + watchdog so a deadlock fails fast instead of
			// hanging the whole suite (important for the empty-input case).
			done := make(chan map[string]int, 1)
			go func() { done <- CheckURLs(tc.urls) }()

			select {
			case got := <-done:
				if len(got) != len(tc.urls) {
					t.Fatalf("got %d results, want %d: %v", len(got), len(tc.urls), got)
				}
				for _, u := range tc.urls {
					if got[u] != 200 {
						t.Errorf("status[%q] = %d, want 200", u, got[u])
					}
				}
			case <-time.After(2 * time.Second):
				t.Fatal("CheckURLs deadlocked (no result after 2s) — check the channel close logic / empty input")
			}
		})
	}
}

func TestCheckURLsManyUnique(t *testing.T) {
	urls := make([]string, 25)
	for i := range urls {
		urls[i] = fmt.Sprintf("https://host-%d.example", i)
	}
	got := CheckURLs(urls)
	if len(got) != len(urls) {
		t.Fatalf("got %d unique results, want %d", len(got), len(urls))
	}
	for _, u := range urls {
		if got[u] != 200 {
			t.Errorf("status[%q] = %d, want 200", u, got[u])
		}
	}
}
