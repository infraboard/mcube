package mongo_test

import (
	"os"
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/mongo"
)

func TestGetClientGetter(t *testing.T) {
	m := mongo.Client()
	t.Log(m)
}

func init() {
	os.Setenv("MONGO_ENDPOINTS", "127.0.0.1:27017")
	os.Setenv("MONGO_USERNAME", "xxx")
	os.Setenv("MONGO_PASSWORD", "xxx")
	os.Setenv("MONGO_DATABASE", "xxx")
	os.Setenv("MONGO_AUTH_DB", "xxx")
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
