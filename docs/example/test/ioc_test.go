package test_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/datasource"
)

var (
	ctx = context.Background()
)

func TestIocLoad(t *testing.T) {
	// 查询注册的对象列表, 通过导入datasource库 完成注册
	t.Log(ioc.Config().List())

	// 加载配置
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "../etc/application.toml"
	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}

	// 使用ioc对象(datasource配置 衍生对象)
	// ioc.Config().Get(DATASOURCE).(*dataSource).db
	t.Log(datasource.DB(ctx))
}
