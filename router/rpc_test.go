package router

import (
	"encoding/gob"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPC(t *testing.T) {
	r, _ := New()
	err := r.Start()
	if err != nil {
		t.Fatal("router start:", err)
	}
	defer r.Stop()

	expectedReply, _ := genRandomBytesSlice(4096)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expectedReply)
	}))
	defer ts.Close()

	sName := "TestService"
	localAddr := "127.0.0.1:8080"
	remoteAddr := ts.Listener.Addr().String()

	err = prepareRouterByRPC(sName, localAddr, remoteAddr)
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

func prepareRouterByRPC(sName, localAddr, remoteAddr string) error {
	mapTo, err := NewInstance("test instance", remoteAddr)
	if err != nil {
		return err
	}

	client, err := rpc.DialHTTP("tcp", ":14817")
	if err != nil {
		return err
	}

	sArgs := &ServiceArgs{
		ServiceName: sName,
		LocalAddr:   localAddr,
		Policy:      NewRandomSelectPolicy(),
	}
	var sReply ServiceReply

	gob.Register(sArgs.Policy)

	err = client.Call("RouterRPC.AddService", sArgs, &sReply)
	if err != nil {
		return err
	}

	iArgs := &InstanceArgs{
		ServiceName: sName,
		Instance:    mapTo,
	}
	var iReply InstanceReply

	gob.Register(iArgs.Instance)

	err = client.Call("RouterRPC.AddServiceInstance", iArgs, &iReply)
	if err != nil {
		return err
	}

	return nil
}
