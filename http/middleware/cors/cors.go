package cors

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/infraboard/mcube/logger"
)

// Cors interface
type Cors interface {
	AllowedMethods() []string
	AllowedHeaders() []string
	AllowedOriginsAll() bool

	IsMethodAllowed(method string) bool
	Handler(h http.Handler) http.Handler
	HandlePreflight(w http.ResponseWriter, r *http.Request)
	HandleActualRequest(w http.ResponseWriter, r *http.Request)
}

// cors http handler
type cors struct {
	// Debug logger
	debug bool
	log   logger.FormatLogger
	// Normalized list of plain allowed origins
	allowedOrigins []string
	// List of allowed origins containing wildcards
	allowedWOrigins []wildcard
	// Optional origin validator function
	allowOriginFunc func(origin string) bool
	// Optional origin validator (with request) function
	allowOriginRequestFunc func(r *http.Request, origin string) bool
	// Normalized list of allowed headers
	allowedHeaders []string
	// Normalized list of allowed methods
	allowedMethods []string
	// Normalized list of exposed headers
	exposedHeaders []string
	maxAge         int
	// Set to true when allowed origins contains a "*"
	allowedOriginsAll bool
	// Set to true when allowed headers contains a "*"
	allowedHeadersAll bool
	allowCredentials  bool
	optionPassthrough bool
}

func (c *cors) AllowedMethods() []string {
	return c.allowedMethods
}

func (c *cors) AllowedHeaders() []string {
	return c.allowedHeaders
}

func (c *cors) AllowedOriginsAll() bool {
	return c.allowedOriginsAll
}

// Default creates a new Cors handler with default options.
func Default() Cors {
	return New(Options{})
}

// AllowAll create a new Cors handler with permissive configuration allowing all
// origins with all standard methods with any header and credentials.
func AllowAll() Cors {
	return New(Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	})
}

// New creates a new Cors handler with the provided options.
func New(options Options) Cors {
	c := &cors{
		exposedHeaders:         convert(options.ExposedHeaders, http.CanonicalHeaderKey),
		allowOriginFunc:        options.AllowOriginFunc,
		allowOriginRequestFunc: options.AllowOriginRequestFunc,
		allowCredentials:       options.AllowCredentials,
		maxAge:                 options.MaxAge,
		optionPassthrough:      options.OptionsPassthrough,
		debug:                  options.Debug,
	}

	// Normalize options
	// Note: for origins and methods matching, the spec requires a case-sensitive matching.
	// As it may error prone, we chose to ignore the spec here.

	// Allowed Origins
	if len(options.AllowedOrigins) == 0 {
		if options.AllowOriginFunc == nil && options.AllowOriginRequestFunc == nil {
			// Default is all origins
			c.allowedOriginsAll = true
		}
	} else {
		c.allowedOrigins = []string{}
		c.allowedWOrigins = []wildcard{}
		for _, origin := range options.AllowedOrigins {
			// Normalize
			origin = strings.ToLower(origin)
			if origin == "*" {
				// If "*" is present in the list, turn the whole list into a match all
				c.allowedOriginsAll = true
				c.allowedOrigins = nil
				c.allowedWOrigins = nil
				break
			} else if i := strings.IndexByte(origin, '*'); i >= 0 {
				// Split the origin in two: start and end string without the *
				w := wildcard{origin[0:i], origin[i+1:]}
				c.allowedWOrigins = append(c.allowedWOrigins, w)
			} else {
				c.allowedOrigins = append(c.allowedOrigins, origin)
			}
		}
	}

	// Allowed Headers
	if len(options.AllowedHeaders) == 0 {
		// Use sensible defaults
		c.allowedHeaders = []string{"Origin", "Accept", "Content-Type", "X-Requested-With"}
	} else {
		// Origin is always appended as some browsers will always request for this header at preflight
		c.allowedHeaders = convert(append(options.AllowedHeaders, "Origin"), http.CanonicalHeaderKey)
		for _, h := range options.AllowedHeaders {
			if h == "*" {
				c.allowedHeadersAll = true
				c.allowedHeaders = nil
				break
			}
		}
	}

	// Allowed Methods
	if len(options.AllowedMethods) == 0 {
		// Default is spec's "simple" methods
		c.allowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead}
	} else {
		c.allowedMethods = convert(options.AllowedMethods, strings.ToUpper)
	}

	return c
}

// Handler apply the CORS specification on the request, and add relevant CORS headers
// as necessary.
func (c *cors) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
			c.logf("Handler: Preflight request")
			c.HandlePreflight(w, r)
			// Preflight requests are standalone and should stop the chain as some other
			// middleware may not handle OPTIONS requests correctly. One typical example
			// is authentication middleware ; OPTIONS requests won't carry authentication
			// headers (see #1)
			if c.optionPassthrough {
				h.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			c.logf("Handler: Actual request")
			c.HandleActualRequest(w, r)
			h.ServeHTTP(w, r)
		}
	})
}

// HandlePreflight handles pre-flight CORS requests
func (c *cors) HandlePreflight(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	if r.Method != http.MethodOptions {
		c.logf("  Preflight aborted: %s!=OPTIONS", r.Method)
		return
	}
	// Always set Vary headers
	// see https://github.com/rs/cors/issues/10,
	//     https://github.com/rs/cors/commit/dbdca4d95feaa7511a46e6f1efb3b3aa505bc43f#commitcomment-12352001
	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")

	if origin == "" {
		c.logf("  Preflight aborted: empty origin")
		return
	}
	if !c.isOriginAllowed(r, origin) {
		c.logf("  Preflight aborted: origin '%s' not allowed", origin)
		return
	}

	reqMethod := r.Header.Get("Access-Control-Request-Method")
	if !c.IsMethodAllowed(reqMethod) {
		c.logf("  Preflight aborted: method '%s' not allowed", reqMethod)
		return
	}
	reqHeaders := parseHeaderList(r.Header.Get("Access-Control-Request-Headers"))
	if !c.areHeadersAllowed(reqHeaders) {
		c.logf("  Preflight aborted: headers '%v' not allowed", reqHeaders)
		return
	}
	if c.allowedOriginsAll {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Set("Access-Control-Allow-Origin", origin)
	}
	// Spec says: Since the list of methods can be unbounded, simply returning the method indicated
	// by Access-Control-Request-Method (if supported) can be enough
	headers.Set("Access-Control-Allow-Methods", strings.ToUpper(reqMethod))
	if len(reqHeaders) > 0 {

		// Spec says: Since the list of headers can be unbounded, simply returning supported headers
		// from Access-Control-Request-Headers can be enough
		headers.Set("Access-Control-Allow-Headers", strings.Join(reqHeaders, ", "))
	}
	if c.allowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	if c.maxAge > 0 {
		headers.Set("Access-Control-Max-Age", strconv.Itoa(c.maxAge))
	}
	c.logf("  Preflight response headers: %v", headers)
}

// HandleActualRequest handles simple cross-origin requests, actual request or redirects
func (c *cors) HandleActualRequest(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()
	origin := r.Header.Get("Origin")

	// Always set Vary, see https://github.com/rs/cors/issues/10
	headers.Add("Vary", "Origin")
	if origin == "" {
		c.logf("  Actual request no headers added: missing origin")
		return
	}
	if !c.isOriginAllowed(r, origin) {
		c.logf("  Actual request no headers added: origin '%s' not allowed", origin)
		return
	}

	// Note that spec does define a way to specifically disallow a simple method like GET or
	// POST. Access-Control-Allow-Methods is only used for pre-flight requests and the
	// spec doesn't instruct to check the allowed methods for simple cross-origin requests.
	// We think it's a nice feature to be able to have control on those methods though.
	if !c.IsMethodAllowed(r.Method) {
		c.logf("  Actual request no headers added: method '%s' not allowed", r.Method)

		return
	}
	if c.allowedOriginsAll {
		headers.Set("Access-Control-Allow-Origin", "*")
	} else {
		headers.Set("Access-Control-Allow-Origin", origin)
	}
	if len(c.exposedHeaders) > 0 {
		headers.Set("Access-Control-Expose-Headers", strings.Join(c.exposedHeaders, ", "))
	}
	if c.allowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
	c.logf("  Actual response added headers: %v", headers)
}

// convenience method. checks if a logger is set.
func (c *cors) logf(format string, a ...interface{}) {
	if c.debug && c.log != nil {
		c.log.Debugf(format, a...)
		return
	}

	if c.debug {
		log.Printf(format, a...)
	}
}

// isOriginAllowed checks if a given origin is allowed to perform cross-domain requests
// on the endpoint
func (c *cors) isOriginAllowed(r *http.Request, origin string) bool {
	if c.allowOriginRequestFunc != nil {
		return c.allowOriginRequestFunc(r, origin)
	}
	if c.allowOriginFunc != nil {
		return c.allowOriginFunc(origin)
	}
	if c.allowedOriginsAll {
		return true
	}
	origin = strings.ToLower(origin)
	for _, o := range c.allowedOrigins {
		if o == origin {
			return true
		}
	}
	for _, w := range c.allowedWOrigins {
		if w.match(origin) {
			return true
		}
	}
	return false
}

// IsMethodAllowed checks if a given method can be used as part of a cross-domain request
// on the endpoing
func (c *cors) IsMethodAllowed(method string) bool {
	if len(c.allowedMethods) == 0 {
		// If no method allowed, always return false, even for preflight request
		return false
	}
	method = strings.ToUpper(method)
	if method == http.MethodOptions {
		// Always allow preflight requests
		return true
	}
	for _, m := range c.allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

// areHeadersAllowed checks if a given list of headers are allowed to used within
// a cross-domain request.
func (c *cors) areHeadersAllowed(requestedHeaders []string) bool {
	if c.allowedHeadersAll || len(requestedHeaders) == 0 {
		return true
	}
	for _, header := range requestedHeaders {
		header = http.CanonicalHeaderKey(header)
		found := false
		for _, h := range c.allowedHeaders {
			if h == header {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
