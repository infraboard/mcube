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

const (
	MONGODB = "mongodb"
)

func init() {
	ioc.Config().Registry(&MongoDB{
		UserName:  "mcube",
		Password:  "123456",
		Database:  "mcube",
		AuthDB:    "",
		Endpoints: []string{"127.0.0.1:27017"},
	})
}

func GetMongoDB() *MongoDB {
	return ioc.Config().Get(MONGODB).(*MongoDB)
}

type MongoDB struct {
	ioc.IocObjectImpl

	Endpoints []string `toml:"endpoints" env:"MONGO_ENDPOINTS" envSeparator:","`
	UserName  string   `toml:"username" env:"MONGO_USERNAME"`
	Password  string   `toml:"password" env:"MONGO_PASSWORD"`
	Database  string   `toml:"database" env:"MONGO_DATABASE"`
	AuthDB    string   `toml:"auth_db" env:"MONGO_AUTH_DB"`

	client *mongo.Client
}

func (m *MongoDB) Name() string {
	return MONGODB
}

func (m *MongoDB) Init() error {
	if m.client == nil {
		conn, err := m.getClient()
		if err != nil {
			return err
		}
		m.client = conn
	}
	return nil
}

func (m *MongoDB) GetAuthDB() string {
	if m.AuthDB != "" {
		return m.AuthDB
	}

	return m.Database
}

func (m *MongoDB) GetDB() (*mongo.Database, error) {
	return m.client.Database(m.Database), nil
}

// 关闭数据库连接
func (m *MongoDB) Close(ctx context.Context) error {
	if m.client == nil {
		return nil
	}

	return m.client.Disconnect(ctx)
}

// Client 获取一个全局的mongodb客户端连接
func (m *MongoDB) Client() *mongo.Client {
	return m.client
}

func (m *MongoDB) getClient() (*mongo.Client, error) {
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
	opts.Monitor = otelmongo.NewMonitor(
		otelmongo.WithCommandAttributeDisabled(true),
	)

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
