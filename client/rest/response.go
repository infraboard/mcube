package rest

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/infraboard/mcube/client/compressor"
	"github.com/infraboard/mcube/client/negotiator"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
)

func NewResponse(c *RESTClient) *Response {
	return &Response{
		log: logger.Sub("http.response"),
	}
}

// Response contains the result of calling Request.Do().
type Response struct {
	body        io.ReadCloser
	headers     http.Header
	statusCode  int
	err         error
	bf          []byte
	contentType string
	isRead      bool

	log  *zerolog.Logger
	lock sync.Mutex
}

func (r *Response) withBody(body io.ReadCloser) {
	r.body = body
}

func (r *Response) withHeader(headers http.Header) {
	r.headers = headers
}

func (r *Response) withStatusCode(code int) {
	r.statusCode = code
}

func (r *Response) readBody() {
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
	r.log.Debug().Msgf("Status Code: %d", r.statusCode)

	r.log.Debug().Msgf("Response Headers:")
	for k, v := range r.headers {
		r.log.Debug().Msgf("%s=%s", k, strings.Join(v, ","))
	}
	r.log.Debug().Msgf("Body: %s", string(body))
}

func (r *Response) Header(header string, v *string) *Response {
	*v = r.headers.Get(header)
	return r
}

// 请求正常的情况下, 获取返回的数据, 不做解析
func (r *Response) Raw() ([]byte, error) {
	if err := r.Error(); err != nil {
		return nil, err
	}
	return r.bf, r.err
}

// 直接返回stream, 常用于websocket
func (r *Response) Stream() (io.ReadCloser, error) {
	return r.body, r.err
}

// 默认不设置通过Content-Type获取, 如果设定, 以设定为准
func (r *Response) ContentType(m negotiator.MIME) *Response {
	r.contentType = string(m)
	return r
}

// 请求正常的情况下, 获取返回的数据, 会根据Content-Type做解析
func (r *Response) Into(v any) error {
	if err := r.Error(); err != nil {
		return err
	}

	// 解析数据
	if r.contentType == "" {
		r.contentType = HeaderFilterFlags(r.headers.Get(CONTENT_TYPE_HEADER))
	}

	nt := negotiator.GetNegotiator(r.contentType)
	return nt.Decode(r.bf, v)
}

// 不处理返回, 直接判断请求是否正常
func (r *Response) Error() error {
	r.readBody()

	// 判断status code
	if r.statusCode/100 != 2 {
		r.err = fmt.Errorf("status code is %d, not 2xx, response: %s", r.statusCode, string(r.bf))
	}

	return r.err
}
