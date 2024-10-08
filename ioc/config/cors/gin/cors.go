package gin

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/infraboard/mcube/v2/ioc"
	ioc_cors "github.com/infraboard/mcube/v2/ioc/config/cors"
	ioc_gin "github.com/infraboard/mcube/v2/ioc/config/gin"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(&CORS{
		CORS: &ioc_cors.CORS{
			Enabled:        true,
			AllowedHeaders: []string{"*"},
			AllowedOrigins: []string{"*"},
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
		m.AllowedOrigins = append(m.AllowedOrigins, "*")
	}
	if len(m.AllowedHeaders) == 0 {
		m.AllowedHeaders = append(m.AllowedHeaders, "*")
	}

	// 将中间件添加到Router中
	if m.Enabled {
		r := ioc_gin.RootRouter()
		r.Use(cors.New(cors.Config{
			AllowOrigins:     m.AllowedOrigins,
			AllowMethods:     m.AllowedMethods,
			AllowHeaders:     m.AllowedHeaders,
			ExposeHeaders:    m.ExposeHeaders,
			AllowCredentials: m.AllowCookies,
			MaxAge:           time.Duration(m.MaxAge) * time.Second,
			AllowWildcard:    true,
		}))
		m.log.Info().Msg("cors enabled")
	}

	return nil
}

func (i *CORS) Priority() int {
	return 288
}
