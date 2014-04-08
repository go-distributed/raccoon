package controller

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
	return "NewRouterEvent"
}
