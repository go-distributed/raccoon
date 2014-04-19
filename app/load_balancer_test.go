package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-distributed/raccoon/controller"
	"github.com/go-distributed/raccoon/instance"
	"github.com/go-distributed/raccoon/router"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestLBAddRouterListener(t *testing.T) {
	testLBFunction(t, 0)
}

func TestLBAddInstanceListener(t *testing.T) {
	testLBFunction(t, 1)
}

func testLBFunction(t *testing.T, option int) {
	rName := "test router"
	sName := "test service"
	expectedReply := []byte("hello, world")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expectedReply)
	}))
	defer ts.Close()
	remoteAddr := ts.Listener.Addr().String()

	mapTo, err := instance.NewInstance("test instance", sName, remoteAddr)
	if err != nil {
		t.Fatal(err)
	}

	cAddr := "127.0.0.1:14818"
	rAddr := "127.0.0.1:14817"

	r, err := router.New(rName, rAddr, cAddr)
	if err != nil {
		t.Fatal(err)
	}

	r.Start()
	defer r.Stop()

	c, err := controller.New(cAddr)
	if err != nil {
		t.Fatal(err)
	}

	c.Start()
	defer c.Stop()

	cr, err := controller.NewCRouter(rName, rAddr)
	if err != nil {
		t.Fatal(err)
	}

	lb := NewLoadBalancer(c)
	switch option {
	case 0:
		c.AddListener(controller.AddRouterEventType, lb.AddRouterListener)

		instances := make([]*instance.Instance, 0)
		c.ServiceInstances[sName] = append(instances, mapTo)

		c.RegisterRouter(cr)
	case 1:
		c.AddListener(controller.AddInstanceEventType, lb.AddInstanceListener)

		c.Routers[rName] = cr

		c.RegisterServiceInstance(mapTo)
	default:
		t.Fatal("Unknown option")
	}

	// http test
	resp, err := http.Get("http://127.0.0.1" + ServicePortMap[sName] + "/")
	if err != nil {
		t.Fatal(err)
	}

	reply, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, reply, expectedReply)

	c.AddListener(controller.FailureInstanceEventType, lb.FailureInstanceListener)

	err = r.ReportFailure(mapTo)
	assert.Nil(t, err)
}
