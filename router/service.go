package router

type service struct {
	name    string
	policy  routePolicy
	manager *serviceManager
	proxy   *proxy
}

func newService(name string, policy routePolicy, localAddr string) (*service, error) {
	var err error

	s := &service{
		name:   name,
		policy: policy,
	}

	selector, err := newSelector(policy)
	if err != nil {
		return nil, err
	}

	s.manager, err = newServiceManager(localAddr, selector)
	if err != nil {
		return nil, err
	}

	s.proxy, err = newProxy(localAddr, s.manager)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *service) start() error {
	err := s.proxy.start()
	return err
}

func (s *service) stop() error {
	err := s.proxy.stop()
	return err
}
