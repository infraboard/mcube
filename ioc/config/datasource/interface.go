package datasource

import (
	"context"

	"github.com/infraboard/mcube/ioc"
	"gorm.io/gorm"
)

const (
	DATASOURCE = "datasource"
)

func DB(ctx context.Context) *gorm.DB {
	return ioc.Config().Get(DATASOURCE).(*dataSource).GetTransactionOrDB(ctx)
}
