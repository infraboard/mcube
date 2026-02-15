package ioc_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
)

func TestCompareVersion(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		// 相等情况
		{name: "equal versions", v1: "1.0.0", v2: "1.0.0", expected: 0},
		{name: "equal major versions", v1: "2.0.0", v2: "2.0.0", expected: 0},
		{name: "equal complex versions", v1: "1.2.3", v2: "1.2.3", expected: 0},

		// v1 > v2
		{name: "major version greater", v1: "2.0.0", v2: "1.0.0", expected: 1},
		{name: "minor version greater", v1: "1.2.0", v2: "1.1.0", expected: 1},
		{name: "patch version greater", v1: "1.0.3", v2: "1.0.2", expected: 1},
		{name: "complex greater", v1: "2.1.5", v2: "2.1.4", expected: 1},
		{name: "major dominates", v1: "2.0.0", v2: "1.9.9", expected: 1},

		// v1 < v2
		{name: "major version less", v1: "1.0.0", v2: "2.0.0", expected: -1},
		{name: "minor version less", v1: "1.1.0", v2: "1.2.0", expected: -1},
		{name: "patch version less", v1: "1.0.2", v2: "1.0.3", expected: -1},
		{name: "complex less", v1: "2.1.4", v2: "2.1.5", expected: -1},

		// 边界情况
		{name: "zero versions", v1: "0.0.0", v2: "0.0.0", expected: 0},
		{name: "zero vs one", v1: "0.0.1", v2: "0.0.0", expected: 1},
		{name: "large numbers", v1: "10.20.30", v2: "10.20.29", expected: 1},
		{name: "very large major", v1: "100.0.0", v2: "99.99.99", expected: 1},

		// 不规范格式（容错处理）
		{name: "missing patch v1", v1: "1.0", v2: "1.0.0", expected: 0},
		{name: "missing patch v2", v1: "1.0.0", v2: "1.0", expected: 0},
		{name: "single number", v1: "2", v2: "1.0.0", expected: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ioc.CompareVersion(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("CompareVersion(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestRegistryVersionOverwrite(t *testing.T) {
	// 测试版本覆盖策略
	t.Run("higher version overwrites lower version", func(t *testing.T) {
		store := ioc.DefaultStore.Namespace("test_version_overwrite_high")

		// 注册 1.0.0 版本
		obj1 := &TestVersionObject{name: "test_obj_v1", version: "1.0.0"}
		store.Registry(obj1)

		// 注册 2.0.0 版本（应该成功覆盖）
		obj2 := &TestVersionObject{name: "test_obj_v1", version: "2.0.0"}
		store.Registry(obj2)

		// 验证已被覆盖为新版本
		got := store.Get("test_obj_v1")
		if got.Version() != "2.0.0" {
			t.Errorf("Expected version 2.0.0, got %s", got.Version())
		}
	})

	t.Run("same version panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for duplicate version")
			}
		}()

		store := ioc.DefaultStore.Namespace("test_same_version")
		obj1 := &TestVersionObject{name: "test_obj_v2", version: "1.0.0"}
		store.Registry(obj1)

		// 注册相同版本（应该 panic）
		obj2 := &TestVersionObject{name: "test_obj_v2", version: "1.0.0"}
		store.Registry(obj2)
	})

	t.Run("lower version panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for lower version")
			}
		}()

		store := ioc.DefaultStore.Namespace("test_lower_version")

		// 注册 2.0.0 版本
		obj1 := &TestVersionObject{name: "test_obj_v3", version: "2.0.0"}
		store.Registry(obj1)

		// 注册 1.0.0 版本（应该 panic）
		obj2 := &TestVersionObject{name: "test_obj_v3", version: "1.0.0"}
		store.Registry(obj2)
	})
}

// TestVersionObject 用于版本测试的对象
type TestVersionObject struct {
	ioc.ObjectImpl
	name    string
	version string
}

func (o *TestVersionObject) Name() string {
	return o.name
}

func (o *TestVersionObject) Version() string {
	return o.version
}
