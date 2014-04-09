package controller

import "github.com/go-distributed/raccoon/router"

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
	Instance *router.Instance
}

type RegInstanceReply struct {
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
