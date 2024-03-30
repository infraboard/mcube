package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/infraboard/mcube/v2/grpc/middleware/recovery"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Grpc{
	Host:           "127.0.0.1",
	Port:           18080,
	EnableRecovery: true,
	EnableTrace:    true,
}

type Grpc struct {
	ioc.ObjectImpl

	// 开启GRPC服务
	Enable *bool  `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	Host   string `json:"host" yaml:"host" toml:"host" env:"HOST"`
	Port   int    `json:"port" yaml:"port" toml:"port" env:"PORT"`

	EnableSSL bool   `json:"enable_ssl" yaml:"enable_ssl" toml:"enable_ssl" env:"ENABLE_SSL"`
	CertFile  string `json:"cert_file" yaml:"cert_file" toml:"cert_file" env:"CERT_FILE"`
	KeyFile   string `json:"key_file" yaml:"key_file" toml:"key_file" env:"KEY_FILE"`

	// 开启recovery恢复
	EnableRecovery bool `json:"enable_recovery" yaml:"enable_recovery" toml:"enable_recovery" env:"ENABLE_RECOVERY"`
	// 开启Trace
	EnableTrace bool `json:"enable_trace" yaml:"enable_trace" toml:"enable_trace" env:"ENABLE_TRACE"`

	// 解析后的数据
	interceptors []grpc.UnaryServerInterceptor
	svr          *grpc.Server
	log          *zerolog.Logger

	// 启动后执行
	PostStart func(context.Context) error `json:"-" yaml:"-" toml:"-" env:"-"`
	// 关闭前执行
	PreStop func(context.Context) error `json:"-" yaml:"-" toml:"-" env:"-"`
}

func (g *Grpc) Name() string {
	return AppName
}

func (g *Grpc) IsEnable() bool {
	if g.Enable == nil {
		return len(g.svr.GetServiceInfo()) > 0
	}

	return *g.Enable
}

func (i *Grpc) Priority() int {
	return -89
}

func (g *Grpc) Addr() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}

func (g *Grpc) Server() *grpc.Server {
	if g.svr == nil {
		panic("gprc server not initital")
	}
	return g.svr
}

func (g *Grpc) Init() error {
	g.log = log.Sub("grpc")
	g.svr = grpc.NewServer(g.ServerOpts()...)
	return nil
}

func (g *Grpc) AddInterceptors(interceptors ...grpc.UnaryServerInterceptor) {
	g.interceptors = append(g.interceptors, interceptors...)
}

func (g *Grpc) Interceptors() (interceptors []grpc.UnaryServerInterceptor) {
	if g.EnableRecovery {
		interceptors = append(interceptors,
			recovery.NewInterceptor(recovery.NewZeroLogRecoveryHandler()).
				UnaryServerInterceptor())
	}

	interceptors = append(interceptors, g.interceptors...)
	return
}

type ServiceInfoCtxKey struct{}

func (g *Grpc) ServerOpts() []grpc.ServerOption {
	opts := []grpc.ServerOption{}
	// 补充Trace选项
	if trace.Get().Enable && g.EnableTrace {
		g.log.Info().Msg("enable mongodb trace")
		otelgrpc.NewServerHandler()
		opts = append(opts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	}
	// 补充中间件
	opts = append(opts, grpc.ChainUnaryInterceptor(g.Interceptors()...))
	return opts
}

func (g *Grpc) Start(ctx context.Context) {
	// 启动GRPC服务
	lis, err := net.Listen("tcp", g.Addr())
	if err != nil {

		g.log.Error().Msgf("listen grpc tcp conn error, %s", err)
		return
	}

	// 启动后勾子
	ctx = context.WithValue(ctx, ServiceInfoCtxKey{}, g.svr.GetServiceInfo())
	if g.PostStart != nil {
		if err := g.PostStart(ctx); err != nil {
			g.log.Error().Msg(err.Error())
			return
		}
	}

	g.log.Info().Msgf("GRPC 服务监听地址: %s", g.Addr())

	if err := g.svr.Serve(lis); err != nil {
		g.log.Error().Msg(err.Error())
	}
}

func (g *Grpc) Stop(ctx context.Context) error {
	// 停止之前的Hook
	if g.PreStop != nil {
		if err := g.PreStop(ctx); err != nil {
			return err
		}
	}

	g.svr.GracefulStop()
	return nil
}
