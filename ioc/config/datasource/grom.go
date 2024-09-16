package datasource

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	"github.com/rs/zerolog"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &dataSource{
	Provider: PROVIDER_MYSQL,
	Host:     "127.0.0.1",
	Port:     3306,
	DB:       application.Get().Name(),
	Debug:    false,
	Trace:    true,
}

type dataSource struct {
	ioc.ObjectImpl
	Provider PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"PROVIDER"`
	Host     string   `json:"host" yaml:"host" toml:"host" env:"HOST"`
	Port     int      `json:"port" yaml:"port" toml:"port" env:"PORT"`
	DB       string   `json:"database" yaml:"database" toml:"database" env:"DB"`
	Username string   `json:"username" yaml:"username" toml:"username" env:"USERNAME"`
	Password string   `json:"password" yaml:"password" toml:"password" env:"PASSWORD"`
	Debug    bool     `json:"debug" yaml:"debug" toml:"debug" env:"DEBUG"`
	Trace    bool     `toml:"trace" json:"trace" yaml:"trace"  env:"TRACE"`

	db  *gorm.DB
	log *zerolog.Logger
}

func (m *dataSource) Name() string {
	return AppName
}

func (i *dataSource) Priority() int {
	return 699
}

func (m *dataSource) Init() error {
	m.log = log.Sub(m.Name())
	db, err := gorm.Open(m.Dialector(), &gorm.Config{})
	if err != nil {
		return err
	}

	if trace.Get().Enable && m.Trace {
		m.log.Info().Msg("enable gorm trace")
		if err := db.Use(otelgorm.NewPlugin()); err != nil {
			return err
		}
	}

	if m.Debug {
		db = db.Debug()
	}

	m.db = db
	return nil
}

// 关闭数据库连接
func (m *dataSource) Close(ctx context.Context) error {
	if m.db == nil {
		return nil
	}

	d, err := m.db.DB()
	if err != nil {
		return err
	}
	return d.Close()
}

// 从上下文中获取事物, 如果获取不到则直接返回 无事物的DB对象
func (m *dataSource) GetTransactionOrDB(ctx context.Context) *gorm.DB {
	db := GetTransactionFromCtx(ctx)
	if db != nil {
		return db
	}
	return m.db
}

func (m *dataSource) Dialector() gorm.Dialector {
	switch m.Provider {
	case PROVIDER_POSTGRES:
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			m.Host,
			m.Username,
			m.Password,
			m.DB,
			m.Port,
		)
		return postgres.Open(dsn)
	default:
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			m.Username,
			m.Password,
			m.Host,
			m.Port,
			m.DB,
		)
		return mysql.Open(dsn)
	}
}
