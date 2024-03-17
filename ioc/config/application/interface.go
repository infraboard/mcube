package application

import (
	"github.com/infraboard/mcube/v2/ioc"
)

const (
	AppName = "app"
)

func Get() *Application {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Application)
}
