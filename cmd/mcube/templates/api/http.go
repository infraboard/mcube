package api

// HTTPTemplate api模板
const HTTPTemplate = `package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/infraboard/mcube/http/middleware/accesslog"
	"github.com/infraboard/mcube/http/middleware/cors"
	"github.com/infraboard/mcube/http/middleware/recovery"
	"github.com/infraboard/mcube/http/router"
	"github.com/infraboard/mcube/http/router/httprouter"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"{{.PKG}}/client"
	"{{.PKG}}/conf"
	"{{.PKG}}/pkg"
)

// NewHTTPService 构建函数
func NewHTTPService() *HTTPService {
	r := httprouter.New()
	r.Use(recovery.NewWithLogger(zap.L().Named("Recovery")))
	r.Use(accesslog.NewWithLogger(zap.L().Named("AccessLog")))
	r.Use(cors.AllowAll())
	r.EnableAPIRoot()
	server := &http.Server{
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      25 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Addr:              conf.C().HTTP.Addr(),
		Handler:           r,
	}
	return &HTTPService{
		r:      r,
		server: server,
		l:      zap.L().Named("HTTP Service"),
		c:      conf.C(),
	}
}

// HTTPService http服务
type HTTPService struct {
	r      router.Router
	l      logger.Logger
	c      *conf.Config
	server *http.Server
}

// Start 启动服务
func (s *HTTPService) Start() error {
	hc := s.c.HTTP

	// 初始化GRPC客户端
	if err := s.initGRPCClient(); err != nil {
		return err
	}

	// 装置子服务路由
	if err := pkg.InitV1HTTPAPI(s.c.App.Name, s.r); err != nil {
		return err
	}
	// 启动HTTPS服务
	if hc.EnableSSL {
		// 安全的算法挑选标准依赖: https://wiki.mozilla.org/Security/Server_Side_TLS
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256, tls.X25519},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				// tls 1.2
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				// tls 1.3
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
			},
		}
		s.server.TLSConfig = cfg
		s.server.TLSNextProto = make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)
		s.l.Infof("HTTPS服务启动成功, 监听地址: %s", s.server.Addr)
		if err := s.server.ListenAndServeTLS(hc.CertFile, hc.KeyFile); err != nil {
			if err == http.ErrServerClosed {
				s.l.Info("service is stopped")
			}
			return fmt.Errorf("start service error, %s", err.Error())
		}
	}
	// 启动 HTTP服务
	s.l.Infof("HTTP服务启动成功, 监听地址: %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			s.l.Info("service is stopped")
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}

// Stop 停止server
func (s *HTTPService) Stop() error {
	s.l.Info("start graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	// 优雅关闭HTTP服务
	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Errorf("graceful shutdown timeout, force exit")
	}
	return nil
}

// InitGRPCClient 初始化grpc客户端
func (s *HTTPService) initGRPCClient() error {
	cf := client.NewDefaultConfig()
	cf.SetAddress(s.c.GRPC.Addr())
	cf.SetClientCredentials(s.c.Keyauth.ClientID, s.c.Keyauth.ClientSecret)
	cli, err := client.NewClient(cf)
	if err != nil {
		return err
	}
	client.SetGlobal(cli)
	return err
}
`
