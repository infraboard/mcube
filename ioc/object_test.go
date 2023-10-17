package ioc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/infraboard/mcube/ioc"
)

var (
	obj *TestObject
)

func TestObjectLoad(t *testing.T) {
	conf := ioc.Default()
	conf.Registry(&TestObject{})
	fmt.Println(ioc.Default().List())

	os.Setenv("ATTR1", "a1")
	os.Setenv("ATTR2", "a2")
	ioc.DevelopmentSetup()

	obj = ioc.Default().Get("*ioc_test.TestObject").(*TestObject)
	t.Log(obj)
}

type TestObject struct {
	Attr1 string `toml:"attr1" env:"ATTR1"`
	Attr2 string `toml:"attr2" env:"ATTR2"`
	ioc.ObjectImpl
}

func (t *TestObject) Hello() string {
	return "hello"
}

type TestService interface {
	Hello() string
}
