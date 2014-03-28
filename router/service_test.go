package router

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestService(t *testing.T) {
	localAddr := "127.0.0.1:8080"
	remoteAddr := "127.0.0.1:8081"

	expectedReply, err := genRandomBytesSlice(4096)
	if err != nil {
		t.Fatal(err)
	}

	go startHTTPServer(remoteAddr, expectedReply)

	s, err := newService("name", "policy", localAddr)
	if err != nil {
		t.Fatal(err)
	}

	err = s.manager.addServiceInstance("testServiceInstance", remoteAddr)
	if err != nil {
		t.Fatal(err)
	}

	go s.start()

	resp, err := http.Get("http://" + localAddr + "/")
	if err != nil {
		t.Fatal(err)
	}
	reply, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, reply, expectedReply)

	s.stop()
}

func startHTTPServer(hostPort string, writeBack []byte) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(writeBack)
	})

	http.ListenAndServe(hostPort, nil)
}

func genRandomBytesSlice(size int) ([]byte, error) {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}
