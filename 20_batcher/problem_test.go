package batcher

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestFlushesBySize(t *testing.T) {
	var mu sync.Mutex
	var all []int
	var sizes []int

	b := NewBatcher(5, 10*time.Second, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		all = append(all, batch...)
		sizes = append(sizes, len(batch))
	})

	for i := 1; i <= 12; i++ {
		b.Add(i)
	}
	b.Stop()

	mu.Lock()
	defer mu.Unlock()

	want := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	if !reflect.DeepEqual(all, want) {
		t.Fatalf("flushed items = %v, want %v", all, want)
	}
	if len(sizes) < 2 || sizes[0] != 5 || sizes[1] != 5 {
		t.Fatalf("first two batch sizes = %v, want first two to be size 5 (size-triggered)", sizes)
	}
}

func TestFlushesByInterval(t *testing.T) {
	var mu sync.Mutex
	var all []int

	b := NewBatcher(1000, 30*time.Millisecond, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		all = append(all, batch...)
	})

	b.Add(1)
	b.Add(2)
	b.Add(3)

	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	got := append([]int(nil), all...)
	mu.Unlock()

	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("after interval, flushed = %v, want %v (interval-triggered)", got, want)
	}

	b.Stop()
}

func TestStopFlushesRemainder(t *testing.T) {
	var mu sync.Mutex
	var all []int

	b := NewBatcher(10, 10*time.Second, func(batch []int) {
		mu.Lock()
		defer mu.Unlock()
		all = append(all, batch...)
	})

	b.Add(1)
	b.Add(2)
	b.Add(3)
	b.Add(4)
	b.Stop()

	mu.Lock()
	defer mu.Unlock()

	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(all, want) {
		t.Fatalf("flushed items = %v, want %v", all, want)
	}
}
