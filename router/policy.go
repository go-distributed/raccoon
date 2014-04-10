package router

type PolicyType uint8

const (
	RandomSelect PolicyType = iota + 1
	RoundRobin              // round robin load balance
	Weighting               // weighted load balance
	PrioritizedSelect
)

type Policy interface {
	// ptype returns the type of the policy
	Type() PolicyType

	// TODO: Define contents
	// contents returns the information used to help setup
	// the selector according to the policy. E.g. priority select
	// policy will give the selector a list of prioritized choices.
}

type SimplePolicy struct {
	PolicyType PolicyType
}

func (p *SimplePolicy) Type() PolicyType {
	return p.PolicyType
}

func NewRandomSelectPolicy() *SimplePolicy {
	return &SimplePolicy{RandomSelect}
}

func NewRoundRobinPolicy() *SimplePolicy {
	return &SimplePolicy{RoundRobin}
}
