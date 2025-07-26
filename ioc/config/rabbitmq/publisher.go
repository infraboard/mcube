package rabbitmq

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher 增强版，支持自动重连
type Publisher struct {
	conn      *RabbitConn
	channel   *amqp.Channel
	mu        sync.Mutex
	closeChan chan struct{}
}

// NewPublisher 创建支持自动重连的Publisher
func NewPublisher() (*Publisher, error) {
	p := &Publisher{
		conn:      GetConn(),
		closeChan: make(chan struct{}),
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
			log.Printf("Failed to reset publisher channel after reconnect: %v", err)
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
