package router

import (
	"net"
)

type service struct {
	name    string
	policy  string
	hosts   []*host
	proxies []*proxy
}
