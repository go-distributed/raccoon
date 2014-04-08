package controller

import (
	"fmt"
	"sync"

	"github.com/go-distributed/raccoon/router"
)

type Controller struct {
	serviceInstances map[string][]*router.Instance
	routers          map[string]router.Router

	// TODO:
	// 1. reader writer lock
	// 2. two locks for each map
	sync.RWMutex
}

func New() *Controller {
	c := &Controller{
		serviceInstances: make(map[string][]*router.Instance),
		routers:          make(map[string]router.Router),
	}

	return c
}

func (c *Controller) RegisterRouter(cr *CRouter) error {
	c.Lock()
	defer c.Unlock()
	_, ok := c.routers[cr.id]
	if ok {
		return fmt.Errorf("router '%s' already exists", cr.id)
	}

	c.routers[cr.id] = cr
	return nil
}

func (c *Controller) RegisterServiceInstance(ins *router.Instance) error {
	c.Lock()
	defer c.Unlock()
	instances := c.serviceInstances[ins.Service]

	for _, instance := range instances {
		if instance.Name == ins.Name {
			return fmt.Errorf("router '%s' already exists", ins.Name)
		}
	}

	c.serviceInstances[ins.Service] = append(instances, ins)
	return nil
}
