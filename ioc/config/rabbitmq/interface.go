package rabbitmq

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	APP_NAME = "rabbitmq"
)

func GetConn() *RabbitConn {
	return ioc.Config().Get(APP_NAME).(*Client).conn
}

func NewQueueMessage(queue string, data []byte) *Message {
	return &Message{
		RoutingKey: queue,
		Body:       data,
	}
}

func NewFanoutMessage(topic string, data []byte) *Message {
	return &Message{
		Exchange: topic,
		Body:     data,
	}
}

func NewTopicMessage(topic, routingKey string, data []byte) *Message {
	return &Message{
		Exchange:   topic,
		RoutingKey: routingKey,
		Body:       data,
	}
}

// Message 消息结构
type Message struct {
	Exchange   string
	RoutingKey string
	Mandatory  bool
	Immediate  bool
	Body       []byte
	Headers    amqp.Table
}

// CallBackHandler 消息处理回调
type CallBackHandler func(ctx context.Context, msg *Message) error
