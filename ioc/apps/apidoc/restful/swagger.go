package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/apidoc"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&SwaggerApiDoc{
		ApiDoc: apidoc.ApiDoc{
			Path: "/apidocs.json",
		},
	})
}

// 等待所有的API接口都加载到Router上后, 提取出所有的接口并生产API Doc
type SwaggerApiDoc struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	apidoc.ApiDoc
}

func (h *SwaggerApiDoc) Name() string {
	return apidoc.AppName
}

func (h *SwaggerApiDoc) Init() error {
	h.log = log.Sub("api_doc")
	h.Registry()
	return nil
}

func (i *SwaggerApiDoc) Priority() int {
	return -100
}

func (h *SwaggerApiDoc) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	meta.CustomPathPrefix = h.Path
	return meta
}

func (h *SwaggerApiDoc) Registry() {
	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
		swagger := restfulspec.BuildSwagger(h.SwaggerDocConfig())
		w.WriteAsJson(swagger)
	}),
	)

	h.log.Info().Msgf("Get the API Doc using %s", http.Get().ApiObjectAddr(h))
}

// API Doc
func (h *SwaggerApiDoc) SwaggerDocConfig() restfulspec.Config {
	return restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       http.Get().ApiObjectPathPrefix(h),
		PostBuildSwaggerObjectHandler: http.Get().SwagerDocs,
		DefinitionNameHandler: func(name string) string {
			if name == "state" || name == "sizeCache" || name == "unknownFields" {
				return ""
			}
			return name
		},
	}
}
