package ioc

import (
	"fmt"
)

// IocObject 内部服务实例, 不需要暴露
type IocObject interface {
	// 对象初始化
	Init() error
	// 对象的名称
	Name() string
	// 对象版本
	Version() string
	// 对象优先级
	Priority() int
	// 对象的销毁
	Destory()
	// 是否允许同名对象被覆盖, 默认不允许被覆盖
	AllowOverwrite() bool
}

func ObjectUid(o IocObject) string {
	return fmt.Sprintf("%s.%s", o.Name(), o.Version())
}

const (
	DEFAULT_VERSION = "v1"
)

type IocObjectImpl struct {
}

func (i *IocObjectImpl) Init() error {
	return nil
}

func (i *IocObjectImpl) Name() string {
	return ""
}

func (i *IocObjectImpl) Destory() {
}

func (i *IocObjectImpl) Version() string {
	return DEFAULT_VERSION
}

func (i *IocObjectImpl) Priority() int {
	return 0
}

func (i *IocObjectImpl) AllowOverwrite() bool {
	return false
}
