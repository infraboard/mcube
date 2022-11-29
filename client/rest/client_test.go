package rest_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/client/rest"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/logger/zap"
)

func TestClient(t *testing.T) {
	c := rest.NewRESTClient()
	c.SetBaseURL("https://www.baidu.com/cmdb/api/v1")
	c.SetBearerTokenAuth("UAoVkI07gDGlfARUTToCA8JW")

	resp := make(map[string]any)
	err := c.Group("group1").Get("host22").
		Do(context.Background()).
		Into(response.NewData(&resp))
	if err != nil {
		t.Fatal(err)
	}

	t.Log(resp)
}

func init() {
	// 设置日志模式
	zap.DevelopmentSetup()
}
