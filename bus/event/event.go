package event

import (
	"github.com/rs/xid"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/infraboard/mcube/types/ftime"
)

// 事件主题定义(由事件类型确定)
// 1. 操作事件

// NewOperateEvent 实例
func NewOperateEvent(e *OperateEventData) (*Event, error) {
	body, err := anypb.New(e)
	if err != nil {
		return nil, err
	}

	obj := &Event{
		Id:     xid.New().String(),
		Type:   Type_Operate,
		Header: NewHeader(),
		Body:   body,
	}

	return obj, nil
}

// NewDefaultEvent todo
func NewDefaultEvent() *Event {
	return &Event{
		Header: NewHeader(),
	}
}

// NewHeader todo
func NewHeader() *Header {
	return &Header{
		Time: ftime.Now().Timestamp(),
		Meta: make(map[string]string),
	}
}

// Validate 校验事件是否合法
func (e *Event) Validate() error {
	return nil
}

// GetMetaKey 获取meta信息
func (e *Event) GetMetaKey(key string) (string, bool) {
	v, ok := e.Header.Meta[key]
	return v, ok
}

// SetMeta 设置meta信息
func (e *Event) SetMeta(key, value string) {
	e.Header.Meta[key] = value
}

// SetLevel 设置事件级别
func (e *Event) SetLevel(l Level) {
	e.Header.Level = l
}

// SetSource 设置事件来源
func (e *Event) SetSource(src string) {
	e.Header.Source = src
}

// ParseBoby todo
func (e *Event) ParseBoby(body proto.Message) error {
	return anypb.UnmarshalTo(e.Body, body, proto.UnmarshalOptions{})
}
