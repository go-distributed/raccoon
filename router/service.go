package router

import (
	"fmt"
	"net"
	"sync"

	rmtService "github.com/go-distributed/raccoon/service"
)

// TODO: discussion needed
// When app calls AddService for router, how does it know which port to use?
var ServicePortMap map[string]string

func init() {
	ServicePortMap = make(map[string]string)
	ServicePortMap["test service"] = ":8080"
}

type instance struct {
	remote  *rmtService.Instance
	netAddr *net.TCPAddr
}

type service struct {
	name      string
	policy    Policy
	proxy     *proxy
	instances []*instance

	selector
	sync.RWMutex
}

func newInstance(remote *rmtService.Instance) (*instance, error) {
	netAddr, err := net.ResolveTCPAddr("tcp", remote.Addr)
	if err != nil {
		return nil, err
	}

	ins := &instance{
		remote:  remote,
		netAddr: netAddr,
	}
	return ins, nil
}

func newService(name string, localAddr string, policy Policy) (s *service, err error) {
	selector, err := newSelector(policy)
	if err != nil {
		return nil, err
	}

	s = &service{
		name:      name,
		instances: make([]*instance, 0),
		selector:  selector,
	}

	s.proxy, err = newProxy(localAddr, s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *service) addInstance(remote *rmtService.Instance) error {
	s.Lock()
	defer s.Unlock()

	if s.isInstanceExist(remote) {
		return fmt.Errorf("Instance %s already exists", remote.Name)
	}

	ins, err := newInstance(remote)
	if err != nil {
		return err
	}

	s.instances = append(s.instances, ins)
	return nil
}

func (s *service) removeInstance(remote *rmtService.Instance) error {
	s.Lock()
	defer s.Unlock()

	if s.isInstanceExist(remote) {
		return fmt.Errorf("Instance %s already exists", remote.Name)
	}

	newInstances := make([]*instance, len(s.instances)-1)
	i := 0
	for _, ours := range s.instances {
		if ours.remote.Name != remote.Name {
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

func (s *service) isInstanceExist(rmtInstance *rmtService.Instance) bool {
	for _, ours := range s.instances {
		if ours.remote.Name == rmtInstance.Name {
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
