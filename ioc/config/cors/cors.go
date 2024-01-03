package cors

import (
	"net/http"
	"time"

	"github.com/emicklei/go-restful/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	ioc_http "github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(&CORS{
		Enabled:        true,
		AllowedHeaders: []string{".*"},
		AllowedDomains: []string{".*"},
		AllowedMethods: []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
		MaxAge:         12 * 60 * 60,
	})
}

type CORS struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	Enabled        bool     `toml:"enabled" json:"enabled" yaml:"enabled"  env:"HTTP_CORS_ENABLED"`
	AllowedHeaders []string `json:"cors_allowed_headers" yaml:"cors_allowed_headers" toml:"cors_allowed_headers" env:"HTTP_CORS_ALLOWED_HEADERS" envSeparator:","`
	AllowedDomains []string `json:"cors_allowed_domains" yaml:"cors_allowed_domains" toml:"cors_allowed_domains" env:"HTTP_CORS_ALLOWED_DOMAINS" envSeparator:","`
	AllowedMethods []string `json:"cors_allowed_methods" yaml:"cors_allowed_methods" toml:"cors_allowed_methods" env:"HTTP_CORS_ALLOWED_METHODS" envSeparator:","`
	ExposeHeaders  []string `json:"cors_expose_headers" yaml:"cors_expose_headers" toml:"cors_expose_headers" env:"HTTP_CORS_EXPOSE_HEADERS" envSeparator:","`
	AllowCookies   bool     `toml:"cors_allow_cookies" json:"cors_allow_cookies" yaml:"cors_allow_cookies"  env:"HTTP_CORS_ALLOW_COOKIES"`
	// 单位秒, 默认12小时
	MaxAge int `toml:"max_age" json:"max_age" yaml:"max_age"  env:"HTTP_CORS_MAX_AGE"`
}

func (m *CORS) Name() string {
	return AppName
}

func (m *CORS) goRestfulDefault() {
	if len(m.AllowedDomains) == 0 {
		m.AllowedDomains = append(m.AllowedDomains, ".*")
	}
	if len(m.AllowedHeaders) == 0 {
		m.AllowedHeaders = append(m.AllowedHeaders, ".*")
	}
}

func (m *CORS) ginRestfulDefault() {
	if len(m.AllowedDomains) == 0 {
		m.AllowedDomains = append(m.AllowedDomains, "*")
	}
	if len(m.AllowedHeaders) == 0 {
		m.AllowedHeaders = append(m.AllowedHeaders, "*")
	}
}

func (m *CORS) Init() error {
	m.log = log.Sub("cors")
	rb := ioc_http.Get().GetRouterBuilder()

	// 将中间件添加到Router中
	if m.Enabled {
		rb.BeforeLoadHooks(func(h http.Handler) {
			switch r := h.(type) {
			case *restful.Container:
				m.goRestfulDefault()
				cors := restful.CrossOriginResourceSharing{
					AllowedHeaders: m.AllowedHeaders,
					AllowedDomains: m.AllowedDomains,
					AllowedMethods: m.AllowedMethods,
					ExposeHeaders:  m.ExposeHeaders,
					CookiesAllowed: m.AllowCookies,
					MaxAge:         m.MaxAge,
					Container:      r,
				}
				r.Filter(cors.Filter)
			case *gin.Engine:
				m.ginRestfulDefault()
				r.Use(cors.New(cors.Config{
					AllowOrigins:     m.AllowedDomains,
					AllowMethods:     m.AllowedMethods,
					AllowHeaders:     m.AllowedHeaders,
					ExposeHeaders:    m.ExposeHeaders,
					AllowCredentials: m.AllowCookies,
					MaxAge:           time.Duration(m.MaxAge) * time.Second,
					AllowWildcard:    true,
				}))
			}
			m.log.Info().Msg("cors enabled")
		})
	}

	return nil
}

func (i *CORS) Priority() int {
	return -10
}
