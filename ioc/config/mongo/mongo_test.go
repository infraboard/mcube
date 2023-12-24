package mongo_test

import (
	"os"
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/mongo"
	"github.com/infraboard/mcube/v2/tools/file"
)

func TestGetClientGetter(t *testing.T) {
	// 获取mongodb 客户端对象
	client := mongo.Client()
	t.Log(client)

	// 获取DB对象
	db := mongo.DB()
	t.Log(db)
}

func TestDefaultConfig(t *testing.T) {
	file.MustToToml(
		mongo.AppName,
		ioc.Config().Get(mongo.AppName),
		"test/default.toml",
	)
}

func init() {
	os.Setenv("MONGO_ENDPOINTS", "127.0.0.1:27017")
	os.Setenv("MONGO_USERNAME", "admin")
	os.Setenv("MONGO_PASSWORD", "123456")
	os.Setenv("MONGO_DATABASE", "admin")
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
