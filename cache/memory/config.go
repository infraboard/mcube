package memory

import "time"

// NewDefaultConfig 使用默认配置
func NewDefaultConfig() *Config {
	return &Config{
		TTL:  300,
		Size: 500,
	}
}

// Config 配置选项
type Config struct {
	TTL  int `json:"ttl,omitempty" yaml:"ttl" toml:"ttl" env:"MCUBE_CACHE_TTL"` // 默认秒
	Size int `json:"size,omitempty" yaml:"size" toml:"size" env:"MCUBE_CACHE_SIZE"`
}

func (c *Config) GetTTL() time.Duration {
	return time.Duration(c.TTL) * time.Second
}
