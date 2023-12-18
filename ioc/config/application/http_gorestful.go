package application

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/apidoc"
	"github.com/infraboard/mcube/v2/ioc/apps/health"
	"github.com/infraboard/mcube/v2/ioc/apps/metric"
	"github.com/infraboard/mcube/v2/ioc/config/logger"
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
	cors := App().HTTP.Cors
	if App().HTTP.Cors.Enabled {
		cors := restful.CrossOriginResourceSharing{
			AllowedHeaders: cors.AllowedHeaders,
			AllowedDomains: cors.AllowedDomains,
			AllowedMethods: cors.AllowedMethods,
			CookiesAllowed: false,
			Container:      r,
		}
		r.Filter(cors.Filter)
	}

	// trace中间件
	if App().Trace.Enable && App().HTTP.EnableTrace {
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

	// Metric
	mc := App().Metric
	if mc.Enable {
		r.Add(metric.NewMetric(mc.Endpoint).WebService())
		log.Info().Msgf("Get the Metric using http://%s%s", App().HTTP.Addr(), mc.Endpoint)
	}

	// API Doc
	doc := App().HTTP.ApiDoc
	if doc.Enabled {
		r.Add(apidoc.APIDocs(doc.DocPath, App().SwagerDocs))
		log.Info().Msgf("Get the API Doc using http://%s%s", App().HTTP.Addr(), doc.DocPath)
	}

	// HealthCheck
	if App().HTTP.HealthCheck.Enabled {
		hc := health.NewDefaultHealthChecker()
		r.Add(hc.WebService())
		log.Info().Msgf("健康检查地址: http://%s%s", App().HTTP.Addr(), hc.HealthCheckPath)
	}

	return r, nil
}
