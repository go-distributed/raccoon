package router

import (
	"fmt"
	"net"
	"sync"
)

type serviceManager struct {
	serviceInstances []*serviceInstance
	selector         selector
	sync.RWMutex
}

func newServiceManager(localAddr string, selector selector) (*serviceManager, error) {
	sm := &serviceManager{
		serviceInstances: make([]*serviceInstance, 0),
		selector:         selector,
	}
	return sm, nil
}

func (sm *serviceManager) addServiceInstance(name, addrStr string) error {
	sm.Lock()
	defer sm.Unlock()

	if sm.isInstanceExist(name) {
		return fmt.Errorf("%s already exists", name)
	}

	si, err := newServiceInstance(name, addrStr)
	if err != nil {
		return err
	}

	sm.serviceInstances = append(sm.serviceInstances, si)
	return nil
}

func (sm *serviceManager) removeServiceInstance(name string) error {
	sm.Lock()
	defer sm.Unlock()

	if !sm.isInstanceExist(name) {
		return fmt.Errorf("%s does not exist", name)
	}

	newInstances := make([]*serviceInstance, len(sm.serviceInstances)-1)
	i := 0
	for _, si := range sm.serviceInstances {
		if si.name != name {
			newInstances[i] = si
			i++
		}
	}
	return nil
}

func (sm *serviceManager) selectServiceAddr() (*net.TCPAddr, error) {
	sm.RLock()
	defer sm.RUnlock()

	raddr, err := sm.selector(sm.serviceInstances)
	if err != nil {
		return nil, err
	}

	return raddr, nil
}

func (sm *serviceManager) isInstanceExist(name string) bool {
	for i := range sm.serviceInstances {
		if sm.serviceInstances[i].name == name {
			return true
		}
	}
	return false
}
