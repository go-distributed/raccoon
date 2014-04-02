package router

import (
	"net/rpc"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRPC(t *testing.T) {
	r, _ := New()
	err := r.Start()
	if err != nil {
		t.Fatal("router start:", err)
	}
	defer r.Stop()

	time.Sleep(time.Millisecond * 300)

	client, err := rpc.DialHTTP("tcp", ":14817")
	if err != nil {
		t.Fatal("dialing:", err)
	}

	reply := new(Reply)
	err = client.Call("RouterRPC.Echo", "hello router!", reply)
	if err != nil {
		t.Fatal("router rpc:", err)
	}

	assert.Equal(t, reply.Value, "hello router!")
}
