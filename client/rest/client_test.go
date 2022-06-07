package rest_test

import (
	"testing"

	"github.com/infraboard/mcube/client/rest"
)

func TestClient(t *testing.T) {
	c := rest.NewRESTClient()
	c.Get("").Do(nil)
}
