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
	Controller    *controller.Controller
	failureRecord map[string][]string
}

func NewDemoApp(c *controller.Controller) *DemoApp {
	return &DemoApp{c, make(map[string][]string)}
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

		// TODO: debugging, delete it
		if r.Id() == "rr" {
			port = ":8081"
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
		// TODO: debugging, delete it
		if r.Id() == "rr" {
			port = ":8081"
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

	log.Printf("Failure Instance: '%v', '%v'\n", e.Reporter, e.Instance)

	list := da.failureRecord[e.Instance.Name]
	for _, reporter := range list {
		if reporter == e.Reporter {
			log.Println(reporter, "has reported failure of", e.Instance.Name, "before")
			return
		}
	}

	list = append(list, e.Reporter)
	da.failureRecord[e.Instance.Name] = list

	if len(list) == len(da.Controller.Routers) {
		for _, r := range da.Controller.Routers {
			instance := e.Instance

			err := r.RemoveServiceInstance(instance.Service, instance)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
