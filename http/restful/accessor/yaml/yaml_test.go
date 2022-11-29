package yaml_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/restful/accessor/yaml"
)

type Book struct {
	Title         string
	Author        string
	PublishedYear int
}

func TestYAMLEncoding(t *testing.T) {
	b := Book{"Singing for Dummies", "john doe", 2015}

	// Write
	httpWriter := httptest.NewRecorder()
	resp := restful.NewResponse(httpWriter)
	resp.SetRequestAccepts(yaml.MIME_YAML)
	resp.WriteEntity(b)
	t.Log(httpWriter.Body.String())

	// Read
	bodyReader := bytes.NewReader(httpWriter.Body.Bytes())
	httpRequest, _ := http.NewRequest("GET", "/test", bodyReader)
	httpRequest.Header.Set("Content-Type", yaml.MIME_YAML)
	request := restful.NewRequest(httpRequest)

	b = Book{}
	err := request.ReadEntity(&b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
}
