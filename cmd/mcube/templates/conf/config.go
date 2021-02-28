package conf

// Template todo
const Template = `package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/infraboard/mcube/cache/memory"
	"github.com/infraboard/mcube/cache/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db *sql.DB
	mgoclient *mongo.Client
)

func newConfig() *Config {
	return &Config{
		App:   newDefaultAPP(),
		Log:   newDefaultLog(),
		MySQL: newDefaultMySQL(),
		Mongo: newDefaultMongoDB(),
		Cache: newDefaultCache(),
		Auth:  newDefaultAuth(),
	}
}

// Config 应用配置
type Config struct {
	App   *app   {{.Backquote}}toml:"app"{{.Backquote}}
	Log   *log   {{.Backquote}}toml:"log"{{.Backquote}}
	MySQL *mysql {{.Backquote}}toml:"mysql"{{.Backquote}}
	Mongo *mongodb {{.Backquote}}toml:"mongodb"{{.Backquote}}
	Auth  *auth  {{.Backquote}}toml:"auth"{{.Backquote}}
	Cache *_cache  {{.Backquote}}toml:"cache"{{.Backquote}}
}

type app struct {
	Name      string {{.Backquote}}toml:"name" env:"MCUBE_APP_NAME"{{.Backquote}}
	Host      string {{.Backquote}}toml:"host" env:"MCUBE_APP_HOST"{{.Backquote}}
	Port      string {{.Backquote}}toml:"port" env:"MCUBE_APP_PORT"{{.Backquote}}
	Key       string {{.Backquote}}toml:"key" env:"MCUBE_APP_KEY"{{.Backquote}}
	EnableSSL bool   {{.Backquote}}toml:"enable_ssl" env:"MCUBE_APP_ENABLE_SSL"{{.Backquote}}
	CertFile  string {{.Backquote}}toml:"cert_file" env:"MCUBE_APP_CERT_FILE"{{.Backquote}}
	KeyFile   string {{.Backquote}}toml:"key_file" env:"MCUBE_APP_KEY_FILE"{{.Backquote}}
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
	Level   string    {{.Backquote}}toml:"level" env:"MCUBE_LOG_LEVEL"{{.Backquote}}
	PathDir string    {{.Backquote}}toml:"path_dir" env:"MCUBE_LOG_PATH_DIR"{{.Backquote}}
	Format  LogFormat {{.Backquote}}toml:"format" env:"MCUBE_LOG_FORMAT"{{.Backquote}}
	To      LogTo     {{.Backquote}}toml:"to" env:"MCUBE_LOG_TO"{{.Backquote}}
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
	Address string {{.Backquote}}toml:"address" env:"MCUBE_AUTH_ADDRESS"{{.Backquote}}
}

func newDefaultAuth() *auth {
	return &auth{}
}

func newDefaultMongoDB() *mongodb {
	return &mongodb{
		Database:  "",
		Endpoints: []string{"127.0.0.1:27017"},
	}
}

type mongodb struct {
	Endpoints []string {{.Backquote}}toml:"endpoints" env:"MCUBE_MONGO_ENDPOINTS" envSeparator:","{{.Backquote}}
	UserName  string   {{.Backquote}}toml:"username" env:"MCUBE_MONGO_USERNAME"{{.Backquote}}
	Password  string   {{.Backquote}}toml:"password" env:"MCUBE_MONGO_PASSWORD"{{.Backquote}}
	Database  string   {{.Backquote}}toml:"database" env:"MCUBE_MONGO_DATABASE"{{.Backquote}}
}

// Client 获取一个全局的mongodb客户端连接
func (m *mongodb) Client() *mongo.Client {
	if mgoclient == nil {
		panic("please load mongo client first")
	}

	return mgoclient
}

func (m *mongodb) GetDB() *mongo.Database {
	return m.Client().Database(m.Database)
}

func (m *mongodb) getClient() (*mongo.Client, error) {
	opts := options.Client()

	cred := options.Credential{
		AuthSource: m.Database,
	}

	if m.UserName != "" && m.Password != "" {
		cred.Username = m.UserName
		cred.Password = m.Password
		cred.PasswordSet = true
		opts.SetAuth(cred)
	}
	opts.SetHosts(m.Endpoints)
	opts.SetConnectTimeout(5 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("new mongodb client error, %s", err)
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("ping mongodb server(%s) error, %s", m.Endpoints, err)
	}

	return client, nil
}

type mysql struct {
	Host        string {{.Backquote}}toml:"host" env:"MCUBE_MYSQL_HOST"{{.Backquote}}
	Port        string {{.Backquote}}toml:"port" env:"MCUBE_MYSQL_PORT"{{.Backquote}}
	UserName    string {{.Backquote}}toml:"username" env:"MCUBE_MYSQL_USERNAME"{{.Backquote}}
	Password    string {{.Backquote}}toml:"password" env:"MCUBE_MYSQL_PASSWORD"{{.Backquote}}
	Database    string {{.Backquote}}toml:"database" env:"MCUBE_MYSQL_DATABASE"{{.Backquote}}
	MaxOpenConn int    {{.Backquote}}toml:"max_open_conn" env:"MCUBE_MYSQL_MAX_OPEN_CONN"{{.Backquote}}
	MaxIdleConn int    {{.Backquote}}toml:"max_idle_conn" env:"MCUBE_MYSQL_MAX_IDLE_CONN"{{.Backquote}}
	MaxLifeTime int    {{.Backquote}}toml:"max_life_time" env:"MCUBE_MYSQL_MAX_LIFE_TIME"{{.Backquote}}
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
}

func newDefaultCache() *_cache {
	return &_cache{
		Type:   "memory",
		Memory: memory.NewDefaultConfig(),
		Redis:  redis.NewDefaultConfig(),
	}
}

type _cache struct {
	Type   string         {{.Backquote}}toml:"type" json:"type" yaml:"type" env:"MCUBE_CACHE_TYPE"{{.Backquote}}
	Memory *memory.Config {{.Backquote}}toml:"memory" json:"memory" yaml:"memory"{{.Backquote}}
	Redis  *redis.Config  {{.Backquote}}toml:"redis" json:"redis" yaml:"redis"{{.Backquote}}
}`
