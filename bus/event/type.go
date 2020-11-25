//go:generate mcube enum -m

package event

const (
	// OperateEventType (operate) 资源操作事件
	OperateEventType Type = iota
	// StatusEventType (status) 告警事件
	StatusEventType
)

// Type 事件类型
type Type uint
