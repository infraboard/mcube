package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/xid"

	"github.com/infraboard/mcube/types/ftime"
)

// 事件主题定义(由事件类型确定)
// 1. 操作事件

// NewOperateEvent 实例
func NewOperateEvent(e *OperateEvent) *Event {
	return &Event{
		ID:   xid.New().String(),
		Time: ftime.Now(),
		Type: OperateEventType,
		Body: e,
	}
}

// NewDefaultEvent todo
func NewDefaultEvent() *Event {
	return &Event{
		Time: ftime.Now(),
		Meta: make(map[string]string),
		Body: &json.RawMessage{},
	}
}

// Event 事件数据结构
type Event struct {
	ID     string            `bson:"_id" json:"id"`        // 事件ID
	Time   ftime.Time        `bson:"time" json:"time"`     // 事件发生时间(毫秒)
	Type   Type              `bson:"type" json:"type"`     // 事件类型
	Source string            `bson:"source" json:"source"` // 事件来源, 比如cmdb
	Level  Level             `bson:"level" json:"level"`   // 事件等级
	Meta   map[string]string `bson:"label" json:"label"`   // 标签
	Body   interface{}       `bson:"body" json:"body"`     // 事件数据

	ctx context.Context
}

// Validate 校验事件是否合法
func (e *Event) Validate() error {
	return nil
}

// GetMeta 获取meta信息
func (e *Event) GetMeta(key string) string {
	if v, ok := e.Meta[key]; ok {
		return v
	}

	return ""
}

// SetMeta 设置meta信息
func (e *Event) SetMeta(key, value string) {
	e.Meta[key] = value
}

// SetLevel 设置事件级别
func (e *Event) SetLevel(l Level) {
	e.Level = l
}

// SetSource 设置事件来源
func (e *Event) SetSource(src string) {
	e.Source = src
}

// ParseBody todo
func (e *Event) ParseBody() error {
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
	default:
		return fmt.Errorf("unknown event type: %s", e.Type)
	}
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
	Session      string      `bson:"session" json:"session"`             // 回话ID
	Account      string      `bson:"account" json:"account"`             // 操作人
	RequestID    string      `bson:"request_id" json:"request_id"`       // 请求ID
	IPAddress    string      `bson:"ip_address" json:"ip_address"`       // 操作者IP
	UserAgent    string      `bson:"user_agent" json:"user_agent"`       // 用户UA
	UserName     string      `bson:"user_name" json:"user_name"`         // 用户名称
	UserType     string      `bson:"user_type" json:"user_type"`         // 用户类型
	UserDomain   string      `bson:"user_domain" json:"user_domain"`     // 操作的域
	ServiceName  string      `bson:"service_name" json:"service_name"`   // 服务名称
	FeaturePath  string      `bson:"feature_path" json:"feature_path"`   // 功能路径
	ResourceType string      `bson:"resource_type" json:"resource_type"` // 资源类型
	Action       string      `bson:"action" json:"action"`               // 操作动作
	Request      interface{} `bson:"request" json:"request,omitempty"`   // 事件数据
	Response     interface{} `bson:"response" json:"response,omitempty"` // 事件数据
}
