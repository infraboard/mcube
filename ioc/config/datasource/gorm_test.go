package datasource_test

import (
	"context"
	"os"
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/datasource"
)

var (
	ctx = context.Background()
)

func TestGetDB(t *testing.T) {
	m := datasource.DB(ctx)
	t.Log(m)

	tx := m.Begin().WithContext(ctx)
	defer datasource.EndTransaction(tx, nil)

	// tx 业务处理

	// 调用其他服务接口
	txCtx := datasource.WithTransactionCtx(ctx, tx)
	// svc.XXX(newCtx)
	// datasource.DB(ctx)
	t.Log(txCtx)
}

func init() {
	os.Setenv("DATASOURCE_HOST", "127.0.0.1")
	os.Setenv("DATASOURCE_PORT", "3306")
	os.Setenv("DATASOURCE_DB", "xxx")
	os.Setenv("DATASOURCE_USERNAME", "xxx")
	os.Setenv("DATASOURCE_PASSWORD", "xxx")
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
