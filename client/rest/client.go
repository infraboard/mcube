package rest

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/infraboard/mcube/client/negotiator"
	"github.com/infraboard/mcube/flowcontrol"
	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

// NewRESTClient creates a new RESTClient. This client performs generic REST functions
// such as Get, Put, Post, and Delete on specified paths.
func NewRESTClient() *RESTClient {
	client := http.DefaultClient
	// 保存Transport信息, 便于修改
	transport := http.DefaultTransport.(*http.Transport)
	client.Transport = transport
	return &RESTClient{
		rateLimiter: tokenbucket.NewBucketWithRate(10, 10),
		client:      client,
		log:         logger.Sub("client.rest"),
		headers:     NewDefaultHeader(),
		transport:   transport,
	}
}

func NewDefaultHeader() http.Header {
	header := http.Header{}
	header.Add("Accept-Encoding", "gzip")
	header.Add(CONTENT_TYPE_HEADER, string(negotiator.MIME_JSON))
	header.Add(ACCEPT_HEADER, string(negotiator.MIME_JSON))
	return header
}

type RESTClient struct {
	rateLimiter flowcontrol.RateLimiter
	transport   *http.Transport
	client      *http.Client
	cookies     []*http.Cookie
	headers     http.Header
	log         *zerolog.Logger
	baseURL     string

	authType AuthType
	user     *User
	token    string

	provider    oteltrace.TracerProvider
	propagators propagation.TextMapPropagator
	tr          oteltrace.Tracer
}

func (c *RESTClient) SetBaseURL(url string) *RESTClient {
	c.baseURL = strings.TrimRight(url, "/")
	return c
}

func (c *RESTClient) SetTLSConfig(conf *tls.Config) *RESTClient {
	c.transport.TLSClientConfig = conf
	return c
}

func (c *RESTClient) Clone() *RESTClient {
	cloned := &RESTClient{}
	*cloned = *c
	return cloned
}

func (c *RESTClient) Group(urlPath string) *RESTClient {
	cloned := c.Clone()

	baseURL, err := url.JoinPath(c.baseURL, urlPath)
	if err != nil {
		panic(err)
	}
	cloned.baseURL = baseURL
	return cloned
}

func (c *RESTClient) SetTimeout(t time.Duration) *RESTClient {
	c.client.Timeout = t
	return c
}

// SetRequestRate 设置全局请求速率, rate: 速率, capacity: 容量(控制并发量)
func (c *RESTClient) SetRequestRate(rate float64, capacity int64) *RESTClient {
	c.rateLimiter = tokenbucket.NewBucketWithRate(rate, capacity)
	return c
}

func (c *RESTClient) SetHeader(key string, values ...string) *RESTClient {
	if c.headers == nil {
		c.headers = http.Header{}
	}
	c.headers.Del(key)
	for _, value := range values {
		c.headers.Add(key, value)
	}
	return c
}

// SetCookie method sets an array of cookies in the client instance.
// These cookies will be added to all the request raised from this client instance.
//
//	cookies := []*http.Cookie{
//		&http.Cookie{
//			Name:"key-1",
//			Value:"This is cookie 1 value",
//		},
//		&http.Cookie{
//			Name:"key2-2",
//			Value:"This is cookie 2 value",
//		},
//	}
//
//	// Setting a cookies into
//	client.SetCookie(cookies...)
func (c *RESTClient) SetCookie(cs ...*http.Cookie) *RESTClient {
	if c.cookies == nil {
		c.cookies = make([]*http.Cookie, 0)
	}
	c.cookies = append(c.cookies, cs...)
	return c
}

// SetContentType set the Content-Type header of the request.
// application/json
func (c *RESTClient) SetContentType(contentType string) {
	c.SetHeader(CONTENT_TYPE_HEADER, contentType)
}

// SetBasicAuth method sets the basic authentication header in the current HTTP request.
//
// For Example:
//
//	Authorization: Basic <base64-encoded-value>
//
// To set the header for username "go-resty" and password "welcome"
//
//	client.SetBasicAuth("mcube", "welcome")
//
// This method overrides the credentials set by method `Client.SetBasicAuth`.
func (c *RESTClient) SetBasicAuth(username, password string) *RESTClient {
	c.authType = BasicAuth
	c.user = &User{Username: username, Password: password}
	return c
}

// SetAuthToken method sets the auth token header(Default Scheme: Bearer) in the current HTTP request. Header example:
//
//	Authorization: Bearer <auth-token-value-comes-here>
//
// For Example: To set bearer token BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F
//
//	client.SetBearerToken("BC594900518B4F7EAC75BD37F019E08FBC594900518B4F7EAC75BD37F019E08F")
//
// This method overrides the Auth token set by method `Client.SetAuthToken`.
func (c *RESTClient) SetBearerTokenAuth(token string) *RESTClient {
	c.authType = BearerToken
	c.token = token
	return c
}

// Verb begins a request with a verb (GET, POST, PUT, DELETE).
//
// Example usage of RESTClient's request building interface:
// c, err := NewRESTClient(...)
// if err != nil { ... }
// resp, err := c.Verb("GET").
//
//	Path("pods").
//	SelectorParam("labels", "area=staging").
//	Timeout(10*time.Second).
//	Do()
//
// if err != nil { ... }
// list, ok := resp.(*api.PodList)
func (c *RESTClient) Method(verb string) *Request {
	return NewRequest(c).Method(verb)
}

// Post begins a POST request. Short for c.Verb("POST").
func (c *RESTClient) Post(path string) *Request {
	return c.Method("POST").URL(path)
}

// Put begins a PUT request. Short for c.Verb("PUT").
func (c *RESTClient) Put(path string) *Request {
	return c.Method("PUT").URL(path)
}

// Patch begins a PATCH request. Short for c.Verb("Patch").
func (c *RESTClient) Patch(path string) *Request {
	return c.Method("PATCH").URL(path)
}

// Get begins a GET request. Short for c.Verb("GET").
func (c *RESTClient) Get(path string) *Request {
	return c.Method("GET").URL(path)
}

// Delete begins a DELETE request. Short for c.Verb("DELETE").
func (c *RESTClient) Delete(path string) *Request {
	return c.Method("DELETE").URL(path)
}
