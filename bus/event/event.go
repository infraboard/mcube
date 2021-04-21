package event

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/rs/xid"

	"github.com/infraboard/mcube/types/ftime"
)

// 事件主题定义(由事件类型确定)
// 1. 操作事件

// NewOperateEvent 实例
func NewOperateEvent(e *OperateEventData) (*Event, error) {
	any, err := ptypes.MarshalAny(e)
	if err != nil {
		return nil, err
	}
	return &Event{
		Id:   xid.New().String(),
		Time: ftime.Now().Timestamp(),
		Type: Type_OperateEvent,
		Data: any,
	}, nil
}

// NewDefaultEvent todo
func NewDefaultEvent() *Event {
	return &Event{
		Time: ftime.Now().Timestamp(),
		Meta: make(map[string]string),
	}
}

// Validate 校验事件是否合法
func (e *Event) Validate() error {
	return nil
}

// GetMeta 获取meta信息
func (e *Event) GetMetaKey(key string) string {
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
func (e *Event) ParseBody(data proto.Message) error {
	return ptypes.UnmarshalAny(e.Data, data)
}
