package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/mcube/v2/ioc/config/application"
)

// 正常请求数据返回
func Success(c *gin.Context, data any) {
	// 是否需要脱敏
	if v, ok := data.(response.DesenseObj); ok {
		v.Desense()
	}

	c.JSON(http.StatusOK, data)
}

// 异常情况的数据返回, 返回我们的业务Exception
func Failed(c *gin.Context, err error) {
	// 如果出现多个Handler， 需要通过手动abord
	defer c.Abort()

	var e *exception.APIException
	if v, ok := err.(*exception.APIException); ok {
		e = v
	} else {
		// 非可以预期, 没有定义业务的情况
		e = exception.NewAPIException(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		).WithMessage(err.Error())
		e.HttpCode = http.StatusInternalServerError

	}

	if e.Namespace == "" {
		e.WithNamespace(application.Get().AppName)
	}

	c.JSON(e.HttpCode, e)
}
