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

	// Unauthorized 未认证
	Unauthorized = http.StatusUnauthorized
	// BadRequest 请求不合法
	BadRequest = http.StatusBadRequest
	// InternalServerError 服务端内部错误
	InternalServerError = http.StatusInternalServerError
	// Forbidden 无权限
	Forbidden = http.StatusForbidden
	// NotFound 接口未找到
	NotFound = http.StatusNotFound

	// UnKnownException 未知异常
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
