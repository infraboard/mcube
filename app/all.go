package app

func InitAllApp() error {
	// 优先初始化内部app
	for _, api := range internalApps {
		if err := api.Config(); err != nil {
			return err
		}
	}

	for _, api := range grpcApps {
		if err := api.Config(); err != nil {
			return err
		}
	}

	for _, api := range restfulApps {
		if err := api.Config(); err != nil {
			return err
		}
	}

	for _, api := range httpApps {
		if err := api.Config(); err != nil {
			return err
		}
	}

	for _, api := range ginApps {
		if err := api.Config(); err != nil {
			return err
		}
	}

	return nil
}
