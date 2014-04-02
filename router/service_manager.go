package router

import (
	"fmt"
	"net"
	"sync"
)

type serviceManager struct {
	instances []*Instance
	selector
	sync.RWMutex
}

func newServiceManager(selector selector) (*serviceManager, error) {
	return &serviceManager{
		instances: make([]*Instance, 0),
		selector:  selector,
	}, nil
}

func (sm *serviceManager) addInstance(instance *Instance) error {
	sm.Lock()
	defer sm.Unlock()

	if sm.isInstanceExist(instance) {
		return fmt.Errorf("Instance %s already exists", instance.name)
	}

	sm.instances = append(sm.instances, instance)
	return nil
}

func (sm *serviceManager) removeInstance(instance *Instance) error {
	sm.Lock()
	defer sm.Unlock()

	if !sm.isInstanceExist(instance) {
		return fmt.Errorf("Instance %s does not exist", instance.name)
	}

	newInstances := make([]*Instance, len(sm.instances)-1)
	i := 0
	for _, ours := range sm.instances {
		if ours.name != instance.name {
			newInstances[i] = ours
			i++
		}
	}

	sm.instances = newInstances
	return nil
}

func (sm *serviceManager) selectServiceAddr() (*net.TCPAddr, error) {
	sm.RLock()
	defer sm.RUnlock()

	raddr, err := sm.doSelection(sm.instances)
	if err != nil {
		return nil, err
	}

	return raddr, nil
}

func (sm *serviceManager) isInstanceExist(instance *Instance) bool {
	for _, ours := range sm.instances {
		if ours.name == instance.name {
			return true
		}
	}
	return false
}
