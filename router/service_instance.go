package router

import (
	"net"
)

type serviceInstance struct {
	addr *net.TCPAddr
	name string
}

func newServiceInstance(name, addrStr string) (*serviceInstance, error) {
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}
	si := &serviceInstance{
		name: name,
		addr: addr,
	}
	return si, nil
}
