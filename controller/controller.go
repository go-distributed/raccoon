package controller

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strings"
	"sync"

	"github.com/go-distributed/raccoon/instance"
	"github.com/go-distributed/raccoon/router"
)

type Controller struct {
	ServiceInstances map[string][]*instance.Instance
	Routers          map[string]router.Router
	dispatcher       *dispatcher
	listener         net.Listener
	addr             *net.TCPAddr
	// TODO:
	// 1. reader writer lock
	// 2. two locks for each map
	sync.RWMutex
}

func New(addrStr string) (*Controller, error) {
	c := &Controller{
		ServiceInstances: make(map[string][]*instance.Instance),
		Routers:          make(map[string]router.Router),
		dispatcher:       newDispatcher(),
	}

	var err error
	c.addr, err = net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Controller) Start() (err error) {
	s := rpc.NewServer()

	if err = s.Register(newControllerRPC(c)); err != nil {
		return
	}

	c.listener, err = net.ListenTCP("tcp", c.addr)
	if err != nil {
		return
	}

	go func() {
		closeErr := "use of closed network connection"
		for {
			if conn, err := c.listener.Accept(); err != nil {
				if strings.Contains(err.Error(), closeErr) {
					return
				}
				log.Fatal(err)
			} else {
				go s.ServeConn(conn)
			}
		}
	}()

	return
}

func (c *Controller) Stop() error {
	err := c.listener.Close()
	return err
}

// add a cRouter into router list and dispatch add-router-event
func (c *Controller) RegisterRouter(cr *CRouter) error {
	c.Lock()
	defer c.Unlock()

	_, ok := c.Routers[cr.id]
	if ok {
		return fmt.Errorf("router '%s' already exists", cr.id)
	}

	log.Println("Router registered:", cr.id, cr.addr)
	c.Routers[cr.id] = cr

	c.dispatcher.dispatch(NewAddRouterEvent(cr.id, cr.addr))
	return nil
}

func (c *Controller) RegisterServiceInstance(ins *instance.Instance) error {
	c.Lock()
	defer c.Unlock()

	instances := c.ServiceInstances[ins.Service]

	for _, instance := range instances {
		if instance.Name == ins.Name {
			return fmt.Errorf("router '%s' already exists", ins.Name)
		}
	}

	c.ServiceInstances[ins.Service] = append(instances, ins)

	c.dispatcher.dispatch(NewAddInstanceEvent(ins))
	return nil
}

func (c *Controller) AddListener(typ string, listener EventListener) {
	c.dispatcher.addListener(typ, listener)
}
