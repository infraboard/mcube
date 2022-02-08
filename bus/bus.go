package bus

import (
	"github.com/infraboard/mcube/pb/event"
)

// PubManager 带管理
type PubManager interface {
	Manager
	Publisher
}

// Publisher 发送事件
type Publisher interface {
	// 发送事件
	Pub(topic string, e *event.Event) error
}

// SubManager 带管理
type SubManager interface {
	Manager
	Subscriber
}

// Subscriber 订阅事件
type Subscriber interface {
	Sub(topic string, h EventHandler) error
}

// Manager 管理端
type Manager interface {
	Connect() error
	Disconnect() error
}

// EventHandler is used to process messages via a subscription of a topic.
type EventHandler func(topic string, e *event.Event) error
