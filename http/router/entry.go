package router

import "fmt"

// Entry 路由条目
type Entry struct {
	Path      string            `json:"path,omitempty"`
	Method    string            `json:"method,omitempty"`
	Resource  string            `json:"resource,omitempty"`
	Protected bool              `json:"protected"`
	Labels    map[string]string `json:"labels,omitempty"`
}

func gernateEntryID(path, mothod string) string {
	return path + "." + mothod
}

// ID entry 唯一标识符
func (e *Entry) ID() string {
	return e.Path + "." + e.Method
}

// NewEntrySet 实例化一个集合
func NewEntrySet() *EntrySet {
	return &EntrySet{
		items: make(map[string]*Entry),
	}
}

// EntrySet 路由集合
type EntrySet struct {
	items map[string]*Entry
}

// AddEntry 添加理由条目
func (s *EntrySet) AddEntry(es ...*Entry) error {
	for _, e := range es {
		if v, ok := s.items[e.ID()]; ok {
			return fmt.Errorf("router entry %s has exist", v.ID())
		}
		s.items[e.ID()] = e
	}

	return nil
}

// ShowEntries 显示理由条目
func (s *EntrySet) ShowEntries() []Entry {
	entries := make([]Entry, 0, len(s.items))
	for _, v := range s.items {
		entries = append(entries, *v)
	}

	return entries
}

// FindEntry 通过函数地址找到对于的路由条目
// 由于该接口是高频接口, 为了不影响性能 采用指针输出, 请不要修改Entry
func (s *EntrySet) FindEntry(path, method string) *Entry {
	id := gernateEntryID(path, method)
	entry, ok := s.items[id]
	if !ok {
		return nil
	}

	return entry
}
