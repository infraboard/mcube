package router

// Entry 路由条目
type Entry struct {
	Path             string            `json:"path,omitempty"`
	Method           string            `json:"method,omitempty"`
	FunctionName     string            `json:"function_name,omitempty"`
	Resource         string            `json:"resource,omitempty"`
	AuthEnable       bool              `json:"auth_enable"`
	PermissionEnable bool              `json:"permission_enable"`
	Labels           map[string]string `json:"labels,omitempty"`
}

// AddLabel 添加Label
func (e *Entry) AddLabel(labels ...*Label) EntryDecorator {
	for i := range labels {
		e.Labels[labels[i].Key()] = labels[i].Value()
	}

	return e
}

// EnableAuth 启动身份验证
func (e *Entry) EnableAuth() {
	e.AuthEnable = true
}

// EnablePermission 启用授权验证
func (e *Entry) EnablePermission() {
	e.PermissionEnable = true
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
