package router

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestProxy(t *testing.T) {
	writeBack, err := genRandomBytesSlice(4096)
	if err != nil {
		t.Fatal(err)
	}

	serverHostPort := "127.0.0.1:8080"
	proxyHostPort := "127.0.0.1:8081"

	go startHTTPServer(serverHostPort, writeBack)
	go startProxy(t, proxyHostPort, serverHostPort)
	time.Sleep(time.Millisecond * 50)

	resp, _ := http.Get("http://" + serverHostPort + "/")
	body, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, body, writeBack)
}

func startHTTPServer(hostPort string, writeBack []byte) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(writeBack)
	})

	http.ListenAndServe(hostPort, nil)
}

func startProxy(t *testing.T, proxyHostPort, serverHostPort string) {
	l, err := net.Listen("tcp", proxyHostPort)
	if err != nil {
		t.Fatal(err)
	}

	one, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}

	other, err := net.Dial("tcp", serverHostPort)
	if err != nil {
		t.Fatal(err)
	}

	p := newProxy(one, other)
	p.start()
}

func genRandomBytesSlice(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
