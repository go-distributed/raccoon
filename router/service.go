package router

type service struct {
	name    string
	policy  string
	proxies []*proxy
}
