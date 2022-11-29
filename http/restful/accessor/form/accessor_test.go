package form_test

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/http/restful/accessor/form"
)

type Book struct {
	Title         string `form:"title"`
	Author        string `form:"author"`
	PublishedYear int    `form:"published_year"`
}

func TestPostFormEncoding(t *testing.T) {
	formData := make(url.Values)
	formData.Add("title", "test")
	formData.Add("author", "test")
	formData.Add("published_year", "2022")

	t.Log(formData.Encode())

	// Read
	bodyReader := bytes.NewReader([]byte(formData.Encode()))
	httpRequest, _ := http.NewRequest("POST", "/test", bodyReader)
	httpRequest.Header.Set("Content-Type", form.MIME_POST_FORM)
	request := restful.NewRequest(httpRequest)

	b := &Book{}
	err := request.ReadEntity(b)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(b)
}
