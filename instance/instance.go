package instance

import (
	"log"
	"net"
	"net/rpc"
)

type Instance struct {
	Addr    *net.TCPAddr
	Name    string
	Service string
	Stats   *Stats
}

type RegInstanceArgs struct {
	Instance *Instance
}

type RegInstanceReply struct {
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
		Stats:   new(Stats),
	}, nil
}

func (i *Instance) NewStats() {
	i.Stats = new(Stats)
}

func (i *Instance) RegisterOnCtler(ctlAddr string) error {
	c, err := rpc.Dial("tcp", ctlAddr)
	if err != nil {
		log.Println(err)
		return err
	}
	defer c.Close()

	args := &RegInstanceArgs{
		Instance: i,
	}
	err = c.Call("ControllerRPC.RegisterServiceInstance", args, nil)

	return err
}
