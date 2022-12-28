package rest

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/infraboard/mcube/client/compressor"
	"github.com/infraboard/mcube/client/negotiator"
	"github.com/infraboard/mcube/logger"
)

func NewResponse(c *RESTClient) *Response {
	return &Response{
		log: c.log.Named("response"),
	}
}

// Response contains the result of calling Request.Do().
type Response struct {
	body       io.ReadCloser
	headers    http.Header
	statusCode int
	err        error
	bf         []byte
	isRead     bool

	log  logger.Logger
	lock sync.Mutex
}

func (r *Response) withBody(body io.ReadCloser) {
	r.body = body
	r.read()
}

func (r *Response) withHeader(headers http.Header) {
	r.headers = headers
}

func (r *Response) withStatusCode(code int) {
	r.statusCode = code
}

func (r *Response) read() {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.body == nil || r.isRead {
		return
	}

	r.isRead = true
	defer r.body.Close()

	var bodyReader io.Reader = r.body

	// 解压缩
	et := HeaderFilterFlags(r.headers.Get(CONTENT_ENCODING_HEADER))
	if et != "" {
		cp := compressor.GetCompressor(et)
		dp, err := cp.Decompress(r.body)
		if err != nil {
			r.err = err
			return
		}
		bodyReader = dp
	}

	// 读取数据
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		r.err = err
		return
	}
	r.debug(body)

	r.bf = body
}

func (r *Response) debug(body []byte) {
	r.log.Debugf("Status Code: %d", r.statusCode)

	r.log.Debugf("Response Headers:")
	for k, v := range r.headers {
		r.log.Debugf("%s=%s", k, strings.Join(v, ","))
	}
	r.log.Debugf("Body: %s", string(body))
}

func (r *Response) Into(v any) error {
	if err := r.Error(); err != nil {
		return err
	}

	// 解析数据
	ct := HeaderFilterFlags(r.headers.Get(CONTENT_TYPE_HEADER))
	nt := negotiator.GetNegotiator(ct)
	return nt.Decode(r.bf, v)
}

func (r *Response) Error() error {
	// 判断status code
	if r.statusCode/100 != 2 {
		r.err = fmt.Errorf("status code is %d, not 2xx, response: %s", r.statusCode, string(r.bf))
	}

	return r.err
}
