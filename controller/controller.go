package controller

import (
	"github.com/go-distributed/raccoon/router"
)

type controller struct {
	serviceInstances map[string][]*router.Instance
}
