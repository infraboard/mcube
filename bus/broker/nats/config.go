package nats

import "time"

const (
	// DefaultDrainTimeout 默认超时时间
	DefaultDrainTimeout = 300
	// DefaultConnectTimeout 建立连接超时时间
	DefaultConnectTimeout = 5
	// DefaultReconnectWait 重连超时时间
	DefaultReconnectWait = 2
	// DefaultMaxReconnect 重试的次数
	DefaultMaxReconnect = 60
)

// NewDefaultConfig 默认配置
func NewDefaultConfig() *Config {
	return &Config{
		Servers:        []string{"nats://127.0.0.1:4222"},
		DrainTimeout:   DefaultDrainTimeout,
		ConnectTimeout: DefaultConnectTimeout,
		ReconnectWait:  DefaultReconnectWait,
		MaxReconnect:   DefaultMaxReconnect,
	}
}

// Config 配置
type Config struct {
	Servers        []string `json:"servers,omitempty" yaml:"servers" toml:"servers" env:"NATS_SERVERS" envSeparator:","`
	DrainTimeout   int      `json:"drain_timeout" yaml:"drain_timeout" toml:"drain_timeout" env:"NATS_DRAIN_TIMEOUT"`         // 单位秒
	ConnectTimeout int      `json:"connect_timeout" yaml:"connect_timeout" toml:"connect_timeout" env:"NATS_CONNECT_TIMEOUT"` // 单位秒
	ReconnectWait  int      `json:"reconnect_wait" ymal:"reconnect_wait" toml:"reconnect_wait" env:"NATS_RECONNECT_WAIT"`     // 单位秒
	MaxReconnect   int      `json:"max_reconnect" yaml:"max_reconnect" toml:"max_reconnect" env:"NATS_MAX_RECONNECT"`         // 最大重连次数
}

// GetDrainTimeout todo
func (c *Config) GetDrainTimeout() time.Duration {
	var dt = DefaultDrainTimeout
	if c.DrainTimeout != 0 {
		dt = c.DrainTimeout
	}

	return time.Duration(dt) * time.Second
}

// GetConnectTimeout todo
func (c *Config) GetConnectTimeout() time.Duration {
	var ct = DefaultConnectTimeout
	if c.ConnectTimeout != 0 {
		ct = c.ConnectTimeout
	}

	return time.Duration(ct) * time.Second
}

// GetReconnectWait todo
func (c *Config) GetReconnectWait() time.Duration {
	var rt = DefaultReconnectWait
	if c.ConnectTimeout != 0 {
		rt = c.ReconnectWait
	}

	return time.Duration(rt) * time.Second
}

// GetMaxReconnect todo
func (c *Config) GetMaxReconnect() int {
	if c.MaxReconnect != 0 {
		return c.MaxReconnect
	}

	return DefaultMaxReconnect
}
