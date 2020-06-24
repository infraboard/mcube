package nats

import "time"

const (
	// DefaultDrainTimeout 默认超时时间
	DefaultDrainTimeout = 300
	// DefaultConnectTimeout 建立连接超时时间
	DefaultConnectTimeout = 60
)

// NewDefaultConfig 默认配置
func NewDefaultConfig() *Config {
	return &Config{
		Servers:        []string{"nats://127.0.0.1:4222"},
		DrainTimeout:   DefaultDrainTimeout,
		ConnectTimeout: DefaultConnectTimeout,
	}
}

// Config 配置
type Config struct {
	Servers        []string `json:"servers,omitempty" yaml:"servers" toml:"servers" env:"NATS_SERVERS" envSeparator:","`
	DrainTimeout   int      `json:"drain_timeout" yaml:"drain_timeout" toml:"drain_timeout" env:"NATS_DRAIN_TIMEOUT"`
	ConnectTimeout int      `json:"connect_timeout" yaml:"connect_timeout" toml:"connect_timeout" env:"NATS_CONNECT_TIMEOUT"`
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
