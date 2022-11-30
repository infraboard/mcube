package response

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/exception"
	"github.com/infraboard/mcube/http/response"
	"github.com/infraboard/mcube/logger/zap"
)

// Failed use to response error messge
func Failed(w *restful.Response, err error, opts ...response.Option) {
	var (
		errCode  int
		httpCode int
		ns       string
		reason   string
		data     interface{}
		meta     interface{}
	)

	switch t := err.(type) {
	case exception.APIException:
		errCode = t.ErrorCode()
		reason = t.Reason()
		data = t.Data()
		meta = t.Meta()
		ns = t.Namespace()
	default:
		errCode = exception.UnKnownException
	}

	// 映射http status code 1xx - 5xx
	// 如果为其他errCode, 统一成200
	if errCode/100 >= 1 && errCode/100 <= 5 {
		httpCode = errCode
	} else {
		httpCode = http.StatusOK
	}

	resp := response.Data{
		Code:      &errCode,
		Namespace: ns,
		Reason:    reason,
		Message:   err.Error(),
		Data:      data,
		Meta:      meta,
	}

	for _, opt := range opts {
		opt.Apply(&resp)
	}

	err = w.WriteHeaderAndEntity(httpCode, resp)
	if err != nil {
		zap.L().Errorf("send failed response error, %s", err)
	}
}

// Success use to response success data
func Success(w *restful.Response, data interface{}, opts ...response.Option) {
	err := w.WriteEntity(data)

	if err != nil {
		zap.L().Errorf("send success response error, %s", err)
	}
}
