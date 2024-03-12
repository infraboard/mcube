package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/server"

	// 引入生成好的API Doc代码
	_ "github.com/infraboard/mcube/v2/examples/http_gin/docs"
	// 引入集成工程
	_ "github.com/infraboard/mcube/v2/ioc/apps/apidoc/swaggo"

	// 开启Health健康检查
	_ "github.com/infraboard/mcube/v2/ioc/apps/health/gin"
	// 开启Metric
	_ "github.com/infraboard/mcube/v2/ioc/apps/metric/gin"
	// 开启CORS, 允许资源跨域共享
	_ "github.com/infraboard/mcube/v2/ioc/config/cors/gin"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// 注册HTTP接口类
	ioc.Api().Registry(&HelloServiceApiHandler{})

	// 启动应用
	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}

type HelloServiceApiHandler struct {
	// 继承自Ioc对象
	ioc.ObjectImpl
}

// 模块的名称, 会作为路径的一部分比如: /mcube_service/api/v1/hello_module/
// 路径构成规则 <service_name>/<path_prefix>/<service_version>/<module_name>
func (h *HelloServiceApiHandler) Name() string {
	return "hello_module"
}

func (h *HelloServiceApiHandler) Version() string {
	return "v1"
}

// API路由
func (h *HelloServiceApiHandler) Registry(r gin.IRouter) {
	r.GET("/", h.Hello)
}

// @Summary 修改文章标签
// @Description  修改文章标签
// @Tags         文章管理
// @Produce  json
// @Param id path int true "ID"
// @Param name query string true "ID"
// @Param state query int false "State"
// @Param modified_by query string true "ModifiedBy"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /api/v1/tags/{id} [put]
func (h *HelloServiceApiHandler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data": "hello mcube",
	})
}
