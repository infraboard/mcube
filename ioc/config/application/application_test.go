package application_test

import (
	"os"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
)

func TestGetClientGetter(t *testing.T) {
	m := application.App()
	t.Log(m.HTTP.EnableTrace)

}

func TestDefaultConfig(t *testing.T) {
	f, err := os.OpenFile("test/default.toml", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	appConf := map[string]any{application.AppName: application.App()}
	toml.NewEncoder(f).Encode(appConf)
}

func init() {
	os.Setenv("HTTP_ENABLE_TRACE", "false")
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "test/application.toml"
	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}
