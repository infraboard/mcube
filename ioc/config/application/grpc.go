package application

import (
	"context"
	"fmt"
	"net"

	"github.com/infraboard/mcube/grpc/middleware/recovery"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func NewDefaultGrpc() *Grpc {
	return &Grpc{
		Host:           "127.0.0.1",
		Port:           18080,
		EnableRecovery: true,
	}
}

type Grpc struct {
	// 开启GRPC服务
	Enable *bool  `json:"enable" yaml:"enable" toml:"enable" env:"GRPC_ENABLE"`
	Host   string `json:"host" yaml:"host" toml:"host" env:"GRPC_HOST"`
	Port   int    `json:"port" yaml:"port" toml:"port" env:"GRPC_PORT"`

	EnableSSL bool   `json:"enable_ssl" yaml:"enable_ssl" toml:"enable_ssl" env:"GRPC_ENABLE_SSL"`
	CertFile  string `json:"cert_file" yaml:"cert_file" toml:"cert_file" env:"GRPC_CERT_FILE"`
	KeyFile   string `json:"key_file" yaml:"key_file" toml:"key_file" env:"GRPC_KEY_FILE"`

	// 开启recovery恢复
	EnableRecovery bool `json:"enable_recovery" yaml:"enable_recovery" toml:"enable_recovery" env:"GRPC_ENABLE_RECOVERY"`
	// 开启Trace
	EnableTrace bool `json:"enable_trace" yaml:"enable_trace" toml:"enable_trace" env:"GRPC_ENABLE_TRACE"`

	// 解析后的数据
	interceptors []grpc.UnaryServerInterceptor
	svr          *grpc.Server
	log          *zerolog.Logger

	// 启动后执行
	PostStart func(context.Context) error `json:"-" yaml:"-" toml:"-" env:"-"`
	// 关闭前执行
	PreStop func(context.Context) error `json:"-" yaml:"-" toml:"-" env:"-"`
}

func (g *Grpc) setEnable(v bool) {
	g.Enable = &v
}

func (g *Grpc) Addr() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}

func (g *Grpc) Parse() error {
	g.log = logger.Sub("grpc")

	if g.Enable == nil {
		g.setEnable(ioc.GrpcControllerCount() > 0)
	}
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
	if g.EnableTrace {
		interceptors = append(interceptors, otelgrpc.UnaryServerInterceptor())
	}

	interceptors = append(interceptors, g.interceptors...)
	return
}

type ServiceInfoCtxKey struct{}

func (g *Grpc) Start(ctx context.Context, cb ErrHandler) {
	g.svr = grpc.NewServer(grpc.ChainUnaryInterceptor(g.Interceptors()...))

	// 装载所有GRPC服务
	ioc.LoadGrpcController(g.svr)

	// 启动GRPC服务
	lis, err := net.Listen("tcp", g.Addr())
	if err != nil {
		cb(fmt.Errorf("listen grpc tcp conn error, %s", err))
		return
	}

	// 启动后勾子
	ctx = context.WithValue(ctx, ServiceInfoCtxKey{}, g.svr.GetServiceInfo())
	if g.PostStart != nil {
		if err := g.PostStart(ctx); err != nil {
			cb(err)
			return
		}
	}

	g.log.Info().Msgf("GRPC 服务监听地址: %s", g.Addr())
	cb(g.svr.Serve(lis))
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
