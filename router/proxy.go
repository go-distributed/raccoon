package router

import (
	"fmt"
	"net"
	"sync"
)

type proxy struct {
	status     string
	connectors []*connector
	localAddr  *net.TCPAddr
	remoteAddr *net.TCPAddr
	listener   net.Listener
	sync.Mutex
}

const (
	initialized = "initialized"
	running     = "running"
	stopped     = "stopped"
)

func newProxy(laddrStr, raddrStr string) (*proxy, error) {
	laddr, err := net.ResolveTCPAddr("tcp", laddrStr)
	if err != nil {
		return nil, err
	}

	raddr, err := net.ResolveTCPAddr("tcp", raddrStr)
	if err != nil {
		return nil, err
	}

	p := &proxy{
		connectors: make([]*connector, 0),
		localAddr:  laddr,
		remoteAddr: raddr,
		status:     initialized,
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
			other, err := net.Dial("tcp", p.remoteAddr.String())
			if err != nil {
				return
			}

			c := newConnector(one, other)
			if err := p.addConnector(c); err != nil {
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
