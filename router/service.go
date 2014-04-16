package router

import (
	"fmt"
	"net"
	"sync"

	"github.com/go-distributed/raccoon/instance"
)

type service struct {
	name      string
	policy    Policy
	proxy     *proxy
	instances []*instance.Instance

	selector
	sync.RWMutex
}

func newService(name string, localAddr string, policy Policy) (s *service, err error) {
	selector, err := newSelector(policy)
	if err != nil {
		return nil, err
	}

	s = &service{
		name:      name,
		instances: make([]*instance.Instance, 0),
		selector:  selector,
	}

	s.proxy, err = newProxy(localAddr, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *service) addInstance(remote *instance.Instance) error {
	s.Lock()
	defer s.Unlock()

	if s.isInstanceExist(remote) {
		return fmt.Errorf("instance '%s' already exists", remote.Name)
	}

	s.instances = append(s.instances, remote)
	return nil
}

func (s *service) removeInstance(remote *instance.Instance) error {
	s.Lock()
	defer s.Unlock()

	if !s.isInstanceExist(remote) {
		return fmt.Errorf("instance '%s' does not exist", remote.Name)
	}

	newInstances := make([]*instance.Instance, len(s.instances)-1)
	i := 0
	for _, ours := range s.instances {
		if ours.Name != remote.Name {
			newInstances[i] = ours
			i++
		}
	}

	s.instances = newInstances
	return nil
}

func (s *service) selectInstanceAddr() (*net.TCPAddr, error) {
	s.RLock()
	defer s.RUnlock()

	raddr, err := s.doSelection(s.instances)
	if err != nil {
		return nil, err
	}

	return raddr, nil
}

func (s *service) isInstanceExist(ins *instance.Instance) bool {
	for _, ours := range s.instances {
		if ours.Name == ins.Name {
			return true
		}
	}
	return false
}

func (s *service) setPolicy(policy Policy) error {
	selector, err := newSelector(policy)
	if err != nil {
		return err
	}

	s.selector = selector
	s.policy = policy
	return nil
}

func (s *service) start() error {
	err := s.proxy.start()
	return err
}

func (s *service) stop() error {
	err := s.proxy.stop()
	return err
}
