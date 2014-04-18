package controller

import "github.com/go-distributed/raccoon/instance"

const (
	AddRouterEventType       = "AddRouterEvent"
	AddInstanceEventType     = "AddInstanceEvent"
	RmInstanceEventType      = "RemoveInstanceEvent"
	FailureInstanceEventType = "FailureInstanceEvent"
)

type Event interface {
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
	Instance *instance.Instance
}

func NewAddInstanceEvent(i *instance.Instance) *AddInstanceEvent {
	return &AddInstanceEvent{i}
}

func (e *AddInstanceEvent) Type() string {
	return AddInstanceEventType
}

type RmInstanceEvent struct {
	Instance *instance.Instance
}

func NewRmInstanceEvent(i *instance.Instance) *RmInstanceEvent {
	return &RmInstanceEvent{i}
}

func (e *RmInstanceEvent) Type() string {
	return RmInstanceEventType
}

type FailureInstanceEvnet struct {
	reporter string
	instance *instance.Instance
}

func NewFailureInstanceEvent(reporter string, i *instance.Instance) *FailureInstanceEvnet {
	return &FailureInstanceEvnet{reporter, i}
}

func (e *FailureInstanceEvnet) Type() string {
	return FailureInstanceEventType
}
