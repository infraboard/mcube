package datasource

import (
	"fmt"

	"github.com/infraboard/mcube/ioc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	ioc.Config().Registry(&dataSource{
		Provider: PROVIDER_MYSQL,
		Host:     "127.0.0.1",
		Port:     3306,
		DB:       "mcube",
		Username: "root",
		Password: "123456",
	})
}

type dataSource struct {
	Provider PROVIDER `json:"provider" yaml:"provider" toml:"provider" env:"DATASOURCE_PROVIDER"`
	Host     string   `json:"host" yaml:"host" toml:"host" env:"DATASOURCE_HOST"`
	Port     int      `json:"port" yaml:"port" toml:"port" env:"DATASOURCE_PORT"`
	DB       string   `json:"database" yaml:"database" toml:"database" env:"DATASOURCE_DB"`
	Username string   `json:"username" yaml:"username" toml:"username" env:"DATASOURCE_USERNAME"`
	Password string   `json:"password" yaml:"password" toml:"password" env:"DATASOURCE_PASSWORD"`

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
	m.db = db
	return nil
}

func (m *dataSource) GetDB() *gorm.DB {
	return m.db
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
