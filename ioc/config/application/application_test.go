package application_test

import (
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/application"
)

func TestGetClientGetter(t *testing.T) {
	m := application.App()
	t.Log(m)
}

func init() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	req.ConfigFile.Path = "test/application.toml"
	err := ioc.ConfigIocObject(req)
	if err != nil {
		panic(err)
	}
}
