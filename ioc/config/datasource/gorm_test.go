package datasource_test

import (
	"context"
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
)

var (
	ctx = context.Background()
)

func TestGetDB(t *testing.T) {
	m := datasource.DB(ctx)
	t.Log(m)

	//
	tx := datasource.BeginTransaction(ctx)
	defer datasource.EndTransaction(tx, nil)

	// tx 业务处理

	// 调用其他服务接口
	txCtx := datasource.WithTransactionCtx(ctx, tx)
	// svc.XXX(newCtx)
	// datasource.DB(ctx)
	t.Log(txCtx)
}

func TestDefaultConfig(t *testing.T) {
	f, err := os.OpenFile("test/default.toml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	appConf := map[string]any{datasource.AppName: ioc.Config().Get(datasource.AppName)}
	toml.NewEncoder(f).Encode(appConf)
}

func init() {
	os.Setenv("DATASOURCE_HOST", "127.0.0.1")
	os.Setenv("DATASOURCE_PORT", "3306")
	os.Setenv("DATASOURCE_DB", "test")
	os.Setenv("DATASOURCE_USERNAME", "root")
	os.Setenv("DATASOURCE_PASSWORD", "123456")
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
