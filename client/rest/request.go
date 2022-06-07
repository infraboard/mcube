package rest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/infraboard/mcube/client/negotiator"
	"github.com/infraboard/mcube/flowcontrol"
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
	}

	return r
}

// Request allows for building up a request to a server in a chained fashion.
// Any errors are stored until the end of your call, so you only have to
// check once.
type Request struct {
	c *RESTClient

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

func (r *Request) URL(p string) *Request {
	u, err := url.Parse(r.url)
	if err != nil {
		r.err = err
		return r
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

	ct := FilterFlags(r.headers.Get(CONTENT_TYPE_HEADER))
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
	resp := NewResponse()

	// 准备请求
	req, err := http.NewRequestWithContext(ctx, r.method, r.url, r.body)
	req.Header = r.headers
	req.URL.RawQuery = r.params.Encode()
	r.buildAuth(req)
	if err != nil {
		resp.err = err
		return resp
	}

	// 发起请求
	raw, err := r.c.client.Do(req)
	if err != nil {
		resp.err = err
		return resp
	}
	resp.statusCode = raw.StatusCode
	resp.headers = raw.Header
	resp.body = raw.Body
	return resp
}

func (r *Request) buildAuth(req *http.Request) {
	req.BasicAuth()
	switch r.authType {
	case BasicAuth:
		req.SetBasicAuth(r.user.Username, r.user.Password)
	case BearerToken:
		r.Header(AUTHORIZATION_HEADER, "Bearer "+r.token)
	}
}
