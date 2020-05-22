package conf

// Template todo
const Template = `package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db *sql.DB
)

func newConfig() *Config {
	return &Config{
		App:   newDefaultAPP(),
		Log:   newDefaultLog(),
		MySQL: newDefaultMySQL(),
		Auth:  newDefaultAuth(),
	}
}

// Config 应用配置
type Config struct {
	App   *app   {{.Backquote}}toml:"app"{{.Backquote}}
	Log   *log   {{.Backquote}}toml:"log"{{.Backquote}}
	MySQL *mysql {{.Backquote}}toml:"mysql"{{.Backquote}}
	Auth  *auth  {{.Backquote}}toml:"auth"{{.Backquote}}
}

type app struct {
	Name      string {{.Backquote}}toml:"name" env:"D_APP_NAME"{{.Backquote}}
	Host      string {{.Backquote}}toml:"host" env:"D_APP_HOST"{{.Backquote}}
	Port      string {{.Backquote}}toml:"port" env:"D_APP_PORT"{{.Backquote}}
	Key       string {{.Backquote}}toml:"key" env:"D_APP_KEY"{{.Backquote}}
	EnableSSL bool   {{.Backquote}}toml:"enable_ssl" env:"D_APP_ENABLE_SSL"{{.Backquote}}
	CertFile  string {{.Backquote}}toml:"cert_file" env:"D_APP_CERT_FILE"{{.Backquote}}
	KeyFile   string {{.Backquote}}toml:"key_file" env:"D_APP_KEY_FILE"{{.Backquote}}
}

func (a *app) Addr() string {
	return a.Host + ":" + a.Port
}
func newDefaultAPP() *app {
	return &app{
		Name: "{{.Name}}",
		Host: "127.0.0.1",
		Port: "8050",
		Key:  "default",
	}
}

type log struct {
	Level   string    {{.Backquote}}toml:"level" env:"D_LOG_LEVEL"{{.Backquote}}
	PathDir string    {{.Backquote}}toml:"path_dir" env:"D_LOG_PATH_DIR"{{.Backquote}}
	Format  LogFormat {{.Backquote}}toml:"format" env:"D_LOG_FORMAT"{{.Backquote}}
	To      LogTo     {{.Backquote}}toml:"to" env:"D_LOG_TO"{{.Backquote}}
}

func newDefaultLog() *log {
	return &log{
		Level:   "debug",
		PathDir: "logs",
		Format:  "text",
		To:      "stdout",
	}
}

// Auth auth 配置
type auth struct {
	Address string {{.Backquote}}toml:"address" env:"D_AUTH_ADDRESS"{{.Backquote}}
}

func newDefaultAuth() *auth {
	return &auth{}
}

type mysql struct {
	Host        string {{.Backquote}}toml:"host" env:"D_MYSQL_HOST"{{.Backquote}}
	Port        string {{.Backquote}}toml:"port" env:"D_MYSQL_PORT"{{.Backquote}}
	UserName    string {{.Backquote}}toml:"username" env:"D_MYSQL_USERNAME"{{.Backquote}}
	Password    string {{.Backquote}}toml:"password" env:"D_MYSQL_PASSWORD"{{.Backquote}}
	Database    string {{.Backquote}}toml:"database" env:"D_MYSQL_DATABASE"{{.Backquote}}
	MaxOpenConn int    {{.Backquote}}toml:"max_open_conn" env:"D_MYSQL_MAX_OPEN_CONN"{{.Backquote}}
	MaxIdleConn int    {{.Backquote}}toml:"max_idle_conn" env:"D_MYSQL_MAX_IDLE_CONN"{{.Backquote}}
	MaxLifeTime int    {{.Backquote}}toml:"max_life_time" env:"D_MYSQL_MAX_LIFE_TIME"{{.Backquote}}
	lock        sync.Mutex
}

func newDefaultMySQL() *mysql {
	return &mysql{
		Database:    "{{.Name}}",
		Host:        "127.0.0.1",
		Port:        "3306",
		MaxOpenConn: 200,
		MaxIdleConn: 16,
		MaxLifeTime: 300,
	}
}

// getDBConn use to get db connection pool
func (m *mysql) getDBConn() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.UserName, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	db.SetConnMaxLifetime(time.Minute * time.Duration(m.MaxLifeTime))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	return db, nil
}
func (m *mysql) GetDB() (*sql.DB, error) {
	// 加载全局数据量单例
	m.lock.Lock()
	defer m.lock.Unlock()
	if db == nil {
		conn, err := m.getDBConn()
		if err != nil {
			return nil, err
		}
		db = conn
	}
	return db, nil
}`
