package controller

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-distributed/raccoon/instance"
	"github.com/go-distributed/raccoon/router"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestCRouter(t *testing.T) {
	rId := "test cRouter"
	r, _ := router.New(rId, ":14817", "")
	err := r.Start()
	if err != nil {
		t.Fatal("router start:", err)
	}
	defer r.Stop()

	cr, err := NewCRouter(rId, ":14817")
	if err != nil {
		t.Fatal(err)
	}

	expectedReply, _ := genRandomBytesSlice(4096)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expectedReply)
	}))
	defer ts.Close()

	sName := "test service"
	localAddr := "127.0.0.1:8080"
	remoteAddr := ts.Listener.Addr().String()

	mapTo, err := instance.NewInstance("test instance", sName, remoteAddr)
	if err != nil {
		t.Fatal(err)
	}

	// setting up service
	err = cr.AddService(sName, localAddr, router.NewRandomSelectPolicy())
	if err != nil {
		t.Fatal(err)
	}

	err = cr.AddServiceInstance(sName, mapTo)
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

func genRandomBytesSlice(size int) ([]byte, error) {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
