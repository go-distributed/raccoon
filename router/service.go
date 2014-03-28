package router

type service struct {
	name    string
	policy  string
	manager *serviceManager
}

func newService(name, policy, localAddr string, raddrStrs []string) (*service, error) {
	var err error

	s := &service{
		name:   name,
		policy: policy,
	}

	// TODO: the selector should be chosen based on the policy
	s.manager, err = newServiceManager(localAddr, raddrStrs, defaultSelector)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *service) start() error {
	err := s.manager.proxy.start()
	return err
}

func (s *service) stop() error {
	err := s.manager.proxy.stop()
	return err
}
