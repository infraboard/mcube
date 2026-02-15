package ioc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
)

func TestStructAutowire(t *testing.T) {
	c := ioc.Default()
	c.Registry(&InjectObject{})
	// 只在对象不存在时才注册，避免与其他测试冲突
	if c.Get("*ioc_test.TestObject") == nil {
		c.Registry(&TestObject{})
	}
	fmt.Println(ioc.Default().List())

	os.Setenv("ATTR1", "a1")
	os.Setenv("ATTR2", "a2")
	ioc.DevelopmentSetup()

	v := ioc.Default().Get("*ioc_test.InjectObject").(*InjectObject)
	fmt.Println(v.Obj)
}

type InjectObject struct {
	Attr1 string `toml:"attr1" env:"ATTR1"`
	Attr2 string `toml:"attr2" env:"ATTR2"`

	// 支持的标签参数: autowire=true;namespace=default;name=*ioc_test.InjectObject;version=1.0.0
	Obj *TestObject `ioc:"autowire=true;namespace=default"`
	ioc.ObjectImpl
}

func TestInterfaceAutowire(t *testing.T) {
	c := ioc.Default()
	c.Registry(&InjectInterface{})
	// 只在对象不存在时才注册，避免与其他测试冲突
	if c.Get("*ioc_test.TestObject") == nil {
		c.Registry(&TestObject{})
	}
	fmt.Println(ioc.Default().List())

	os.Setenv("ATTR1", "a1")
	os.Setenv("ATTR2", "a2")
	ioc.DevelopmentSetup()

	v := ioc.Default().Get("*ioc_test.InjectInterface").(*InjectInterface)
	fmt.Println(v.Obj)
}

type InjectInterface struct {
	Attr1 string `toml:"attr1" env:"ATTR1"`
	Attr2 string `toml:"attr2" env:"ATTR2"`

	Obj TestService `ioc:"autowire=true;namespace=default"`
	ioc.ObjectImpl
}
