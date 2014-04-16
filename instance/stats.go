package instance

import (
	"sync"
)

type Stats struct {
	totalCount   int
	currentCount int
	sync.Mutex
}

func (s *Stats) IncTotal() {
	s.Lock()
	defer s.Unlock()

	s.totalCount++
}

func (s *Stats) IncCurr(d int) {
	s.Lock()
	defer s.Unlock()

	s.currentCount = s.currentCount + d
}
