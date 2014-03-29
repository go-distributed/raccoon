package router

type routePolicyType uint8

const (
	randomSelect routePolicyType = iota + 1
	rRLoadBalance
	weightedLoadBalance
	prioritySelect
)

type RoutePolicy interface {
	// getType returns the type of the policy
	getType() routePolicyType
	// getContents returns the information used to help setup
	// the selector according to the policy. E.g. priority select
	// policy will tell the selector a list of prioritized choices.
	getContents() interface{}
}

type simplePolicy struct {
	Type routePolicyType
}

func (p *simplePolicy) getType() routePolicyType {
	return p.Type
}

func (p *simplePolicy) getContents() interface{} {
	return nil
}

// *************************
// ****  RANDOM SELECT  ****
// *************************

func NewRandomSelectPolicy() *simplePolicy {
	return &simplePolicy{randomSelect}
}

// *************************
// **** RANDOM SELECTOR ****
// *************************

func NewRRLoadBalancePolicy() *simplePolicy {
	return &simplePolicy{rRLoadBalance}
}
