package rest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/infraboard/mcube/client/negotiator"
	"github.com/infraboard/mcube/flowcontrol"
	"github.com/infraboard/mcube/flowcontrol/tokenbucket"
	"github.com/infraboard/mcube/http/queryparams"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/otel/trace"
)

// NewRequest creates a new request helper object.
func NewRequest(c *RESTClient) *Request {
	r := &Request{
		c:           c,
		rateLimiter: c.rateLimiter,
		timeout:     c.client.Timeout,
		basePath:    c.baseURL,
		headers:     c.headers,
		cookies:     c.cookies,
		authType:    c.authType,
		user:        c.user,
		token:       c.token,
		log:         logger.Sub("http.request"),
	}

	return r
}

// Request allows for building up a request to a server in a chained fashion.
// Any errors are stored until the end of your call, so you only have to
// check once.
type Request struct {
	c *RESTClient

	log         *zerolog.Logger
	rateLimiter flowcontrol.RateLimiter
	timeout     time.Duration

	authType AuthType
	user     *User
	token    string

	// generic components accessible via method setters
	method   string
	prePath  string
	subPath  string
	basePath string
	reqPath  string
	isAbs    bool
	cookies  []*http.Cookie
	headers  http.Header
	params   url.Values
	body     io.Reader

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

// Prefix adds segments to the relative beginning to the request path. These
// items will be placed before the optional Namespace, Resource, or Name sections.
// Setting AbsPath will clear any previously set Prefix segments
func (r *Request) Prefix(segments ...string) *Request {
	if r.err != nil {
		return r
	}
	r.prePath = path.Join(r.prePath, path.Join(segments...))
	return r
}

// Suffix appends segments to the end of the path. These items will be placed after the prefix and optional
// Namespace, Resource, or Name sections.
func (r *Request) Suffix(segments ...string) *Request {
	if r.err != nil {
		return r
	}
	r.subPath = path.Join(r.subPath, path.Join(segments...))
	return r
}

// 设置请求的URL, 会补充前缀与后缀
func (r *Request) URL(p string) *Request {
	u, err := url.Parse(p)
	if err != nil {
		r.err = err
		return r
	}

	r.reqPath = u.String()
	return r
}

// 设置请求的URL的绝对路径
func (r *Request) AbsURL(p string) *Request {
	r.URL(p)
	r.isAbs = true
	return r
}

func (r *Request) url() string {
	u, err := url.Parse(r.reqPath)
	if err != nil {
		r.err = err
		return ""
	}

	if !r.isAbs {
		if r.prePath != "" {
			u.Path = path.Join(r.prePath, u.Path)
		}
		if r.subPath != "" {
			u.Path = path.Join(u.Path, r.subPath)
		}
	}

	url, err := url.JoinPath(r.basePath, u.Path)
	if err != nil {
		r.err = err
		return ""
	}
	return url
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

// Param creates a query parameter with the given json obj to url.Values.
func (r *Request) ParamJson(obj any) *Request {
	if r.err != nil {
		return r
	}
	if r.params == nil {
		r.params = make(url.Values)
	}

	values, err := queryparams.Convert(obj)
	if err != nil {
		r.err = err
		return r
	}

	for k, v := range values {
		r.params[k] = append(r.params[k], v...)
	}

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
	// Trace
	var span trace.Span
	if r.c.tr != nil {
		ctx, span = r.c.tr.Start(ctx, r.url())
		defer span.End()

		ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
	}

	// 请求速率控制
	r.rateLimiter.Wait(1)

	// 请求响应对象
	resp := NewResponse(r.c)

	// 准备请求
	req, err := http.NewRequestWithContext(ctx, r.method, r.url(), r.body)
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
	r.log.Debug().Msgf("[%s] %s", req.Method, req.URL.String())
	r.log.Debug().Msgf("Request Headers:")
	for k, v := range req.Header {
		r.log.Debug().Msgf("%s=%s", k, strings.Join(v, ","))
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
