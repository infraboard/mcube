package redis

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Redis{
	DB:        0,
	Endpoints: []string{"127.0.0.1:6379"},
	Trace:     true,
}

type Redis struct {
	ioc.ObjectImpl
	Endpoints []string `toml:"endpoints" json:"endpoints" yaml:"endpoints" env:"ENDPOINTS" envSeparator:","`
	DB        int      `toml:"db" json:"db" yaml:"db"  env:"DB"`
	UserName  string   `toml:"username" json:"username" yaml:"username"  env:"USERNAME"`
	Password  string   `toml:"password" json:"password" yaml:"password"  env:"PASSWORD"`
	Trace     bool     `toml:"trace" json:"trace" yaml:"trace"  env:"TRACE"`
	Metric    bool     `toml:"metric" json:"metric" yaml:"metric"  env:"METRIC"`

	client redis.UniversalClient
	log    *zerolog.Logger
}

func (m *Redis) Name() string {
	return AppName
}

func (i *Redis) Priority() int {
	return 697
}

// https://opentelemetry.io/ecosystem/registry/?s=redis&component=&language=go
// https://github.com/redis/go-redis/tree/master/extra/redisotel
func (m *Redis) Init() error {
	m.log = log.Sub(m.Name())
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    m.Endpoints,
		DB:       m.DB,
		Username: m.UserName,
		Password: m.Password,
	})

	if trace.Get().Enable && m.Trace {
		m.log.Info().Msg("enable redis trace")
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			return err
		}
	}

	if m.Metric {
		if err := redisotel.InstrumentMetrics(rdb); err != nil {
			return err
		}
	}

	m.client = rdb
	return nil
}

// 关闭数据库连接
func (m *Redis) Close(ctx context.Context) {
	if m.client == nil {
		return
	}

	err := m.client.Close()
	if err != nil {
		m.log.Error().Msgf("close redis client error, %s", err)
	}
	return
}
