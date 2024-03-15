package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/http"
)

func init() {
	ioc.Config().Registry(&GinFramework{})
}

type GinFramework struct {
	ioc.ObjectImpl
	Engine *gin.Engine
}

func (g *GinFramework) Init() error {
	g.Engine = gin.Default()

	// 注册给Http服务器
	http.Get().SetRouter(g.Engine)
	return nil
}

func (i *GinFramework) Priority() int {
	return 99
}

func (g *GinFramework) Name() string {
	return AppName
}
