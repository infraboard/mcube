package application

import (
	"fmt"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/tools/pretty"
)

func init() {
	ioc.Config().Registry(&Application{
		AppName:      "mcube_app",
		EncryptKey:   "defualt app encrypt key",
		CipherPrefix: "@ciphered@",
		HTTP: &Http{
			Host: "127.0.0.1",
			Port: 8080,
		},
		GRPC: &Grpc{
			Host: "127.0.0.1",
			Port: 18080,
		},
	})
}

type Application struct {
	AppName      string `json:"name" yaml:"name" toml:"name" env:"APP_NAME"`
	EncryptKey   string `json:"encrypt_key" yaml:"encrypt_key" toml:"encrypt_key" env:"APP_ENCRYPT_KEY"`
	CipherPrefix string `json:"cipher_prefix" yaml:"cipher_prefix" toml:"cipher_prefix" env:"APP_CIPHER_PREFIX"`
	HTTP         *Http  `json:"http" yaml:"http"  toml:"http"`
	GRPC         *Grpc  `json:"grpc" yaml:"grpc"  toml:"grpc"`

	ioc.ObjectImpl
}

func (a *Application) HTTPPrefix() string {
	return fmt.Sprintf("/%s/api", a.AppName)
}

func (m *Application) String() string {
	return pretty.ToJSON(m)
}

func (m *Application) Name() string {
	return APPLICATION
}

func (m *Application) Init() error {
	return nil
}

type Http struct {
	Host      string `json:"size" yaml:"size" toml:"size" env:"HTTP_HOST"`
	Port      int    `json:"port" yaml:"port" toml:"port" env:"HTTP_PORT"`
	EnableSSL bool   `json:"enable_ssl" yaml:"enable_ssl" toml:"enable_ssl" env:"HTTP_ENABLE_SSL"`
	CertFile  string `json:"cert_file" yaml:"cert_file" toml:"cert_file" env:"HTTP_CERT_FILE"`
	KeyFile   string `json:"key_file" yaml:"key_file" toml:"key_file" env:"HTTP_KEY_FILE"`
}

func (a *Http) Addr() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
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
