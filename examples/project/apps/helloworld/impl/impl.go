package impl

import (
	"errors"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/datasource"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
)

func init() {
	ioc.Controller().Registry(&HelloServiceImpl{})
}

// 业务逻辑实现类
type HelloServiceImpl struct {
	db       *gorm.DB
	colector *EventCollect

	ioc.ObjectImpl
}

// 控制器初始化
func (i *HelloServiceImpl) Init() error {
	// 从Ioc总获取GORM DB对象, GORM相关配置已经托管给Ioc
	// Ioc会负责GORM的配置读取和为你初始化DB对象实例,以及关闭
	i.db = datasource.DB()

	// 将采集器注册到默认注册表
	i.colector = NewEventCollect()
	prometheus.MustRegister(i.colector)
	return nil
}

// 具体业务逻辑
func (i *HelloServiceImpl) Hello() string {

	// 模拟存储失败报错, 然后调用采集器纪录状态
	err := errors.New("save event error")
	if err != nil {
		i.colector.Inc()
	}

	return "hello world"
}
