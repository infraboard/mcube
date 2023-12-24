package datasource

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"gorm.io/gorm"
)

const (
	AppName = "datasource"
)

func DB() *gorm.DB {
	return ioc.Config().Get(AppName).(*dataSource).db
}

// 从上下文中获取事物, 如果获取不到则直接返回 无事物的DB对象
func DBFromCtx(ctx context.Context) *gorm.DB {
	return ioc.Config().Get(AppName).(*dataSource).GetTransactionOrDB(ctx)
}
