package rest

import (
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/infraboard/mcube/client/compressor"
	"github.com/infraboard/mcube/client/negotiator"
)

func NewResponse() *Response {
	return &Response{}
}

// Response contains the result of calling Request.Do().
type Response struct {
	body       io.ReadCloser
	headers    http.Header
	statusCode int
	err        error
	bf         []byte
	isRead     bool

	sync.Mutex
}

func (r *Response) read() {
	r.Lock()
	defer r.Unlock()

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

	r.bf = body
}

func (r *Response) Into(v any) error {
	if r.err != nil {
		return r.err
	}

	// 读取body里面的数据
	r.read()

	// 判断status code
	if r.statusCode/100 != 2 {
		return fmt.Errorf("status code is %d, not 2xx, response: %s", r.statusCode, string(r.bf))
	}

	// 解析数据
	ct := HeaderFilterFlags(r.headers.Get(CONTENT_TYPE_HEADER))
	nt := negotiator.GetNegotiator(ct)

	if err := nt.Decode(r.bf, v); err != nil {
		return fmt.Errorf("decode err: %s, data: %s", err, string(r.bf))
	}
	return nil
}
