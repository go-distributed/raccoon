package router

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	rmtService "github.com/go-distributed/raccoon/service"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestService(t *testing.T) {
	localAddr := "127.0.0.1:8080"

	expectedReply, err := genRandomBytesSlice(4096)
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(expectedReply)
	}))

	defer ts.Close()

	remoteAddr := ts.Listener.Addr().String()

	mapTo := rmtService.NewInstance("test instance", "test", remoteAddr)

	s, err := newService("name", localAddr, NewRandomSelectPolicy())
	if err != nil {
		t.Fatal(err)
	}

	err = s.addInstance(mapTo)
	if err != nil {
		t.Fatal(err)
	}

	go s.start()

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

	err = s.stop()
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
