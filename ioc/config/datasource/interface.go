package datasource

import (
	"context"

	"github.com/infraboard/mcube/ioc"
	"gorm.io/gorm"
)

const (
	AppName = "datasource"
)

func DB(ctx context.Context) *gorm.DB {
	return ioc.Config().Get(AppName).(*dataSource).GetTransactionOrDB(ctx)
}
