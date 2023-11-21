package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/infraboard/mcube/ioc"
	"github.com/infraboard/mcube/ioc/config/logger"
	"github.com/rs/zerolog"
)

func NewDefaultHttp() *Http {
	return &Http{
		Host:                    "127.0.0.1",
		Port:                    8080,
		ReadHeaderTimeoutSecond: 30,
		ReadTimeoutSecond:       60,
		WriteTimeoutSecond:      60,
		IdleTimeoutSecond:       300,
		MaxHeaderSize:           "16kb",
		WEB_FRAMEWORK:           WEB_FRAMEWORK_GIN,
		routerBuilders: map[WEB_FRAMEWORK]RouterBuilder{
			WEB_FRAMEWORK_GO_RESTFUL: NewGoRestfulRouterBuilder(),
			WEB_FRAMEWORK_GIN:        NewGinRouterBuilder(),
		},
		RouterBuildConfig: &BuildConfig{},
	}
}

type Http struct {
	// 默认根据
	Enable *bool `json:"enable" yaml:"enable" toml:"enable" env:"HTTP_ENABLE"`
	// HTTP服务Host
	Host string `json:"size" yaml:"size" toml:"size" env:"HTTP_HOST"`
	// HTTP服务端口
	Port int `json:"port" yaml:"port" toml:"port" env:"HTTP_PORT"`

	// 使用的http框架, 启用后会自动从ioc中加载 该框架的hanlder
	WEB_FRAMEWORK WEB_FRAMEWORK `json:"web_framework" yaml:"web_framework" toml:"web_framework" env:"HTTP_WEB_FRAMEWORK"`

	// HTTP服务器参数
	ReadHeaderTimeoutSecond int `json:"read_header_timeout" yaml:"read_header_timeout" toml:"read_header_timeout" env:"HTTP_READ_HEADER_TIMEOUT"`
	// 读取HTTP整个请求时的参数
	ReadTimeoutSecond int `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout" env:"HTTP_READ_TIMEOUT"`
	// 响应超时事件
	WriteTimeoutSecond int `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout" env:"HTTP_WRITE_TIMEOUT"`
	// 启用了KeepAlive时 复用TCP链接的超时时间
	IdleTimeoutSecond int `json:"idle_timeout" yaml:"idle_timeout" toml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT"`
	// header最大大小
	MaxHeaderSize string `json:"max_header_size" yaml:"max_header_size" toml:"max_header_size" env:"HTTP_MAX_HEADER_SIZE"`

	// SSL启用参数
	EnableSSL bool   `json:"enable_ssl" yaml:"enable_ssl" toml:"enable_ssl" env:"HTTP_ENABLE_SSL"`
	CertFile  string `json:"cert_file" yaml:"cert_file" toml:"cert_file" env:"HTTP_CERT_FILE"`
	KeyFile   string `json:"key_file" yaml:"key_file" toml:"key_file" env:"HTTP_KEY_FILE"`

	// 开启Trace
	EnableTrace bool `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace"  env:"HTTP_ENABLE_TRACE"`
	// 开启HTTP健康检查
	EnableHealthCheck bool `toml:"enable_health_check" json:"enable_health_check" yaml:"enable_health_check"  env:"HTTP_ENABLE_HEALTH_CHECK"`
	// 开启跨越允许
	EnableCors bool `toml:"enable_cors" json:"enable_cors" yaml:"enable_cors"  env:"HTTP_ENABLE_CORS"`

	// 是否开启API Doc
	EnableApiDoc bool   `json:"enable_api_doc" yaml:"enable_api_doc" toml:"enable_api_doc" env:"HTTP_ENABLE_API_DOC"`
	ApiDocPath   string `json:"api_doc_path" yaml:"api_doc_path" toml:"api_doc_path" env:"HTTP_API_DOC_PATH"`

	// 解析后的数据
	maxHeaderBytes    uint64
	log               *zerolog.Logger
	server            *http.Server
	routerBuilders    map[WEB_FRAMEWORK]RouterBuilder `json:"-" yaml:"-" toml:"-" env:"-"`
	RouterBuildConfig *BuildConfig
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
