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

	// 打印ioc当前托管的对象
	err = application.App().Start(context.Background())
	if err != nil {
		panic(err)
	}
}
