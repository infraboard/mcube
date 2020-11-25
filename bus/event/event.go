package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/xid"

	"github.com/infraboard/mcube/types/ftime"
)

// 事件主题定义(由事件类型确定)
// 1. 资源变更事件 (变更成功和变更失败)
// 2. 资源告警事件 ( 更加metric计算触发)

// NewOperateEvent 实例
func NewOperateEvent(header *Header, e *OperateEvent) *Event {
	return &Event{
		ID:     xid.New().String(),
		Time:   ftime.Now(),
		Type:   OperateEventType,
		Header: header,
		Body:   e,
	}
}

// NewDefaultEvent todo
func NewDefaultEvent() *Event {
	return &Event{
		Header: NewHeader(),
		Body:   &json.RawMessage{},
	}
}

// Event 事件数据结构
type Event struct {
	ID      string           `bson:"_id" json:"id"`    // 事件ID
	Time    ftime.Time       `bson:"time" json:"time"` // 事件发生时间(毫秒)
	Type    Type             `bson:"type" json:"type"` // 事件类型
	*Header `bson:",inline"` // 事件头
	Body    interface{}      `bson:"body" json:"body"` // 事件数据

	ctx    context.Context
	parsed bool
}

// Validate 校验事件是否合法
func (e *Event) Validate() error {
	return nil
}

// ParseBody todo
func (e *Event) ParseBody() error {
	if e.parsed {
		return nil
	}

	body, err := e.getBytesBody()
	if err != nil {
		return err
	}

	switch e.Type {
	case OperateEventType:
		e.Body, err = ParseOperateEventFromBytes(body)
		if err != nil {
			return err
		}
	case StatusEventType:
		e.Body, err = ParseStatusEventFromBytes(body)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown event type: %s", e.Type)
	}

	e.parsed = true
	return nil
}

// WithContext 添加上下文
func (e *Event) WithContext(ctx context.Context) {
	e.ctx = ctx
}

// Context 返回事件的上下文
func (e *Event) Context() context.Context {
	return e.ctx
}

func (e *Event) getBytesBody() ([]byte, error) {
	switch v := e.Body.(type) {
	case []byte:
		return v, nil
	case json.RawMessage:
		return v, nil
	case *json.RawMessage:
		return *v, nil
	default:
		return nil, fmt.Errorf("body type is not []byte or json.RawMessage")
	}
}

// NewHeader todo
func NewHeader() *Header {
	return &Header{
		Label: make(map[string]string),
	}
}

// Header todo
type Header struct {
	Source string            `bson:"source" json:"source"` // 事件来源, 比如cmdb
	Level  Level             `bson:"level" json:"level"`   // 事件等级
	Label  map[string]string `bson:"label" json:"label"`   // 标签
	Meta   Meta              `bson:"meta" json:"meta"`     // 事件的元数据
}

// ParseOperateEventFromBytes todo
func ParseOperateEventFromBytes(data []byte) (*OperateEvent, error) {
	oe := &OperateEvent{}
	if err := json.Unmarshal(data, oe); err != nil {
		return nil, err
	}
	return oe, nil
}

// OperateEvent 事件具体数据
type OperateEvent struct {
	OperateSession string      `bson:"operate_session" json:"operate_session"` // 回话ID
	OperateUser    string      `bson:"operate_user" json:"operate_user"`       // 操作人
	ResourceType   string      `bson:"resource_type" json:"resource_type"`     // 资源类型,
	ResourceUUID   string      `bson:"resource_uuid" json:"resource_uuid"`     // 资源UUID
	ResourceName   string      `bson:"resource_name" json:"resource_name"`     // 资源名称
	Action         string      `bson:"action" json:"action"`                   // 操作
	Data           interface{} `bson:"data" json:"data,omitempty"`             // 事件数据
}

// ParseStatusEventFromBytes todo
func ParseStatusEventFromBytes(data []byte) (*StatusEvent, error) {
	ae := &StatusEvent{}
	if err := json.Unmarshal(data, ae); err != nil {
		return nil, err
	}
	return ae, nil
}

// StatusEvent 事件具体数据
type StatusEvent struct {
	Reason       string      `bson:"reason" json:"reason,omitempty"`     // 触发原因, 比如 创建/删除/绑定/告警/恢复
	Message      string      `bson:"message" json:"message,omitempty"`   // 事件消息,
	ResourceType string      `bson:"resource_type" json:"resource_type"` // 资源类型,
	ResourceUUID string      `bson:"resource_uuid" json:"resource_uuid"` // 资源UUID
	ResourceName string      `bson:"resource_name" json:"resource_name"` // 资源名称
	Data         interface{} `bson:"data" json:"data,omitempty"`         // 事件数据
}

// Meta todo
type Meta interface {
	// Put associates the specified value with the specified key. If the map
	// previously contained a mapping for the key, the old value is replaced and
	// returned. The key can be expressed in dot-notation (e.g. x.y) to put a value
	// into a nested map.
	//
	// If you need insert keys containing dots then you must use bracket notation
	// to insert values (e.g. m[key] = value).
	Put(key string, value interface{}) (interface{}, error)
	// Delete deletes the given key from the map.
	Delete(key string) error
	// GetValue 获取值
	Get(key string) (interface{}, error)
}
