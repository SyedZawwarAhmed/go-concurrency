package pubsub

import (
	"fmt"
	"testing"
	"time"
)

func recv(t *testing.T, ch <-chan string, want string) {
	t.Helper()
	select {
	case got := <-ch:
		if got != want {
			t.Fatalf("received %q, want %q", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timed out waiting to receive %q", want)
	}
}

func recvClosed(t *testing.T, ch <-chan string) {
	t.Helper()
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("channel delivered a value, want closed (ok==false)")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for channel to be closed")
	}
}

func TestFanOutToAllSubscribers(t *testing.T) {
	b := NewBroker()
	sub1 := b.Subscribe(1)
	sub2 := b.Subscribe(1)
	sub3 := b.Subscribe(1)

	b.Publish("hello")

	recv(t, sub1, "hello")
	recv(t, sub2, "hello")
	recv(t, sub3, "hello")
}

func TestSlowSubscriberDoesNotBlockPublish(t *testing.T) {
	b := NewBroker()
	_ = b.Subscribe(1)    // slow: buffered to 1, never drained
	fast := b.Subscribe(10)

	done := make(chan struct{})
	go func() {
		for i := 0; i < 5; i++ {
			b.Publish(fmt.Sprintf("m%d", i))
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Publish blocked on a slow subscriber")
	}

	// fast has a buffer of 10, so all 5 must have fit.
	for i := 0; i < 5; i++ {
		recv(t, fast, fmt.Sprintf("m%d", i))
	}
	// The slow subscriber dropping messages is fine; we don't assert on it.
}

func TestUnsubscribeStopsDelivery(t *testing.T) {
	b := NewBroker()
	sub := b.Subscribe(1)
	b.Unsubscribe(sub)

	recvClosed(t, sub)

	// Publishing after the only subscriber left must not panic.
	b.Publish("x")
}

func TestCloseClosesAllSubscribers(t *testing.T) {
	b := NewBroker()
	a := b.Subscribe(1)
	c := b.Subscribe(1)

	b.Close()

	recvClosed(t, a)
	recvClosed(t, c)

	// Publish after Close is a safe no-op.
	b.Publish("y")
}
