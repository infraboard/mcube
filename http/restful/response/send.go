package response

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/exception"
	"github.com/infraboard/mcube/v2/http/response"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/log"
)

// Failed use to response error messge
func Failed(w *restful.Response, err error, opts ...response.Option) {
	var e *exception.APIException
	if v, ok := err.(*exception.APIException); ok {
		e = v
	} else {
		// 非可以预期, 没有定义业务的情况
		e = exception.NewAPIException(
			http.StatusInternalServerError,
			http.StatusText(http.StatusInternalServerError),
			"%s",
			err.Error(),
		)
		e.HttpCode = http.StatusInternalServerError
	}

	if e.Namespace == "" {
		e.WithNamespace(application.Get().AppName)
	}

	err = w.WriteHeaderAndEntity(e.HttpCode, e)
	if err != nil {
		log.L().Error().Msgf("send failed response error, %s", err)
	}
}

// Success use to response success data
func Success(w *restful.Response, data any, opts ...response.Option) {
	// 是否需要脱敏
	if v, ok := data.(response.DesenseObj); ok {
		v.Desense()
	}

	err := w.WriteEntity(data)
	if err != nil {
		log.L().Error().Msgf("send success response error, %s", err)
	}
}
