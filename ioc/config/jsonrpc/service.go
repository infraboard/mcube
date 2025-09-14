package jsonrpc

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sync"

	"github.com/emicklei/go-restful/v3"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/infraboard/mcube/v2/ioc/config/trace"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
)

func init() {
	ioc.Api().Registry(&JsonRpc{
		Host:       "127.0.0.1",
		Port:       9090,
		PathPrefix: "jsonrpc",
		methods:    map[string]*MethodInfo{},
	})
}

type JsonRpc struct {
	ioc.ObjectImpl

	// 是否开启HTTP Server, 默认会根据是否有注册得有API对象来自动开启
	Enable *bool `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	// HTTP服务Host
	Host string `json:"host" yaml:"host" toml:"host" env:"HOST"`
	// HTTP服务端口
	Port int `json:"port" yaml:"port" toml:"port" env:"PORT"`
	// API接口前缀
	PathPrefix string `json:"path_prefix" yaml:"path_prefix" toml:"path_prefix" env:"PATH_PREFIX"`
	// 开启Trace
	Trace bool `toml:"trace" json:"trace" yaml:"trace" env:"TRACE"`
	// 访问日志
	AccessLog bool `toml:"access_log" json:"access_log" yaml:"access_log" env:"ACCESS_LOG"`

	EnableSSL bool   `json:"enable_ssl" yaml:"enable_ssl" toml:"enable_ssl" env:"ENABLE_SSL"`
	CertFile  string `json:"cert_file" yaml:"cert_file" toml:"cert_file" env:"CERT_FILE"`
	KeyFile   string `json:"key_file" yaml:"key_file" toml:"key_file" env:"KEY_FILE"`

	server    *http.Server
	Container *restful.Container
	mu        sync.RWMutex
	log       *zerolog.Logger
	methods   map[string]*MethodInfo
}

func (h *JsonRpc) Addr() string {
	return fmt.Sprintf("%s:%d", h.Host, h.Port)
}

func (j *JsonRpc) Priority() int {
	return -89
}

func (j *JsonRpc) Name() string {
	return APP_NAME
}

func (h *JsonRpc) HTTPPrefix() string {
	u, err := url.JoinPath("/", h.PathPrefix, application.Get().AppName, h.Version())
	if err != nil {
		return fmt.Sprintf("/%s/%s/%s", application.Get().AppName, h.PathPrefix, h.Version())
	}
	return u
}

func (h *JsonRpc) RPCURL() string {
	return fmt.Sprintf("http://%s%s", h.Addr(), h.HTTPPrefix())
}

func (h *JsonRpc) Start(ctx context.Context) {
	h.log.Info().Msgf("JSON RPC服务启动成功, 监听地址: %s", h.RPCURL())
	if err := h.server.ListenAndServe(); err != nil {
		h.log.Error().Msg(err.Error())
	}
}

// Stop 停止server
func (h *JsonRpc) Stop(ctx context.Context) error {
	h.log.Info().Msg("start graceful shutdown")
	// 优雅关闭HTTP服务
	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http graceful shutdown timeout, force exit")
	}
	return nil
}

func (h *JsonRpc) IsEnable() bool {
	if h.Enable == nil {
		return len(h.methods) > 0
	}

	return *h.Enable
}

func (j *JsonRpc) Init() error {
	j.log = log.Sub(j.Name())

	if len(j.methods) == 0 {
		j.log.Info().Msgf("no reigstry service")
		return nil
	}

	// 在Init函数中修改循环部分
	j.PrintMethods()

	j.Container = restful.DefaultContainer
	restful.DefaultResponseContentType(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)

	// 注册路由
	if j.Trace && trace.Get().Enable {
		j.log.Info().Msg("enable jsonrpc trace")
		j.Container.Filter(otelrestful.OTelFilter(application.Get().GetAppName()))
	}

	// RPC的服务架设在“/jsonrpc”路径，
	// 在处理函数中基于http.ResponseWriter和http.Request类型的参数构造一个io.ReadWriteCloser类型的conn通道。
	// 然后基于conn构建针对服务端的json编码解码器。
	// 最后通过rpc.ServeRequest函数为每次请求处理一次RPC方法调用
	ws := new(restful.WebService)
	ws.Path(j.HTTPPrefix()).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Route(ws.POST("").To(j.HandleRequest))
	// 添加到Root Container
	RootRouter().Add(ws)

	j.server = &http.Server{
		Addr:    j.Addr(),
		Handler: j.Container,
	}
	return nil
}

// 打印所有注册的方法信息
func (j *JsonRpc) PrintMethods() {
	j.mu.RLock()
	defer j.mu.RUnlock()

	for name, info := range j.methods {
		j.log.Info().Msgf("method: %s --> %s(%s)", name, info.FuncName, info.ParamType.String())
	}
}

// 方法信息结构
type MethodInfo struct {
	Name      string       // 方法名
	Handler   HandlerFunc  // 处理器函数
	FuncName  string       // 原始函数名
	ParamType reflect.Type // 参数类型
}
