package router

import (
	"fmt"
	"net"
	"sync"
)

type serviceManager struct {
	proxy            *proxy
	serviceInstances []*serviceInstance
	selector         selector
	sync.RWMutex
}

func newServiceManager(localAddr string, raddrStrs []string, selector selector) (*serviceManager, error) {
	var err error

	if len(raddrStrs) == 0 {
		return nil, fmt.Errorf("no remote address is given")
	}

	sm := &serviceManager{
		serviceInstances: make([]*serviceInstance, 0),
		selector:         selector,
	}

	sm.proxy, err = newProxy(localAddr, sm)
	if err != nil {
		return nil, err
	}

	for _, raddrStr := range raddrStrs {
		if err = sm.addServiceInstance(raddrStr); err != nil {
			return nil, err
		}
	}
	return sm, nil
}

func (sm *serviceManager) addServiceInstance(addrStr string) error {
	sm.Lock()
	defer sm.Unlock()

	servInst, err := newServiceInstance(addrStr)
	if err != nil {
		return err
	}

	sm.serviceInstances = append(sm.serviceInstances, servInst)

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
