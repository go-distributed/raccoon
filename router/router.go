package router

import "fmt"

type Router struct {
	services map[string]*service
}

func NewRouter() (*Router, error) {
	r := &Router{
		make(map[string]*service),
	}
	return r, nil
}

func (r *Router) CreateService(sName, localAddr string, policy routePolicy, remoteAddrs ...string) error {
	_, ok := r.services[sName]
	if ok {
		return fmt.Errorf("name '%s' used already", sName)
	}

	s, err := newService(sName, policy, localAddr)
	if err != nil {
		return err
	}

	for _, addr := range remoteAddrs {
		s.manager.addServiceInstance("", addr)
	}

	go s.start()

	return nil
}

func (r *Router) DeleteService(sName string) error {
	panic("Not Implemented")
}

func (r *Router) SetServicePolicy(sName string, policy routePolicy) error {
	panic("Not Implemented")
}

func (r *Router) AddServiceMapping(sName, remoteAddr, siName string) error {
	panic("Not Implemented")
}
