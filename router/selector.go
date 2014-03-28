package router

import (
	"math/rand"
	"net"
)

type selector func(remoteAddrs []*net.TCPAddr) *net.TCPAddr

// defaultSelector randomly select a remote address from the proxy
// remote address list.
func defaultSelector(remoteAddrs []*net.TCPAddr) *net.TCPAddr {
	which := rand.Int() % len(remoteAddrs)
	return remoteAddrs[which]
}
