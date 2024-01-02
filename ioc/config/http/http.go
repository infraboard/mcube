package http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/go-openapi/spec"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Config().Registry(&Http{
		Host:                    "127.0.0.1",
		Port:                    8080,
		PathPrefix:              "api",
		ReadHeaderTimeoutSecond: 30,
		ReadTimeoutSecond:       60,
		WriteTimeoutSecond:      60,
		IdleTimeoutSecond:       300,
		MaxHeaderSize:           "16kb",
		EnableTrace:             true,
		WEB_FRAMEWORK:           WEB_FRAMEWORK_AUTO,
		routerBuilders: map[WEB_FRAMEWORK]RouterBuilder{
			WEB_FRAMEWORK_GO_RESTFUL: NewGoRestfulRouterBuilder(),
			WEB_FRAMEWORK_GIN:        NewGinRouterBuilder(),
		},
		handlerCount: map[WEB_FRAMEWORK]int{
			WEB_FRAMEWORK_GO_RESTFUL: 0,
			WEB_FRAMEWORK_GIN:        0,
		},
	})
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

	// 使用的http框架, 默认会根据当前注册的API对象,自动选择合适的框架
	WEB_FRAMEWORK WEB_FRAMEWORK `json:"web_framework" yaml:"web_framework" toml:"web_framework" env:"WEB_FRAMEWORK"`

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

	// 开启Trace
	EnableTrace bool `toml:"enable_trace" json:"enable_trace" yaml:"enable_trace" env:"ENABLE_TRACE"`

	// 解析后的数据
	maxHeaderBytes uint64
	log            *zerolog.Logger
	server         *http.Server
	routerBuilders map[WEB_FRAMEWORK]RouterBuilder `json:"-" yaml:"-" toml:"-" env:"-"`
	handlerCount   map[WEB_FRAMEWORK]int           `json:"-" yaml:"-" toml:"-" env:"-"`
}

func (h *Http) HTTPPrefix() string {
	u, err := url.JoinPath("/"+application.Get().AppName, h.PathPrefix)
	if err != nil {
		return fmt.Sprintf("/%s/%s", application.Get().AppName, h.PathPrefix)
	}
	return u
}

type WEB_FRAMEWORK string

const (
	// 根据ioc当前加载的对象自动判断使用那种框架
	WEB_FRAMEWORK_AUTO       WEB_FRAMEWORK = ""
	WEB_FRAMEWORK_GO_RESTFUL WEB_FRAMEWORK = "go-restful"
	WEB_FRAMEWORK_GIN        WEB_FRAMEWORK = "gin"
)

func (h *Http) setEnable(v bool) {
	h.Enable = &v
}

func (h *Http) DetectAndSetWebFramework() {
	if h.Enable == nil && h.WEB_FRAMEWORK == WEB_FRAMEWORK_AUTO {
		ioc.Api().ForEach(func(w *ioc.ObjectWrapper) {
			switch w.Value.(type) {
			case ioc.GoRestfulApiObject:
				h.handlerCount[WEB_FRAMEWORK_GO_RESTFUL]++
			case ioc.GinApiObject:
				h.handlerCount[WEB_FRAMEWORK_GIN]++
			}
		})
		if wf, count := h.maxHandlerCount(); count > 0 {
			h.setEnable(true)
			h.WEB_FRAMEWORK = wf
		}
	}

	if h.Enable == nil {
		h.setEnable(false)
	}
}

func (h *Http) maxHandlerCount() (maxKey WEB_FRAMEWORK, maxValue int) {
	for k, v := range h.handlerCount {
		if v > maxValue {
			maxKey = k
			maxValue = v
		}
	}
	return
}

func (h *Http) Name() string {
	return AppName
}

// 配置数据解析
func (h *Http) Init() error {
	h.log = log.Sub("http")

	h.DetectAndSetWebFramework()
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

	r, err := rb.Build()
	if err != nil {
		return err
	}

	h.server.Handler = r
	return nil
}

// Start 启动服务
func (h *Http) Start(ctx context.Context) {
	if err := h.BuildRouter(); err != nil {
		h.log.Error().Msgf("build http router error, %s", err)
		return
	}

	// 启动 HTTP服务
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
}
