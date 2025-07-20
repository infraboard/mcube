package nats

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(&Client{
		URL: nats.DefaultURL,
	})
}

type Client struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	//
	URL   string `toml:"url" json:"url" yaml:"url"  env:"URL"`
	Token string `toml:"token" json:"token" yaml:"token"  env:"TOKEN"`

	conn *nats.Conn
}

func (c *Client) Name() string {
	return APP_NAME
}

func (c *Client) Options() (opts []nats.Option) {
	if c.Token != "" {
		opts = append(opts, nats.Token(c.Token))
	}
	return
}

func (c *Client) Init() error {
	c.log = log.Sub(c.Name())
	conn, err := nats.Connect(c.URL, c.Options()...)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Close(ctx context.Context) {
	if c.conn != nil {
		c.conn.Drain()
	}
}
