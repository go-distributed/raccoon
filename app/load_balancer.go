package app

import (
	"github.com/go-distributed/raccoon/controller"
	"github.com/go-distributed/raccoon/router"
	"github.com/go-distributed/raccoon/service"
)

type LoadBalancer struct {
	Controller *controller.Controller
}

// AddRouterListener makes the Load Balancer know of the newly added router.
// It adds known service instances to the router and set service policy as round robin.
// TODO: It assumes the new router is newly created: which means it has no services.
func (lb *LoadBalancer) AddRouterListener(event *controller.AddRouterEvent) {
	cr := lb.Controller.Routers[event.Id]

	for service, instances := range lb.Controller.ServiceInstances {
		cr.AddService(service, instances[0].Addr, router.NewRoundRobinPolicy())
		for _, instance := range instances {
			cr.AddServiceInstance(service, instance)
		}
	}
}

// AddInstanceListener makes the Load Balancer know of the newly added service instance.
// It adds the service instances to the routers it knew.
// TODO: It assumes the new instance is newly created: which means no other router had it before.
func (lb *LoadBalancer) AddInstanceListener(event *controller.AddInstanceEvent) {
	for _, r := range lb.Controller.Routers {
		instance := &service.Instance{
			Addr:    event.Addr,
			Name:    event.Name,
			Service: event.Service,
		}
		r.AddServiceInstance(event.Service, instance)
	}
}
