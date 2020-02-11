package router

// Entry 路由条目
type Entry struct {
	Path         string            `json:"path,omitempty"`
	Method       string            `json:"method,omitempty"`
	FunctionName string            `json:"function_name,omitempty"`
	Resource     string            `json:"resource,omitempty"`
	Protected    bool              `json:"protected"`
	Labels       map[string]string `json:"labels,omitempty"`
}

// NewEntrySet 实例
func NewEntrySet() *EntrySet {
	return &EntrySet{}
}

// EntrySet 路由条目集
type EntrySet struct {
	Items []*Entry `json:"items"`
}

// AddEntry 添加Entry
func (s *EntrySet) AddEntry(e Entry) {
	s.Items = append(s.Items, &e)
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
