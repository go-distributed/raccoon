package router

import (
	"net"
)

type proxy struct {
	connectors []*connector
	localAddr  *net.TCPAddr
	remoteAddr *net.TCPAddr
	listener   net.Listener
}

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
	}

	return p, nil
}

func (p *proxy) start() error {
	var err error
	p.listener, err = net.Listen("tcp", p.localAddr.String())
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
			c.connect()
		}(one)

	}
}

func (p *proxy) stop() {

}
