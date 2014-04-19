package router

import (
	"github.com/go-distributed/raccoon/instance"
)

type Status int

const (
	Ok Status = iota + 1
	NotOk
)

type RouterRPC struct {
	router Router
}

type Reply struct {
	Value string
}

type ServiceArgs struct {
	ServiceName string
	LocalAddr   string
	Policy      Policy
}

type ServiceReply struct {
}

type InstanceArgs struct {
	Instance *instance.Instance
}

type InstanceReply struct {
}

func newRouterRPC(router Router) *RouterRPC {
	return &RouterRPC{
		router: router,
	}
}

func (rpc *RouterRPC) Echo(arg string, reply *Reply) error {
	reply.Value = arg
	return nil
}

func (rpc *RouterRPC) AddService(args *ServiceArgs, reply *ServiceReply) error {
	return rpc.router.AddService(args.ServiceName, args.LocalAddr, args.Policy)
}

func (rpc *RouterRPC) RemoveService(args *ServiceArgs, reply *ServiceReply) error {
	return rpc.router.RemoveService(args.ServiceName)
}

func (rpc *RouterRPC) SetServicePolicy(args *ServiceArgs, reply *ServiceReply) error {
	return rpc.router.SetServicePolicy(args.ServiceName, args.Policy)
}

func (rpc *RouterRPC) AddServiceInstance(args *InstanceArgs, reply *InstanceReply) error {
	return rpc.router.AddServiceInstance(args.Instance.Service, args.Instance)
}

func (rpc *RouterRPC) RemoveServiceInstance(args *InstanceArgs, reply *InstanceReply) error {
	return rpc.router.RemoveServiceInstance(args.Instance.Service, args.Instance)
}

func (rpc *RouterRPC) GetServiceInstances(args string, reply *[]*instance.Instance) error {
	is, err := rpc.router.GetServiceInstances(args)
	if err != nil {
		return err
	}
	reply = is

	return nil
}

type ReportFailureArgs struct {
	Reporter string
	Instance *instance.Instance
}
