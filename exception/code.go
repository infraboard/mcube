package exception

import "net/http"

const (
	// AccessTokenExpired token过期
	AccessTokenExpired = 1000
	// RefreshTokenExpired token过期
	RefreshTokenExpired = 1001

	// 1xx - 5xx copy from http status code
	Unauthorized        = http.StatusUnauthorized
	BadRequest          = http.StatusBadRequest
	InternalServerError = http.StatusInternalServerError
	Forbidden           = http.StatusForbidden
	NotFound            = http.StatusNotFound

	UnKnownException = 9999
)

var (
	reasonMap = map[int]string{
		Unauthorized:        "认证失败",
		NotFound:            "资源未找到",
		BadRequest:          "请求不合法",
		InternalServerError: "系统内部错误",
		Forbidden:           "访问未授权",
		UnKnownException:    "未知异常",
		AccessTokenExpired:  "访问过期, 请刷新",
		RefreshTokenExpired: "刷新过期, 请登录",
	}
)

func codeReason(code int) string {
	v, ok := reasonMap[code]
	if !ok {
		v = reasonMap[UnKnownException]
	}

	return v
}
