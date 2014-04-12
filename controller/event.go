package controller

import "github.com/go-distributed/raccoon/service"

const (
	AddRouterEventType   = "AddRouterEvent"
	AddInstanceEventType = "AddInstanceEvent"
)

type event interface {
	Type() string
}

type AddRouterEvent struct {
	Id   string
	Addr string
}

func NewAddRouterEvent(id, addr string) *AddRouterEvent {
	return &AddRouterEvent{
		Id:   id,
		Addr: addr,
	}
}

func (e *AddRouterEvent) Type() string {
	return AddRouterEventType
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
		Addr:    i.Addr,
	}
}

func (e *AddInstanceEvent) Type() string {
	return AddInstanceEventType
}
