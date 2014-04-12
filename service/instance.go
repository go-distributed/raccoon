package service

type Instance struct {
	Addr    string
	Name    string
	Service string
}

func NewInstance(name, service, addr string) *Instance {
	return &Instance{
		Name:    name,
		Service: service,
		Addr:    addr,
	}
}
