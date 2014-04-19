package controller

import "github.com/go-distributed/raccoon/instance"

type Status int

type ControllerRPC struct {
	controller *Controller
}

type RegRouterArgs struct {
	Id   string
	Addr string
}

type RegRouterReply struct {
}

type RegInstanceArgs struct {
	Instance *instance.Instance
}

type RegInstanceReply struct {
}

type ReportFailureArgs struct {
	Reporter string
	Instance *instance.Instance
}

func newControllerRPC(controller *Controller) *ControllerRPC {
	return &ControllerRPC{
		controller: controller,
	}
}

func (rpc *ControllerRPC) RegisterRouter(args *RegRouterArgs, reply *RegRouterReply) error {
	cr, err := NewCRouter(args.Id, args.Addr)
	if err != nil {
		return err
	}
	return rpc.controller.RegisterRouter(cr)
}

func (rpc *ControllerRPC) RegisterServiceInstance(args *RegInstanceArgs, reply *RegInstanceReply) error {
	return rpc.controller.RegisterServiceInstance(args.Instance)
}

func (rpc *ControllerRPC) ReportFailure(args *ReportFailureArgs, reply *struct{}) error {
	rpc.controller.ReportFailure(args.Reporter, args.Instance)
	return nil
}
