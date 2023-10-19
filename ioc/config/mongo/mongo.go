package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/ioc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func init() {
	ioc.Config().Registry(&mongoDB{
		Database:  "mcube",
		AuthDB:    "admin",
		Endpoints: []string{"127.0.0.1:27017"},
	})
}

type mongoDB struct {
	Endpoints   []string `toml:"endpoints" json:"endpoints" yaml:"endpoints" env:"MONGO_ENDPOINTS" envSeparator:","`
	UserName    string   `toml:"username" json:"username" yaml:"username"  env:"MONGO_USERNAME"`
	Password    string   `toml:"password" json:"password" yaml:"password"  env:"MONGO_PASSWORD"`
	Database    string   `toml:"database" json:"database" yaml:"database"  env:"MONGO_DATABASE"`
	AuthDB      string   `toml:"auth_db" json:"auth_db" yaml:"auth_db"  env:"MONGO_AUTH_DB"`
	EnableTrace bool     `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace"  env:"MONGO_ENABLE_TRACE"`

	client *mongo.Client
	ioc.ObjectImpl
}

func (m *mongoDB) Name() string {
	return MONGODB
}

func (m *mongoDB) Init() error {
	conn, err := m.getClient()
	if err != nil {
		return err
	}
	m.client = conn
	return nil
}

// 关闭数据库连接
func (m *mongoDB) Close(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	return m.client.Disconnect(ctx)
}

func (m *mongoDB) GetAuthDB() string {
	if m.AuthDB != "" {
		return m.AuthDB
	}

	return m.Database
}

func (m *mongoDB) GetDB() *mongo.Database {
	return m.client.Database(m.Database)
}

// Client 获取一个全局的mongodb客户端连接
func (m *mongoDB) Client() *mongo.Client {
	return m.client
}

func (m *mongoDB) getClient() (*mongo.Client, error) {
	opts := options.Client()

	if m.UserName != "" && m.Password != "" {
		cred := options.Credential{
			AuthSource: m.GetAuthDB(),
		}

		cred.Username = m.UserName
		cred.Password = m.Password
		cred.PasswordSet = true
		opts.SetAuth(cred)
	}
	opts.SetHosts(m.Endpoints)
	opts.SetConnectTimeout(5 * time.Second)
	if m.EnableTrace {
		opts.Monitor = otelmongo.NewMonitor(
			otelmongo.WithCommandAttributeDisabled(true),
		)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*5))
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("new mongodb client error, %s", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongodb server(%s) error, %s", m.Endpoints, err)
	}

	return client, nil
}
