package application

import (
	"context"
	"fmt"
	"net/http"
	"time"

	humanize "github.com/dustin/go-humanize"
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
		WEB_FRAMEWORK:           WEB_FRAMEWORK_GO_RESTFUL,
		log:                     logger.Sub("http"),
	}
}

type Http struct {
	Host string `json:"size" yaml:"size" toml:"size" env:"HTTP_HOST"`
	Port int    `json:"port" yaml:"port" toml:"port" env:"HTTP_PORT"`

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

	// 是否开启API Doc
	EnableApiDoc bool   `json:"enable_api_doc" yaml:"enable_api_doc" toml:"enable_api_doc" env:"HTTP_ENABLE_API_DOC"`
	ApiDocPath   string `json:"api_doc_path" yaml:"api_doc_path" toml:"api_doc_path" env:"HTTP_API_DOC_PATH"`

	// 解析后的数据
	maxHeaderBytes uint64
	log            *zerolog.Logger
	server         *http.Server
}

type WEB_FRAMEWORK string

const (
	WEB_FRAMEWORK_GO_RESTFUL WEB_FRAMEWORK = "go-restful"
	WEB_FRAMEWORK_GIN        WEB_FRAMEWORK = "gin"
)

// 配置数据解析
func (h *Http) Parse() error {
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

func (h *Http) BuildRouter() error {
	rb, ok := RouterBuilderStore[h.WEB_FRAMEWORK]
	if !ok {
		return fmt.Errorf("router builder for web framework %s not found", h.WEB_FRAMEWORK)
	}

	if err := rb.Init(); err != nil {
		return err
	}

	h.server.Handler = rb.GetRouter()
	return nil
}

// Start 启动服务
func (h *Http) Start(ctx context.Context) {
	h.BuildRouter()

	// // 注册路由条目
	// // s.RegistryEndpoint(ctx)

	// 启动 HTTP服务
	h.log.Info().Msgf("HTTP服务启动成功, 监听地址: %s", h.Addr())
	if err := h.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			h.log.Info().Msg("service is stopped")
			return
		}
		h.log.Error().Msgf("start service error, %s", err.Error())
	}
}

// Stop 停止server
func (h *Http) Stop(ctx context.Context) error {
	h.log.Info().Msg("start graceful shutdown")
	// 优雅关闭HTTP服务
	if err := h.server.Shutdown(ctx); err != nil {
		h.log.Error().Msg("graceful shutdown timeout, force exit")
	}
	return nil
}
