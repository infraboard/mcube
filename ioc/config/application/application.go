package application

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-openapi/spec"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/infraboard/mcube/tools/pretty"
	"github.com/rs/zerolog"
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

	ch     chan os.Signal
	log    *zerolog.Logger
	ctx    context.Context
	cancle context.CancelFunc
}

func (a *Application) UseGoRestful() {
	a.HTTP.WEB_FRAMEWORK = WEB_FRAMEWORK_GO_RESTFUL
}

func (a *Application) HTTPPrefix() string {
	return fmt.Sprintf("/%s/api", a.AppName)
}

func (a *Application) String() string {
	return pretty.ToJSON(a)
}

func (a *Application) Name() string {
	return APPLICATION
}

func (a *Application) Init() error {
	// 处理信号量
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	a.ch = ch
	a.log = logger.Sub("application")
	a.ctx, a.cancle = context.WithCancel(context.Background())

	if err := a.HTTP.Parse(); err != nil {
		return err
	}
	if err := a.GRPC.Parse(); err != nil {
		return err
	}
	return nil
}

func (a *Application) SwagerDocs(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       a.AppName,
			Description: a.AppDescription,
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

func (a *Application) Start(ctx context.Context) error {
	a.log.Info().Msgf("loaded configs: %s", ioc.Config().List())
	a.log.Info().Msgf("loaded controllers: %s", ioc.Controller().List())
	a.log.Info().Msgf("loaded apis: %s", ioc.Api().List())

	if *a.HTTP.Enable {
		go a.HTTP.Start(ctx, a.HandleError)
	}
	if *a.GRPC.Enable {
		go a.GRPC.Start(ctx, a.HandleError)
	}

	a.waitSign()
	return nil
}

func (a *Application) HandleError(err error) {
	if err != nil {
		a.log.Error().Msg(err.Error())
	}
}

func (a *Application) waitSign() {
	defer a.cancle()

	for sg := range a.ch {
		switch v := sg.(type) {
		default:
			a.log.Info().Msgf("receive signal '%v', start graceful shutdown", v.String())

			if *a.GRPC.Enable {
				if err := a.GRPC.Stop(a.ctx); err != nil {
					a.log.Error().Msgf("grpc graceful shutdown err: %s, force exit", err)
				} else {
					a.log.Info().Msg("grpc service stop complete")
				}
			}

			if *a.HTTP.Enable {
				if err := a.HTTP.Stop(a.ctx); err != nil {
					a.log.Error().Msgf("http graceful shutdown err: %s, force exit", err)
				} else {
					a.log.Info().Msgf("http service stop complete")
				}
			}
			return
		}
	}
}
