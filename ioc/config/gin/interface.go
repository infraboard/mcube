package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
)

const (
	AppName = "gin_webframework"
)

func Get() *gin.Engine {
	return ioc.Config().Get(AppName).(*GinFramework).Engine
}
