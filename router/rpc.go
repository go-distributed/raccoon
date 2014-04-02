package router

type RouterRPC struct {
	router *Router
}

type Reply struct {
	Value string
}

func newRouterRPC(router *Router) *RouterRPC {
	return &RouterRPC{
		router: router,
	}
}

func (rr *RouterRPC) Echo(arg string, reply *Reply) error {
	reply.Value = arg
	return nil
}
