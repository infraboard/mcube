package logger

import (
	"github.com/infraboard/mcube/ioc"
	"github.com/rs/zerolog"
)

const (
	AppName = "log"
)

const (
	SUB_LOGGER_KEY = "component"
)

func Sub(name string) *zerolog.Logger {
	return ioc.Config().Get(AppName).(*Config).Logger(name)
}

func L() *zerolog.Logger {
	return ioc.Config().Get(AppName).(*Config).root
}
