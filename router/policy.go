package router

type routePolicyType uint8

const (
	randomSelect routePolicyType = iota + 1
	roundRobinLB                 // round robin load balance
	weightedLB                   // weighted load balance
	prioritizedSelect
)

type routePolicy interface {
	// ptype returns the type of the policy
	pType() routePolicyType

	// TODO: Define contents
	// contents returns the information used to help setup
	// the selector according to the policy. E.g. priority select
	// policy will give the selector a list of prioritized choices.
}

type simplePolicy struct {
	policyType routePolicyType
}

func (p *simplePolicy) pType() routePolicyType {
	return p.policyType
}

func NewRandomSelectPolicy() *simplePolicy {
	return &simplePolicy{randomSelect}
}

func NewRoundRobinLBPolicy() *simplePolicy {
	return &simplePolicy{roundRobinLB}
}
