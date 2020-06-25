package nats

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/infraboard/mcube/bus"
	"github.com/infraboard/mcube/bus/event"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"github.com/nats-io/nats.go"
)

// NewBroker todo
func NewBroker(conf *Config) *Broker {
	b := &Broker{
		conf:    conf,
		closeCh: make(chan error),
		l:       zap.L().Named("Nats Bus"),
	}

	b.opts = nats.GetDefaultOptions()

	b.opts.Servers = conf.Servers
	b.opts.DrainTimeout = conf.GetDrainTimeout()
	b.opts.Timeout = conf.GetConnectTimeout()
	b.opts.ReconnectWait = conf.GetReconnectWait()
	b.opts.MaxReconnect = conf.GetMaxReconnect()

	b.opts.ClosedCB = b.closeHandler
	b.opts.AsyncErrorCB = b.asyncErrorHandler
	b.opts.DisconnectedErrCB = b.disconnectedErrorHandler
	b.opts.ReconnectedCB = b.reconnectHandler
	return b
}

// Broker todo
type Broker struct {
	connected bool
	conf      *Config
	conn      *nats.EncodedConn
	opts      nats.Options
	closeCh   chan error
	l         logger.Logger
	subs      []*nats.Subscription

	lock sync.RWMutex
}

// Debug 设置Logger
func (b *Broker) Debug(l logger.Logger) {
	b.l = l
}

// Connect 实现连接
func (b *Broker) Connect() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.connected {
		return nil
	}

	status := nats.CLOSED
	if b.conn != nil {
		status = b.conn.Conn.Status()
	}

	switch status {
	case nats.CONNECTED, nats.RECONNECTING, nats.CONNECTING:
		b.connected = true
		return nil
	default:
		c, err := b.opts.Connect()
		if err != nil {
			return err
		}
		ec, err := nats.NewEncodedConn(c, nats.JSON_ENCODER)
		if err != nil {
			return err
		}
		b.conn = ec
		b.connected = true
		return nil
	}
}

// Disconnect 断开连接
func (b *Broker) Disconnect() error {
	b.lock.Lock()
	defer b.lock.Unlock()

	if b.conn == nil {
		return errors.New("not connected")
	}

	// 优雅关闭客户端, 通过回调确认是否关闭完成
	if err := b.conn.Drain(); err != nil {
		return err
	}

	// set not connected
	b.connected = false

	return <-b.closeCh
}

// Pub 发布事件
func (b *Broker) Pub(topic string, e *event.Event) error {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.conn == nil {
		return errors.New("not connected")
	}

	if err := b.conn.Publish(topic, e); err != nil {
		return err
	}
	return nil
}

// Sub 订阅事件
func (b *Broker) Sub(topic string, h bus.EventHandler) error {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.conn == nil {
		return errors.New("not connected")
	}

	fn := func(msg *nats.Msg) {
		b.l.Debugf("receive an message from %s, data: %s", msg.Subject, string(msg.Data))
		e := event.NewEvent()
		if err := json.Unmarshal(msg.Data, e); err != nil {
			b.l.Errorf("unmarshal data to event error, %s", err)
			return
		}

		if err := h(e); err != nil {
			b.l.Errorf("call back to return event error, %s", err)
		}
	}

	sub, err := b.conn.Subscribe(topic, fn)
	if err != nil {
		return nil
	}

	b.subs = append(b.subs, sub)
	return nil
}

func (b *Broker) reconnectHandler(nc *nats.Conn) {
	b.l.Debugf("reconnected [%s]", nc.ConnectedUrl())
}

func (b *Broker) closeHandler(nc *nats.Conn) {
	b.l.Debugf("exiting: %v", nc.LastError())
	b.closeCh <- nc.LastError()
	return
}

func (b *Broker) asyncErrorHandler(conn *nats.Conn, sub *nats.Subscription, err error) {
	if err != nil {
		b.l.Error("async error, %s", err)
	}

	return
}

func (b *Broker) disconnectedErrorHandler(conn *nats.Conn, err error) {
	b.l.Errorf("disconnected due to:%v, will attempt reconnects for %ds", err, b.conf.ReconnectWait)
	return
}
