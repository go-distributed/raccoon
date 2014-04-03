package router

import (
	"net"
)

type Instance struct {
	Addr *net.TCPAddr
	Name string
}

func NewInstance(name, addrStr string) (*Instance, error) {
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}
	instance := &Instance{
		Name: name,
		Addr: addr,
	}
	return instance, nil
}
