package application

import (
	"net/http"
	"sync"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/apps/apidoc"
	"github.com/infraboard/mcube/ioc/apps/health"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
)

func NewGoRestfulRouterBuilder() *GoRestfulRouterBuilder {
	return &GoRestfulRouterBuilder{
		log: logger.Sub("http"),
	}
}

type GoRestfulRouterBuilder struct {
	Router *restful.Container
	log    *zerolog.Logger
	lock   sync.Mutex
}

func (b *GoRestfulRouterBuilder) Init() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.Router == nil {
		r := restful.DefaultContainer
		restful.DefaultResponseContentType(restful.MIME_JSON)
		restful.DefaultRequestContentType(restful.MIME_JSON)

		// CORS中间件
		cors := restful.CrossOriginResourceSharing{
			AllowedHeaders: []string{"*"},
			AllowedDomains: []string{"*"},
			AllowedMethods: []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
			CookiesAllowed: false,
			Container:      r,
		}
		r.Filter(cors.Filter)
		// trace中间件
		filter := otelrestful.OTelFilter(App().AppName)
		restful.DefaultContainer.Filter(filter)

		// 装置子服务路由
		ioc.LoadGoRestfulApi(App().HTTPPrefix(), r)

		// API Doc
		r.Add(apidoc.APIDocs(App().HTTP.ApiDocPath, App().SwagerDocs))
		b.log.Info().Msgf("Get the API Doc using http://%s%s", App().HTTP.Addr(), App().HTTP.ApiDocPath)

		// HealthCheck
		hc := health.NewDefaultHealthChecker()
		r.Add(hc.WebService())
		b.log.Info().Msgf("健康检查地址: http://%s%s", App().HTTP.Addr(), hc.HealthCheckPath)
		b.Router = r
	}

	return nil
}

func (b *GoRestfulRouterBuilder) GetRouter() http.Handler {
	return b.Router
}
