package restful

import (
	"encoding/json"
	"net/http"
	"strings"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/apps/apidoc"
	"github.com/infraboard/mcube/v2/ioc/config/application"
	"github.com/infraboard/mcube/v2/ioc/config/gorestful"
	httpconf "github.com/infraboard/mcube/v2/ioc/config/http"
	"github.com/infraboard/mcube/v2/ioc/config/log"
	"github.com/rs/zerolog"
)

func init() {
	ioc.Api().Registry(&SwaggerApiDoc{
		ApiDoc: apidoc.DefaultApiDoc(),
	})
}

// SwaggerApiDoc 在全部 API 注册完成后，提取路由并产出 API Doc。
type SwaggerApiDoc struct {
	ioc.ObjectImpl
	log *zerolog.Logger

	*apidoc.ApiDoc
	cache apidoc.DocCache
}

func (h *SwaggerApiDoc) Name() string {
	return apidoc.AppName
}

func (h *SwaggerApiDoc) Init() error {
	h.log = log.Sub("api_doc")
	h.normalizeConfig()
	h.Registry()
	return nil
}

func (h *SwaggerApiDoc) Priority() int {
	return -100
}

func (h *SwaggerApiDoc) Meta() ioc.ObjectMeta {
	meta := ioc.DefaultObjectMeta()
	if h.BasePath != "" {
		meta.CustomPathPrefix = h.BasePath
	}
	return meta
}

func (h *SwaggerApiDoc) normalizeConfig() {
	if h.ApiDoc == nil {
		h.ApiDoc = apidoc.DefaultApiDoc()
	}
	def := apidoc.DefaultApiDoc()
	if strings.TrimSpace(h.JsonPath) == "" {
		h.JsonPath = def.JsonPath
	}
	if strings.TrimSpace(h.OpenAPIPath) == "" {
		h.OpenAPIPath = def.OpenAPIPath
	}
	if strings.TrimSpace(h.UIPath) == "" {
		h.UIPath = def.UIPath
	}
	if strings.TrimSpace(h.TryItUIPath) == "" {
		h.TryItUIPath = def.TryItUIPath
	}
}

func (h *SwaggerApiDoc) ApiDocPath() string {
	if application.Get().AppAddress != "" {
		return strings.TrimRight(application.Get().AppAddress, "/") + h.ApiDocPathRel()
	}
	return httpconf.Get().ApiObjectAddr(h) + h.JsonPath
}

func (h *SwaggerApiDoc) ApiUIPath() string {
	if application.Get().AppAddress != "" {
		return strings.TrimRight(application.Get().AppAddress, "/") + h.ApiUIPathRel()
	}
	return httpconf.Get().ApiObjectAddr(h) + h.UIPath
}

func (h *SwaggerApiDoc) ApiDocPathRel() string {
	return apidoc.JoinURLPath(httpconf.Get().ApiObjectPathPrefix(h), h.JsonPath)
}

func (h *SwaggerApiDoc) ApiOpenAPIPathRel() string {
	return apidoc.JoinURLPath(httpconf.Get().ApiObjectPathPrefix(h), h.OpenAPIPath)
}

func (h *SwaggerApiDoc) ApiUIPathRel() string {
	return apidoc.JoinURLPath(httpconf.Get().ApiObjectPathPrefix(h), h.UIPath)
}

func (h *SwaggerApiDoc) ApiTryItPathRel() string {
	return apidoc.JoinURLPath(httpconf.Get().ApiObjectPathPrefix(h), h.TryItUIPath)
}

func (h *SwaggerApiDoc) Registry() {
	tags := []string{"API 文档"}
	ws := gorestful.ObjectRouter(h)

	ws.Route(ws.GET(h.JsonPath).To(h.ServeSwaggerJSON).
		Doc("Swagger 2.0 JSON").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	ws.Route(ws.HEAD(h.JsonPath).To(h.ServeSwaggerHEAD).
		Doc("Swagger 2.0 HEAD").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)
	ws.Route(ws.GET(h.OpenAPIPath).To(h.ServeOpenAPIJSON).
		Doc("OpenAPI 3.0 JSON").
		Metadata(restfulspec.KeyOpenAPITags, tags),
	)

	engine := h.EffectiveUIEngine()
	if engine == apidoc.UIEngineScalar || engine == apidoc.UIEngineBoth {
		ws.Route(ws.GET(h.UIPath).To(h.ServeScalarUI).
			Doc("Scalar UI").
			Metadata(restfulspec.KeyOpenAPITags, tags),
		)
	}
	if engine == apidoc.UIEngineSwagger || engine == apidoc.UIEngineBoth {
		ws.Route(ws.GET(h.TryItUIPath).To(h.ServeSwaggerUI).
			Doc("Swagger UI (Try it)").
			Metadata(restfulspec.KeyOpenAPITags, tags),
		)
	}

	h.log.Info().Msgf("API Doc swagger=%s openapi=%s scalar=%s try_it=%s engine=%s",
		h.ApiDocPathRel(), h.ApiOpenAPIPathRel(), h.ApiUIPathRel(), h.ApiTryItPathRel(), engine)
}

func (h *SwaggerApiDoc) disabled(w *restful.Response) bool {
	if h.Enabled {
		return false
	}
	_ = w.WriteErrorString(http.StatusNotFound, "api doc disabled")
	return true
}

func (h *SwaggerApiDoc) ServeSwaggerHEAD(_ *restful.Request, w *restful.Response) {
	if h.disabled(w) {
		return
	}
	b, err := h.bundle()
	if err != nil {
		_ = w.WriteError(500, err)
		return
	}
	w.Header().Set("ETag", b.ETag)
	w.WriteHeader(http.StatusOK)
}

func (h *SwaggerApiDoc) ServeSwaggerJSON(r *restful.Request, w *restful.Response) {
	if h.disabled(w) {
		return
	}
	b, err := h.bundle()
	if err != nil {
		_ = w.WriteError(500, err)
		return
	}
	if match := r.Request.Header.Get("If-None-Match"); match != "" && match == b.ETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	data := b.SwaggerJSON
	if r.QueryParameter("open_to_api_key") == "true" {
		filtered := apidoc.FilterSwaggerByOpenAPIKey(b.Swagger)
		data, err = json.Marshal(filtered)
		if err != nil {
			_ = w.WriteError(500, err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("ETag", b.ETag)
	w.Header().Set("Cache-Control", "public, max-age=30")
	_, _ = w.Write(data)
}

func (h *SwaggerApiDoc) ServeOpenAPIJSON(r *restful.Request, w *restful.Response) {
	if h.disabled(w) {
		return
	}
	b, err := h.bundle()
	if err != nil {
		_ = w.WriteError(500, err)
		return
	}
	if match := r.Request.Header.Get("If-None-Match"); match != "" && match == b.ETag {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	data := b.OpenAPIJSON
	if r.QueryParameter("open_to_api_key") == "true" {
		filtered := apidoc.FilterSwaggerByOpenAPIKey(b.Swagger)
		data, err = json.Marshal(apidoc.Swagger2ToOpenAPI3(filtered))
		if err != nil {
			_ = w.WriteError(500, err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("ETag", b.ETag)
	w.Header().Set("Cache-Control", "public, max-age=30")
	_, _ = w.Write(data)
}

func (h *SwaggerApiDoc) ServeScalarUI(_ *restful.Request, w *restful.Response) {
	if h.disabled(w) {
		return
	}
	switchHref := ""
	if h.EffectiveUIEngine() == apidoc.UIEngineBoth {
		switchHref = h.ApiTryItPathRel()
	}
	// Scalar 吃 OAS3
	html := apidoc.RenderScalarHTML(switchHref, h.ApiOpenAPIPathRel(), h.ScalarJS())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func (h *SwaggerApiDoc) ServeSwaggerUI(_ *restful.Request, w *restful.Response) {
	if h.disabled(w) {
		return
	}
	switchHref := ""
	if h.EffectiveUIEngine() == apidoc.UIEngineBoth {
		switchHref = h.ApiUIPathRel()
	}
	// 试调同样用 OAS3，与 Scalar 同源契约
	html := apidoc.RenderSwaggerUIHTML(switchHref, h.ApiOpenAPIPathRel(), h.SwaggerUIBase())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

func (h *SwaggerApiDoc) bundle() (*apidoc.DocBundle, error) {
	return h.cache.Get(h.CacheTTLSecond, func() (*apidoc.DocBundle, error) {
		swagger := restfulspec.BuildSwagger(h.SwaggerDocConfig())
		if swagger.Host == "" {
			swagger.Host = application.Get().Host()
		}
		openapi := apidoc.Swagger2ToOpenAPI3(swagger)
		return apidoc.NewDocBundle(swagger, openapi)
	})
}

func (h *SwaggerApiDoc) SwaggerDocConfig() restfulspec.Config {
	return restfulspec.Config{
		Host:        application.Get().Host(),
		WebServices: restful.RegisteredWebServices(),
		APIPath:     httpconf.Get().ApiObjectPathPrefix(h),
		PostBuildSwaggerObjectHandler: func(swo *spec.Swagger) {
			httpconf.Get().SwagerDocs(swo)
			apidoc.SanitizeSwaggerDefinitions(swo)
			apidoc.EnrichSwaggerFromRoutes(swo)
			apidoc.ApplyTagDescriptions(swo)
			apidoc.RunPostBuildHooks(swo)
		},
		DefinitionNameHandler: func(name string) string {
			if name == "state" || name == "sizeCache" || name == "unknownFields" {
				return ""
			}
			return apidoc.SanitizeDefinitionName(name)
		},
		Schemes: []string{"http", "https"},
	}
}
