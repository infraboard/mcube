package application

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/apps/apidoc"
	"github.com/infraboard/mcube/ioc/apps/health"
	"github.com/infraboard/mcube/ioc/config/logger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
)

func NewGoRestfulRouterBuilder() *GoRestfulRouterBuilder {
	return &GoRestfulRouterBuilder{
		conf: &BuildConfig{},
	}
}

type GoRestfulRouterBuilder struct {
	conf *BuildConfig
}

func (b *GoRestfulRouterBuilder) Config(c *BuildConfig) {
	b.conf = c
}

func (b *GoRestfulRouterBuilder) Build() (http.Handler, error) {
	log := logger.Sub("go-restful")

	r := restful.DefaultContainer
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)

	// CORS中间件
	if App().HTTP.EnableCors {
		cors := restful.CrossOriginResourceSharing{
			AllowedHeaders: []string{"*"},
			AllowedDomains: []string{"*"},
			AllowedMethods: []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
			CookiesAllowed: false,
			Container:      r,
		}
		r.Filter(cors.Filter)
	}

	// trace中间件
	if App().HTTP.EnableTrace {
		filter := otelrestful.OTelFilter(App().AppName)
		restful.DefaultContainer.Filter(filter)
	}

	// 装载Ioc路由之前
	if b.conf.BeforeLoad != nil {
		b.conf.BeforeLoad(r)
	}

	// 装置子服务路由
	ioc.LoadGoRestfulApi(App().HTTPPrefix(), r)

	// 装载Ioc路由之后
	if b.conf.AfterLoad != nil {
		b.conf.AfterLoad(r)
	}

	// API Doc
	if App().HTTP.EnableApiDoc {
		r.Add(apidoc.APIDocs(App().HTTP.ApiDocPath, App().SwagerDocs))
		log.Info().Msgf("Get the API Doc using http://%s%s", App().HTTP.Addr(), App().HTTP.ApiDocPath)
	}

	// HealthCheck
	if App().HTTP.EnableHealthCheck {
		hc := health.NewDefaultHealthChecker()
		r.Add(hc.WebService())
		log.Info().Msgf("健康检查地址: http://%s%s", App().HTTP.Addr(), hc.HealthCheckPath)
	}

	return r, nil
}
