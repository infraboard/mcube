package rabbitmq_test

import (
	"testing"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/rabbitmq"
	"github.com/rabbitmq/amqp091-go"
)

func TestReconnect(t *testing.T) {
	rabbitmq.GetConn().RegisterReconnectCallback(func(c *amqp091.Connection) {
		t.Log(c.ConnectionState())
	})

	time.Sleep(60 * time.Second)
}
