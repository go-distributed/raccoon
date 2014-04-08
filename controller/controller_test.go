package controller

import (
	"fmt"
	"testing"

	"github.com/go-distributed/raccoon/router"
	"github.com/stretchr/testify/assert"
)

var _ = fmt.Printf
var _ = router.NewInstance

func TestRegisterRouter(t *testing.T) {
	r, err := router.New(":14817")
	if err != nil {
		t.Fatal(err)
	}

	err = r.Start()
	if err != nil {
		t.Fatal("router start:", err)
	}

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
}
