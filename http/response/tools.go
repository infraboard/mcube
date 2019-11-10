package response

import (
	"encoding/json"
	"net/http"

	"github.com/infraboard/mcube/exception"
)

// Data to be used by controllers.
type Data struct {
	Code    *int        `json:"code"`              // 自定义返回码  0:表示正常
	Type    string      `json:"type,omitempty"`    // 数据类型, 可以缺省
	Message string      `json:"message,omitempty"` // 关于这次响应的说明信息
	Data    interface{} `json:"data,omitempty"`    // 返回的具体数据
}

// PageData 数据分页数据
type PageData struct {
	PageSize   int64       `json:"page_size"`   // 总共多少页
	TotalCount int64       `json:"total_count"` // 总共多少条
	PageIndex  int64       `json:"page_index"`  // 当前页
	List       interface{} `json:"list"`        // 页面数据
}

// Failed use to response error messge
func Failed(w http.ResponseWriter, err error) {
	var (
		errCode  int
		httpCode int
	)

	switch t := err.(type) {
	case exception.APIException:
		errCode = t.ErrorCode()
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

	resp := Data{
		Code:    &errCode,
		Message: err.Error(),
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
	resp := Data{
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
