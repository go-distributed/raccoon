package router

import (
	"net"
)

type serviceInstance struct {
	addr *net.TCPAddr
}

func newServiceInstance(addrStr string) (*serviceInstance, error) {
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}
	servInst := &serviceInstance{
		addr: addr,
	}
	return servInst, nil
}
