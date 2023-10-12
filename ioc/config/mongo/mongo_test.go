package mongo_test

import (
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/mongo"
)

func TestGetMongoDB(t *testing.T) {
	m := mongo.GetClientGetter()
	t.Log(m)
}

func init() {
	ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
}
