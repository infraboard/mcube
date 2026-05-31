package server

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/jsonrpc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

var (
	DefaultConfig = ioc.NewLoadConfigRequest()

	// defaultShutdownTimeout 首次收到退出信号后，优雅关闭的最长等待时间。
	defaultShutdownTimeout = 30 * time.Second
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

	http    *http.Http
	grpc    *grpc.Grpc
	jsonrpc *jsonrpc.JsonRpc

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
	// 缓冲避免 shutdown 期间再次 Ctrl+C 时阻塞终端信号投递
	s.ch = make(chan os.Signal, 2)
	signal.Notify(s.ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, syscall.SIGQUIT)
	s.ctx, s.cancle = context.WithCancel(context.Background())

	s.http = http.Get()
	s.grpc = grpc.Get()
	s.jsonrpc = jsonrpc.Get()
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

	// 打印对象
	ioc.DefaultStore.ForEatch(func(ns *ioc.NamespaceStore) {
		s.log.Info().Msgf("loaded %s: %s", ns.Namespace, ns.List())
	})

	if s.http.IsEnable() {
		go s.http.Start(ctx)
	}
	if s.grpc.IsEnable() {
		go s.grpc.Start(ctx)
	}
	if s.jsonrpc.IsEnable() {
		go s.jsonrpc.Start(ctx)
	}
	s.waitSign()

	// TODO: 判断所有服务启动成功
	return nil
}

func (s *Server) HandleError(err error) {
	if err != nil {
		s.log.Error().Msg(err.Error())
	}
}

func (s *Server) waitSign() {
	defer s.cancle()
	defer signal.Stop(s.ch)

	for sg := range s.ch {
		switch sg.(type) {
		default:
			s.shutdown(sg)
			return
		}
	}
}

// shutdown 处理退出信号。首次信号触发带超时的优雅关闭；再次收到信号则取消等待并强制停止。
func (s *Server) shutdown(sig os.Signal) {
	s.log.Info().Msgf("receive signal '%v', start graceful shutdown (timeout %s)", sig, defaultShutdownTimeout)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	var forced atomic.Bool
	s.watchForceShutdown(shutdownCtx, cancel, &forced)

	s.stopServices(shutdownCtx)

	ioc.DefaultStore.Stop(shutdownCtx)

	if forced.Load() {
		s.log.Warn().Msg("forced shutdown, exit now")
		os.Exit(1)
	}
}

// watchForceShutdown 在优雅关闭进行中监听后续信号；再次 Ctrl+C 则取消 shutdownCtx 以打断阻塞的 Stop。
func (s *Server) watchForceShutdown(shutdownCtx context.Context, cancel context.CancelFunc, forced *atomic.Bool) {
	go func() {
		select {
		case sig := <-s.ch:
			s.log.Warn().Msgf("receive signal '%v' again, force shutdown", sig)
			forced.Store(true)
			cancel()
		case <-shutdownCtx.Done():
		}
	}()
}

func (s *Server) stopServices(ctx context.Context) {
	if s.grpc.IsEnable() {
		if err := s.grpc.Stop(ctx); err != nil {
			s.log.Error().Msgf("grpc shutdown err: %s", err)
		} else {
			s.log.Info().Msg("grpc service stop complete")
		}
	}

	if s.http.IsEnable() {
		if err := s.http.Stop(ctx); err != nil {
			s.log.Error().Msgf("http shutdown err: %s", err)
		} else {
			s.log.Info().Msg("http service stop complete")
		}
	}

	if s.jsonrpc.IsEnable() {
		if err := s.jsonrpc.Stop(ctx); err != nil {
			s.log.Error().Msgf("jsonrpc shutdown err: %s", err)
		} else {
			s.log.Info().Msg("jsonrpc service stop complete")
		}
	}
}
