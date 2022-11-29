package rest_test

import (
	"context"
	"testing"

	"github.com/infraboard/mcube/client/rest"
	"github.com/infraboard/mcube/http/response"
)

func TestClient(t *testing.T) {
	c := rest.NewRESTClient()
	c.SetBaseURL("http://127.0.0.1:8060/cmdb/api/v1")
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
