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
	// DEFAULT_VERSION 为对象默认版本，推荐使用 v 前缀短形式（如 v1、v1.1.0），也兼容纯数字形式（如 1.0.0）。
	DEFAULT_VERSION = "v1"
)

// CompareVersion 比较两个语义化版本号。
// 支持: v1、v1.0.0、V1.1.1（v 前缀且下一字符为数字时去掉前缀）、以及 1.0.0 等形式。
// 返回值: 1 表示 v1 > v2, -1 表示 v1 < v2, 0 表示 v1 == v2
func CompareVersion(v1, v2 string) int {
	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	for i := range 3 {
		if parts1[i] > parts2[i] {
			return 1
		}
		if parts1[i] < parts2[i] {
			return -1
		}
	}
	return 0
}

// stripSemverVPrefix 去掉 Go 模块风格的 v 前缀：仅当首字符为 v/V 且下一字符为数字时剥离（如 v1、v1.0.0）。
// 不会剥离 "v1beta" 这类非 semver 形式（下一字符非数字则保留原样，解析时数字段可能为 0）。
func stripSemverVPrefix(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 1 && (s[0] == 'v' || s[0] == 'V') {
		c := s[1]
		if c >= '0' && c <= '9' {
			return s[1:]
		}
	}
	return s
}

// parseVersion 解析版本号为 [major, minor, patch]
func parseVersion(version string) [3]int {
	var result [3]int
	version = stripSemverVPrefix(version)
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
