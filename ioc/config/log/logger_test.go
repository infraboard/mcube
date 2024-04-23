package log_test

import (
	"context"
	"os"
	"runtime/debug"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
)

func TestGetClientGetter(t *testing.T) {
	sub := log.Sub("module_a")
	log.T("module_a").Trace(context.Background())
	sub.Debug().Msgf("hello %s", "a")
}

func TestDefaultConfig(t *testing.T) {
	f, err := os.OpenFile("test/default.toml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	appConf := map[string]any{log.AppName: ioc.Config().Get(log.AppName).(*log.Config)}
	toml.NewEncoder(f).Encode(appConf)
}

func TestPanicStack(t *testing.T) {
	// 捕获 panic
	defer func() {
		if err := recover(); err != nil {
			log.L().Error().Stack().Msgf("Panic occurred: %v\n%s", err, debug.Stack())
		}
	}()

	// 代码中可能发生 panic 的地方
	panic("Something went wrong!")
}

func TestAppEnv(t *testing.T) {
	t.Log(log.Get().Level)
	t.Log(log.Get().File.FilePath)
}

func init() {
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("LOG_FILE_PATH", "/test")
	req := ioc.NewLoadConfigRequest()
	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}
