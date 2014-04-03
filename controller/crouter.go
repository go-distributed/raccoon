package controller

import (
	"encoding/gob"
	"net"
	"net/rpc"

	"github.com/go-distributed/raccoon/router"
)

type CRouter struct {
	client *rpc.Client
}

func NewCRouter(addr string) (*CRouter, error) {
	client, err := rpc.DialHTTP("tcp", net.JoinHostPort(addr, "14817"))
	if err != nil {
		return nil, err
	}

	// register policies for gob encoding in RPC calls
	gob.Register(new(router.SimplePolicy))

	return &CRouter{
		client: client,
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
