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
		CORS: ioc_cors.Default(),
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

	if len(m.AllowedDomains) == 0 {
		m.AllowedDomains = append(m.AllowedDomains, ".*")
	}
	if len(m.AllowedHeaders) == 0 {
		m.AllowedHeaders = append(m.AllowedHeaders, ".*")
	}

	// 将中间件添加到Router中
	r := gorestful.RootRouter()
	if m.Enabled {
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
		m.log.Info().Msg("cors enabled")
	}

	return nil
}

func (i *CORS) Priority() int {
	return -10
}
