package controller

import (
	"encoding/gob"
	"net/rpc"

	"github.com/go-distributed/raccoon/router"
)

type CRouter struct {
	id     string
	client *rpc.Client
	addr   string
}

func NewCRouter(id, addr string) (*CRouter, error) {
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	// register policies for gob encoding in RPC calls
	gob.Register(new(router.SimplePolicy))

	return &CRouter{
		id:     id,
		client: client,
		addr:   addr,
	}, nil
}

func (cr *CRouter) AddService(sName, localAddr string, policy router.RoutePolicy) error {
	sArgs := &router.ServiceArgs{
		ServiceName: sName,
		LocalAddr:   localAddr,
		Policy:      policy,
	}
	return cr.client.Call("RouterRPC.AddService", sArgs, nil)
}

func (cr *CRouter) RemoveService(sName string) error {
	sArgs := &router.ServiceArgs{
		ServiceName: sName,
	}
	return cr.client.Call("RouterRPC.RemoveService", sArgs, nil)
}

func (cr *CRouter) SetServicePolicy(sName string, policy router.RoutePolicy) error {
	sArgs := &router.ServiceArgs{
		ServiceName: sName,
		Policy:      policy,
	}
	return cr.client.Call("RouterRPC.SetServicePolicy", sArgs, nil)
}

func (cr *CRouter) AddServiceInstance(sName string, instance *router.Instance) error {
	iArgs := &router.InstanceArgs{
		ServiceName: sName,
		Instance:    instance,
	}
	return cr.client.Call("RouterRPC.AddServiceInstance", iArgs, nil)
}

func (cr *CRouter) RemoveServiceInstance(sName string, instance *router.Instance) error {
	iArgs := &router.InstanceArgs{
		ServiceName: sName,
		Instance:    instance,
	}
	return cr.client.Call("RouterRPC.RemoveServiceInstance", iArgs, nil)
}
