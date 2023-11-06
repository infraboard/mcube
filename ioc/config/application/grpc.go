package application

import "fmt"

func NewDefaultGrpc() *Grpc {
	return &Grpc{
		Host: "127.0.0.1",
		Port: 18080,
	}
}

type Grpc struct {
	Host      string `json:"host" yaml:"host" toml:"host" env:"GRPC_HOST"`
	Port      int    `json:"port" yaml:"port" toml:"port" env:"GRPC_PORT"`
	EnableSSL bool   `json:"enable_ssl" yaml:"enable_ssl" toml:"enable_ssl" env:"GRPC_ENABLE_SSL"`
	CertFile  string `json:"cert_file" yaml:"cert_file" toml:"cert_file" env:"GRPC_CERT_FILE"`
	KeyFile   string `json:"key_file" yaml:"key_file" toml:"key_file" env:"GRPC_KEY_FILE"`
}

func (a *Grpc) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
