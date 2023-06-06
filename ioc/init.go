package ioc

func InitIocObject() error {
	// 优先初始化内部app
	for _, api := range controllers {
		if err := api.Init(); err != nil {
			return err
		}
	}

	for _, api := range grpcServers {
		if err := api.Init(); err != nil {
			return err
		}
	}

	for _, api := range goRestfulApis {
		if err := api.Init(); err != nil {
			return err
		}
	}

	for _, api := range httpApis {
		if err := api.Init(); err != nil {
			return err
		}
	}

	for _, api := range ginApis {
		if err := api.Init(); err != nil {
			return err
		}
	}

	return nil
}
