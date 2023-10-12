package datasource

import (
	"github.com/infraboard/mcube/ioc"
	"gorm.io/gorm"
)

const (
	DATASOURCE = "datasource"
)

type ClientGetter interface {
	// 获取DB
	GetDB() *gorm.DB
}

func GetClientGetter() ClientGetter {
	return ioc.Config().Get(DATASOURCE).(ClientGetter)
}
