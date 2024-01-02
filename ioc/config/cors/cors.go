package cors

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	ioc_http "github.com/infraboard/mcube/v2/ioc/config/http"
)

func init() {
	ioc.Config().Registry(&CORS{
		Enabled:        true,
		AllowedHeaders: []string{"*"},
		AllowedDomains: []string{"*"},
		AllowedMethods: []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
	})
}

type CORS struct {
	ioc.ObjectImpl

	Enabled        bool     `toml:"enabled" json:"enabled" yaml:"enabled"  env:"CORS_ENABLED"`
	AllowedHeaders []string `json:"cors_allowed_headers" yaml:"cors_allowed_headers" toml:"cors_allowed_headers" env:"CORS_ALLOWED_HEADERS" envSeparator:","`
	AllowedDomains []string `json:"cors_allowed_domains" yaml:"cors_allowed_domains" toml:"cors_allowed_domains" env:"CORS_ALLOWED_DOMAINS" envSeparator:","`
	AllowedMethods []string `json:"cors_allowed_methods" yaml:"cors_allowed_methods" toml:"cors_allowed_methods" env:"CORS_ALLOWED_METHODS" envSeparator:","`
}

func (m *CORS) Name() string {
	return AppName
}

func (m *CORS) Init() error {
	rb := ioc_http.Get().GetRouterBuilder()

	// 将中间件添加到Router中
	if m.Enabled {
		rb.BeforeLoadHooks(func(h http.Handler) {
			switch r := h.(type) {
			case *restful.Container:
				cors := restful.CrossOriginResourceSharing{
					AllowedHeaders: m.AllowedHeaders,
					AllowedDomains: m.AllowedDomains,
					AllowedMethods: m.AllowedMethods,
					CookiesAllowed: false,
					Container:      r,
				}
				r.Filter(cors.Filter)
			}
		})
	}
	return nil
}

func (i *CORS) Priority() int {
	return -10
}
