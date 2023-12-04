package validator

type Config struct {
}

func (m *Config) Name() string {
	return AppName
}

func (m *Config) Init() error {
	return nil
}
