package client

// ClientProxyTemplate todo
const ClientProxyTemplate = `package client

import (
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"

	"{{.PKG}}/pkg/example"
)

var (
	client *Client
)

// SetGlobal todo
func SetGlobal(cli *Client) {
	client = cli
}

// C Global
func C() *Client {
	return client
}

// NewClient todo
func NewClient(conf *Config) (*Client, error) {
	zap.DevelopmentSetup()
	log := zap.L()

	conn, err := grpc.Dial(conf.address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(conf.Authentication))
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
		log:  log,
	}, nil
}

// Client 客户端
type Client struct {
	conn *grpc.ClientConn
	log  logger.Logger
}

// Example todo
func (c *Client) Example() example.ServiceClient {
	return example.NewServiceClient(c.conn)
}`
