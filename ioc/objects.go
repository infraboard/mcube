package ioc

import (
	"fmt"
)

func ObjectUid(o Object) string {
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
