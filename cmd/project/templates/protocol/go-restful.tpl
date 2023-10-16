package protocol

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
{{ if $.EnableMcenter -}}
	"github.com/infraboard/keyauth/apps/endpoint"
	httpb "github.com/infraboard/mcube/pb/http"
	"github.com/infraboard/mcube/http/label"
{{- end }}
	"github.com/infraboard/mcube/app"

	"{{.PKG}}/conf"
	"{{.PKG}}/swagger"
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

	r := restful.DefaultContainer
	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open http://localhost:8080/apidocs/?url=http://localhost:8080/apidocs.json
	// http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("/Users/emicklei/Projects/swagger-ui/dist"))))

	// Optionally, you may need to enable CORS for the UI to work.
	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"*"},
		AllowedMethods: []string{"*"},
		CookiesAllowed: false,
		Container:      r}
	r.Filter(cors.Filter)

	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, // 1M
		Addr:              conf.C().App.HTTP.Addr(),
		Handler:           r,
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
	r      *restful.Container
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
	app.LoadRESTfulApp(s.PathPrefix(), s.r)

	// API Doc
	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: swagger.Docs}
	s.r.Add(restfulspec.NewOpenAPIService(config))
	s.l.Infof("Get the API using http://%s%s", s.c.App.HTTP.Addr(), config.APIPath)

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

        entries := []*httpb.Entry{}
        wss := s.r.RegisteredWebServices()
        for i := range wss {
                for _, r := range wss[i].Routes() {
                        m := label.Meta(r.Metadata)
                        entries = append(entries, &httpb.Entry{
                                FunctionName:     r.Operation,
                                Path:             fmt.Sprintf("%s.%s", r.Method, r.Path),
                                Method:           r.Method,
                                Resource:         m.Resource(),
                                AuthEnable:       m.AuthEnable(),
                                PermissionEnable: m.PermissionEnable(),
                                Allow:            m.Allow(),
                                AuditLog:         m.AuditEnable(),
                                Labels: map[string]string{
                                        label.Action: m.Action(),
                                },
                        })
                }
        }

        req := endpoint.NewRegistryRequest(version.Short(), entries)
        _, err := s.endpoint.RegistryEndpoint(context.Background(), req)
        if err != nil {
                s.l.Warnf("registry endpoints error, %s", err)
        } else {
                s.l.Debug("service endpoints registry success")
        }
}
{{- end }}