package app

import (
	"log"

	"github.com/go-distributed/raccoon/controller"
	"github.com/go-distributed/raccoon/router"
)

var rIdPolicyMap map[string]router.Policy

func init() {
	rIdPolicyMap = make(map[string]router.Policy)
	rIdPolicyMap["rr"] = router.NewRoundRobinPolicy()
	rIdPolicyMap["rand"] = router.NewRandomSelectPolicy()
}

type DemoApp struct {
	Controller *controller.Controller
}

func NewDemoApp(c *controller.Controller) *DemoApp {
	return &DemoApp{c}
}

func (da *DemoApp) AddRouterListener(event controller.Event) {
	if event.Type() != controller.AddRouterEventType {
		panic("")
	}

	e := event.(*controller.AddRouterEvent)

	r := da.Controller.Routers[e.Id]

	for service, instances := range da.Controller.ServiceInstances {
		port, ok := ServicePortMap[service]
		if !ok {
			log.Println("Unknown port for service:", service)
			continue
		}
		err := r.AddService(service, port, rIdPolicyMap[e.Id])
		if err != nil {
			log.Println(err)
		}

		for _, instance := range instances {
			err := r.AddServiceInstance(service, instance)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (da *DemoApp) AddInstanceListener(event controller.Event) {
	if event.Type() != controller.AddInstanceEventType {
		panic("")
	}

	e := event.(*controller.AddInstanceEvent)
	for _, r := range da.Controller.Routers {
		instance := e.Instance

		port, ok := ServicePortMap[instance.Service]
		if !ok {
			log.Println("Unknown port for service:", instance.Service)
			continue
		}

		err := r.AddService(instance.Service, port, rIdPolicyMap[r.Id()])
		if err != nil {
			log.Println(err)
		}

		err = r.AddServiceInstance(instance.Service, instance)
		if err != nil {
			log.Println(err)
		}
	}
}

func (da *DemoApp) RmInstanceListener(event controller.Event) {
	if event.Type() != controller.RmInstanceEventType {
		panic("")
	}

	e := event.(*controller.RmInstanceEvent)

	for _, r := range da.Controller.Routers {
		instance := e.Instance

		err := r.RemoveServiceInstance(instance.Service, instance)
		if err != nil {
			log.Println(err)
		}
	}
}

func (da *DemoApp) FailureInstanceListener(event controller.Event) {
	if event.Type() != controller.FailureInstanceEventType {
		panic("")
	}

	e := event.(*controller.FailureInstanceEvent)

	log.Printf("Failure Instance: '%v', '%v'", e.Reporter, e.Instance)
}
