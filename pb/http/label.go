package http

// NewLable label实例
func NewLable(k, v string) *Label {
	return &Label{
		key:   k,
		value: v,
	}
}

// Label 路由标签
type Label struct {
	key   string
	value string
}

// Key 健
func (l *Label) Key() string {
	return l.key
}

// Value 值
func (l *Label) Value() string {
	return l.value
}
