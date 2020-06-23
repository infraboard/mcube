package redis

// NewDefaultConfig 默认配置
func NewDefaultConfig() *Config {
	return &Config{
		Prefix:     "",
		Address:    "127.0.0.1:6379",
		DB:         0,
		Password:   "",
		DefaultTTL: 300,
	}
}

// Config 配置
type Config struct {
	Prefix     string `json:"prefix,omitempty" yaml:"prefix" toml:"prefix" env:"MCUBE_CACHE_PREFIX"`
	Address    string `json:"address,omitempty" yaml:"address" toml:"address" env:"MCUBE_CACHE_ADDRESS"`
	DB         int    `json:"db,omitempty" yaml:"db" toml:"db" env:"MCUBE_CACHE_DB"`
	Password   string `json:"password,omitempty" yaml:"password" toml:"password" env:"MCUBE_CACHE_PASSWORD"`
	DefaultTTL int    `json:"default_ttl,omitempty" yaml:"default_ttl" toml:"default_ttl" env:"MCUBE_CACHE_TTL"`
}
