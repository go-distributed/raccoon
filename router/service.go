package router

type service struct {
	name    string
	policy  string
	manager *serviceManager
	proxy   *proxy
}

func newService(name, policy, localAddr string) (*service, error) {
	var err error

	s := &service{
		name:   name,
		policy: policy,
	}

	// TODO: the selector should be chosen based on the policy
	s.manager, err = newServiceManager(localAddr, defaultSelector)
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
