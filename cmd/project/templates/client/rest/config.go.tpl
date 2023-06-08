package rest

func NewDefaultConfig() *Config {
	return &Config{
		Address:    "http://127.0.0.1:8080",
		PathPrefix: "/mpaas/api/v1/",
	}
}

type Config struct {
	Token      string `json:"token" toml:"token" yaml:"token" env:"MCENTER_TOKEN"`
	Address    string `json:"address" toml:"address" yaml:"address" env:"MCENTER_ADDRESS"`
	PathPrefix string `json:"path_prefix" toml:"path_prefix" yaml:"path_prefix" env:"MCENTER_PATH_PREFIX"`
}
