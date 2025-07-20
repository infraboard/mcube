package nats

import (
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/nats-io/nats.go"
)

const (
	APP_NAME = "nats"
)

func Get() *nats.Conn {
	return ioc.Config().Get(APP_NAME).(*Client).conn
}
