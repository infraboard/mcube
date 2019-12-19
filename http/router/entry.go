package router

// Entry 路由条目
type Entry struct {
	Name      string            `json:"name,omitempty"`
	Resource  string            `json:"resource,omitempty"`
	Path      string            `json:"path,omitempty"`
	Method    string            `json:"method,omitempty"`
	Protected bool              `json:"protected"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// NewEntrySet 实例化一个集合
func NewEntrySet() *EntrySet {
	return &EntrySet{}
}

// EntrySet 路由集合
type EntrySet struct {
	items []Entry
}

// AddEntry 添加理由条目
func (s *EntrySet) AddEntry(e ...Entry) {
	s.items = append(s.items, e...)
}

// ShowEntries 显示理由条目
func (s *EntrySet) ShowEntries() []Entry {
	return s.items
}
