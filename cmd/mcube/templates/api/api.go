package api

// Template api模板
const Template = `package api

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

	"{{.PKG}}/conf"
	"{{.PKG}}/pkg"
	"{{.PKG}}/version"
)

// NewHTTPService 构建函数
func NewHTTPService() *HTTPService {
	r := httprouter.New()
	r.Use(recovery.NewWithLogger(zap.L().Named("Recovery")))
	r.Use(accesslog.NewWithLogger(zap.L().Named("AccessLog")))
	r.Use(cors.AllowAll())
	r.EnableAPIRoot()
	r.SetAuther(pkg.NewAuther())
	r.Auth(true)
	server := &http.Server{
		ReadHeaderTimeout: 20 * time.Second,
		ReadTimeout:       20 * time.Second,
		WriteTimeout:      25 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Addr:              conf.C().App.Addr(),
		Handler:           r,
	}
	return &HTTPService{
		r:      r,
		server: server,
		l:      zap.L().Named("API"),
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
	app := s.c.App
	// 装置子服务路由
	if err := pkg.InitV1HTTPAPI(app.Name, s.r); err != nil {
		return err
	}
	if err := s.RegistryHTTPEndpoints(); err != nil {
		s.l.Errorf("registry http endpoint failed, %s", err)
	} else {
		s.l.Info("http endpoint registry success")
	}
	// 启动HTTPS服务
	if app.EnableSSL {
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
		if err := s.server.ListenAndServeTLS(app.CertFile, app.KeyFile); err != nil {
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

// RegistryHTTPEndpoints 服务所有的路由条目注册
func (s *HTTPService) RegistryHTTPEndpoints() error {
	fmt.Println(version.FullVersion())
	return nil
	// return pkg.Auth.Registry(version.ServiceName, version.GIT_TAG, s.r.GetEndpoints())
}`
