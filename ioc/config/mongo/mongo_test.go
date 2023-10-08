package mongo_test

import (
	"testing"

	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/mongo"
)

func TestGetMongoDB(t *testing.T) {
	m := mongo.GetMongoDB()
	t.Log(m)
}

func init() {
	ioc.LoadConfig(ioc.NewLoadConfigRequest())
}
