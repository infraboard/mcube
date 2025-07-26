package rabbitmq

import (
	"sync"
	"time"

	"github.com/infraboard/mcube/v2/ioc/config/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
)

func NewRabbitConn(conf RabbitConnConfig) (*RabbitConn, error) {
	rc := &RabbitConn{
		conf: conf,
		log:  log.Sub(APP_NAME),
	}

	if err := rc.connect(); err != nil {
		return nil, err
	}

	rc.log.Debug().Msgf("connect rabbitmq success")
	return rc, nil
}

type RabbitConnConfig struct {
	// RabbitMQ连接URL
	URL string `toml:"url" json:"url" yaml:"url"  env:"URL"`
	// 连接超时时间
	Timeout time.Duration `toml:"timeout" json:"timeout" yaml:"timeout"  env:"TIMEOUT"`
	// 心跳间隔
	Heartbeat time.Duration `toml:"heart_beat" json:"heart_beat" yaml:"heart_beat"  env:"HEART_BEAT"`
	// 重连配置
	ReconnectInterval time.Duration `toml:"reconect_interval" json:"reconect_interval" yaml:"reconect_interval"  env:"RECONNECT_INTERVAL"`
	// 最大重连次数, 为0 时表示无限重连
	MaxReconnectAttempts int `toml:"max_reconnect_attempts" json:"max_reconnect_attempts" yaml:"max_reconnect_attempts"  env:"MAX_RECONNECT_ATTEMPTS"`
}

type RabbitConn struct {
	conf RabbitConnConfig

	conn      *amqp.Connection
	mu        sync.RWMutex
	closing   bool
	closeChan chan struct{}

	log *zerolog.Logger
	// 重连回调
	reconnectCallbacks []func(*amqp.Connection)
}

// RegisterReconnectCallback 注册重连回调
func (rc *RabbitConn) RegisterReconnectCallback(cb func(*amqp.Connection)) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.reconnectCallbacks = append(rc.reconnectCallbacks, cb)
}

// connect 内部连接方法
func (rc *RabbitConn) connect() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	conn, err := amqp.DialConfig(rc.conf.URL, amqp.Config{
		Dial:      amqp.DefaultDial(rc.conf.Timeout),
		Heartbeat: rc.conf.Heartbeat,
	})
	if err != nil {
		return err
	}

	if rc.conn != nil && !rc.conn.IsClosed() {
		_ = rc.conn.Close()
	}

	rc.conn = conn
	go rc.monitorConnection()

	// 执行重连回调
	for _, cb := range rc.reconnectCallbacks {
		go cb(conn)
	}
	return nil
}

// monitorConnection 监控连接状态
func (rc *RabbitConn) monitorConnection() {
	notifyClose := make(chan *amqp.Error)
	rc.conn.NotifyClose(notifyClose)

	select {
	case err := <-notifyClose:
		if err != nil {
			rc.log.Info().Msgf("RabbitMQ connection closed: %v", err)
		}
		rc.handleDisconnect()
	case <-rc.closeChan:
		return
	}
}

// handleDisconnect 处理连接断开
func (rc *RabbitConn) handleDisconnect() {
	rc.mu.RLock()
	closing := rc.closing
	rc.mu.RUnlock()

	if closing {
		return
	}

	rc.log.Debug().Msg("Attempting to reconnect to RabbitMQ...")
	var attempts int
	for {
		attempts++

		if rc.conf.MaxReconnectAttempts > 0 && attempts > rc.conf.MaxReconnectAttempts {
			rc.log.Debug().Msgf("Max reconnect attempts (%d) reached", rc.conf.MaxReconnectAttempts)
			return
		}

		if err := rc.connect(); err == nil {
			rc.log.Debug().Msg("Successfully reconnected to RabbitMQ")
			return
		}

		rc.log.Debug().Msgf("Reconnect attempt %d failed, retrying in %v", attempts, rc.conf.ReconnectInterval)
		time.Sleep(rc.conf.ReconnectInterval)
	}
}

// Close 关闭连接
func (rc *RabbitConn) Close() error {
	rc.mu.Lock()
	defer rc.mu.Unlock()

	if rc.closing {
		return nil
	}

	rc.closing = true
	if rc.closeChan != nil {
		close(rc.closeChan)
		rc.closeChan = nil
	}

	if rc.conn != nil && !rc.conn.IsClosed() {
		return rc.conn.Close()
	}
	return nil
}

// IsClosed 检查连接是否关闭
func (rc *RabbitConn) IsClosed() bool {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.closing || rc.conn == nil || rc.conn.IsClosed()
}

// GetConnection 获取当前连接（线程安全）
func (rc *RabbitConn) GetConnection() *amqp.Connection {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.conn
}
