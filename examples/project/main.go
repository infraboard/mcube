package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/v2/ioc/config/grpc"
	ioc_http "github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/server"

	// metric 模块
	_ "github.com/infraboard/mcube/v2/ioc/apps/metric/gin"

	// 加载业务模块
	_ "github.com/infraboard/mcube/v2/examples/project/apps/helloworld/api"
	_ "github.com/infraboard/mcube/v2/examples/project/apps/helloworld/impl"
)

func main() {
	server.DefaultConfig.ConfigFile.Enabled = true
	server.DefaultConfig.ConfigFile.Path = "etc/application.toml"

	// 直接启动
	// err := server.Run(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// 启动前 设置
	server.SetUp(func() {
		// 开启GRPC服务注册与注销
		grpc.Get().PostStart = func(ctx context.Context) error {
			return nil
		}
		grpc.Get().PreStop = func(ctx context.Context) error {
			return nil
		}

		// HTTP业务路由加载之前
		rb := ioc_http.Get().GetRouterBuilder()
		rb.BeforeLoadHooks(func(r http.Handler) {
			// Gin框架
			if router, ok := r.(*restful.Container); ok {
				// GoRestful Router对象
				fmt.Println(router)
			}

			// GoRestful 框架
			if router, ok := r.(*gin.Engine); ok {
				// Gin Engine对象
				fmt.Println(router)
			}
		})
		// HTTP业务路由加载之后
		rb.AfterLoadHooks(func(r http.Handler) {
			// Gin框架
			if router, ok := r.(*restful.Container); ok {
				// GoRestful Router对象
				fmt.Println(router)
			}

			// GoRestful 框架
			if router, ok := r.(*gin.Engine); ok {
				// Gin Engine对象
				fmt.Println(router)
			}
		})

		// 补充Grpc 中间件
		grpc.Get().AddInterceptors()
	}).Run(context.Background())
}
