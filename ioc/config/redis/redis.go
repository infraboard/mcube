package redis

import (
	"context"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

func init() {
	ioc.Config().Registry(&Redist{
		Database:    0,
		Endpoints:   []string{"127.0.0.1:6379"},
		EnableTrace: true,
	})
}

type Redist struct {
	Endpoints     []string `toml:"endpoints" json:"endpoints" yaml:"endpoints" env:"REDIS_ENDPOINTS" envSeparator:","`
	Database      int      `toml:"database" json:"database" yaml:"database"  env:"REDIS_DATABASE"`
	UserName      string   `toml:"username" json:"username" yaml:"username"  env:"REDIS_USERNAME"`
	Password      string   `toml:"password" json:"password" yaml:"password"  env:"REDIS_PASSWORD"`
	EnableTrace   bool     `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace"  env:"REDIS_ENABLE_TRACE"`
	EnableMetrics bool     `toml:"enable_metrics" json:"enable_metrics" yaml:"enable_metrics"  env:"REDIS_ENABLE_METRICS"`

	client redis.UniversalClient
	ioc.ObjectImpl
}

func (m *Redist) Name() string {
	return AppName
}

// https://opentelemetry.io/ecosystem/registry/?s=redis&component=&language=go
// https://github.com/redis/go-redis/tree/master/extra/redisotel
func (m *Redist) Init() error {
	rdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    m.Endpoints,
		DB:       m.Database,
		Username: m.UserName,
		Password: m.Password,
	})

	if application.App().Trace.Enable && m.EnableTrace {
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			return err
		}
	}

	if m.EnableMetrics {
		if err := redisotel.InstrumentMetrics(rdb); err != nil {
			return err
		}
	}

	m.client = rdb
	return nil
}

// 关闭数据库连接
func (m *Redist) Close(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	return m.Close(ctx)
}
