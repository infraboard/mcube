package application_test

import (
	"os"
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
)

func TestGetClientGetter(t *testing.T) {
	m := application.App()
	t.Log(m.HTTP.EnableTrace)
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
