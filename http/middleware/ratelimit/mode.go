package ratelimit

const (
	// GlobalMode 限制所有的请求
	GlobalMode Mode = iota
	// RemoteIPMode 限制远程IP的速率
	RemoteIPMode
	// HeaderKeyMode 根据Header中的特定Key进行限制
	HeaderKeyMode
	// CookieKeyMode 更加cookie中特定的Key进行限制
	CookieKeyMode
)

// Mode 限制模式
type Mode uint
