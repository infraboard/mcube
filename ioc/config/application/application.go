package application

import (
	"fmt"

	"github.com/go-openapi/spec"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/tools/pretty"
)

func init() {
	ioc.Config().Registry(&Application{
		AppName:      "mcube_app",
		EncryptKey:   "defualt app encrypt key",
		CipherPrefix: "@ciphered@",
		HTTP:         NewDefaultHttp(),
		GRPC:         NewDefaultGrpc(),
	})
}

type Application struct {
	AppName        string `json:"name" yaml:"name" toml:"name" env:"APP_NAME"`
	AppDescription string `json:"description" yaml:"description" toml:"description" env:"APP_DESCRIPTION"`
	EncryptKey     string `json:"encrypt_key" yaml:"encrypt_key" toml:"encrypt_key" env:"APP_ENCRYPT_KEY"`
	CipherPrefix   string `json:"cipher_prefix" yaml:"cipher_prefix" toml:"cipher_prefix" env:"APP_CIPHER_PREFIX"`
	HTTP           *Http  `json:"http" yaml:"http"  toml:"http"`
	GRPC           *Grpc  `json:"grpc" yaml:"grpc"  toml:"grpc"`

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

func (m *Application) SwagerDocs(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       m.AppName,
			Description: m.AppDescription,
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: Short(),
		},
	}
}

func (m *Application) Init() error {
	return nil
}
