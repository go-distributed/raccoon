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

func (s *Stats) IncCurr(d int) {
	s.l.Lock()
	defer s.l.Unlock()

	s.CurrentCount = s.CurrentCount + d
}
