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
	return Get().root
}

func Sub(name string) *zerolog.Logger {
	return Get().Logger(name)
}

func T(name string) *TraceLogger {
	return Get().TraceLogger(name)
}

func Get() *Config {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return ioc.Config().Get(AppName).(*Config)
}
