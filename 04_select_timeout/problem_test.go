package selecttimeout

import (
	"errors"
	"testing"
	"time"
)

// sendAfter returns a channel that delivers v after d. Buffered so the sender
// goroutine never leaks if no one receives.
func sendAfter(v string, d time.Duration) <-chan string {
	ch := make(chan string, 1)
	go func() {
		time.Sleep(d)
		ch <- v
	}()
	return ch
}

func TestFetchResilient(t *testing.T) {
	const timeout = 150 * time.Millisecond

	tests := []struct {
		name    string
		primary <-chan string
		replica <-chan string
		want    string
		wantErr error
	}{
		{
			name:    "primary wins",
			primary: sendAfter("primary-data", 10*time.Millisecond),
			replica: sendAfter("replica-data", 100*time.Millisecond),
			want:    "primary-data",
		},
		{
			name:    "replica wins",
			primary: sendAfter("primary-data", 100*time.Millisecond),
			replica: sendAfter("replica-data", 10*time.Millisecond),
			want:    "replica-data",
		},
		{
			name:    "both timeout",
			primary: sendAfter("primary-data", 500*time.Millisecond),
			replica: sendAfter("replica-data", 500*time.Millisecond),
			wantErr: ErrTimeout,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := FetchResilient(tc.primary, tc.replica, timeout)

			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
