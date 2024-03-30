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
	ioc.Config().Registry(&GinFramework{})
}

type GinFramework struct {
	ioc.ObjectImpl
	Engine *gin.Engine
	log    *zerolog.Logger
}

func (g *GinFramework) Init() error {
	g.log = log.Sub(g.Name())
	g.Engine = gin.Default()

	if http.Get().EnableTrace && trace.Get().Enable {
		g.log.Info().Msg("enable gin trace")
		g.Engine.Use(otelgin.Middleware(application.Get().GetAppNameWithDefault("default")))
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
