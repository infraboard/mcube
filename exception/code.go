package exception

import "net/http"

const (
	// AccessTokenIllegal 访问token不合法
	AccessTokenIllegal = 50008
	// RefreshTokenIllegal 刷新token不合法
	RefreshTokenIllegal = 50009
	// OtherClientsLoggedIn 启动端登录
	OtherClientsLoggedIn = 50012
	// AccessTokenExpired token过期
	AccessTokenExpired = 50014
	// RefreshTokenExpired token过期
	RefreshTokenExpired = 50015

	// 1xx - 5xx copy from http status code
	Unauthorized        = http.StatusUnauthorized
	BadRequest          = http.StatusBadRequest
	InternalServerError = http.StatusInternalServerError
	Forbidden           = http.StatusForbidden
	NotFound            = http.StatusNotFound

	UnKnownException = 99999
)

var (
	reasonMap = map[int]string{
		Unauthorized:         "认证失败",
		NotFound:             "资源未找到",
		BadRequest:           "请求不合法",
		InternalServerError:  "系统内部错误",
		Forbidden:            "访问未授权",
		UnKnownException:     "未知异常",
		AccessTokenIllegal:   "访问令牌不合法",
		RefreshTokenIllegal:  "刷新令牌不合法",
		OtherClientsLoggedIn: "用户已经通过其他端登录",
		AccessTokenExpired:   "访问过期, 请刷新",
		RefreshTokenExpired:  "刷新过期, 请登录",
	}
)

func codeReason(code int) string {
	v, ok := reasonMap[code]
	if !ok {
		v = reasonMap[UnKnownException]
	}

	return v
}
