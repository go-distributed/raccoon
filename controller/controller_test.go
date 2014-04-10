package controller

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-distributed/raccoon/router"
	"github.com/go-distributed/raccoon/service"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf

func TestRegisterRouter(t *testing.T) {
	r, err := router.New(":14817")
	if err != nil {
		t.Fatal(err)
	}

	err = r.Start()
	if err != nil {
		t.Fatal("router start:", err)
	}
	defer func() {
		r.Stop()
		time.Sleep(time.Millisecond * 50)
	}()

	cr, err := NewCRouter("test router", ":14817")
	if err != nil {
		t.Fatal(err)
	}

	cAddr := "127.0.0.1:14818"
	c, err := New(cAddr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Empty(t, c.routers)

	err = c.RegisterRouter(cr)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(c.routers), 1)
	assert.Equal(t, c.routers["test router"], cr)

	err = c.RegisterRouter(cr)
	assert.NotNil(t, err)
}

func TestRegisterServiceInstance(t *testing.T) {
	r, err := router.New(":14817")
	if err != nil {
		t.Fatal(err)
	}

	err = r.Start()
	if err != nil {
		t.Fatal("router start:", err)
	}
	defer func() {
		r.Stop()
		time.Sleep(time.Millisecond * 50)
	}()

	ins, err := service.NewInstance("test instance", "test service", ":8888")
	if err != nil {
		t.Fatal(err)
	}

	cAddr := "127.0.0.1:14818"
	c, err := New(cAddr)
	if err != nil {
		t.Fatal(err)
	}

	err = c.RegisterServiceInstance(ins)
	if err != nil {
		t.Fatal(err)
	}

	err = c.RegisterServiceInstance(ins)
	assert.NotNil(t, err)
}
