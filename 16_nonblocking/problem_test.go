package nonblocking

import (
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTryPushReportsFull(t *testing.T) {
	q := NewQueue(3)

	if !q.TryPush(1) {
		t.Fatal("TryPush(1) returned false, want true (room available)")
	}
	if !q.TryPush(2) {
		t.Fatal("TryPush(2) returned false, want true (room available)")
	}
	if !q.TryPush(3) {
		t.Fatal("TryPush(3) returned false, want true (room available)")
	}

	if q.TryPush(4) {
		t.Fatal("TryPush(4) returned true on a full queue, want false")
	}

	got := q.Drain()
	want := []int{1, 2, 3}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Drain() = %v, want %v", got, want)
	}

	// The 5th element was never accepted, so a fresh Drain is empty.
	if rest := q.Drain(); len(rest) != 0 {
		t.Fatalf("Drain() after first drain = %v, want empty", rest)
	}
}

func TestDrainEmptyIsNonBlocking(t *testing.T) {
	q := NewQueue(2)

	done := make(chan []int, 1)
	go func() {
		done <- q.Drain()
	}()

	select {
	case got := <-done:
		if len(got) != 0 {
			t.Fatalf("Drain() on empty queue = %v, want len 0", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Drain() on empty queue blocked for 2s, want prompt return")
	}
}

func TestConcurrentProducersNeverBlock(t *testing.T) {
	const (
		producers      = 50
		pushPerProducer = 1000
	)
	q := NewQueue(100)

	var accepted int64
	var producersWG sync.WaitGroup
	producersWG.Add(producers)

	for p := 0; p < producers; p++ {
		go func(base int) {
			defer producersWG.Done()
			for i := 0; i < pushPerProducer; i++ {
				if q.TryPush(base*pushPerProducer + i) {
					atomic.AddInt64(&accepted, 1)
				}
			}
		}(p)
	}

	producersDone := make(chan struct{})
	go func() {
		producersWG.Wait()
		close(producersDone)
	}()

	var drained int64
	consumerDone := make(chan struct{})
	go func() {
		defer close(consumerDone)
		for {
			select {
			case <-producersDone:
				// Final sweep to collect anything left buffered.
				drained += int64(len(q.Drain()))
				return
			default:
				drained += int64(len(q.Drain()))
			}
		}
	}()

	select {
	case <-consumerDone:
	case <-time.After(10 * time.Second):
		t.Fatal("concurrent producers/consumer did not finish within 10s, possible deadlock")
	}

	if got := atomic.LoadInt64(&drained); got != accepted {
		t.Fatalf("drained %d items, want %d (= accepted); accepted+dropped bookkeeping is off", got, accepted)
	}
}
