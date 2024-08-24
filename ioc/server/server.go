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
	return NewServer().Run(ctx)
}

func SetUp(cb func()) *Server {
	return NewServer().WithSetUp(cb)
}

func NewServer() *Server {
	return &Server{}
}

type Server struct {
	ioc.ObjectImpl
	setupHook func()

	http *http.Http
	grpc *grpc.Grpc

	ch     chan os.Signal
	log    *zerolog.Logger
	ctx    context.Context
	cancle context.CancelFunc
}

func (a *Server) WithSetUp(setup func()) *Server {
	a.setupHook = setup
	return a
}

func (s *Server) setup() {
	// 处理信号量
	s.ch = make(chan os.Signal, 1)
	signal.Notify(s.ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	s.ctx, s.cancle = context.WithCancel(context.Background())

	s.http = http.Get()
	s.grpc = grpc.Get()
	s.log = log.Sub("server")
	if s.setupHook != nil {
		s.setupHook()
	}
}

func (s *Server) Run(ctx context.Context) error {
	// 初始化ioc
	err := ioc.ConfigIocObject(DefaultConfig)
	if err != nil {
		return err
	}

	// ioc setup
	s.setup()

	s.log.Info().Msgf("loaded configs: %s", ioc.Config().List())
	s.log.Info().Msgf("loaded controllers: %s", ioc.Controller().List())
	s.log.Info().Msgf("loaded apis: %s", ioc.Api().List())
	s.log.Info().Msgf("loaded defaults: %s", ioc.Default().List())

	if s.http.IsEnable() {
		go s.http.Start(ctx)
	}
	if s.grpc.IsEnable() {
		go s.grpc.Start(ctx)
	}
	s.waitSign()
	return nil
}

func (s *Server) HandleError(err error) {
	if err != nil {
		s.log.Error().Msg(err.Error())
	}
}

func (s *Server) waitSign() {
	defer s.cancle()

	for sg := range s.ch {
		switch v := sg.(type) {
		default:
			s.log.Info().Msgf("receive signal '%v', start graceful shutdown", v.String())

			if s.grpc.IsEnable() {
				if err := s.grpc.Stop(s.ctx); err != nil {
					s.log.Error().Msgf("grpc graceful shutdown err: %s, force exit", err)
				} else {
					s.log.Info().Msg("grpc service stop complete")
				}
			}

			if s.http.IsEnable() {
				if err := s.http.Stop(s.ctx); err != nil {
					s.log.Error().Msgf("http graceful shutdown err: %s, force exit", err)
				} else {
					s.log.Info().Msgf("http service stop complete")
				}
			}
			return
		}
	}
}
