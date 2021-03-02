package httprouter

import (
	"fmt"
	"net/http"

	httppb "github.com/infraboard/mcube/pb/http"
)

// Entry 路由条目
type entry struct {
	*httppb.Entry
	h http.HandlerFunc
}

func newEntrySet() *entrySet {
	return &entrySet{
		items: map[string]*entry{},
	}
}

// EntrySet 路由集合
type entrySet struct {
	order []string
	items map[string]*entry
}

// AddEntry 添加理由条目
func (s *entrySet) AddEntry(es ...*entry) error {
	for _, e := range es {
		key := s.ID(e.Path, e.Method)
		if _, ok := s.items[key]; ok {
			return fmt.Errorf("router entry %s has exist", key)
		}
		s.items[key] = e
		s.order = append(s.order, key)
	}

	return nil
}

// ShowEntries 显示理由条目
func (s *entrySet) EntrieSet() *httppb.EntrySet {
	es := httppb.NewEntrySet()
	for _, key := range s.order {
		es.AddEntry(*s.items[key].Entry)
	}

	return es
}

// FindEntry 通过函数地址找到对于的路由条目
// 由于该接口是高频接口, 为了不影响性能 采用指针输出, 请不要修改Entry
func (s *entrySet) FindEntry(path, method string) *entry {
	id := s.ID(path, method)
	entry, ok := s.items[id]
	if !ok {
		return nil
	}

	return entry
}

func (s *entrySet) ID(path, mothod string) string {
	return path + "." + mothod
}
