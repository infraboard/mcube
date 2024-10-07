package swaggo

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/apidoc"
	ioc_gin "github.com/infraboard/mcube/v2/ioc/config/gin"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
	"github.com/swaggo/swag"
)

func init() {
	ioc.Api().Registry(&SwaggerApiDoc{
		ApiDoc:       apidoc.ApiDoc{},
		InstanceName: "swagger",
	})
}

type SwaggerApiDoc struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	apidoc.ApiDoc
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
	meta.CustomPathPrefix = h.BasePath
	return meta
}

func (h *SwaggerApiDoc) Registry() {
	r := ioc_gin.ObjectRouter(h)
	r.GET("/", func(c *gin.Context) {
		c.Writer.WriteString(swag.GetSwagger(h.InstanceName).ReadDoc())
	})

	h.log.Info().Msgf("Get the API Doc using %s", http.Get().ApiObjectAddr(h))
}
