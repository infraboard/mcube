package test

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/logger/zap"

	"{{.PKG}}/conf"

	// 注册所有服务
	_ "{{.PKG}}/apps"
)

func DevelopmentSetup() {
	// 初始化日志实例
	zap.DevelopmentSetup()

	// 针对http handler的测试需要提前设置默认数据格式
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)

	// 初始化配置, 提前配置好/etc/unit_test.env
	err := conf.LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}

	// 加载单元测试的变量
	LoadConfigFromEnv()

	// 初始化全局app
	if err := ioc.InitIocObject(); err != nil {
		panic(err)
	}
}
