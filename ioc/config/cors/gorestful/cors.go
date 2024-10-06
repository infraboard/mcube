package gorestful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	ioc_cors "github.com/infraboard/mcube/v2/ioc/config/cors"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(&CORS{
		CORS: &ioc_cors.CORS{
			Enabled:        true,
			AllowedHeaders: []string{".*"},
			AllowedOrigins: []string{".*"},
			AllowedMethods: []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
			MaxAge:         12 * 60 * 60,
		},
	})
}

type CORS struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	*ioc_cors.CORS
}

func (m *CORS) Name() string {
	return AppName
}

func (m *CORS) Init() error {
	m.log = log.Sub("cors")

	if len(m.AllowedOrigins) == 0 {
		m.AllowedOrigins = append(m.AllowedOrigins, ".*")
	}
	if len(m.AllowedHeaders) == 0 {
		m.AllowedHeaders = append(m.AllowedHeaders, ".*")
	}

	// 将中间件添加到Router中
	r := gorestful.RootRouter()
	if m.Enabled {
		cors := restful.CrossOriginResourceSharing{
			AllowedHeaders: m.AllowedHeaders,
			AllowedDomains: m.AllowedOrigins,
			AllowedMethods: m.AllowedMethods,
			ExposeHeaders:  m.ExposeHeaders,
			CookiesAllowed: m.AllowCookies,
			MaxAge:         m.MaxAge,
			Container:      r,
		}
		r.Filter(cors.Filter)
		m.log.Info().Msg("cors enabled")
	}

	return nil
}

func (i *CORS) Priority() int {
	return 289
}
