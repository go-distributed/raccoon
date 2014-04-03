package router

type Status int

const (
	Ok Status = iota + 1
	NotOk
)

type RouterRPC struct {
	router *Router
}

type Reply struct {
	Value string
}

type ServiceArgs struct {
	ServiceName string
	LocalAddr   string
	Policy      RoutePolicy
}

type ServiceReply struct {
	Status Status
}

type InstanceArgs struct {
	ServiceName string
	Instance    *Instance
}

type InstanceReply struct {
	Status Status
}

func newRouterRPC(router *Router) *RouterRPC {
	return &RouterRPC{
		router: router,
	}
}

func (rpc *RouterRPC) Echo(arg string, reply *Reply) error {
	reply.Value = arg
	return nil
}

func (rpc *RouterRPC) AddService(args *ServiceArgs, reply *ServiceReply) error {
	err := rpc.router.AddService(args.ServiceName, args.LocalAddr, args.Policy)
	if err != nil {
		reply.Status = NotOk
	}

	return err
}

func (rpc *RouterRPC) RemoveService(args *ServiceArgs, reply *ServiceReply) error {
	err := rpc.router.RemoveService(args.ServiceName)
	if err != nil {
		reply.Status = NotOk
	}

	return err
}

func (rpc *RouterRPC) SetServicePolicy(args *ServiceArgs, reply *ServiceReply) error {
	err := rpc.router.SetServicePolicy(args.ServiceName, args.Policy)
	if err != nil {
		reply.Status = NotOk
	}

	return err
}

func (rpc *RouterRPC) AddServiceInstance(args *InstanceArgs, reply *InstanceReply) error {
	err := rpc.router.AddServiceInstance(args.ServiceName, args.Instance)
	if err != nil {
		reply.Status = NotOk
	}

	return err
}

func (rpc *RouterRPC) RemoveServiceInstance(args *InstanceArgs, reply *InstanceReply) error {
	err := rpc.router.RemoveServiceInstance(args.ServiceName, args.Instance)
	if err != nil {
		reply.Status = NotOk
	}

	return err
}
