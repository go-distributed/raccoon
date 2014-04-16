package router

import (
	"fmt"
	"math/rand"

	"github.com/go-distributed/raccoon/instance"
)

type selector interface {
	doSelection(instances []*instance.Instance) (*instance.Instance, error)
}

func newSelector(policy Policy) (selector, error) {
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
func (s *randomSelector) doSelection(instances []*instance.Instance) (*instance.Instance, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("No service instance exists")
	}
	which := rand.Int() % len(instances)
	return instances[which], nil
}

type roundRounbinSelector struct {
	counter uint32
}

func newRoundbinSelector() *roundRounbinSelector {
	return new(roundRounbinSelector)
}

func (s *roundRounbinSelector) doSelection(instances []*instance.Instance) (*instance.Instance, error) {
	if len(instances) == 0 {
		return nil, fmt.Errorf("No service instance exists")
	}

	which := s.counter % uint32(len(instances))
	s.counter++
	return instances[which], nil
}
