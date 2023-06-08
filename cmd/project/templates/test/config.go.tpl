package test

import "github.com/caarlos0/env/v6"

var (
	C *Config
)

// LoadConfigFromEnv 从环境变量中加载配置
func LoadConfigFromEnv() {
	C = &Config{}
	if err := env.Parse(C); err != nil {
		panic(err)
	}
}

type Config struct {
	MCENTER_GRPC_ADDRESS  string `env:"MCENTER_GRPC_ADDRESS"`
	MCENTER_CLINET_ID     string `env:"MCENTER_CLINET_ID"`
	MCENTER_CLIENT_SECRET string `env:"MCENTER_CLIENT_SECRET"`
}
