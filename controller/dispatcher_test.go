package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testEvent int

func (e *testEvent) Type() string {
	return "testListener"
}

func TestDispatcher(t *testing.T) {
	d := newDispatcher()

	called := false
	l := func(e event) {
		called = true
	}

	d.addListener("testListener", l)
	d.dispatch(new(testEvent))

	assert.Equal(t, called, true)
}
