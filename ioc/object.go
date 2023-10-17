package ioc

import (
	"context"
	"fmt"
)

func ObjectUid(o *ObjectWrapper) string {
	return fmt.Sprintf("%s.%s", o.Name, o.Version)
}

const (
	DEFAULT_VERSION = "v1"
)

type ObjectImpl struct {
}

func (i *ObjectImpl) Init() error {
	return nil
}

func (i *ObjectImpl) Name() string {
	return ""
}

func (i *ObjectImpl) Close(ctx context.Context) error {
	return nil
}

func (i *ObjectImpl) Version() string {
	return DEFAULT_VERSION
}

func (i *ObjectImpl) Priority() int {
	return 0
}

func (i *ObjectImpl) AllowOverwrite() bool {
	return false
}
