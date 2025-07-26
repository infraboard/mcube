package rabbitmq

import (
	"context"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
)

func init() {
	ioc.Config().Registry(&Client{
		RabbitConnConfig: RabbitConnConfig{
			URL:               "amqp://guest:guest@localhost:5672/",
			Timeout:           5 * time.Second,
			Heartbeat:         10 * time.Second,
			ReconnectInterval: 10 * time.Second,
		},
	})
}

type Client struct {
	ioc.ObjectImpl

	// 连接配置
	RabbitConnConfig

	conn *RabbitConn
}

func (c *Client) Name() string {
	return APP_NAME
}

func (c *Client) Init() error {
	conn, err := NewRabbitConn(c.RabbitConnConfig)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Close(ctx context.Context) {
	if c.conn != nil {
		c.conn.Close()
	}
}
