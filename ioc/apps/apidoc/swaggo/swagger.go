package swaggo

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/apidoc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	ioc_gin "github.com/infraboard/mcube/v2/ioc/config/gin"
	"github.com/infraboard/mcube/v2/ioc/config/http"
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
	h.log.Info().Msgf("Get the API JSON data using %s", h.ApiDocPath())

	r.GET(h.UIPath, h.SwaggerUI)
	h.log.Info().Msgf("Get the API UI using %s", h.ApiUIPath())
}

func (h *SwaggerApiDoc) ApiDocPath() string {
	if application.Get().Domain != "" {
		return application.Get().Endpoint() + filepath.Join(http.Get().ApiObjectPathPrefix(h), h.JsonPath)
	}

	return http.Get().ApiObjectAddr(h) + h.JsonPath
}

func (h *SwaggerApiDoc) ApiUIPath() string {
	if application.Get().Domain != "" {
		return application.Get().Endpoint() + filepath.Join(http.Get().ApiObjectPathPrefix(h), h.UIPath)
	}

	return http.Get().ApiObjectAddr(h) + h.UIPath
}

func (h *SwaggerApiDoc) SwaggerJson(ctx *gin.Context) {
	ctx.Writer.WriteString(swag.GetSwagger(h.InstanceName).ReadDoc())
}

func (h *SwaggerApiDoc) SwaggerUI(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/html")
	ctx.Writer.WriteString(fmt.Sprintf(apidoc.HTML_REDOC, h.ApiDocPath()))
}
