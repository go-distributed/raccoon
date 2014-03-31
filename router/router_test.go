package router

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Fprintf
var _ = assert.Empty

// TestRouter tests whether the router could route service correctly.
func TestRouter(t *testing.T) {
	sName := "TestService"
	localAddr := "127.0.0.1:8080"
	remoteAddr := "127.0.0.1:8081"

	expectedReply, err := genRandomBytesSlice(4096)
	if err != nil {
		t.Fatal(err)
	}

	go startHTTPServer(remoteAddr, expectedReply)

	r, err := NewRouter()
	if err != nil {
		t.Fatal(err)
	}

	// setting up service
	err = r.CreateService(sName, localAddr, NewRandomSelectPolicy(), remoteAddr)
	if err != nil {
		t.Fatal(err)
	}

	// testing service routing
	resp, err := http.Get("http://" + localAddr + "/")
	if err != nil {
		t.Fatal(err)
	}
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, reply, expectedReply)
}
