//go:generate mcube enum -m

package event

const (
	// OperateEventType (operate) 资源操作事件
	OperateEventType Type = iota
	// AlertEventType (alert) 告警事件
	AlertEventType
)

// Type 事件类型
type Type uint
