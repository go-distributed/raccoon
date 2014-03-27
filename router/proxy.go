package router

import (
	"io"
)

type proxy struct {
	one   io.ReadWriteCloser
	other io.ReadWriteCloser
}

func newProxy(one io.ReadWriteCloser, other io.ReadWriteCloser) *proxy {
	return &proxy{
		one:   one,
		other: other,
	}
}

func (p *proxy) start() {
	go io.Copy(p.one, p.other)
	go io.Copy(p.other, p.one)
}

func (p *proxy) stop() {
	p.one.Close()
	p.other.Close()
}
