package application

import (
	"github.com/infraboard/mcube/v2/ioc"
)

const (
	AppName = "app"
)

func App() *Application {
	return ioc.Config().Get(AppName).(*Application)
}
