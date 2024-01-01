package restful

import (
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/apidoc"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&SwaggerApiDoc{
		ApiDoc: apidoc.ApiDoc{
			Path: apidoc.DEFAUL_API_DOC_PATH,
		},
	})
}

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

func (h *SwaggerApiDoc) Registry(ws *restful.WebService) {
	ws.Path(h.Path)
	ws.Produces(restful.MIME_JSON)

	swagger := restfulspec.BuildSwagger(h.SwaggerDocConfig())
	ws.Route(ws.GET("/").To(func(r *restful.Request, w *restful.Response) {
		w.WriteAsJson(swagger)
	}))
	h.log.Info().Msgf("Get the Health using http://%s%s", http.Get().Addr(), h.Path)
}

// API Doc
func (h *SwaggerApiDoc) SwaggerDocConfig() restfulspec.Config {
	return restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       h.Path,
		PostBuildSwaggerObjectHandler: http.Get().SwagerDocs,
		DefinitionNameHandler: func(name string) string {
			if name == "state" || name == "sizeCache" || name == "unknownFields" {
				return ""
			}
			return name
		},
	}
}
