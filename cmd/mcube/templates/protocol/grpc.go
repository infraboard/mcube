package protocol

// GRPCTemplate api模板
const GRPCTemplate = `package protocol

import (
	"context"
	"fmt"
	"log"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/infraboard/keyauth/pkg/endpoint"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"{{.PKG}}/conf"
	"{{.PKG}}/pkg"
	"{{.PKG}}/version"
)

// NewGRPCService todo
func NewGRPCService(interceptors ...grpc.UnaryServerInterceptor) *GRPCService {
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		interceptors...,
	)))

	return &GRPCService{
		svr: grpcServer,
		l:   zap.L().Named("GRPC"),
		c:   conf.C(),
	}
}

// GRPCService grpc服务
type GRPCService struct {
	svr *grpc.Server
	l   logger.Logger
	c   *conf.Config
}

// Start 启动GRPC服务
func (s *GRPCService) Start() error {
	// 装载所有GRPC服务
	pkg.InitV1GRPCAPI(s.svr)

	// 启动HTTP服务
	s.l.Infof("GRPC 开始启动, 监听地址: %s", s.c.GRPC.Addr())
	lis, err := net.Listen("tcp", s.c.GRPC.Addr())
	if err != nil {
		log.Fatal(err)
	}
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			s.l.Info("service is stopped")
		}

		return fmt.Errorf("start service error, %s", err.Error())
	}

	return nil
}

// RegistryEndpoints 注册条目
func (s *GRPCService) RegistryEndpoints() error {
	kc, err := s.c.Keyauth.Client()
	if err != nil {
		return err
	}
	req := endpoint.NewRegistryRequest(version.Short(), pkg.HTTPEntry().Items)
	_, err = kc.Endpoint().Registry(context.Background(), req)
	return err
}

// Stop 停止GRPC服务
func (s *GRPCService) Stop() error {
	s.l.Info("start graceful shutdown")

	// 优雅关闭HTTP服务
	s.svr.GracefulStop()

	return nil
}
`
