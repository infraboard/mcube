package ioc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/infraboard/mcube/ioc"
)

func TestObjectLoad(t *testing.T) {
	// 将*TestObject对象 注册到默认空间
	conf := ioc.Default()
	conf.Registry(&TestObject{})
	fmt.Println(ioc.Default().List())

	// 通过环境变量配置TestObject对象
	os.Setenv("ATTR1", "a1")
	os.Setenv("ATTR2", "a2")
	ioc.DevelopmentSetup()

	// 除了采用Get直接获取对象, 也可以通过Load动态加载, 等价于获取后赋值
	// ioc.Default().Get("*ioc_test.TestObject").(*TestObject)
	obj := &TestObject{}
	err := ioc.Default().Load(obj)
	if err != nil {
		t.Fatal(err)
	}
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
