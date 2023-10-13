package datasource

import (
	"github.com/infraboard/mcube/ioc"
	"gorm.io/gorm"
)

const (
	DATASOURCE = "datasource"
)

func DB() *gorm.DB {
	return ioc.Config().Get(DATASOURCE).(*dataSource).db
}
