package impl

import (
	"context"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/datasource"
	"gorm.io/gorm"
)

func init() {
	ioc.Controller().Registry(&HelloServiceImpl{})
}

// 业务逻辑实现类
type HelloServiceImpl struct {
	db *gorm.DB

	ioc.ObjectImpl
}

// 控制器初始化
func (i *HelloServiceImpl) Init() error {
	// 从Ioc总获取GORM DB对象, GORM相关配置已经托管给Ioc
	// Ioc会负责GORM的配置读取和为你初始化DB对象实例,以及关闭
	i.db = datasource.DB(context.Background())
	return nil
}

// 具体业务逻辑
func (i *HelloServiceImpl) Hello() string {
	return "hello world"
}
