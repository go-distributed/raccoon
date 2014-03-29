package router

import (
	"fmt"
	"math/rand"
	"net"
)

type selector interface {
	doSelection(servInstances []*serviceInstance) (*net.TCPAddr, error)
}

type defaultSelector struct{}

// defaultSelector randomly select a remote address from the proxy
// remote address list.
func (s *defaultSelector) doSelection(servInstances []*serviceInstance) (*net.TCPAddr, error) {
	if len(servInstances) == 0 {
		return nil, fmt.Errorf("No service instance exists")
	}
	which := rand.Int() % len(servInstances)
	return servInstances[which].addr, nil
}

type roundRounbinSelector struct {
	counter uint32
}

func newRoundbinSelector() *roundRounbinSelector {
	return new(roundRounbinSelector)
}

func (s *roundRounbinSelector) doSelection(serviceInstances []*serviceInstance) (*net.TCPAddr, error) {
	if len(serviceInstances) == 0 {
		return nil, fmt.Errorf("No service instance exists")
	}

	which := s.counter % uint32(len(serviceInstances))
	s.counter++
	return serviceInstances[which].addr, nil
}
