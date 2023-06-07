package client_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"{{.PKG}}/apps/book"
	"{{.PKG}}/client"
)

func TestBookQuery(t *testing.T) {
	should := assert.New(t)

	conf := client.NewDefaultConfig()
	// 设置GRPC服务地址
	// conf.SetAddress("127.0.0.1:8050")
	// 携带认证信息
	// conf.SetClientCredentials("secret_id", "secret_key")
	c, err := client.NewClient(conf)
	if should.NoError(err) {
		resp, err := c.Book().QueryBook(
			context.Background(),
			book.NewQueryBookRequest(),
		)
		should.NoError(err)
		fmt.Println(resp.Items)
	}
}
