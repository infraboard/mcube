package datasource_test

import (
	"os"
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/datasource"
)

func TestGetDB(t *testing.T) {
	m := datasource.DB()
	t.Log(m)
}

func init() {
	os.Setenv("DATASOURCE_HOST", "127.0.0.1")
	os.Setenv("DATASOURCE_PORT", "3306")
	os.Setenv("DATASOURCE_DB", "xxx")
	os.Setenv("DATASOURCE_USERNAME", "xxx")
	os.Setenv("DATASOURCE_PASSWORD", "xxx")
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}