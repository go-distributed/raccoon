package router

import (
	"fmt"
	"math/rand"
	"net"
)

type selector interface {
	doSelection(instances []*Instance) (*net.TCPAddr, error)
}

func newSelector(policy RoutePolicy) (selector, error) {
	var sel selector
	switch policy.Type() {
	case RandomSelect:
		sel = new(randomSelector)
	case RoundRobin:
		sel = new(roundRounbinSelector)
	default:
		return nil, fmt.Errorf("unknown route policy")
	}
	return sel, nil
}

type randomSelector struct{}

// defaultSelector randomly select a remote address from the proxy
// remote address list.
func (s *randomSelector) doSelection(instances []*Instance) (*net.TCPAddr, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("No service instance exists")
	}
	which := rand.Int() % len(instances)
	return instances[which].Addr, nil
}

type roundRounbinSelector struct {
	counter uint32
}

func newRoundbinSelector() *roundRounbinSelector {
	return new(roundRounbinSelector)
}

func (s *roundRounbinSelector) doSelection(instances []*Instance) (*net.TCPAddr, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("No service instance exists")
	}

	which := s.counter % uint32(len(instances))
	s.counter++
	return instances[which].Addr, nil
}
