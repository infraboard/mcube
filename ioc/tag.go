package ioc

import "strings"

// autowire=true
func ParseInjectTag(v string) *InjectTag {
	ins := NewInjectTag()

	v = strings.TrimSpace(v)
	items := strings.Split(v, ";")
	for i := range items {
		kv := strings.Split(items[i], "=")
		switch kv[0] {
		case "autowire":
			ins.Autowire = true
			if len(kv) > 1 {
				ins.Autowire = kv[1] == "true"
			}
		case "namespace":
			ins.Namespace = defaultNamespace
			if len(kv) > 1 {
				v := strings.Join(kv[1:], "")
				if v != "" {
					ins.Namespace = v
				}
			}
		case "name":
			if len(kv) > 1 {
				ins.Name = strings.Join(kv[1:], "")
			}
		}
	}
	return ins
}

func NewInjectTag() *InjectTag {
	return &InjectTag{}
}

type InjectTag struct {
	// 是否自动注入
	Autowire bool
	// 空间
	Namespace string
	// 注入对象的名称
	Name string
}
