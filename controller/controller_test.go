package controller

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-distributed/raccoon/router"
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

	cr, err := NewCRouter("test cRouter", ":14817")
	if err != nil {
		t.Fatal(err)
	}

	c := New()

	err = c.RegisterRouter(cr)
	if err != nil {
		t.Fatal(err)
	}

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

	ins, err := router.NewInstance("test instance", "test service", ":8888")
	if err != nil {
		t.Fatal(err)
	}

	c := New()

	err = c.RegisterServiceInstance(ins)
	if err != nil {
		t.Fatal(err)
	}

	err = c.RegisterServiceInstance(ins)
	assert.NotNil(t, err)
}
