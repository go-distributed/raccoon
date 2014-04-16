package instance

import (
	"net"
)

type Instance struct {
	Addr    *net.TCPAddr
	Name    string
	Service string
}

func NewInstance(name, service, addrStr string) (*Instance, error) {
	addr, err := net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}

	return &Instance{
		Name:    name,
		Service: service,
		Addr:    addr,
	}, nil
}
