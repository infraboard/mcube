package protocol

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
{{ if $.EnableMcenter -}}
	"github.com/infraboard/keyauth/apps/endpoint"
{{- end }}
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/http/middleware/cors"

	"{{.PKG}}/conf"
{{ if $.EnableMcenter -}}
	"{{.PKG}}/version"
{{- end }}
)

// NewHTTPService 构建函数
func NewHTTPService() *HTTPService {
{{ if $.EnableMcenter -}}
	c, err := conf.C().Keyauth.Client()
	if err != nil {
		panic(err)
	}
{{- end }}

	r := gin.New()

	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.HTTP.Addr(),
		Handler:           cors.AllowAll().Handler(r),
	}
	return &HTTPService{
		r:        r,
		server:   server,
		l:        logger.Sub("HTTP Service"),
		c:        conf.C(),
{{ if $.EnableMcenter -}}
		endpoint: c.Endpoint(),
{{- end }}
	}
}

// HTTPService http服务
type HTTPService struct {
	r      gin.IRouter
	l      logger.Logger
	c      *conf.Config
	server *http.Server
{{ if $.EnableMcenter -}}
	endpoint endpoint.ServiceClient
{{- end }}
}

func (s *HTTPService) PathPrefix() string {
	return fmt.Sprintf("/%s/api", s.c.App.Name)
}

// Start 启动服务
func (s *HTTPService) Start() error {
	// 装置子服务路由
	app.LoadGinApp(s.PathPrefix(), s.r)

{{ if $.EnableMcenter -}}
	// 注册路由条目
	s.RegistryEndpoint()
{{- end }}

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

{{ if $.EnableMcenter -}}
func (s *HTTPService) RegistryEndpoint() {
	// 注册服务权限条目
	s.l.Info("start registry endpoints ...")

	req := endpoint.NewRegistryRequest(version.Short(), nil)
	_, err := s.endpoint.RegistryEndpoint(context.Background(), req)
	if err != nil {
		s.l.Warnf("registry endpoints error, %s", err)
	} else {
		s.l.Debug("service endpoints registry success")
	}
}
{{- end }}