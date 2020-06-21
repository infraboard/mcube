package bus

import (
	"github.com/infraboard/mcube/bus/event"
)

// Publisher 发送事件
type Publisher interface {
	// 发送事件
	Pub(*event.Event) error
	Connect() error
	Close() error
}

// Subscriber 订阅事件
type Subscriber interface {
	Sub() (<-chan *event.Event, error)
	Connect() error
	Close() error
}
