package router

import (
	"net"
)

type Instance struct {
	addr *net.TCPAddr
	name string
	cpu  int
}

func NewInstance(name, addrStr string) (*Instance, error) {
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}
	instance := &Instance{
		name: name,
		addr: addr,
	}
	return instance, nil
}
