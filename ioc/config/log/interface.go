package log

import (
	"fmt"
	"runtime"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/rs/zerolog"
)

const (
	AppName = "log"
)

const (
	SUB_LOGGER_KEY = "logger"
)

func L() *zerolog.Logger {
	return Get().root
}

func Sub(name string) *zerolog.Logger {
	return Get().Logger(name)
}

func Get() *Config {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return ioc.Config().Get(AppName).(*Config)
}

const RECOVERY_BUF_SIZE = 64 << 1

func Recover() {
	if r := recover(); r != nil {
		buf := make([]byte, RECOVERY_BUF_SIZE)
		buf = buf[:runtime.Stack(buf, false)]
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		L().Error().Msgf("pannic: %s ... \n %s", err, string(buf))
	}
}
