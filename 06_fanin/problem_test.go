package fanin

import (
	"sort"
	"testing"
	"time"
)

// streamOf returns a channel that emits vals then closes.
func streamOf(vals ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range vals {
			ch <- v
		}
	}()
	return ch
}

// collect drains ch into a slice, with a watchdog so a never-closed output
// fails the test instead of hanging forever.
func collect(t *testing.T, ch <-chan int) []int {
	t.Helper()
	var got []int
	done := make(chan struct{})
	go func() {
		defer close(done)
		for v := range ch {
			got = append(got, v)
		}
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Merge output never closed (collect timed out) — check the close logic")
	}
	return got
}

func equalMultiset(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	x := append([]int(nil), a...)
	y := append([]int(nil), b...)
	sort.Ints(x)
	sort.Ints(y)
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func TestMerge(t *testing.T) {
	a := streamOf(1, 2, 3)
	b := streamOf(4, 5)
	c := streamOf(6, 7, 8, 9)

	got := collect(t, Merge(a, b, c))
	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !equalMultiset(got, want) {
		t.Errorf("merged = %v, want (in any order) %v", got, want)
	}
}

func TestMergeSingle(t *testing.T) {
	got := collect(t, Merge(streamOf(10, 20, 30)))
	if !equalMultiset(got, []int{10, 20, 30}) {
		t.Errorf("merged = %v, want [10 20 30]", got)
	}
}

func TestMergeZeroChannels(t *testing.T) {
	got := collect(t, Merge())
	if len(got) != 0 {
		t.Errorf("merging zero channels should yield nothing, got %v", got)
	}
}
