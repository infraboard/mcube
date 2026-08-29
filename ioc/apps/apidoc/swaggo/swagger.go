package swaggo

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/apidoc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	ioc_gin "github.com/infraboard/mcube/v2/ioc/config/gin"
	httpconf "github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
	"github.com/swaggo/swag"
)

func init() {
	ioc.Api().Registry(&SwaggerApiDoc{
		ApiDoc:       apidoc.DefaultApiDoc(),
		InstanceName: "swagger",
	})
}

type SwaggerApiDoc struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	*apidoc.ApiDoc
	InstanceName string `json:"instance_name" yaml:"instance_name" toml:"instance_name" env:"SWAGGER_INSTANCE_NAME"`
}

func (h *SwaggerApiDoc) Name() string {
	return apidoc.AppName
}

func (h *SwaggerApiDoc) Init() error {
	h.log = log.Sub("api_doc")
	if h.ApiDoc == nil {
		h.ApiDoc = apidoc.DefaultApiDoc()
	}
	def := apidoc.DefaultApiDoc()
	if strings.TrimSpace(h.JsonPath) == "" {
		h.JsonPath = def.JsonPath
	}
	if strings.TrimSpace(h.OpenAPIPath) == "" {
		h.OpenAPIPath = def.OpenAPIPath
	}
	if strings.TrimSpace(h.UIPath) == "" {
		h.UIPath = def.UIPath
	}
	if strings.TrimSpace(h.TryItUIPath) == "" {
		h.TryItUIPath = def.TryItUIPath
	}
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

func (h *SwaggerApiDoc) Registry() {
	r := ioc_gin.ObjectRouter(h)
	r.GET(h.JsonPath, h.SwaggerJson)
	r.GET(h.UIPath, h.SwaggerUI)
	r.GET(h.TryItUIPath, h.SwaggerTryItUI)
	h.log.Info().Msgf("Get the API JSON data using %s", h.ApiDocPath())
	h.log.Info().Msgf("Get the API UI using %s", h.ApiUIPath())
}

func (h *SwaggerApiDoc) ApiDocPath() string {
	if application.Get().AppAddress != "" {
		return strings.TrimRight(application.Get().AppAddress, "/") + h.ApiDocPathRel()
	}
	return httpconf.Get().ApiObjectAddr(h) + h.JsonPath
}

func (h *SwaggerApiDoc) ApiUIPath() string {
	if application.Get().AppAddress != "" {
		return strings.TrimRight(application.Get().AppAddress, "/") + h.ApiUIPathRel()
	}
	return httpconf.Get().ApiObjectAddr(h) + h.UIPath
}

func (h *SwaggerApiDoc) ApiDocPathRel() string {
	return apidoc.JoinURLPath(httpconf.Get().ApiObjectPathPrefix(h), h.JsonPath)
}

func (h *SwaggerApiDoc) ApiOpenAPIPathRel() string {
	return apidoc.JoinURLPath(httpconf.Get().ApiObjectPathPrefix(h), h.OpenAPIPath)
}

func (h *SwaggerApiDoc) ApiUIPathRel() string {
	return apidoc.JoinURLPath(httpconf.Get().ApiObjectPathPrefix(h), h.UIPath)
}

func (h *SwaggerApiDoc) ApiTryItPathRel() string {
	return apidoc.JoinURLPath(httpconf.Get().ApiObjectPathPrefix(h), h.TryItUIPath)
}

func (h *SwaggerApiDoc) SwaggerJson(ctx *gin.Context) {
	if !h.Enabled {
		ctx.String(http.StatusNotFound, "api doc disabled")
		return
	}
	ctx.Header("Content-Type", "application/json; charset=utf-8")
	ctx.Writer.WriteString(swag.GetSwagger(h.InstanceName).ReadDoc())
}

func (h *SwaggerApiDoc) SwaggerUI(ctx *gin.Context) {
	if !h.Enabled {
		ctx.String(http.StatusNotFound, "api doc disabled")
		return
	}
	switchHref := ""
	if h.EffectiveUIEngine() == apidoc.UIEngineBoth {
		switchHref = h.ApiTryItPathRel()
	}
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Writer.WriteString(apidoc.RenderScalarHTML(switchHref, h.ApiOpenAPIPathRel(), h.ScalarJS()))
}

func (h *SwaggerApiDoc) SwaggerTryItUI(ctx *gin.Context) {
	if !h.Enabled {
		ctx.String(http.StatusNotFound, "api doc disabled")
		return
	}
	switchHref := ""
	if h.EffectiveUIEngine() == apidoc.UIEngineBoth {
		switchHref = h.ApiUIPathRel()
	}
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Writer.WriteString(apidoc.RenderSwaggerUIHTML(switchHref, h.ApiOpenAPIPathRel(), h.SwaggerUIBase()))
}
