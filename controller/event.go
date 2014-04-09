package controller

import "github.com/go-distributed/raccoon/service"

type event interface {
	Type() string
}

type AddRouterEvent struct {
	ID   string
	Addr string
}

func NewAddRouterEvent(id, addr string) *AddRouterEvent {
	return &AddRouterEvent{
		ID:   id,
		Addr: addr,
	}
}

func (e *AddRouterEvent) Type() string {
	return "AddRouterEvent"
}

type AddInstanceEvent struct {
	Name    string
	Service string
	Addr    string
}

func NewAddInstanceEvent(i *service.Instance) *AddInstanceEvent {
	return &AddInstanceEvent{
		Name:    i.Name,
		Service: i.Service,
		Addr:    i.Addr.String(),
	}
}

func (e *AddInstanceEvent) Type() string {
	return "AddInstanceEvent"
}
