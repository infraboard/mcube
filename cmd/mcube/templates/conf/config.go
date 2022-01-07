package conf

// Template todo
const Template = `package conf

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	kc "github.com/infraboard/keyauth/client"

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
		App:     newDefaultAPP(),
		Log:     newDefaultLog(),
		MySQL:   newDefaultMySQL(),
		Mongo:   newDefaultMongoDB(),
		Cache:   newDefaultCache(),
		Keyauth: newDefaultKeyauth(),
	}
}

// Config 应用配置
type Config struct {
	App   *app   {{.Backquote}}toml:"app"{{.Backquote}}
	Log   *log   {{.Backquote}}toml:"log"{{.Backquote}}
	MySQL *mysql {{.Backquote}}toml:"mysql"{{.Backquote}}
	Mongo *mongodb {{.Backquote}}toml:"mongodb"{{.Backquote}}
	Keyauth  *keyauth  {{.Backquote}}toml:"keyauth"{{.Backquote}}
	Cache *_cache  {{.Backquote}}toml:"cache"{{.Backquote}}
}

type app struct {
	Name       string {{.Backquote}}toml:"name" env:"APP_NAME"{{.Backquote}}
	EncryptKey string {{.Backquote}}toml:"encrypt_key" env:"APP_ENCRYPT_KEY"{{.Backquote}}
	HTTP       *http  {{.Backquote}}toml:"http"{{.Backquote}}
	GRPC       *grpc  {{.Backquote}}toml:"grpc"{{.Backquote}}
}

func newDefaultAPP() *app {
	return &app{
		Name:       "cmdb",
		EncryptKey: "defualt app encrypt key",
		HTTP:       newDefaultHTTP(),
		GRPC:       newDefaultGRPC(),
	}
}

type http struct {
	Host      string {{.Backquote}}toml:"host" env:"HTTP_HOST"{{.Backquote}}
	Port      string {{.Backquote}}toml:"port" env:"HTTP_PORT"{{.Backquote}}
	EnableSSL bool   {{.Backquote}}toml:"enable_ssl" env:"HTTP_ENABLE_SSL"{{.Backquote}}
	CertFile  string {{.Backquote}}toml:"cert_file" env:"HTTP_CERT_FILE"{{.Backquote}}
	KeyFile   string {{.Backquote}}toml:"key_file" env:"HTTP_KEY_FILE"{{.Backquote}}
}

func (a *http) Addr() string {
	return a.Host + ":" + a.Port
}

func newDefaultHTTP() *http {
	return &http{
		Host: "127.0.0.1",
		Port: "8050",
	}
}

type grpc struct {
	Host      string {{.Backquote}}toml:"host" env:"GRPC_HOST"{{.Backquote}}
	Port      string {{.Backquote}}toml:"port" env:"GRPC_PORT"{{.Backquote}}
	EnableSSL bool   {{.Backquote}}toml:"enable_ssl" env:"GRPC_ENABLE_SSL"{{.Backquote}}
	CertFile  string {{.Backquote}}toml:"cert_file" env:"GRPC_CERT_FILE"{{.Backquote}}
	KeyFile   string {{.Backquote}}toml:"key_file" env:"GRPC_KEY_FILE"{{.Backquote}}
}

func (a *grpc) Addr() string {
	return a.Host + ":" + a.Port
}

func newDefaultGRPC() *grpc {
	return &grpc{
		Host: "127.0.0.1",
		Port: "18050",
	}
}

type log struct {
	Level   string    {{.Backquote}}toml:"level" env:"LOG_LEVEL"{{.Backquote}}
	PathDir string    {{.Backquote}}toml:"path_dir" env:"LOG_PATH_DIR"{{.Backquote}}
	Format  LogFormat {{.Backquote}}toml:"format" env:"LOG_FORMAT"{{.Backquote}}
	To      LogTo     {{.Backquote}}toml:"to" env:"LOG_TO"{{.Backquote}}
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
type keyauth struct {
	Host      string {{.Backquote}}toml:"host" env:"KEYAUTH_HOST"{{.Backquote}}
	Port      string {{.Backquote}}toml:"port" env:"KEYAUTH_PORT"{{.Backquote}}
	ClientID string {{.Backquote}}toml:"client_id" env:"KEYAUTH_CLIENT_ID"{{.Backquote}}
	ClientSecret string {{.Backquote}}toml:"client_secret" env:"KEYAUTH_CLIENT_SECRET"{{.Backquote}}
}

func (a *keyauth) Addr() string {
	return a.Host + ":" + a.Port
}

func (a *keyauth) Client() (*kc.Client, error) {
	if kc.C() == nil {
		conf := kc.NewDefaultConfig()
		conf.SetAddress(a.Addr())
		conf.SetClientCredentials(a.ClientID, a.ClientSecret)
		client, err := kc.NewClient(conf)
		if err != nil {
			return nil, err
		}
		kc.SetGlobal(client)
	}

	return kc.C(), nil
}

func newDefaultKeyauth() *keyauth {
	return &keyauth{}
}

func newDefaultMongoDB() *mongodb {
	return &mongodb{
		Database:  "",
		Endpoints: []string{"127.0.0.1:27017"},
	}
}

type mongodb struct {
	Endpoints []string {{.Backquote}}toml:"endpoints" env:"MONGO_ENDPOINTS" envSeparator:","{{.Backquote}}
	UserName  string   {{.Backquote}}toml:"username" env:"MONGO_USERNAME"{{.Backquote}}
	Password  string   {{.Backquote}}toml:"password" env:"MONGO_PASSWORD"{{.Backquote}}
	Database  string   {{.Backquote}}toml:"database" env:"MONGO_DATABASE"{{.Backquote}}
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
	Host        string {{.Backquote}}toml:"host" env:"MYSQL_HOST"{{.Backquote}}
	Port        string {{.Backquote}}toml:"port" env:"MYSQL_PORT"{{.Backquote}}
	UserName    string {{.Backquote}}toml:"username" env:"MYSQL_USERNAME"{{.Backquote}}
	Password    string {{.Backquote}}toml:"password" env:"MYSQL_PASSWORD"{{.Backquote}}
	Database    string {{.Backquote}}toml:"database" env:"MYSQL_DATABASE"{{.Backquote}}
	MaxOpenConn int    {{.Backquote}}toml:"max_open_conn" env:"MYSQL_MAX_OPEN_CONN"{{.Backquote}}
	MaxIdleConn int    {{.Backquote}}toml:"max_idle_conn" env:"MYSQL_MAX_IDLE_CONN"{{.Backquote}}
	MaxLifeTime int    {{.Backquote}}toml:"max_life_time" env:"MYSQL_MAX_LIFE_TIME"{{.Backquote}}
	MaxIdleTime int    {{.Backquote}}toml:"max_idle_time" env:"MYSQL_MAX_IDLE_TIME"{{.Backquote}}
	lock        sync.Mutex
}

func newDefaultMySQL() *mysql {
	return &mysql{
		Database:    "{{.Name}}",
		Host:        "127.0.0.1",
		Port:        "3306",
		MaxOpenConn: 200,
		MaxIdleConn: 100,
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
	if m.MaxLifeTime != 0 {
		db.SetConnMaxLifetime(time.Second * time.Duration(m.MaxLifeTime))
	}
	if m.MaxIdleConn != 0 {
		db.SetConnMaxIdleTime(time.Second * time.Duration(m.MaxIdleTime))
	}

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
	Type   string         {{.Backquote}}toml:"type" json:"type" yaml:"type" env:"CACHE_TYPE"{{.Backquote}}
	Memory *memory.Config {{.Backquote}}toml:"memory" json:"memory" yaml:"memory"{{.Backquote}}
	Redis  *redis.Config  {{.Backquote}}toml:"redis" json:"redis" yaml:"redis"{{.Backquote}}
}`
