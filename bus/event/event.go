package event

import "github.com/infraboard/mcube/types/ftime"

// Event 事件数据结构
type Event struct {
	Header
	Body
}

// Header 事件元数据
type Header struct {
	ID     string     `bson:"id" json:"id,omitempty"`         // 事件ID
	Time   ftime.Time `bson:"time" json:"time,omitempty"`     // 事件发生时间(毫秒)
	Source string     `bson:"source" json:"source,omitempty"` // 事件来源, 比如cmdb
	Level  Level      `bson:"level" json:"level,omitempty"`   // 事件等级
	Meta   string     `bson:"meta" json:"meta,omitempty"`     // 事件的元数据
}

// Body 事件具体数据
type Body struct {
	Reason       string      `bson:"reason" json:"reason,omitempty"`     // 触发原因, 比如 创建/删除/绑定/告警/恢复
	Message      string      `bson:"message" json:"message,omitempty"`   // 事件消息,
	ResourceType string      `bson:"resource_type" json:"resource_type"` // 资源类型,
	ResourceUUID string      `bson:"resource_uuid" json:"resource_uuid"` // 资源UUID,
	Data         interface{} `bson:"data" json:"data,omitempty"`         // 事件具体数据
}
