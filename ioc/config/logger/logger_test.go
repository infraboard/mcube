package logger_test

import (
	"os"
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
)

func TestGetClientGetter(t *testing.T) {
	m := logger.L("module_a")
	m.Debug().Msgf("hello %s", "a")
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
