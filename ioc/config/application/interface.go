package application

import (
	"github.com/infraboard/mcube/ioc"
)

const (
	APPLICATION = "app"
)

func App() *Application {
	return ioc.Config().Get(APPLICATION).(*Application)
}
