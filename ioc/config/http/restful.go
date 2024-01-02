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
		RouterBuilderHooks: NewRouterBuilderHooks(),
	}
}

type GoRestfulRouterBuilder struct {
	*RouterBuilderHooks
}

func (b *GoRestfulRouterBuilder) Build() (http.Handler, error) {
	r := restful.DefaultContainer
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)

	// trace中间件
	if trace.Get().Enable && Get().EnableTrace {
		filter := otelrestful.OTelFilter(application.Get().AppName)
		r.Filter(filter)
	}

	// 装载Ioc路由之前
	for _, fn := range b.Before {
		fn(r)
	}

	// 装置子服务路由
	ioc.LoadGoRestfulApi(Get().HTTPPrefix(), r)

	// 装载Ioc路由之后
	for _, fn := range b.After {
		fn(r)
	}

	return r, nil
}
