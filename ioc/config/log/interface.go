package log

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"
)

const (
	AppName = "log"
)

const (
	SUB_LOGGER_KEY = "component"
)

func L() *zerolog.Logger {
	return ioc.Config().Get(AppName).(*Config).root
}

func Sub(name string) *zerolog.Logger {
	return ioc.Config().Get(AppName).(*Config).Logger(name)
}

func T(name string) *TraceLogger {
	return ioc.Config().Get(AppName).(*Config).TraceLogger(name)
}
