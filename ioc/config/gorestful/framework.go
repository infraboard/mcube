package gorestful

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/http"
)

func init() {
	ioc.Config().Registry(&GoRestfulFramework{})
}

type GoRestfulFramework struct {
	ioc.ObjectImpl

	Container *restful.Container
}

func (i *GoRestfulFramework) Priority() int {
	return 99
}

func (g *GoRestfulFramework) Init() error {
	g.Container = restful.DefaultContainer
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)

	// 注册给Http服务器
	http.Get().SetRouter(g.Container)
	return nil
}

func (g *GoRestfulFramework) Name() string {
	return AppName
}
