package controller

import (
	"sync"
)

type dispatcher struct {
	listeners map[string][]EventListener
	sync.RWMutex
}

func newDispatcher() *dispatcher {
	d := &dispatcher{
		listeners: make(map[string][]EventListener),
	}

	return d
}

func (d *dispatcher) addListener(typ string, listener EventListener) {
	d.Lock()
	defer d.Unlock()

	listeners := d.listeners[typ]
	d.listeners[typ] = append(listeners, listener)
}

func (d *dispatcher) dispatch(e Event) {
	d.RLock()
	defer d.RUnlock()

	listeners := d.listeners[e.Type()]
	for _, l := range listeners {
		l(e)
	}
}
