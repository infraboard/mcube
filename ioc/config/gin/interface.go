package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/http"
)

const (
	AppName = "gin_webframework"
)

func RootRouter() *gin.Engine {
	return ioc.Config().Get(AppName).(*GinFramework).Engine
}

// 添加对象前缀路径
func ObjectRouter(obj ioc.Object) gin.IRouter {
	modulePath := http.Get().ApiObjectPathPrefix(obj)
	return ioc.Config().Get(AppName).(*GinFramework).Engine.Group(modulePath)
}
