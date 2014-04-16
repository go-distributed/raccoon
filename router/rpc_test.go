package router

import (
	"encoding/gob"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"testing"
	"time"

	"github.com/go-distributed/raccoon/instance"
	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	routerAddr := "127.0.0.1:14817"

	r, _ := New(routerAddr)
	err := r.Start()
	if err != nil {
		t.Fatal("router start:", err)
	}
	defer r.Stop()

	time.Sleep(time.Millisecond * 50)

	client, err := rpc.Dial("tcp", ":14817")
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

func TestRPC(t *testing.T) {
	routerAddr := "127.0.0.1:14817"

	r, _ := New(routerAddr)
	err := r.Start()
	if err != nil {
		t.Fatal("router start:", err)
	}
	defer r.Stop()

	time.Sleep(time.Millisecond * 50)

	expectedReply, _ := genRandomBytesSlice(4096)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expectedReply)
	}))
	defer ts.Close()

	sName := "test service"
	localAddr := "127.0.0.1:8080"
	remoteAddr := ts.Listener.Addr().String()

	err = prepareRouterByRPC(routerAddr, sName, localAddr, remoteAddr)
	if err != nil {
		t.Fatal(err)
	}

	// testing service routing
	resp, err := http.Get("http://" + localAddr + "/")
	if err != nil {
		t.Fatal(err)
	}

	reply, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, reply, expectedReply)

	err = r.RemoveService(sName)
	if err != nil {
		t.Fatal(err)
	}
}

func prepareRouterByRPC(routerAddr, sName, localAddr, remoteAddr string) error {
	mapTo, err := instance.NewInstance("test instance", "test service", remoteAddr)
	if err != nil {
		return err
	}

	//client, err := rpc.DialHTTP("tcp", routerAddr)
	client, err := rpc.Dial("tcp", routerAddr)
	if err != nil {
		return err
	}

	sArgs := &ServiceArgs{
		ServiceName: sName,
		LocalAddr:   localAddr,
		Policy:      NewRandomSelectPolicy(),
	}

	gob.Register(sArgs.Policy)

	err = client.Call("RouterRPC.AddService", sArgs, nil)
	if err != nil {
		return err
	}

	iArgs := &InstanceArgs{
		Instance: mapTo,
	}

	err = client.Call("RouterRPC.AddServiceInstance", iArgs, nil)
	if err != nil {
		return err
	}

	return nil
}
