package exception

import "net/http"

const (
	// CODE_OTHER_PLACE_LGOIN 登录登录
	CODE_OTHER_PLACE_LGOIN = 50010
	// CODE_OTHER_IP_LOGIN 异常IP登录
	CODE_OTHER_IP_LOGIN = 50011
	// CODE_OTHER_CLIENT_LOGIN 用户已经通过其他端登录
	CODE_OTHER_CLIENT_LOGIN = 50012
	// CODE_SESSION_TERMINATED 会话中断
	CODE_SESSION_TERMINATED = 50013
	// CODE_ACESS_TOKEN_EXPIRED token过期
	CODE_ACESS_TOKEN_EXPIRED = 50014
	// CODE_REFRESH_TOKEN_EXPIRED token过期
	CODE_REFRESH_TOKEN_EXPIRED = 50015
	// CODE_ACCESS_TOKEN_ILLEGAL 访问token不合法
	CODE_ACCESS_TOKEN_ILLEGAL = 50016
	// CODE_REFRESH_TOKEN_ILLEGAL 刷新token不合法
	CODE_REFRESH_TOKEN_ILLEGAL = 50017
	// CODE_VERIFY_CODE_REQUIRED 需要验证码
	CODE_VERIFY_CODE_REQUIRED = 50018
	// CODE_PASSWORD_EXPIRED 用户密码过期
	CODE_PASSWORD_EXPIRED = 50019

	// CODE_UNAUTHORIZED 未认证
	CODE_UNAUTHORIZED = http.StatusUnauthorized
	// CODE_BAD_REQUEST 请求不合法
	CODE_BAD_REQUEST = http.StatusBadRequest
	// CODE_INTERNAL_SERVER_ERROR 服务端内部错误
	CODE_INTERNAL_SERVER_ERROR = http.StatusInternalServerError
	// CODE_FORBIDDEN 无权限
	CODE_FORBIDDEN = http.StatusForbidden
	// CODE_NOT_FOUND 接口未找到
	CODE_NOT_FOUND = http.StatusNotFound
	// CODE_CONFLICT 资源冲突, 已经存在
	CODE_CONFLICT = http.StatusConflict

	// CODE_UNKNOWN 未知异常
	CODE_UNKNOWN = 99999
)

var (
	reasonMap = map[int]string{
		CODE_UNAUTHORIZED:          "认证失败",
		CODE_NOT_FOUND:             "资源未找到",
		CODE_CONFLICT:              "资源已经存在",
		CODE_BAD_REQUEST:           "请求不合法",
		CODE_INTERNAL_SERVER_ERROR: "系统内部错误",
		CODE_FORBIDDEN:             "访问未授权",
		CODE_UNKNOWN:               "未知异常",
		CODE_ACCESS_TOKEN_ILLEGAL:  "访问令牌不合法",
		CODE_REFRESH_TOKEN_ILLEGAL: "刷新令牌不合法",
		CODE_OTHER_PLACE_LGOIN:     "异地登录",
		CODE_OTHER_IP_LOGIN:        "异常IP登录",
		CODE_OTHER_CLIENT_LOGIN:    "用户已经通过其他端登录",
		CODE_SESSION_TERMINATED:    "会话结束",
		CODE_ACESS_TOKEN_EXPIRED:   "访问过期, 请刷新",
		CODE_REFRESH_TOKEN_EXPIRED: "刷新过期, 请登录",
		CODE_VERIFY_CODE_REQUIRED:  "异常操作, 需要验证码进行二次确认",
		CODE_PASSWORD_EXPIRED:      "密码过期, 请找回密码或者联系管理员重置",
	}
)

func codeReason(code int) string {
	v, ok := reasonMap[code]
	if !ok {
		v = reasonMap[CODE_UNKNOWN]
	}

	return v
}
