package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/logger"
	"github.com/rs/zerolog"
)

func NewDefaultHttp() *Http {
	return &Http{
		Host:                    "127.0.0.1",
		Port:                    8080,
		PathPrefix:              "api",
		ReadHeaderTimeoutSecond: 30,
		ReadTimeoutSecond:       60,
		WriteTimeoutSecond:      60,
		IdleTimeoutSecond:       300,
		MaxHeaderSize:           "16kb",
		EnableTrace:             true,
		HealthCheck: HealthCheck{
			Enabled: true,
		},
		Cors: CORS{
			Enabled:        false,
			AllowedHeaders: []string{"*"},
			AllowedDomains: []string{"*"},
			AllowedMethods: []string{"HEAD", "OPTIONS", "GET", "POST", "PUT", "PATCH", "DELETE"},
		},
		ApiDoc: ApiDoc{
			Enabled: true,
			DocPath: "apidocs.jso",
		},
		WEB_FRAMEWORK: WEB_FRAMEWORK_GIN,
		routerBuilders: map[WEB_FRAMEWORK]RouterBuilder{
			WEB_FRAMEWORK_GO_RESTFUL: NewGoRestfulRouterBuilder(),
			WEB_FRAMEWORK_GIN:        NewGinRouterBuilder(),
		},
		RouterBuildConfig: &BuildConfig{},
	}
}

type Http struct {
	// 默认根据
	Enable *bool `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	// HTTP服务Host
	Host string `json:"size" yaml:"size" toml:"size" env:"HOST"`
	// HTTP服务端口
	Port int `json:"port" yaml:"port" toml:"port" env:"PORT"`
	// 接口前缀
	PathPrefix string `json:"path_prefix" yaml:"path_prefix" toml:"path_prefix" env:"PATH_PREFIX"`

	// 使用的http框架, 启用后会自动从ioc中加载 该框架的hanlder
	WEB_FRAMEWORK WEB_FRAMEWORK `json:"web_framework" yaml:"web_framework" toml:"web_framework" env:"WEB_FRAMEWORK"`

	// HTTP服务器参数
	ReadHeaderTimeoutSecond int `json:"read_header_timeout" yaml:"read_header_timeout" toml:"read_header_timeout" env:"READ_HEADER_TIMEOUT"`
	// 读取HTTP整个请求时的参数
	ReadTimeoutSecond int `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout" env:"READ_TIMEOUT"`
	// 响应超时事件
	WriteTimeoutSecond int `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout" env:"WRITE_TIMEOUT"`
	// 启用了KeepAlive时 复用TCP链接的超时时间
	IdleTimeoutSecond int `json:"idle_timeout" yaml:"idle_timeout" toml:"idle_timeout" env:"IDLE_TIMEOUT"`
	// header最大大小
	MaxHeaderSize string `json:"max_header_size" yaml:"max_header_size" toml:"max_header_size" env:"MAX_HEADER_SIZE"`

	// SSL启用参数
	EnableSSL bool   `json:"enable_ssl" yaml:"enable_ssl" toml:"enable_ssl" env:"ENABLE_SSL"`
	CertFile  string `json:"cert_file" yaml:"cert_file" toml:"cert_file" env:"CERT_FILE"`
	KeyFile   string `json:"key_file" yaml:"key_file" toml:"key_file" env:"KEY_FILE"`

	// 开启Trace
	EnableTrace bool `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace" env:"ENABLE_TRACE"`
	// 开启HTTP健康检查
	HealthCheck HealthCheck `toml:"health_check" json:"health_check" yaml:"health_check" envPrefix:"HEALTH_CHECK_"`
	// cors配置
	Cors CORS `toml:"cors" json:"cors" yaml:"cors" envPrefix:"CORS_"`
	// API Doc配置 swagger配置
	ApiDoc ApiDoc `json:"api_doc" yaml:"api_doc" toml:"api_doc" envPrefix:"API_DOC_"`

	// 解析后的数据
	maxHeaderBytes    uint64
	log               *zerolog.Logger
	server            *http.Server
	routerBuilders    map[WEB_FRAMEWORK]RouterBuilder `json:"-" yaml:"-" toml:"-" env:"-"`
	RouterBuildConfig *BuildConfig
}

// `envPrefix:"FOO_"`

type HealthCheck struct {
	Enabled bool `toml:"enabled" json:"enabled" yaml:"enabled"  env:"ENABLED"`
}

type CORS struct {
	Enabled        bool     `toml:"enabled" json:"enabled" yaml:"enabled"  env:"ENABLED"`
	AllowedHeaders []string `json:"cors_allowed_headers" yaml:"cors_allowed_headers" toml:"cors_allowed_headers" env:"ALLOWED_HEADERS" envSeparator:","`
	AllowedDomains []string `json:"cors_allowed_domains" yaml:"cors_allowed_domains" toml:"cors_allowed_domains" env:"ALLOWED_DOMAINS" envSeparator:","`
	AllowedMethods []string `json:"cors_allowed_methods" yaml:"cors_allowed_methods" toml:"cors_allowed_methods" env:"ALLOWED_METHODS" envSeparator:","`
}

type ApiDoc struct {
	// 是否开启API Doc
	Enabled bool `json:"enabled" yaml:"enabled" toml:"enabled" env:"HTTP_API_DOC_ENABLED"`
	// Swagger API Doc URL路径
	DocPath string `json:"doc_path" yaml:"doc_path" toml:"doc_path" env:"HTTP_API_DOC_PATH"`
}

type WEB_FRAMEWORK string

const (
	WEB_FRAMEWORK_GO_RESTFUL WEB_FRAMEWORK = "go-restful"
	WEB_FRAMEWORK_GIN        WEB_FRAMEWORK = "gin"
)

type RouterBuilder interface {
	Config(*BuildConfig)
	Build() (http.Handler, error)
}

type BuildHook func(http.Handler)

type BuildConfig struct {
	// 装载Ioc路由之前
	BeforeLoad BuildHook `json:"-" yaml:"-" toml:"-" env:"-"`
	// 装载Ioc路由之后
	AfterLoad BuildHook `json:"-" yaml:"-" toml:"-" env:"-"`
}

func (h *Http) setEnable(v bool) {
	h.Enable = &v
}

// 配置数据解析
func (h *Http) Parse() error {
	h.log = logger.Sub("http")

	if h.Enable == nil {
		h.setEnable(ioc.Api().Count() > 0)
	}

	if !*h.Enable {
		return nil
	}

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
		Handler:           nil,
	}
	return nil
}

func (h *Http) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

func (h *Http) GetRouterBuilder() RouterBuilder {
	return h.routerBuilders[h.WEB_FRAMEWORK]
}

func (h *Http) BuildRouter() error {
	rb, ok := h.routerBuilders[h.WEB_FRAMEWORK]
	if !ok {
		return fmt.Errorf("router builder for web framework %s not found", h.WEB_FRAMEWORK)
	}

	// 传递配置
	rb.Config(h.RouterBuildConfig)

	r, err := rb.Build()
	if err != nil {
		return err
	}

	h.server.Handler = r
	return nil
}

type ErrHandler func(error)

// Start 启动服务
func (h *Http) Start(ctx context.Context, cb ErrHandler) {
	if err := h.BuildRouter(); err != nil {
		cb(fmt.Errorf("build http router error, %s", err))
		return
	}

	// 启动 HTTP服务
	h.log.Info().Msgf("HTTP服务启动成功, 监听地址: %s", h.Addr())
	cb(h.server.ListenAndServe())
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
