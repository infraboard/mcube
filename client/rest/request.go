package rest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/infraboard/mcube/client/negotiator"
	"github.com/infraboard/mcube/flowcontrol"
	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
	"github.com/infraboard/mcube/logger"
)

// NewRequest creates a new request helper object.
func NewRequest(c *RESTClient) *Request {
	r := &Request{
		c:           c,
		rateLimiter: c.rateLimiter,
		timeout:     c.client.Timeout,
		url:         c.baseURL,
		headers:     c.headers,
		cookies:     c.cookies,
		authType:    c.authType,
		user:        c.user,
		token:       c.token,
		log:         c.log.Named("request"),
	}

	return r
}

// Request allows for building up a request to a server in a chained fashion.
// Any errors are stored until the end of your call, so you only have to
// check once.
type Request struct {
	c *RESTClient

	log         logger.Logger
	rateLimiter flowcontrol.RateLimiter
	timeout     time.Duration

	authType AuthType
	user     *User
	token    string

	// generic components accessible via method setters
	method  string
	url     string
	cookies []*http.Cookie
	headers http.Header
	params  url.Values
	body    io.Reader

	err error
}

// Method sets the verb this request will use.
func (r *Request) Method(verb string) *Request {
	r.method = verb
	return r
}

// SetRequestRate 设置请求速率, rate: 速率, capacity: 容量(控制并发量)
func (c *Request) SetRequestRate(rate float64, capacity int64) *Request {
	c.rateLimiter = tokenbucket.NewBucketWithRate(rate, capacity)
	return c
}

func (r *Request) URL(p string) *Request {
	u, err := url.Parse(r.url)
	if err != nil {
		r.err = err
		return r
	}

	for _, group := range r.c.groups {
		u.Path = path.Join(u.Path, group)
	}

	u.Path = path.Join(u.Path, p)
	r.url = u.String()
	return r
}

// Timeout makes the request use the given duration as an overall timeout for the
// request. Additionally, if set passes the value as "timeout" parameter in URL.
func (r *Request) Timeout(d time.Duration) *Request {
	if r.err != nil {
		return r
	}
	r.timeout = d
	return r
}

func (r *Request) Header(key string, values ...string) *Request {
	if r.err != nil {
		return r
	}
	if r.headers == nil {
		r.headers = http.Header{}
	}
	r.headers.Del(key)
	for _, value := range values {
		r.headers.Add(key, value)
	}
	return r
}

func (r *Request) Cookie(cs ...*http.Cookie) *Request {
	if r.err != nil {
		return r
	}
	if r.cookies == nil {
		r.cookies = make([]*http.Cookie, 0)
	}

	r.cookies = append(r.cookies, cs...)
	return r
}

// Param creates a query parameter with the given string value.
func (r *Request) Param(paramName, value string) *Request {
	if r.err != nil {
		return r
	}
	if r.params == nil {
		r.params = make(url.Values)
	}
	r.params[paramName] = append(r.params[paramName], value)
	return r
}

func (r *Request) Body(v any) *Request {
	if r.err != nil {
		return r
	}

	ct := HeaderFilterFlags(r.headers.Get(CONTENT_TYPE_HEADER))
	nt := negotiator.GetNegotiator(ct)

	b, err := nt.Encode(v)
	if err != nil {
		r.err = err
		return r
	}

	r.body = bytes.NewReader(b)
	return r
}

func (r *Request) Do(ctx context.Context) *Response {
	// 请求速率控制
	r.rateLimiter.Wait(1)

	// 请求响应对象
	resp := NewResponse(r.c)

	// 准备请求
	req, err := http.NewRequestWithContext(ctx, r.method, r.url, r.body)
	if err != nil {
		resp.err = err
		return resp
	}
	req.URL.RawQuery = r.params.Encode()

	//补充Header
	for k, vs := range r.headers {
		for i := range vs {
			req.Header.Set(k, vs[i])
		}
	}

	// 补充认证
	r.buildAuth(req)

	// 补充cookie
	for i := range r.cookies {
		req.AddCookie(r.cookies[i])
	}

	// debug信息
	r.debug(req)

	// 发起请求
	raw, err := r.c.client.Do(req)
	if err != nil {
		resp.err = err
		return resp
	}

	// 设置返回
	resp.withStatusCode(raw.StatusCode)
	resp.withHeader(raw.Header)
	resp.withBody(raw.Body)
	return resp
}

func (r *Request) debug(req *http.Request) {
	r.log.Debugf("[%s] %s", req.Method, req.URL.String())
	r.log.Debugf("Request Headers:")
	for k, v := range req.Header {
		r.log.Debugf("%s=%s", k, strings.Join(v, ","))
	}
}

func (r *Request) buildAuth(req *http.Request) {
	switch r.authType {
	case BasicAuth:
		req.SetBasicAuth(r.user.Username, r.user.Password)
	case BearerToken:
		req.Header.Set(AUTHORIZATION_HEADER, "Bearer "+r.token)
	}
}
