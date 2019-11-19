package realip

import (
	"net/http"
	"strings"

	"github.com/infraboard/mcube/http/router"
)

var (
	// DefaultForwareHeaderKey 协商forward ip 的 hander key名称
	defaultScanForwareHeaderKey = []string{"X-Forwarded-For", "X-Real-IP"}
)

// Options 参数
type Options struct {
	forwareHeaderKeys []string
}

// NewDefault 初始化实例
func NewDefault() router.Middleware {
	return New(&Options{
		forwareHeaderKeys: defaultScanForwareHeaderKey,
	})
}

// New 初始化实例
func New(options *Options) router.Middleware {
	r := new(realip)
	r.addScanKey(options.forwareHeaderKeys)
	return r
}

type realip struct {
	scanKeys []string
}

func (r *realip) addScanKey(ks []string) {
	for i := range ks {
		r.scanKeys = append(r.scanKeys, http.CanonicalHeaderKey(ks[i]))
	}
}

// RealIP is a middleware that sets a http.Request's RemoteAddr to the results
// of parsing either the X-Forwarded-For header or the X-Real-IP header (in that
// order).
//
// This middleware should be inserted fairly early in the middleware stack to
// ensure that subsequent layers (e.g., request loggers) which examine the
// RemoteAddr will see the intended value.
//
// You should only use this middleware if you can trust the headers passed to
// you (in particular, the two headers this middleware uses), for example
// because you have placed a reverse proxy like HAProxy or nginx in front of
// chi. If your reverse proxies are configured to pass along arbitrary header
// values from the client, or if you use this middleware without a reverse
// proxy, malicious clients will be able to make you very sad (or, depending on
// how you're using RemoteAddr, vulnerable to an attack of some sort).
func (r *realip) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if rip := r.realIP(req); rip != "" {
			req.RemoteAddr = rip
		}

		next.ServeHTTP(rw, req)
	})
}

func (r *realip) realIP(req *http.Request) string {
	var ip string

	for _, key := range r.scanKeys {
		value := req.Header.Get(key)

		if strings.Contains(value, ", ") {
			i := strings.Index(value, ", ")
			if i == -1 {
				i = len(value)
			}

			ip = value[:i]
			break
		}

		if value != "" {
			ip = value
			break
		}
	}

	return ip
}
