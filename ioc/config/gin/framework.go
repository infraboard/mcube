package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func init() {
	ioc.Config().Registry(&GinFramework{
		Recovery: true,
		Mode:     gin.DebugMode,
		Trace:    true,
	})
}

type GinFramework struct {
	ioc.ObjectImpl
	Engine *gin.Engine
	log    *zerolog.Logger

	// 开启recovery恢复
	Recovery bool `json:"recovery" yaml:"recovery" toml:"recovery" env:"RECOVERY"`
	// Gin模式
	Mode string `toml:"mode" json:"mode" yaml:"mode" env:"Mode"`
	// 开启Trace
	Trace bool `toml:"trace" json:"trace" yaml:"trace" env:"TRACE"`
}

func (g *GinFramework) Init() error {
	g.log = log.Sub(g.Name())
	g.Engine = gin.Default()
	gin.SetMode(g.Mode)

	if g.Recovery {
		g.log.Info().Msg("enable gin recovery")
		g.Engine.Use(gin.Recovery())
	}

	if g.Trace && trace.Get().Enable {
		g.log.Info().Msg("enable gin trace")
		g.Engine.Use(otelgin.Middleware(application.Get().GetAppName()))
	}

	// 注册给Http服务器
	http.Get().SetRouter(g.Engine)
	return nil
}

func (g *GinFramework) Priority() int {
	return 898
}

func (g *GinFramework) Name() string {
	return AppName
}
