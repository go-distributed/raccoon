package router

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type Router struct {
	services map[string]*service
	listener net.Listener
	sync.Mutex
}

func New() (*Router, error) {
	r := &Router{
		services: make(map[string]*service),
	}
	return r, nil
}

func (r *Router) Start() (err error) {
	if err = rpc.Register(newRouterRPC(r)); err != nil {
		return
	}

	rpc.HandleHTTP()
	r.listener, err = net.Listen("tcp", ":14817")
	if err != nil {
		return
	}
	go http.Serve(r.listener, nil)

	return
}

func (r *Router) Stop() error {
	return r.listener.Close()
}

func (r *Router) AddService(sName, localAddr string, policy routePolicy) error {
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

func (r *Router) RemoveService(sName string) error {
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

func (r *Router) SetServicePolicy(sName string, policy routePolicy) error {
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

func (r *Router) AddServiceInstance(sName string, instance *Instance) error {
	r.Lock()
	s, ok := r.services[sName]
	r.Unlock()

	if !ok {
		return fmt.Errorf("service '%s' does not exist", sName)
	}

	err := s.manager.addInstance(instance)
	if err != nil {
		return err
	}

	return nil
}

func (r *Router) RemoveServiceInstance(sName string, instance *Instance) error {
	r.Lock()
	s, ok := r.services[sName]
	r.Unlock()

	if !ok {
		return fmt.Errorf("service '%s' does not exist", sName)
	}

	err := s.manager.removeInstance(instance)
	if err != nil {
		return err
	}

	return nil
}
