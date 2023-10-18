package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/docs/example/helloworld"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&HelloServiceApiHandler{})
}

// 3. 暴露HTTP接口
type HelloServiceApiHandler struct {
	// 依赖业务控制器
	// 使用ioc注解来从自动加载依赖对象, 等同于手动执行:
	// 	h.svc = ioc.Controller().Get("*impl.HelloService").(helloworld.HelloService)
	Svc helloworld.HelloService `ioc:"autowire=true;namespace=controllers"`

	// 日志相关配置已经托管到Ioc中, 由于是私有属性，所有受到注入, 具体见下边初始化方法
	log *zerolog.Logger

	// 继承自Ioc对象
	ioc.ObjectImpl
}

// 对象自定义初始化
func (h *HelloServiceApiHandler) Init() error {
	h.log = logger.Sub("helloworld.api")
	return nil
}

// API路由
func (h *HelloServiceApiHandler) Registry(r gin.IRouter) {
	r.GET("/", h.Hello)
}

// API接口具体实现
func (h *HelloServiceApiHandler) Hello(c *gin.Context) {
	// 业务处理
	resp := h.Svc.Hello()
	h.log.Debug().Msg(resp)

	// 业务响应
	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}
