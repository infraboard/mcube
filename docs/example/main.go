package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/application"

	// 加载业务模块
	_ "github.com/infraboard/mcube/docs/example/helloworld/api"
	_ "github.com/infraboard/mcube/docs/example/helloworld/impl"
)

func main() {
	req := ioc.NewLoadConfigRequest()
	// 配置文件默认路径: etc/applicaiton.toml
	req.ConfigFile.Enabled = true
	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}

	// 打印ioc当前托管的对象
	fmt.Println(ioc.Config().List())
	fmt.Println(ioc.Config().List())
	fmt.Println(ioc.Api().List())

	// 将Ioc所有业务模块的路由 都注册给Gin Root Router
	app := application.App()
	r := gin.Default()
	ioc.LoadGinApi(app.HTTPPrefix(), r)

	// 启动HTTP服务
	// ioc中读取应用相关配置
	r.Run(app.HTTP.Addr())
}
