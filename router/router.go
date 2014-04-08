package router

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type Router interface {
	AddService(sName, localAddr string, policy RoutePolicy) error
	RemoveService(sName string) error
	SetServicePolicy(sName string, policy RoutePolicy) error
	AddServiceInstance(sName string, instance *Instance) error
	RemoveServiceInstance(sName string, instance *Instance) error
}

type router struct {
	services map[string]*service
	listener net.Listener
	addr     *net.TCPAddr
	sync.Mutex
}

func New(addrStr string) (*router, error) {
	r := &router{
		services: make(map[string]*service),
	}

	var err error
	r.addr, err = net.ResolveTCPAddr("tcp", addrStr)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *router) Start() (err error) {
	if err = rpc.Register(newRouterRPC(r)); err != nil {
		return
	}

	rpc.HandleHTTP()
	r.listener, err = net.ListenTCP("tcp", r.addr)
	if err != nil {
		return
	}
	go http.Serve(r.listener, nil)

	return
}

func (r *router) Stop() error {
	return r.listener.Close()
}

func (r *router) AddService(sName, localAddr string, policy RoutePolicy) error {
	r.Lock()
	defer r.Unlock()

	_, ok := r.services[sName]
	if ok {
		return fmt.Errorf("service '%s' already exists", sName)
	}

	s, err := newService(sName, localAddr, policy)
	if err != nil {
		return err
	}

	// TODO: handle error
	go s.start()

	r.services[sName] = s

	return nil
}

func (r *router) RemoveService(sName string) error {
	r.Lock()
	defer r.Unlock()

	s, ok := r.services[sName]
	if !ok {
		return fmt.Errorf("service '%s' does not exist", sName)
	}

	err := s.stop()
	if err != nil {
		return err
	}

	delete(r.services, sName)

	return nil
}

func (r *router) SetServicePolicy(sName string, policy RoutePolicy) error {
	r.Lock()
	s, ok := r.services[sName]
	r.Unlock()

	if !ok {
		return fmt.Errorf("service '%s' does not exist", sName)
	}

	err := s.setPolicy(policy)
	if err != nil {
		return err
	}

	return nil
}

func (r *router) AddServiceInstance(sName string, instance *Instance) error {
	r.Lock()
	s, ok := r.services[sName]
	r.Unlock()

	if !ok {
		return fmt.Errorf("service '%s' does not exist", sName)
	}

	err := s.addInstance(instance)
	if err != nil {
		return err
	}

	return nil
}

func (r *router) RemoveServiceInstance(sName string, instance *Instance) error {
	r.Lock()
	s, ok := r.services[sName]
	r.Unlock()

	if !ok {
		return fmt.Errorf("service '%s' does not exist", sName)
	}

	err := s.removeInstance(instance)
	if err != nil {
		return err
	}

	return nil
}
