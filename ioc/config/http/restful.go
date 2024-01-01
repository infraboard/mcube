package http

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
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
	r := restful.DefaultContainer
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)

	// CORS中间件
	cors := Get().Cors
	if cors.Enabled {
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
	if trace.Get().Enable && Get().EnableTrace {
		filter := otelrestful.OTelFilter(application.Get().AppName)
		restful.DefaultContainer.Filter(filter)
	}

	// 装载Ioc路由之前
	if b.conf.BeforeLoad != nil {
		b.conf.BeforeLoad(r)
	}

	// 装置子服务路由
	ioc.LoadGoRestfulApi(Get().HTTPPrefix(), r)

	// 装载Ioc路由之后
	if b.conf.AfterLoad != nil {
		b.conf.AfterLoad(r)
	}

	return r, nil
}