package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &mongoDB{
	Database:  application.Get().GetAppName(),
	AuthDB:    "admin",
	Endpoints: []string{"127.0.0.1:27017"},
	Trace:     true,
}

type mongoDB struct {
	Endpoints []string `toml:"endpoints" json:"endpoints" yaml:"endpoints" env:"ENDPOINTS" envSeparator:","`
	UserName  string   `toml:"username" json:"username" yaml:"username"  env:"USERNAME"`
	Password  string   `toml:"password" json:"password" yaml:"password"  env:"PASSWORD"`
	Database  string   `toml:"database" json:"database" yaml:"database"  env:"DATABASE"`
	AuthDB    string   `toml:"auth_db" json:"auth_db" yaml:"auth_db"  env:"AUTH_DB"`
	Trace     bool     `toml:"trace" json:"trace" yaml:"trace"  env:"TRACE"`

	client *mongo.Client
	ioc.ObjectImpl
	log *zerolog.Logger
}

func (m *mongoDB) Name() string {
	return AppName
}

func (i *mongoDB) Priority() int {
	return 698
}

func (m *mongoDB) Init() error {
	m.log = log.Sub(m.Name())

	conn, err := m.getClient()
	if err != nil {
		return err
	}
	m.client = conn
	return nil
}

// 关闭数据库连接
func (m *mongoDB) Close(ctx context.Context) {
	if m.client == nil {
		return
	}

	err := m.client.Disconnect(ctx)
	if err != nil {
		m.log.Error().Msgf("close error, %s", err)
	}
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
	if trace.Get().Enable && m.Trace {
		m.log.Info().Msg("enable mongodb trace")
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
