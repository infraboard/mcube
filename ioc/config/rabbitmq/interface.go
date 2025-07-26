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

// Message 消息结构
type Message struct {
	Exchange   string
	RoutingKey string
	Body       []byte
	Headers    amqp.Table
}

// CallBackHandler 消息处理回调
type CallBackHandler func(ctx context.Context, msg *Message) error
