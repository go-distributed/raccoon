package controller

import (
	"net/rpc"
	"testing"
	"time"

	"github.com/go-distributed/raccoon/router"
	"github.com/go-distributed/raccoon/service"
	"github.com/stretchr/testify/assert"
)

func TestRPC(t *testing.T) {
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

	ins, err := service.NewInstance("test instance", "test service", ":8888")
	if err != nil {
		t.Fatal(err)
	}

	regInstanceArgs := &RegInstanceArgs{
		Instance: ins,
	}

	assert.Empty(t, c.serviceInstances)

	err = client.Call("ControllerRPC.RegisterServiceInstance", regInstanceArgs, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(c.serviceInstances), 1)
	assert.NotNil(t, c.serviceInstances["test service"])
	assert.Equal(t, c.serviceInstances["test service"][0], ins)
}
