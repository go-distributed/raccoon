package controller

import (
	"fmt"
	"net/rpc"
	"testing"
	"time"

	"github.com/go-distributed/raccoon/router"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf
var _ = router.NewInstance
var _ = assert.Nil

func TestRegRouterRPC(t *testing.T) {
	r, err := router.New(":14817")
	if err != nil {
		t.Fatal(err)
	}

	err = r.Start()
	if err != nil {
		t.Fatal("router start:", err)
	}
	defer func() {
		r.Stop()
		time.Sleep(time.Millisecond * 50)
	}()

	cAddr := "127.0.0.1:14818"
	c, err := New(cAddr)
	if err != nil {
		t.Fatal(err)
	}

	err = c.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer c.Stop()

	client, err := rpc.Dial("tcp", cAddr)
	if err != nil {
		t.Fatal("dialing:", err)
	}

	regRouterArgs := &RegRouterArgs{
		Id:   "test router",
		Addr: ":14817",
	}

	assert.Empty(t, c.routers)

	err = client.Call("ControllerRPC.RegisterRouter", regRouterArgs, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(c.routers), 1)
	assert.NotNil(t, c.routers["test router"])
}

func TestRegInstanceRPC(t *testing.T) {
}
