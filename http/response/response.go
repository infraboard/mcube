package response

import (
	"encoding/json"
	"net/http"

	"github.com/infraboard/mcube/exception"
)

// Response to be used by controllers.
type Response struct {
	Code       *int        `json:"code"`              // 自定义返回码 0, 表示正常
	Type       string      `json:"type,omitempty"`    // 数据类型, 可以缺省
	Message    string      `json:"message,omitempty"` // 关于这次响应的说明信息
	Data       interface{} `json:"data,omitempty"`    // 返回的具体数据
	StatusCode int         `json:"statusCode,omitempty"`
}

// PageData 数据分页数据
type PageData struct {
	PageSize   int64       `json:"page_size"`       // 总共多少页
	TotalCount int64       `json:"total_count"`     // 总共多少条
	PageIndex  int64       `json:"page_index"`      // 当前页
	Start      int         `json:"start,omitempty"` // 开始位置
	End        int         `json:"end,omitempty"`   // 结束位置
	List       interface{} `json:"list"`            // 页面数据
}

// Failed use to response error messge
func Failed(w http.ResponseWriter, err error) {
	msg := err.Error()
	customCode := 0
	httpCode := http.StatusBadRequest

	switch t := err.(type) {
	case exception.APIException:
		customCode = t.ErrorCode()
	default:
		customCode = exception.UnKnownException
	}

	// 映射http status code 1xx - 5xx
	if customCode/100 >= 1 && customCode/100 <= 5 {
		httpCode = customCode
	}

	resp := Response{
		Code:    &customCode,
		Message: msg,
	}

	// set response heanders
	w.Header().Set("Content-Type", "application/json")

	// if marshal json error, use string to response
	respByt, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error", "message": "encoding to json error"}`))
		return
	}

	w.WriteHeader(httpCode)
	w.Write(respByt)
	return
}

// Success use to response success data
func Success(w http.ResponseWriter, code int, data interface{}) {
	c := 0
	resp := Response{
		Code:    &c,
		Message: "",
		Data:    data,
	}

	// set response heanders
	w.Header().Set("Content-Type", "application/json")

	respByt, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error", "message": "encoding to json error"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(respByt)
	return
}
