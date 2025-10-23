package restful

import (
	"fmt"
	"path/filepath"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/apidoc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&SwaggerApiDoc{
		ApiDoc: apidoc.DefaultApiDoc(),
	})
}

// 等待所有的API接口都加载到Router上后, 提取出所有的接口并生产API Doc
type SwaggerApiDoc struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	*apidoc.ApiDoc
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
	if h.BasePath != "" {
		meta.CustomPathPrefix = h.BasePath
	}
	return meta
}

func (h *SwaggerApiDoc) ApiDocPath() string {
	if application.Get().AppAddress != "" {
		return application.Get().AppAddress + filepath.Join(http.Get().ApiObjectPathPrefix(h), h.JsonPath)
	}

	return http.Get().ApiObjectAddr(h) + h.JsonPath
}

func (h *SwaggerApiDoc) ApiUIPath() string {
	if application.Get().AppAddress != "" {
		return application.Get().AppAddress + filepath.Join(http.Get().ApiObjectPathPrefix(h), h.UIPath)
	}

	return http.Get().ApiObjectAddr(h) + h.UIPath
}

func (h *SwaggerApiDoc) Registry() {
	tags := []string{"API 文档"}

	ws := gorestful.ObjectRouter(h)
	ws.Route(ws.GET(h.JsonPath).To(h.SwaggerApiDoc).
		Doc("Swagger JSON").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	h.log.Info().Msgf("Get the API JSON data using %s", h.ApiDocPath())

	ws.Route(ws.GET(h.UIPath).To(h.SwaggerUI).
		Doc("Swagger UI").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	h.log.Info().Msgf("Get the API UI using %s", h.ApiUIPath())
}

func (h *SwaggerApiDoc) SwaggerApiDoc(r *restful.Request, w *restful.Response) {
	swagger := restfulspec.BuildSwagger(h.SwaggerDocConfig())
	w.WriteAsJson(swagger)
}

func (h *SwaggerApiDoc) SwaggerUI(r *restful.Request, w *restful.Response) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(fmt.Sprintf(apidoc.HTML_REDOC, h.ApiDocPath())))
}

// API Doc
func (h *SwaggerApiDoc) SwaggerDocConfig() restfulspec.Config {
	return restfulspec.Config{
		Host:                          application.Get().Host(),
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       http.Get().ApiObjectPathPrefix(h),
		PostBuildSwaggerObjectHandler: http.Get().SwagerDocs,
		DefinitionNameHandler: func(name string) string {
			if name == "state" || name == "sizeCache" || name == "unknownFields" {
				return ""
			}
			return name
		},
		Schemes: []string{"http", "https"},
	}
}
