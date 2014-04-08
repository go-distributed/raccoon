package controller

import (
	"github.com/go-distributed/raccoon/router"
	"sync"
)

type Controller struct {
	serviceInstances map[string][]*router.Instance
	routers          map[string]*router.Router
	sync.RWMutex
}

func New() *Controller {
	c := &Controller{
		serviceInstances: make(map[string][]*router.Instance),
		routers:          make(map[string]*router.Router),
	}

	return c
}
