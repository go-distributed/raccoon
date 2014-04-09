package router

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	rmtService "github.com/go-distributed/raccoon/service"
	"github.com/stretchr/testify/assert"
)

// TestRouter tests whether the router could route service correctly.
func TestRouter(t *testing.T) {
	sName := "TestService"
	localAddr := "127.0.0.1:8080"
	routerAddr := "127.0.0.1:14817"

	expectedReply, err := genRandomBytesSlice(4096)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expectedReply)
	}))
	defer ts.Close()

	remoteAddr := ts.Listener.Addr().String()

	mapTo, err := rmtService.NewInstance("test instance", "test", remoteAddr)
	if err != nil {
		t.Fatal(err)
	}

	r, err := New(routerAddr)
	if err != nil {
		t.Fatal(err)
	}

	// setting up service
	err = r.AddService(sName, localAddr, NewRandomSelectPolicy())
	if err != nil {
		t.Fatal(err)
	}

	err = r.AddServiceInstance(sName, mapTo)
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
