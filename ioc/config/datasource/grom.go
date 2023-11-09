package datasource

import (
	"context"
	"fmt"

	"github.com/infraboard/mcube/ioc"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	ioc.Config().Registry(&dataSource{
		Provider:    PROVIDER_MYSQL,
		Host:        "127.0.0.1",
		Port:        3306,
		DB:          "mcube",
		Username:    "root",
		Password:    "123456",
		Debug:       false,
		EnableTrace: false,
	})
}

type dataSource struct {
	Provider    PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"DATASOURCE_PROVIDER"`
	Host        string   `json:"host" yaml:"host" toml:"host" env:"DATASOURCE_HOST"`
	Port        int      `json:"port" yaml:"port" toml:"port" env:"DATASOURCE_PORT"`
	DB          string   `json:"database" yaml:"database" toml:"database" env:"DATASOURCE_DB"`
	Username    string   `json:"username" yaml:"username" toml:"username" env:"DATASOURCE_USERNAME"`
	Password    string   `json:"password" yaml:"password" toml:"password" env:"DATASOURCE_PASSWORD"`
	Debug       bool     `json:"debug" yaml:"debug" toml:"debug" env:"DATASOURCE_DEBUG"`
	EnableTrace bool     `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace"  env:"DATASOURCE_ENABLE_TRACE"`

	db *gorm.DB
	ioc.ObjectImpl
}

func (m *dataSource) Name() string {
	return DATASOURCE
}

func (m *dataSource) Init() error {
	db, err := gorm.Open(mysql.Open(m.DSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	if m.EnableTrace {
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
	return m.db.WithContext(ctx)
}

func (m *dataSource) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.DB,
	)
}
