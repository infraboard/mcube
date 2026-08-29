package apidoc

import "strings"

const (
	AppName = "apidoc"

	// Well-known route metadata keys（与业务侧约定一致，写入 OpenAPI x-* 扩展）
	MetaPerm              = "perm"
	MetaOpenToAPIKey      = "open_to_api_key"
	MetaRequiredAuth      = "required_auth"
	MetaRequiredNamespace = "required_namespace"
	MetaResource          = "resource"
	MetaAction            = "action"

	ExtPerm              = "x-perm"
	ExtOpenToAPIKey      = "x-open-to-api-key"
	ExtRequiredAuth      = "x-required-auth"
	ExtRequiredNamespace = "x-required-namespace"
	ExtResource          = "x-resource"
	ExtAction            = "x-action"
)

// UIEngine 文档壳类型
type UIEngine string

const (
	// UIEngineScalar 默认阅读页（产品主文档）
	UIEngineScalar UIEngine = "scalar"
	// UIEngineSwagger 仅 Swagger UI 试调
	UIEngineSwagger UIEngine = "swagger"
	// UIEngineBoth Scalar（读）+ Swagger UI（试），推荐默认
	UIEngineBoth UIEngine = "both"
	// UIEngineRedoc 已废弃，等同 scalar
	UIEngineRedoc UIEngine = "redoc"
)

func DefaultApiDoc() *ApiDoc {
	return &ApiDoc{
		BasePath:         "",
		JsonPath:         "/swagger.json",
		OpenAPIPath:      "/openapi.json",
		UIPath:           "/ui.html",
		TryItUIPath:      "/swagger.html",
		Enabled:          true,
		CacheTTLSecond:   0,
		UIEngine:         string(UIEngineBoth),
		CDNScalarJS:      "https://cdn.jsdelivr.net/npm/@scalar/api-reference",
		CDNSwaggerUIBase: "https://unpkg.com/swagger-ui-dist@5.17.14",
	}
}

type ApiDoc struct {
	BasePath       string `json:"base_path" yaml:"base_path" toml:"base_path" env:"BASE_PATH"`
	JsonPath       string `json:"json_path" yaml:"json_path" toml:"json_path" env:"JSON_PATH"`
	OpenAPIPath    string `json:"openapi_path" yaml:"openapi_path" toml:"openapi_path" env:"OPENAPI_PATH"`
	UIPath         string `json:"ui_path" yaml:"ui_path" toml:"ui_path" env:"UI_PATH"`
	TryItUIPath    string `json:"try_it_ui_path" yaml:"try_it_ui_path" toml:"try_it_ui_path" env:"TRY_IT_UI_PATH"`
	Enabled        bool   `json:"enabled" yaml:"enabled" toml:"enabled" env:"ENABLED"`
	CacheTTLSecond int    `json:"cache_ttl_second" yaml:"cache_ttl_second" toml:"cache_ttl_second" env:"CACHE_TTL_SECOND"`
	// UIEngine scalar | swagger | both（redoc 兼容映射为 scalar）
	UIEngine string `json:"ui_engine" yaml:"ui_engine" toml:"ui_engine" env:"UI_ENGINE"`
	// CDNScalarJS @scalar/api-reference 脚本（内网可改本地）
	CDNScalarJS string `json:"cdn_scalar_js" yaml:"cdn_scalar_js" toml:"cdn_scalar_js" env:"CDN_SCALAR_JS"`
	// CDNSwaggerUIBase swagger-ui-dist 根路径
	CDNSwaggerUIBase string `json:"cdn_swagger_ui_base" yaml:"cdn_swagger_ui_base" toml:"cdn_swagger_ui_base" env:"CDN_SWAGGER_UI_BASE"`
	// CDNRedocJS 已废弃，忽略
	CDNRedocJS string `json:"cdn_redoc_js" yaml:"cdn_redoc_js" toml:"cdn_redoc_js" env:"CDN_REDOC_JS"`
}

func (d *ApiDoc) EffectiveUIEngine() UIEngine {
	if d == nil {
		return UIEngineBoth
	}
	switch UIEngine(strings.ToLower(strings.TrimSpace(d.UIEngine))) {
	case UIEngineScalar, UIEngineRedoc:
		return UIEngineScalar
	case UIEngineSwagger:
		return UIEngineSwagger
	case UIEngineBoth:
		return UIEngineBoth
	default:
		return UIEngineBoth
	}
}

func (d *ApiDoc) ScalarJS() string {
	if d != nil && strings.TrimSpace(d.CDNScalarJS) != "" {
		return strings.TrimSpace(d.CDNScalarJS)
	}
	return DefaultApiDoc().CDNScalarJS
}

func (d *ApiDoc) SwaggerUIBase() string {
	if d != nil && strings.TrimSpace(d.CDNSwaggerUIBase) != "" {
		return strings.TrimRight(strings.TrimSpace(d.CDNSwaggerUIBase), "/")
	}
	return DefaultApiDoc().CDNSwaggerUIBase
}

// JoinURLPath 拼接 API 相对路径，避免 path.Join 把绝对段冲掉前缀。
func JoinURLPath(parts ...string) string {
	out := ""
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if out == "" {
			out = "/" + strings.Trim(p, "/")
			continue
		}
		out = strings.TrimRight(out, "/") + "/" + strings.Trim(p, "/")
	}
	if out == "" {
		return "/"
	}
	return out
}
