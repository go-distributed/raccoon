package router

import (
	"fmt"
	"net"
	"sync"

	rmtService "github.com/go-distributed/raccoon/service"
)

type service struct {
	name      string
	policy    Policy
	proxy     *proxy
	instances []*rmtService.Instance
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
		instances: make([]*rmtService.Instance, 0),
		selector:  selector,
	}

	s.proxy, err = newProxy(localAddr, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *service) addInstance(instance *rmtService.Instance) error {
	s.Lock()
	defer s.Unlock()

	if s.isInstanceExist(instance) {
		return fmt.Errorf("rmtService.Instance %s already exists", instance.Name)
	}

	s.instances = append(s.instances, instance)
	return nil
}

func (s *service) removeInstance(instance *rmtService.Instance) error {
	s.Lock()
	defer s.Unlock()

	if !s.isInstanceExist(instance) {
		return fmt.Errorf("rmtService.Instance %s does not exist", instance.Name)
	}

	newInstances := make([]*rmtService.Instance, len(s.instances)-1)
	i := 0
	for _, ours := range s.instances {
		if ours.Name != instance.Name {
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

func (s *service) isInstanceExist(instance *rmtService.Instance) bool {
	for _, ours := range s.instances {
		if ours.Name == instance.Name {
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
