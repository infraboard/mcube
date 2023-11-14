package main

import (
	"context"

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

	// 启动应用, 应用会自动加载 刚才实现的Gin Api Handler
	err = application.App().Start(context.Background())
	if err != nil {
		panic(err)
	}
}
