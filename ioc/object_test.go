package ioc_test

import (
	"os"
	"testing"

	"github.com/infraboard/mcube/ioc"
)

var (
	obj *TestObject
)

func TestObjectLoad(t *testing.T) {
	os.Setenv("ATTR1", "a1")
	os.Setenv("ATTR2", "a2")

	if err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest()); err != nil {
		panic(err)
	}

	t.Log(obj)
}

type TestObject struct {
	Attr1 string `toml:"attr1" env:"ATTR1"`
	Attr2 string `toml:"attr2" env:"ATTR2"`
	ioc.ObjectImpl
}

func (o *TestObject) Name() string {
	return "test_obj"
}

func init() {
	conf := ioc.Config()
	conf.Registry(&TestObject{})
	obj = ioc.Config().Get("test_obj").(*TestObject)
}
