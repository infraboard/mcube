package response

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/desense"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/log"
)

// Failed use to response error messge
func Failed(w *restful.Response, err error, opts ...response.Option) {
	var e *exception.ApiException
	if v, ok := err.(*exception.ApiException); ok {
		e = v
	} else {
		// 非可以预期, 没有定义业务的情况
		e = exception.NewApiException(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
		).WithMessage(err.Error())
		e.HttpCode = http.StatusInternalServerError
	}

	if e.Service == "" {
		e.WithNamespace(application.Get().AppName)
	}

	statusCode := e.HttpCode
	e.HttpCode = 0
	err = w.WriteHeaderAndEntity(statusCode, e)
	if err != nil {
		log.L().Error().Msgf("send failed response error, %s", err)
	}
}

// Success use to response success data
func Success(w *restful.Response, data any, opts ...response.Option) {
	// 脱敏
	if err := desense.MaskStruct(data); err != nil {
		log.L().Error().Msgf("desense error, %s", err)
	}

	err := w.WriteEntity(data)
	if err != nil {
		log.L().Error().Msgf("send success response error, %s", err)
	}
}
