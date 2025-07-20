package nats

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/bus"
	ioc_nats "github.com/infraboard/mcube/v2/ioc/config/nats"
	"github.com/nats-io/nats.go"
)

func init() {
	ioc.Config().Registry(&BusServiceImpl{})
}

var _ bus.Service = (*BusServiceImpl)(nil)

type BusServiceImpl struct {
	ioc.ObjectImpl
}

func (b *BusServiceImpl) Name() string {
	return bus.APP_NAME
}

// 事件发送
func (b *BusServiceImpl) Publish(ctx context.Context, e *bus.Event) error {
	msg := nats.NewMsg(e.Subject)
	msg.Data = e.Data
	msg.Header = e.Header
	return ioc_nats.Get().PublishMsg(msg)
}

// 订阅事件
func (b *BusServiceImpl) Subscribe(ctx context.Context, subject string, cb bus.EventHandler) error {
	_, err := ioc_nats.Get().Subscribe(subject, func(msg *nats.Msg) {
		cb(&bus.Event{
			Subject: msg.Subject,
			Header:  msg.Header,
			Data:    msg.Data,
		})
	})
	if err != nil {
		return err
	}
	return nil
}
