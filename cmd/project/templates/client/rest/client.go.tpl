package rest

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/client/rest"
)

func NewClient(conf *Config) *ClientSet {
	c := rest.NewRESTClient()
	c.SetBearerTokenAuth(conf.Token)
	c.SetHeader(restful.HEADER_ContentType, restful.MIME_JSON)
	c.SetHeader(restful.HEADER_Accept, restful.MIME_JSON)
	c.SetBaseURL(conf.Address + conf.PathPrefix)
	return &ClientSet{
		c: c,
	}
}

type ClientSet struct {
	c *rest.RESTClient
}