package application

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/grpc/middleware/recovery"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func NewDefaultGrpc() *Grpc {
	rc := recovery.NewInterceptor(recovery.NewZeroLogRecoveryHandler())
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		rc.UnaryServerInterceptor(),
		otelgrpc.UnaryServerInterceptor(),
	))
	fmt.Println(grpcServer)

	return &Grpc{
		Host: "127.0.0.1",
		Port: 18080,
		log:  logger.Sub("grpc"),
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

	// 解析后的数据
	svr *grpc.Server
	log *zerolog.Logger
}

func (g *Grpc) setEnable(v bool) {
	g.Enable = &v
}

func (g *Grpc) Addr() string {
	return fmt.Sprintf("%s:%d", g.Host, g.Port)
}

func (g *Grpc) Parse() error {
	if g.Enable == nil {
		g.setEnable(ioc.GrpcControllerCount() > 0)
	}
	return nil
}

func (g *Grpc) Start(ctx context.Context) error {
	return nil
}

func (g *Grpc) Stop(ctx context.Context) error {
	return nil
}
