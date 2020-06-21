package router

const (
	// ActionLableKey key name
	ActionLableKey = "action"
)

var (
	// GetAction Label
	GetAction = action("get")
	// ListAction label
	ListAction = action("list")
	// CreateAction label
	CreateAction = action("create")
	// UpdateAction label
	UpdateAction = action("update")
	// DeleteAction label
	DeleteAction = action("delete")
)

func action(value string) *Label {
	return NewLable(ActionLableKey, value)
}

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
