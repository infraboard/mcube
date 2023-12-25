package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

var (
	DefaultConfig = ioc.NewLoadConfigRequest()
)

func Run(ctx context.Context) error {
	err := ioc.ConfigIocObject(DefaultConfig)
	if err != nil {
		return err
	}

	return NewServer().Start(ctx)
}

func NewServer() *Server {
	// 处理信号量
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	ctx, cancle := context.WithCancel(context.Background())

	return &Server{
		http:   http.Get(),
		grpc:   grpc.Get(),
		log:    log.Sub("server"),
		ch:     ch,
		ctx:    ctx,
		cancle: cancle,
	}
}

type Server struct {
	ioc.ObjectImpl

	http *http.Http
	grpc *grpc.Grpc

	ch     chan os.Signal
	log    *zerolog.Logger
	ctx    context.Context
	cancle context.CancelFunc
}

func (a *Server) Start(ctx context.Context) error {
	a.log.Info().Msgf("loaded configs: %s", ioc.Config().List())
	a.log.Info().Msgf("loaded controllers: %s", ioc.Controller().List())
	a.log.Info().Msgf("loaded apis: %s", ioc.Api().List())

	if *a.http.Enable {
		go a.http.Start(ctx)
	}
	if *a.grpc.Enable {
		go a.grpc.Start(ctx)
	}

	a.waitSign()
	return nil
}

func (a *Server) HandleError(err error) {
	if err != nil {
		a.log.Error().Msg(err.Error())
	}
}

func (a *Server) waitSign() {
	defer a.cancle()

	for sg := range a.ch {
		switch v := sg.(type) {
		default:
			a.log.Info().Msgf("receive signal '%v', start graceful shutdown", v.String())

			if *a.grpc.Enable {
				if err := a.grpc.Stop(a.ctx); err != nil {
					a.log.Error().Msgf("grpc graceful shutdown err: %s, force exit", err)
				} else {
					a.log.Info().Msg("grpc service stop complete")
				}
			}

			if *a.http.Enable {
				if err := a.http.Stop(a.ctx); err != nil {
					a.log.Error().Msgf("http graceful shutdown err: %s, force exit", err)
				} else {
					a.log.Info().Msgf("http service stop complete")
				}
			}
			return
		}
	}
}
