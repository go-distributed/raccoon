package app

import (
	"log"

	"github.com/go-distributed/raccoon/controller"
	"github.com/go-distributed/raccoon/router"
)

var ServicePortMap map[string]string

func init() {
	ServicePortMap = make(map[string]string)
	ServicePortMap["test service"] = ":8080"
}

type LoadBalancer struct {
	Controller *controller.Controller
}

func NewLoadBalancer(c *controller.Controller) *LoadBalancer {
	return &LoadBalancer{c}
}

// AddRouterListener makes the Load Balancer know of the newly added router.
// It adds known service instances to the router and set service policy as round robin.
// TODO: It assumes the new router is newly created: which means it has no services.
//func (lb *LoadBalancer) AddRouterListener(event *controller.AddRouterEvent) {
func (lb *LoadBalancer) AddRouterListener(event controller.Event) {
	if event.Type() != controller.AddRouterEventType {
		panic("")
	}

	e := event.(*controller.AddRouterEvent)

	r := lb.Controller.Routers[e.Id]

	for service, instances := range lb.Controller.ServiceInstances {
		port, ok := ServicePortMap[service]
		if !ok {
			log.Println("Unknown port for service:", service)
			continue
		}
		err := r.AddService(service, port, router.NewRoundRobinPolicy())
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

// AddInstanceListener makes the Load Balancer know of the newly added service instance.
// It adds the service instances to the routers it knew.
// TODO: It assumes the new instance is newly created: which means no other router had it before.
func (lb *LoadBalancer) AddInstanceListener(event controller.Event) {
	if event.Type() != controller.AddInstanceEventType {
		panic("")
	}

	e := event.(*controller.AddInstanceEvent)
	for _, r := range lb.Controller.Routers {
		instance := e.Instance

		port, ok := ServicePortMap[instance.Service]
		if !ok {
			log.Println("Unknown port for service:", instance.Service)
			continue
		}

		err := r.AddService(instance.Service, port, router.NewRoundRobinPolicy())
		if err != nil {
			log.Println(err)
		}

		err = r.AddServiceInstance(instance.Service, instance)
		if err != nil {
			log.Println(err)
		}
	}
}

func (lb *LoadBalancer) RmInstanceListener(event controller.Event) {
	if event.Type() != controller.RmInstanceEventType {
		panic("")
	}

	e := event.(*controller.RmInstanceEvent)

	for _, r := range lb.Controller.Routers {
		instance := e.Instance

		err := r.RemoveServiceInstance(instance.Service, instance)
		if err != nil {
			log.Println(err)
		}
	}
}

func (lb *LoadBalancer) FailureInstanceListener(event controller.Event) {
	if event.Type() != controller.FailureInstanceEventType {
		panic("")
	}

	e := event.(*controller.FailureInstanceEvent)

	log.Printf("Failure Instance: '%v', '%v'", e.Reporter, e.Instance)
}
