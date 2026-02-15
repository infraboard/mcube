package ioc

import (
	"fmt"
	"strings"
)

// ParseInjectTag 解析IOC注入标签（向后兼容版本，忽略错误）
// 使用示例: `ioc:"autowire=true;namespace=default;name=myService;version=v2"`
//
// Deprecated: 推荐使用 ParseInjectTagWithError 以获得更好的错误处理
func ParseInjectTag(v string) *InjectTag {
	tag, err := ParseInjectTagWithError(v)
	if err != nil {
		// 向后兼容：发生错误时返回默认值而不是nil
		return NewInjectTag()
	}
	return tag
}

// ParseInjectTagWithError 解析IOC注入标签（返回错误版本）
// 支持的标签格式: "key=value;key=value"
//
// 支持的key:
//   - autowire: 是否自动注入 (true/false)
//   - namespace: 对象所在命名空间
//   - name: 注入对象的名称
//   - version: 注入对象的版本
//
// 示例:
//
//	tag, err := ParseInjectTagWithError("autowire=true;namespace=default")
func ParseInjectTagWithError(v string) (*InjectTag, error) {
	ins := NewInjectTag()

	v = strings.TrimSpace(v)
	if v == "" {
		return ins, nil
	}

	items := strings.Split(v, ";")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		// 使用SplitN避免值中包含'='的问题
		kv := strings.SplitN(item, "=", 2)
		key := strings.TrimSpace(kv[0])

		if key == "" {
			return nil, fmt.Errorf("empty key in tag: %q", item)
		}

		var value string
		if len(kv) > 1 {
			value = strings.TrimSpace(kv[1])
		}

		switch key {
		case "autowire":
			// 没有值或值为true时启用
			switch value {
			case "", "true":
				ins.Autowire = true
			case "false":
				ins.Autowire = false
			default:
				return nil, fmt.Errorf("invalid autowire value %q, expected true or false", value)
			}

		case "namespace":
			if value == "" {
				ins.Namespace = DEFAULT_NAMESPACE
			} else {
				ins.Namespace = value
			}

		case "name":
			if value == "" {
				return nil, fmt.Errorf("name value cannot be empty")
			}
			ins.Name = value

		case "version":
			if value == "" {
				ins.Version = DEFAULT_VERSION
			} else {
				ins.Version = value
			}

		default:
			return nil, fmt.Errorf("unknown tag key: %q", key)
		}
	}

	return ins, nil
}

func NewInjectTag() *InjectTag {
	return &InjectTag{
		Namespace: DEFAULT_NAMESPACE,
		Version:   DEFAULT_VERSION,
	}
}

type InjectTag struct {
	// 是否自动注入
	Autowire bool
	// 空间
	Namespace string
	// 注入对象的名称
	Name string
	// 注入对象的版本, 默认v1
	Version string
}
