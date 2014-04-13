package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-distributed/raccoon/controller"
	"github.com/go-distributed/raccoon/router"
	"github.com/go-distributed/raccoon/service"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestLBAddRouterListener(t *testing.T) {
	sName := "test service"
	expectedReply := []byte("hello, world")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expectedReply)
	}))
	defer ts.Close()
	remoteAddr := ts.Listener.Addr().String()

	mapTo := service.NewInstance("test instance", sName, remoteAddr)

	rAddr := "127.0.0.1:14817"
	r, err := router.New(rAddr)
	if err != nil {
		t.Fatal(err)
	}

	r.Start()
	defer r.Stop()

	cAddr := "127.0.0.1:14818"
	c, err := controller.New(cAddr)
	if err != nil {
		t.Fatal(err)
	}

	instances := make([]*service.Instance, 0)
	c.ServiceInstances[sName] = append(instances, mapTo)

	lb := NewLoadBalancer(c)
	c.AddListener(controller.AddRouterEventType, lb.AddRouterListener)

	cr, err := controller.NewCRouter("test router", rAddr)
	if err != nil {
		t.Fatal(err)
	}

	c.RegisterRouter(cr)

	// http test
	resp, err := http.Get("http://127.0.0.1" + router.ServicePortMap[sName] + "/")
	if err != nil {
		t.Fatal(err)
	}

	reply, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, reply, expectedReply)
}

func TestLBAddInstanceListener(t *testing.T) {
	rName := "test router"
	sName := "test service"
	expectedReply := []byte("hello, world")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expectedReply)
	}))
	defer ts.Close()
	remoteAddr := ts.Listener.Addr().String()

	mapTo := service.NewInstance("test instance", sName, remoteAddr)

	rAddr := "127.0.0.1:14817"
	r, err := router.New(rAddr)
	if err != nil {
		t.Fatal(err)
	}

	r.Start()
	defer r.Stop()

	cAddr := "127.0.0.1:14818"
	c, err := controller.New(cAddr)
	if err != nil {
		t.Fatal(err)
	}

	lb := NewLoadBalancer(c)
	c.AddListener(controller.AddInstanceEventType, lb.AddInstanceListener)

	cr, err := controller.NewCRouter(rName, rAddr)
	if err != nil {
		t.Fatal(err)
	}

	c.Routers[rName] = cr

	c.RegisterServiceInstance(mapTo)

	// http test
	resp, err := http.Get("http://127.0.0.1" + router.ServicePortMap[sName] + "/")
	if err != nil {
		t.Fatal(err)
	}

	reply, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, reply, expectedReply)
}
