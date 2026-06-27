// Package pubsub — Problem 23: An in-memory publish/subscribe broker.
//
// CONCEPT: Fan-out with backpressure isolation. A broker holds a set of
// subscriber channels behind a mutex; Publish copies a message to each one with
// a non-blocking `select { case ch <- msg: default: }` so a single slow
// consumer can never stall the publisher or the other subscribers. Closing a
// channel signals "no more messages" to its receiver; the broker owns every
// channel it hands out and is the only one allowed to close them.
//
// SCENARIO: An event bus. Components Subscribe to receive every message
// published after they join, drain at their own pace, and Unsubscribe when
// done. A shutdown Closes the broker, releasing all subscribers at once.
//
// REQUIREMENTS:
//   - Publish fans out to all current subscribers; a full buffer drops the
//     message for that subscriber only (never blocks the publisher).
//   - Subscribe returns a fresh buffered channel; Unsubscribe removes & closes it.
//   - Close closes every subscriber channel; Publish after Close is a no-op.
//   - Guard against double-close and send-on-closed-channel panics.
//   - No data races — go test -race -v ./23_pubsub/
package pubsub

// Broker is a concurrency-safe publish/subscribe hub. Subscribe returns a new
// channel that receives messages published after it subscribes. Publish
// fans a message out to all current subscribers WITHOUT blocking on a slow
// one: if a subscriber's buffer is full, the message is dropped for that
// subscriber only. Unsubscribe removes and closes a subscription. Close shuts
// the broker down and closes every subscriber channel.
type Broker struct{ /* student chooses fields, e.g. sync.Mutex + map */ }

// NewBroker returns an empty broker ready for subscriptions.
func NewBroker() *Broker {
	panic("TODO: implement NewBroker")
}

// Subscribe returns a receive-only channel buffered to `buffer`.
//
// HINT: make(chan string, buffer), add it to the set under the lock, return it.
func (b *Broker) Subscribe(buffer int) <-chan string {
	panic("TODO: implement Subscribe")
}

// Unsubscribe removes ch (as returned by Subscribe) and closes it.
//
// HINT: range the map and match with (<-chan string)(k) == ch, then delete+close.
func (b *Broker) Unsubscribe(ch <-chan string) {
	panic("TODO: implement Unsubscribe")
}

// Publish delivers msg to all current subscribers, non-blocking per subscriber.
//
// HINT: for each ch use select { case ch <- msg: default: } so a full buffer drops.
func (b *Broker) Publish(msg string) {
	panic("TODO: implement Publish")
}

// Close removes and closes all subscriptions; further Publish is a no-op.
//
// HINT: set a closed flag (Publish checks it) before closing every channel.
func (b *Broker) Close() {
	panic("TODO: implement Close")
}
