package main

import (
	"github.com/infraboard/mcube/v2/ioc/server/cmd"

	// metric 模块
	_ "github.com/infraboard/mcube/v2/ioc/apps/metric/gin"

	// 加载业务模块
	_ "github.com/infraboard/mcube/v2/examples/project/apps/helloworld/api"
	_ "github.com/infraboard/mcube/v2/examples/project/apps/helloworld/impl"
)

func main() {
	cmd.Start()
}
