package application_test

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
)

func TestDefaultConfig(t *testing.T) {
	f, err := os.OpenFile("test/default.toml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	appConf := map[string]any{application.AppName: application.Get()}
	toml.NewEncoder(f).Encode(appConf)
}

func TestAppEnv(t *testing.T) {
	t.Log(application.Get().AppName)
}

func init() {
	os.Setenv("HTTP_ENABLE_TRACE", "false")
	os.Setenv("APP_NAME", "test")
	req := ioc.NewLoadConfigRequest()
	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}
