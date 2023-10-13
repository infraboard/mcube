package logger

import (
	"github.com/infraboard/mcube/ioc"
	"github.com/rs/zerolog"
)

const (
	LOG = "log"
)

const (
	SUB_LOGGER_KEY = "component"
)

type LoggerGetter interface {
	// 获取DB
	Logger(name string) *zerolog.Logger
}

func L(name string) *zerolog.Logger {
	return ioc.Config().Get(LOG).(LoggerGetter).Logger(name)
}
