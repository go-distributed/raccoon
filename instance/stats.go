package instance

import (
	"sync"
)

type Stats struct {
	TotalCount   int
	CurrentCount int
	l            sync.Mutex
}

func (s *Stats) IncTotal() {
	s.l.Lock()
	defer s.l.Unlock()

	s.TotalCount++
}

// Current Count is used to keep records of current TCP connection.
// This is useful for long session such as Database.
func (s *Stats) IncCurr(d int) {
	s.l.Lock()
	defer s.l.Unlock()

	s.CurrentCount = s.CurrentCount + d
}
