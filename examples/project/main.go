package main

import (
	"context"

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
	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
