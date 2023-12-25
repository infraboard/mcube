package application

import (
	"github.com/infraboard/mcube/v2/ioc"
)

const (
	AppName = "app"
)

func Get() *Application {
	return ioc.Config().Get(AppName).(*Application)
}
