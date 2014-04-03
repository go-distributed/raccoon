package router

type RoutePolicyType uint8

const (
	RandomSelect RoutePolicyType = iota + 1
	RoundRobin                   // round robin load balance
	Weighting                    // weighted load balance
	PrioritizedSelect
)

type RoutePolicy interface {
	// ptype returns the type of the policy
	Type() RoutePolicyType

	// TODO: Define contents
	// contents returns the information used to help setup
	// the selector according to the policy. E.g. priority select
	// policy will give the selector a list of prioritized choices.
}

type SimplePolicy struct {
	PolicyType RoutePolicyType
}

func (p *SimplePolicy) Type() RoutePolicyType {
	return p.PolicyType
}

func NewRandomSelectPolicy() *SimplePolicy {
	return &SimplePolicy{RandomSelect}
}

func NewRoundRobinPolicy() *SimplePolicy {
	return &SimplePolicy{RoundRobin}
}
