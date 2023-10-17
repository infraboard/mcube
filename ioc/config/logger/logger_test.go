package logger_test

import (
	"os"
	"runtime/debug"
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
)

func TestGetClientGetter(t *testing.T) {
	sub := logger.Sub("module_a")
	sub.Debug().Msgf("hello %s", "a")
}

func TestPanicStack(t *testing.T) {
	// 捕获 panic
	defer func() {
		if err := recover(); err != nil {
			logger.L().Error().Stack().Msgf("Panic occurred: %v\n%s", err, debug.Stack())
		}
	}()

	// 代码中可能发生 panic 的地方
	panic("Something went wrong!")
}

func init() {
	os.Setenv("LOG_CONSOLE_NO_COLOR", "true")
	os.Setenv("LOG_TO_FILE", "true")
	os.Setenv("LOG_CALLER_DEEP", "2")
	os.Setenv("LOG_FILE_PATH", "/tmp/test.log")
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
