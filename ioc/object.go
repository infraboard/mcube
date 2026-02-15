package ioc

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

func ObjectUid(o *ObjectWrapper) string {
	return fmt.Sprintf("%s.%s", o.Name, o.Version)
}

const (
	DEFAULT_VERSION = "1.0.0"
)

// CompareVersion 比较两个语义化版本号
// 返回值: 1 表示 v1 > v2, -1 表示 v1 < v2, 0 表示 v1 == v2
func CompareVersion(v1, v2 string) int {
	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	for i := 0; i < 3; i++ {
		if parts1[i] > parts2[i] {
			return 1
		}
		if parts1[i] < parts2[i] {
			return -1
		}
	}
	return 0
}

// parseVersion 解析版本号为 [major, minor, patch]
func parseVersion(version string) [3]int {
	var result [3]int
	parts := strings.Split(version, ".")

	for i := 0; i < len(parts) && i < 3; i++ {
		if num, err := strconv.Atoi(strings.TrimSpace(parts[i])); err == nil {
			result[i] = num
		}
	}
	return result
}

type ObjectImpl struct {
}

func (i *ObjectImpl) Init() error {
	return nil
}

func (i *ObjectImpl) Name() string {
	return ""
}

func (i *ObjectImpl) Close(ctx context.Context) {
}

func (i *ObjectImpl) Version() string {
	return DEFAULT_VERSION
}

func (i *ObjectImpl) Priority() int {
	return 0
}

func (i *ObjectImpl) Meta() ObjectMeta {
	return DefaultObjectMeta()
}

// 生命周期钩子默认实现

func (i *ObjectImpl) OnPostConfig() error {
	return nil
}

func (i *ObjectImpl) OnPreInit() error {
	return nil
}

func (i *ObjectImpl) OnPostInit() error {
	return nil
}

func (i *ObjectImpl) OnPreStop(ctx context.Context) error {
	return nil
}

func (i *ObjectImpl) OnPostStop(ctx context.Context) error {
	return nil
}

func DefaultObjectMeta() ObjectMeta {
	return ObjectMeta{
		CustomPathPrefix: "",
		Extra:            map[string]string{},
	}
}

type ObjectMeta struct {
	CustomPathPrefix string
	Extra            map[string]string
}
