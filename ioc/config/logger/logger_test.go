package logger_test

import (
	"os"
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
)

func TestGetClientGetter(t *testing.T) {
	sub := logger.Sub("module_a")
	sub.Debug().Msgf("hello %s", "a")

	logger.L().Debug().Msg("test")
}

func init() {
	os.Setenv("LOG_CONSOLE_NO_COLOR", "true")
	os.Setenv("LOG_TO_FILE", "true")
	os.Setenv("LOG_FILE_PATH", "/tmp/test.log")
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
