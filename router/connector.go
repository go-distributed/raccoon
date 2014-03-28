package router

import (
	"io"
)

type connector struct {
	one   io.ReadWriteCloser
	other io.ReadWriteCloser
}

func newConnector(one io.ReadWriteCloser, other io.ReadWriteCloser) *connector {
	return &connector{
		one:   one,
		other: other,
	}
}

func (c *connector) connect() {
	go io.Copy(c.one, c.other)
	go io.Copy(c.other, c.one)
}

func (c *connector) disconnect() {
	c.one.Close()
	c.other.Close()
}
