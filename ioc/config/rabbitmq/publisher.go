package rabbitmq

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// Publisher 增强版，支持自动重连
type Publisher struct {
	conn      *RabbitConn
	channel   *amqp.Channel
	mu        sync.Mutex
	closeChan chan struct{}

	log *zerolog.Logger
}

// NewPublisher 创建支持自动重连的Publisher
func NewPublisher() (*Publisher, error) {
	p := &Publisher{
		conn:      GetConn(),
		closeChan: make(chan struct{}),
		log:       log.Sub(APP_NAME),
	}

	// 初始通道
	if err := p.resetChannel(); err != nil {
		return nil, err
	}

	// 注册重连回调
	p.conn.RegisterReconnectCallback(func(_ *amqp.Connection) {
		p.mu.Lock()
		defer p.mu.Unlock()

		if err := p.resetChannel(); err != nil {
			p.log.Error().Msgf("Failed to reset publisher channel after reconnect: %v", err)
		}
	})

	return p, nil
}

// resetChannel 重置通道
func (p *Publisher) resetChannel() error {
	if p.channel != nil && !p.channel.IsClosed() {
		_ = p.channel.Close()
	}

	conn := p.conn.GetConnection()
	if conn == nil {
		return errors.New("no active connection")
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	p.channel = ch
	return nil
}

// Publish 发布消息（支持自动重连）
func (p *Publisher) Publish(ctx context.Context, msg *Message) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case <-p.closeChan:
		return errors.New("publisher is closed")
	default:
	}

	if p.channel == nil || p.channel.IsClosed() {
		if err := p.resetChannel(); err != nil {
			return err
		}
	}

	if trace.Get().Enable {
		log.FromCtx(ctx).Info().Msg("enable rabbitmq trace")
		tracer := otel.GetTracerProvider().Tracer("rabbitmq-producer")

		// 启动一个Span代表发送消息的操作
		_, span := tracer.Start(ctx, "publish to "+msg.Exchange+"."+msg.RoutingKey)
		defer span.End()

		// 准备消息属性
		headers := make(propagation.MapCarrier)
		// 获取文本映射传播器，并将追踪信息注入到 headers 中
		propagator := otel.GetTextMapPropagator()
		propagator.Inject(ctx, propagation.MapCarrier(headers))

		// 将追踪信息合并到现有headers中，避免覆盖用户自定义参数
		if msg.Headers == nil {
			msg.Headers = make(amqp.Table)
		}
		for k, v := range headers {
			msg.Headers[k] = v
		}
	}

	return p.channel.Publish(
		msg.Exchange,
		msg.RoutingKey,
		msg.Mandatory,
		msg.Immediate,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         msg.Body,
			Headers:      msg.Headers,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

// Close 关闭Publisher
func (p *Publisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case <-p.closeChan:
		return nil
	default:
		close(p.closeChan)
	}

	if p.channel != nil && !p.channel.IsClosed() {
		return p.channel.Close()
	}
	return nil
}
