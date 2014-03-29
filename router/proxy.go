package router

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type proxy struct {
	status         string
	connectors     []*connector
	localAddr      *net.TCPAddr
	listener       net.Listener
	serviceManager *serviceManager
	sync.Mutex
}

const (
	initialized = "initialized"
	running     = "running"
	stopped     = "stopped"
)

func newProxy(laddrStr string, sm *serviceManager) (*proxy, error) {
	var err error

	p := &proxy{
		connectors:     make([]*connector, 0),
		status:         initialized,
		serviceManager: sm,
	}

	if err != nil {
		return nil, err
	}

	p.localAddr, err = net.ResolveTCPAddr("tcp", laddrStr)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *proxy) start() error {
	var err error

	p.Lock()
	if p.status != initialized {
		defer p.Unlock()
		return fmt.Errorf("the status of proxy is not initialized [%s]", p.status)
	}

	p.status = running
	p.listener, err = net.Listen("tcp", p.localAddr.String())
	p.Unlock()

	if err != nil {
		return err
	}

	for {
		one, err := p.listener.Accept()
		if err != nil {
			// handle error
			return err
		}

		go func(one net.Conn) {
			raddr, err := p.serviceManager.selectServiceAddr()

			if err != nil {
				log.Println(err)
				return
			}

			other, err := net.Dial("tcp", raddr.String())
			if err != nil {
				log.Println(err)
				return
			}

			c := newConnector(one, other)
			if err := p.addConnector(c); err != nil {
				log.Println(err)
				return
			}

			c.connect()
		}(one)
	}
}

// todo(xiangli) Graceful shutdown
func (p *proxy) stop() error {
	p.Lock()
	defer p.Unlock()

	if p.status != running {
		return fmt.Errorf("the status of proxy is not running [%s]", p.status)
	}

	p.status = stopped
	p.listener.Close()
	for _, c := range p.connectors {
		c.disconnect()
	}
	p.connectors = nil
	return nil
}

func (p *proxy) addConnector(c *connector) error {
	p.Lock()
	defer p.Unlock()

	if p.status != running {
		return fmt.Errorf("the status of proxy is not running [%s]", p.status)
	}

	p.connectors = append(p.connectors, c)
	return nil
}
