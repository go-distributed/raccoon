package router

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strings"
	"sync"

	rmtService "github.com/go-distributed/raccoon/service"
)

type Router interface {
	AddService(sName, localAddr string, policy Policy) error
	RemoveService(sName string) error
	SetServicePolicy(sName string, policy Policy) error
	// TODO: remove first argument 'sName'.
	// instance contains service name already.
	AddServiceInstance(sName string, instance *rmtService.Instance) error
	RemoveServiceInstance(sName string, instance *rmtService.Instance) error
}

type router struct {
	services map[string]*service
	listener net.Listener
	addr     *net.TCPAddr
	server   *rpc.Server
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
	s := rpc.NewServer()

	if err = s.Register(newRouterRPC(r)); err != nil {
		return
	}

	r.listener, err = net.ListenTCP("tcp", r.addr)
	if err != nil {
		return
	}

	go func() {
		closeErr := "use of closed network connection"
		for {
			if conn, err := r.listener.Accept(); err != nil {
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

func (r *router) Stop() error {
	return r.listener.Close()
}

func (r *router) AddService(sName, localAddr string, policy Policy) error {
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

func (r *router) SetServicePolicy(sName string, policy Policy) error {
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

func (r *router) AddServiceInstance(sName string, instance *rmtService.Instance) error {
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

func (r *router) RemoveServiceInstance(sName string, instance *rmtService.Instance) error {
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
