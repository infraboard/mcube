package http

import (
	"strings"
)

// EntryDecorator 装饰
type EntryDecorator interface {
	// SetLabel 设置子路由标签, 作用于Entry上
	AddLabel(...*Label) EntryDecorator
	EnableAuth() EntryDecorator
	DisableAuth() EntryDecorator
	EnablePermission() EntryDecorator
	DisablePermission() EntryDecorator
}

// NewEntry 行健条目
func NewEntry(path, method, resource string) *Entry {
	return &Entry{
		Path:     path,
		Method:   method,
		Resource: resource,
		Labels:   map[string]string{},
	}
}

// AddLabel 添加Label
func (e *Entry) AddLabel(labels ...*Label) EntryDecorator {
	for i := range labels {
		e.Labels[labels[i].Key()] = labels[i].Value()
	}

	return e
}

// GetLableValue 获取Lable的值
func (e *Entry) GetLableValue(key string) string {
	v, ok := e.Labels[key]
	if ok {
		return v
	}
	return ""
}

// EnableAuth 启动身份验证
func (e *Entry) EnableAuth() EntryDecorator {
	e.AuthEnable = true
	return e
}

// DisableAuth 不启用身份验证
func (e *Entry) DisableAuth() EntryDecorator {
	e.AuthEnable = false
	return e
}

// EnablePermission 启用授权验证
func (e *Entry) EnablePermission() EntryDecorator {
	e.PermissionEnable = true
	return e
}

// DisablePermission 禁用授权验证
func (e *Entry) DisablePermission() EntryDecorator {
	e.PermissionEnable = false
	return e
}

// NewEntrySet 实例
func NewEntrySet() *EntrySet {
	return &EntrySet{}
}

// EntrySet 路由条目集
type EntrySet struct {
	Items []*Entry `json:"items"`
}

// PermissionEnableEntry todo
func (s *EntrySet) PermissionEnableEntry() []*Entry {
	items := []*Entry{}
	for i := range s.Items {
		if s.Items[i].PermissionEnable {
			items = append(items, s.Items[i])
		}
	}

	return items
}

// AuthEnableEntry todo
func (s *EntrySet) AuthEnableEntry() []*Entry {
	items := []*Entry{}
	for i := range s.Items {
		if s.Items[i].AuthEnable {
			items = append(items, s.Items[i])
		}
	}

	return items
}

func (s *EntrySet) String() string {
	strs := []string{}
	for i := range s.Items {
		strs = append(strs, s.Items[i].String())
	}

	return strings.Join(strs, "\n")
}

// Merge todo
func (s *EntrySet) Merge(target *EntrySet) {
	for i := range target.Items {
		s.Items = append(s.Items, target.Items[i])
	}
}

// AddEntry 添加Entry
func (s *EntrySet) AddEntry(es ...Entry) {
	for i := range es {
		s.Items = append(s.Items, &es[i])
	}
}

// GetEntry 获取条目
func (s *EntrySet) GetEntry(path, mothod string) *Entry {
	for i := range s.Items {
		item := s.Items[i]
		if item.Path == path && item.Method == mothod {
			return item
		}
	}

	return nil
}
