package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer RabbitMQ消费者
type Consumer struct {
	conn      *RabbitConn
	channel   *amqp.Channel
	mu        sync.Mutex
	consumers map[string]*consumerContext
	closeChan chan struct{}
}

type consumerContext struct {
	cancel   context.CancelFunc
	config   consumerConfig
	callback CallBackHandler
	msgChan  <-chan amqp.Delivery
}

type consumerConfig struct {
	exchange     string
	queue        string
	exchangeType string
	routingKey   string
	autoAck      bool
	exclusive    bool
	args         amqp.Table
}

// ConsumeOption 消费者配置选项
type ConsumeOption func(*consumerConfig)

// NewConsumer 创建新的消费者实例
func NewConsumer(conn *RabbitConn) (*Consumer, error) {
	c := &Consumer{
		conn:      conn,
		consumers: make(map[string]*consumerContext),
		closeChan: make(chan struct{}),
	}

	if err := c.resetChannel(); err != nil {
		return nil, fmt.Errorf("failed to create initial channel: %w", err)
	}

	// 注册连接重连回调
	conn.RegisterReconnectCallback(func(newConn *amqp.Connection) {
		if err := c.handleReconnect(); err != nil {
			log.Printf("Failed to handle reconnect: %v", err)
		}
	})

	return c, nil
}

// Subscribe 订阅模式 (Pub/Sub)
func (c *Consumer) Subscribe(ctx context.Context, subject string, cb CallBackHandler, opts ...ConsumeOption) error {
	config := consumerConfig{
		exchange:     subject,
		exchangeType: "fanout",
		autoAck:      false,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return c.setupConsumer(ctx, config, cb)
}

// Queue 队列模式 (Competing Consumers)
func (c *Consumer) Queue(ctx context.Context, subject, queue string, cb CallBackHandler, opts ...ConsumeOption) error {
	config := consumerConfig{
		exchange:     subject,
		queue:        queue,
		exchangeType: "direct",
		routingKey:   queue,
		autoAck:      false,
	}

	for _, opt := range opts {
		opt(&config)
	}

	return c.setupConsumer(ctx, config, cb)
}

// Close 关闭消费者
func (c *Consumer) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.closeChan:
		return nil // 已经关闭
	default:
		close(c.closeChan)
	}

	// 取消所有消费者
	for _, ctx := range c.consumers {
		ctx.cancel()
	}

	// 关闭通道
	if c.channel != nil {
		return c.channel.Close()
	}

	return nil
}

// ============== 内部方法 ==============

// setupConsumer 设置消费者
func (c *Consumer) setupConsumer(ctx context.Context, config consumerConfig, cb CallBackHandler) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case <-c.closeChan:
		return errors.New("consumer is closed")
	default:
	}

	// 确保通道可用
	if err := c.ensureChannel(); err != nil {
		return err
	}

	// 声明交换机和队列
	var msgs <-chan amqp.Delivery
	var err error

	if config.exchangeType == "fanout" {
		msgs, err = c.setupFanoutConsumer(config)
	} else {
		msgs, err = c.setupQueueConsumer(config)
	}

	if err != nil {
		return err
	}

	// 创建子上下文
	childCtx, cancel := context.WithCancel(ctx)
	key := consumerKey(config.exchange, config.queue)

	// 保存消费者上下文
	c.consumers[key] = &consumerContext{
		cancel:   cancel,
		config:   config,
		callback: cb,
		msgChan:  msgs,
	}

	// 启动消息消费协程
	go c.consumeMessages(childCtx, msgs, cb, key)

	return nil
}

// setupFanoutConsumer 设置fanout类型消费者
func (c *Consumer) setupFanoutConsumer(config consumerConfig) (<-chan amqp.Delivery, error) {
	// 声明fanout交换机
	err := c.channel.ExchangeDeclare(
		config.exchange,
		"fanout",
		true,  // durable
		false, // autoDelete
		false, // internal
		false, // noWait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// 创建匿名队列
	q, err := c.channel.QueueDeclare(
		"",    // 随机名称
		false, // durable
		true,  // autoDelete
		config.exclusive,
		false, // noWait
		config.args,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// 绑定队列到交换机
	err = c.channel.QueueBind(
		q.Name,
		"", // fanout忽略routing key
		config.exchange,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	// 开始消费
	return c.channel.Consume(
		q.Name,
		"", // consumer tag
		config.autoAck,
		config.exclusive,
		false, // noLocal
		false, // noWait
		config.args,
	)
}

// setupQueueConsumer 设置队列消费者
func (c *Consumer) setupQueueConsumer(config consumerConfig) (<-chan amqp.Delivery, error) {
	// 声明交换机(如果不是默认的direct)
	if config.exchange != "" && config.exchangeType != "direct" {
		err := c.channel.ExchangeDeclare(
			config.exchange,
			config.exchangeType,
			true,  // durable
			false, // autoDelete
			false, // internal
			false, // noWait
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to declare exchange: %w", err)
		}
	}

	// 声明队列
	q, err := c.channel.QueueDeclare(
		config.queue,
		true,  // durable
		false, // autoDelete
		config.exclusive,
		false, // noWait
		config.args,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// 绑定队列到交换机(如果有指定exchange)
	if config.exchange != "" {
		err = c.channel.QueueBind(
			q.Name,
			config.routingKey,
			config.exchange,
			false,
			nil,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to bind queue: %w", err)
		}
	}

	// 开始消费
	return c.channel.Consume(
		q.Name,
		"", // consumer tag
		config.autoAck,
		config.exclusive,
		false, // noLocal
		false, // noWait
		config.args,
	)
}

// consumeMessages 消费消息的核心逻辑
func (c *Consumer) consumeMessages(
	ctx context.Context,
	msgs <-chan amqp.Delivery,
	cb CallBackHandler,
	consumerKey string,
) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Consumer %s panic recovered: %v", consumerKey, r)
		}

		c.mu.Lock()
		delete(c.consumers, consumerKey)
		c.mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Consumer %s context canceled", consumerKey)
			return

		case <-c.closeChan:
			log.Printf("Consumer %s stopped due to client close", consumerKey)
			return

		case d, ok := <-msgs:
			if !ok {
				log.Printf("Message channel closed for consumer %s", consumerKey)
				return
			}

			msg := &Message{
				Exchange:   d.Exchange,
				RoutingKey: d.RoutingKey,
				Body:       d.Body,
				Headers:    d.Headers,
			}

			start := time.Now()
			err := cb(ctx, msg)
			duration := time.Since(start)

			if err != nil {
				log.Printf("Message handling failed (consumer %s, duration %v): %v",
					consumerKey, duration, err)
				// 消息处理失败，拒绝消息并重新入队
				if err := d.Nack(false, true); err != nil {
					log.Printf("Failed to Nack message: %v", err)
				}
				continue
			}

			log.Printf("Message processed successfully (consumer %s, duration %v)",
				consumerKey, duration)
			// 消息处理成功，确认消息
			if err := d.Ack(false); err != nil {
				log.Printf("Failed to Ack message: %v", err)
			}
		}
	}
}

// resetChannel 重置通道
func (c *Consumer) resetChannel() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 关闭旧通道
	if c.channel != nil {
		_ = c.channel.Close()
	}

	// 获取当前连接
	conn := c.conn.GetConnection()
	if conn == nil {
		return errors.New("no active connection")
	}

	// 创建新通道
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}

	// 配置通道
	if err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		_ = ch.Close()
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	c.channel = ch
	return nil
}

// ensureChannel 确保通道可用
func (c *Consumer) ensureChannel() error {
	if c.channel != nil && !c.channel.IsClosed() {
		return nil
	}
	return c.resetChannel()
}

// handleReconnect 处理重新连接
func (c *Consumer) handleReconnect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 1. 重置通道
	if err := c.resetChannel(); err != nil {
		return fmt.Errorf("failed to reset channel: %w", err)
	}

	// 2. 恢复所有消费者
	for key, ctx := range c.consumers {
		// 取消旧的消费者
		ctx.cancel()

		// 重新创建消费者
		var msgs <-chan amqp.Delivery
		var err error

		if ctx.config.exchangeType == "fanout" {
			msgs, err = c.setupFanoutConsumer(ctx.config)
		} else {
			msgs, err = c.setupQueueConsumer(ctx.config)
		}

		if err != nil {
			log.Printf("Failed to recreate consumer %s: %v", key, err)
			delete(c.consumers, key)
			continue
		}

		// 创建新的上下文
		newCtx, cancel := context.WithCancel(context.Background())

		// 更新消费者上下文
		ctx.cancel = cancel
		ctx.msgChan = msgs

		// 启动新的消费协程
		go c.consumeMessages(newCtx, msgs, ctx.callback, key)
	}

	return nil
}

// consumerKey 生成消费者唯一键
func consumerKey(exchange, queue string) string {
	return fmt.Sprintf("%s|%s", exchange, queue)
}

// ============== 消费者选项 ==============

func WithAutoAck(autoAck bool) ConsumeOption {
	return func(c *consumerConfig) {
		c.autoAck = autoAck
	}
}

func WithExclusive(exclusive bool) ConsumeOption {
	return func(c *consumerConfig) {
		c.exclusive = exclusive
	}
}

func WithArgs(args amqp.Table) ConsumeOption {
	return func(c *consumerConfig) {
		c.args = args
	}
}

func WithExchangeType(exchangeType string) ConsumeOption {
	return func(c *consumerConfig) {
		c.exchangeType = exchangeType
	}
}

func WithRoutingKey(routingKey string) ConsumeOption {
	return func(c *consumerConfig) {
		c.routingKey = routingKey
	}
}
