package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Http{
	Host:                    "127.0.0.1",
	Port:                    8080,
	PathPrefix:              "api",
	ReadHeaderTimeoutSecond: 30,
	ReadTimeoutSecond:       60,
	WriteTimeoutSecond:      60,
	IdleTimeoutSecond:       600,
	MaxHeaderSize:           "16kb",
}

type Http struct {
	ioc.ObjectImpl

	// 是否开启HTTP Server, 默认会根据是否有注册得有API对象来自动开启
	Enable *bool `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	// HTTP服务Host
	Host string `json:"host" yaml:"host" toml:"host" env:"HOST"`
	// HTTP服务端口
	Port int `json:"port" yaml:"port" toml:"port" env:"PORT"`
	// API接口前缀
	PathPrefix string `json:"path_prefix" yaml:"path_prefix" toml:"path_prefix" env:"PATH_PREFIX"`

	// HTTP服务器参数
	// HTTP Header读取超时时间
	ReadHeaderTimeoutSecond int `json:"read_header_timeout" yaml:"read_header_timeout" toml:"read_header_timeout" env:"READ_HEADER_TIMEOUT"`
	// 读取HTTP整个请求时的参数
	ReadTimeoutSecond int `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout" env:"READ_TIMEOUT"`
	// 响应超时时间
	WriteTimeoutSecond int `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout" env:"WRITE_TIMEOUT"`
	// 启用了KeepAlive时 复用TCP链接的超时时间
	IdleTimeoutSecond int `json:"idle_timeout" yaml:"idle_timeout" toml:"idle_timeout" env:"IDLE_TIMEOUT"`
	// header最大大小
	MaxHeaderSize string `json:"max_header_size" yaml:"max_header_size" toml:"max_header_size" env:"MAX_HEADER_SIZE"`

	// 解析后的数据
	maxHeaderBytes uint64
	log            *zerolog.Logger
	router         http.Handler
	server         *http.Server
}

func (h *Http) HTTPPrefix() string {
	u, err := url.JoinPath("/", h.PathPrefix, application.Get().AppName)
	if err != nil {
		return fmt.Sprintf("/%s/%s", application.Get().AppName, h.PathPrefix)
	}
	return u
}

func (h *Http) ApiObjectPathPrefix(obj ioc.Object) string {
	cp := obj.Meta().CustomPathPrefix
	if cp != "" {
		return cp
	}

	// 使用对象反射名称
	objName := obj.Name()
	if objName == "" {
		objName = strings.ToLower(strings.TrimLeft(reflect.TypeOf(obj).String(), "*"))
	}

	return fmt.Sprintf("%s/%s/%s",
		h.HTTPPrefix(),
		obj.Version(),
		objName)
}

func (h *Http) ApiObjectAddr(obj ioc.Object) string {
	return fmt.Sprintf("http://%s%s", h.Addr(), h.ApiObjectPathPrefix(obj))
}

func (h *Http) Name() string {
	return AppName
}

func (i *Http) Priority() int {
	return -99
}

// 配置数据解析
func (h *Http) Init() error {
	h.log = log.Sub("http")

	mhz, err := humanize.ParseBytes(h.MaxHeaderSize)
	if err != nil {
		return err
	}
	h.maxHeaderBytes = mhz

	h.server = &http.Server{
		ReadHeaderTimeout: time.Duration(h.ReadHeaderTimeoutSecond) * time.Second,
		ReadTimeout:       time.Duration(h.ReadTimeoutSecond) * time.Second,
		WriteTimeout:      time.Duration(h.WriteTimeoutSecond) * time.Second,
		IdleTimeout:       time.Duration(h.IdleTimeoutSecond) * time.Second,
		MaxHeaderBytes:    int(h.maxHeaderBytes),
		Addr:              h.Addr(),
		Handler:           h.router,
	}
	return nil
}

func (h *Http) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

func (h *Http) SetRouter(r http.Handler) {
	h.router = r
}

func (h *Http) IsEnable() bool {
	if h.Enable == nil {
		return h.router != nil
	}

	return *h.Enable
}

// Start 启动服务
func (h *Http) Start(ctx context.Context) {
	h.log.Info().Msgf("HTTP服务启动成功, 监听地址: %s", h.Addr())
	if err := h.server.ListenAndServe(); err != nil {
		h.log.Error().Msg(err.Error())
	}
}

// Stop 停止server
func (h *Http) Stop(ctx context.Context) error {
	h.log.Info().Msg("start graceful shutdown")
	// 优雅关闭HTTP服务
	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http graceful shutdown timeout, force exit")
	}
	return nil
}

func (a *Http) SwagerDocs(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       application.Get().AppName,
			Description: application.Get().AppDescription,
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT",
					URL:  "http://mit.org",
				},
			},
			Version: application.Short(),
		},
	}
	swo.SecurityDefinitions = map[string]*spec.SecurityScheme{
		// can add more authentication methods here
		"jwt": spec.APIKeyAuth("Authorization", "header"),
	}
	// map routes to security definitions
	a.enrichSwaggerObjectSecurity(swo)
}

const (
	SecurityDefinitionKey = "OAPI_SECURITY_DEFINITION"
)

type OAISecurity struct {
	Name   string   // SecurityDefinition name
	Scopes []string // Scopes for oauth2
}

func (s *OAISecurity) Valid() error {
	switch s.Name {
	case "oauth2":
		return nil
	case "openIdConnect":
		return nil
	default:
		if len(s.Scopes) > 0 {
			return fmt.Errorf("oai Security scopes for scheme '%s' should be empty", s.Name)
		}
	}

	return nil
}

func (h *Http) enrichSwaggerObjectSecurity(swo *spec.Swagger) {

	// loop through all registered web services
	for _, ws := range restful.RegisteredWebServices() {
		for _, route := range ws.Routes() {

			// grab route metadata for a SecurityDefinition
			secdefn, ok := route.Metadata[SecurityDefinitionKey]
			if !ok {
				continue
			}

			// grab pechelper.OAISecurity from the stored interface{}
			var sEntry OAISecurity
			switch v := secdefn.(type) {
			case *OAISecurity:
				sEntry = *v
			case OAISecurity:
				sEntry = v
			default:
				// not valid type
				h.log.Error().Msgf("skipping Security openapi spec for %s:%s, invalid metadata type %v", route.Method, route.Path, v)
				continue
			}

			if _, ok := swo.SecurityDefinitions[sEntry.Name]; !ok {
				h.log.Error().Msgf("skipping Security openapi spec for %s:%s, '%s' not found in SecurityDefinitions", route.Method, route.Path, sEntry.Name)
				continue
			}

			// grab path and path item in openapi spec, need to sanitized rote.Path because swo.Paths have been sanitized
			sanitizedPath, _ := sanitizePath(route.Path)
			path, err := swo.Paths.JSONLookup(sanitizedPath)
			if err != nil {
				h.log.Error().Msgf("skipping Security openapi spec for %s:%s, %s", route.Method, route.Path, err.Error())
				continue
			}
			pItem := path.(*spec.PathItem)

			// Update respective path Option based on method
			var pOption *spec.Operation
			switch method := strings.ToLower(route.Method); method {
			case "get":
				pOption = pItem.Get
			case "post":
				pOption = pItem.Post
			case "patch":
				pOption = pItem.Patch
			case "delete":
				pOption = pItem.Delete
			case "put":
				pOption = pItem.Put
			case "head":
				pOption = pItem.Head
			case "options":
				pOption = pItem.Options
			default:
				// unsupported method
				h.log.Error().Msgf("skipping Security openapi spec for %s:%s, unsupported method '%s'", route.Method, route.Path, route.Method)
				continue
			}

			// update the pOption with security entry
			pOption.SecuredWith(sEntry.Name, sEntry.Scopes...)
		}
	}
}

func sanitizePath(restfulPath string) (string, map[string]string) {
	openapiPath := ""
	patterns := map[string]string{}
	for _, fragment := range strings.Split(restfulPath, "/") {
		if fragment == "" {
			continue
		}
		if strings.HasPrefix(fragment, "{") && strings.Contains(fragment, ":") {
			split := strings.Split(fragment, ":")
			// skip google custom method like `resource/{resource-id}:customVerb`
			if !strings.Contains(split[0], "}") {
				fragment = split[0][1:]
				pattern := split[1][:len(split[1])-1]
				if pattern == "*" { // special case
					pattern = ".*"
				}
				patterns[fragment] = pattern
				fragment = "{" + fragment + "}"
			}
		}
		openapiPath += "/" + fragment
	}
	return openapiPath, patterns
}
