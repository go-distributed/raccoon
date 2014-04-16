package instance

import (
	"net"
)

type Instance struct {
	Addr    *net.TCPAddr
	Name    string
	Service string
	stats   *Stats
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
		stats:   new(Stats),
	}, nil
}

func (i *Instance) Stats() *Stats {
	return i.stats
}

func (i *Instance) NewStats() {
	i.stats = new(Stats)
}
