package gorestful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
)

func init() {
	ioc.Config().Registry(&GoRestfulFramework{
		Trace: true,
	})
}

type GoRestfulFramework struct {
	ioc.ObjectImpl
	Container *restful.Container
	log       *zerolog.Logger

	// 开启Trace
	Trace bool `toml:"trace" json:"trace" yaml:"trace" env:"TRACE"`
}

func (g *GoRestfulFramework) Priority() int {
	return 899
}

func (g *GoRestfulFramework) Name() string {
	return AppName
}

func (g *GoRestfulFramework) Init() error {
	g.log = log.Sub(g.Name())
	g.Container = restful.DefaultContainer
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)

	// 注册路由
	if g.Trace && trace.Get().Enable {
		g.log.Info().Msg("enable go-restful trace")
		g.Container.Filter(otelrestful.OTelFilter(application.Get().GetAppName()))
	}

	// 注册给Http服务器
	http.Get().SetRouter(g.Container)
	return nil
}
