package datasource_test

import (
	"context"
	"os"
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/infraboard/mcube/v2/tools/file"
	"gorm.io/gorm"
)

var (
	ctx = context.Background()
)

func init() {
	os.Setenv("DATASOURCE_HOST", "127.0.0.1")
	os.Setenv("DATASOURCE_PORT", "3306")
	os.Setenv("DATASOURCE_DB", "test")
	os.Setenv("DATASOURCE_USERNAME", "root")
	os.Setenv("DATASOURCE_PASSWORD", "123456")
	os.Setenv("DATASOURCE_DEBUG", "true")
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}

func TestDefaultConfig(t *testing.T) {
	file.MustToToml(
		datasource.AppName,
		ioc.Config().Get(datasource.AppName),
		"test/default.toml",
	)
}

type TestStruct struct {
	Id     string `gorm:"column:id" json:"id"`
	FiledA string `gorm:"column:filed_a" json:"filed_a"`
}

func (s *TestStruct) TableName() string {
	return "test_transactions"
}

func TestTransaction(t *testing.T) {
	m := datasource.DB()
	t.Log(m)

	// 处理事务
	err := m.Transaction(func(tx *gorm.DB) error {
		// 处理自己逻辑
		tx.Save(&TestStruct{Id: "1", FiledA: "test"})

		// 处理其他业务逻辑
		txCtx := datasource.WithTransactionCtx(ctx, tx)
		if err := Tx1(txCtx); err != nil {
			return err
		}
		if err := Tx2(txCtx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// 如果是事务则当前操作在事务中, 如果不是就是一个普通操作,回及时生效
func Tx1(ctx context.Context) error {
	db := datasource.DBFromCtx(ctx)
	db.Save(&TestStruct{Id: "2", FiledA: "test"})
	return nil
}

func Tx2(ctx context.Context) error {
	db := datasource.DBFromCtx(ctx)
	db.Save(&TestStruct{Id: "3", FiledA: "test"})
	return nil
}
