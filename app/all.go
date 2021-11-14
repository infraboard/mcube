package app

func InitAllApp() error {
	for _, api := range grpcApps {
		if err := api.Config(); err != nil {
			return err
		}
	}

	for _, api := range internalApps {
		if err := api.Config(); err != nil {
			return err
		}
	}

	for _, api := range httpApps {
		if err := api.Config(); err != nil {
			return err
		}
	}
	return nil
}
