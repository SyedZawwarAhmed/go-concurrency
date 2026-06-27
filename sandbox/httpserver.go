package sandbox

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"
)

// Server is a real HTTP test server (loopback) that adds latency to every
// response and tracks the peak number of requests it handled simultaneously.
// Use MaxConcurrent in a test to prove a client really bounded its concurrency.
// Remember to Close it.
type Server struct {
	*httptest.Server
	inflight atomic.Int32
	maxSeen  atomic.Int32
}

// NewServer starts a Server that delays each response by latency and echoes the
// request path in the body as "ok:<path>".
func NewServer(latency time.Duration) *Server {
	s := &Server{}
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := s.inflight.Add(1)
		defer s.inflight.Add(-1)
		for {
			m := s.maxSeen.Load()
			if n <= m || s.maxSeen.CompareAndSwap(m, n) {
				break
			}
		}
		if latency > 0 {
			select {
			case <-time.After(latency):
			case <-r.Context().Done():
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "ok:"+r.URL.Path)
	}))
	return s
}

// MaxConcurrent reports the highest number of requests handled at the same time.
func (s *Server) MaxConcurrent() int { return int(s.maxSeen.Load()) }
