package ioc_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/infraboard/mcube/ioc"
)

func TestStructAutowire(t *testing.T) {
	c := ioc.Default()
	c.Registry(&InjectObject{})
	c.Registry(&TestObject{})
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

	// 支持的标签参数: autowire=true;namespace=default;name=*ioc_test.InjectObject;version=v1
	Obj *TestObject `ioc:"autowire=true;namespace=default"`
	ioc.ObjectImpl
}

func TestInterfaceAutowire(t *testing.T) {
	c := ioc.Default()
	c.Registry(&InjectInterface{})
	c.Registry(&TestObject{})
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
