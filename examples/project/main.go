package main

import (
	"context"

	// 非业务模块
	_ "github.com/infraboard/mcube/v2/ioc/apps/metric/gin"
	"github.com/infraboard/mcube/v2/ioc/server"

	// 加载业务模块
	_ "github.com/infraboard/mcube/v2/examples/project/apps/helloworld/api"
	_ "github.com/infraboard/mcube/v2/examples/project/apps/helloworld/impl"
)

func main() {
	server.DefaultConfig.ConfigFile.Enabled = true
	err := server.Run(context.Background())
	if err != nil {
		panic(err)
	}
}
