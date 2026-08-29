package apidoc

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
)

var pathParamPattern = regexp.MustCompile(`{([^}]*)}`)

// EnrichSwaggerFromRoutes 把 restful Route Metadata 写入 operation 的 x-* 扩展。
func EnrichSwaggerFromRoutes(swo *spec.Swagger) {
	if swo == nil || swo.Paths == nil {
		return
	}
	for _, ws := range restful.RegisteredWebServices() {
		if ws == nil {
			continue
		}
		for _, route := range ws.Routes() {
			op := findOperation(swo, route.Method, route.Path)
			if op == nil {
				continue
			}
			setExtString(op, ExtPerm, metaString(route.Metadata, MetaPerm))
			if b, ok := metaBool(route.Metadata, MetaOpenToAPIKey); ok {
				setExtBool(op, ExtOpenToAPIKey, b)
			}
			if b, ok := metaBool(route.Metadata, MetaRequiredAuth); ok {
				setExtBool(op, ExtRequiredAuth, b)
			}
			if b, ok := metaBool(route.Metadata, MetaRequiredNamespace); ok {
				setExtBool(op, ExtRequiredNamespace, b)
			}
			setExtString(op, ExtResource, metaString(route.Metadata, MetaResource))
			setExtString(op, ExtAction, metaString(route.Metadata, MetaAction))
		}
	}
}

func findOperation(swo *spec.Swagger, method, restfulPath string) *spec.Operation {
	sanitized, _ := sanitizeRestfulPath(restfulPath)
	item, ok := swo.Paths.Paths[sanitized]
	if !ok {
		item, ok = swo.Paths.Paths[restfulPath]
		if !ok {
			return nil
		}
	}
	switch strings.ToUpper(method) {
	case "GET":
		return item.Get
	case "POST":
		return item.Post
	case "PUT":
		return item.Put
	case "PATCH":
		return item.Patch
	case "DELETE":
		return item.Delete
	case "HEAD":
		return item.Head
	case "OPTIONS":
		return item.Options
	default:
		return nil
	}
}

func sanitizeRestfulPath(restfulPath string) (string, map[string]string) {
	openapiPath := ""
	patterns := map[string]string{}
	for _, fragment := range strings.Split(restfulPath, "/") {
		if fragment == "" {
			continue
		}
		if !strings.ContainsAny(fragment, "{}:*") {
			openapiPath += "/" + fragment
			continue
		}
		matches := pathParamPattern.FindStringSubmatch(fragment)
		if len(matches) == 2 && matches[1] != "" {
			name := matches[1]
			if strings.Contains(name, ":") {
				parts := strings.SplitN(name, ":", 2)
				name = parts[0]
				patterns[name] = parts[1]
			}
			openapiPath += "/{" + name + "}"
			continue
		}
		openapiPath += "/" + fragment
	}
	if openapiPath == "" {
		openapiPath = "/"
	}
	return openapiPath, patterns
}

func metaString(meta map[string]interface{}, key string) string {
	if meta == nil {
		return ""
	}
	v, ok := meta[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}

func metaBool(meta map[string]interface{}, key string) (bool, bool) {
	if meta == nil {
		return false, false
	}
	v, ok := meta[key]
	if !ok || v == nil {
		return false, false
	}
	b, ok := v.(bool)
	return b, ok
}

func setExtString(op *spec.Operation, key, val string) {
	if op == nil || val == "" {
		return
	}
	if op.Extensions == nil {
		op.Extensions = spec.Extensions{}
	}
	op.Extensions[key] = val
}

func setExtBool(op *spec.Operation, key string, val bool) {
	if op == nil {
		return
	}
	if op.Extensions == nil {
		op.Extensions = spec.Extensions{}
	}
	op.Extensions[key] = val
}

// FilterSwaggerByOpenAPIKey 仅保留 x-open-to-api-key=true 的 operation。
func FilterSwaggerByOpenAPIKey(src *spec.Swagger) *spec.Swagger {
	if src == nil || src.Paths == nil {
		return src
	}
	b, err := json.Marshal(src)
	if err != nil {
		return src
	}
	raw := &spec.Swagger{}
	if err := json.Unmarshal(b, raw); err != nil {
		return src
	}
	next := &spec.Paths{Paths: map[string]spec.PathItem{}}
	for p, item := range raw.Paths.Paths {
		filtered := spec.PathItem{}
		keep := false
		if openAPIKeyOp(item.Get) {
			filtered.Get = item.Get
			keep = true
		}
		if openAPIKeyOp(item.Post) {
			filtered.Post = item.Post
			keep = true
		}
		if openAPIKeyOp(item.Put) {
			filtered.Put = item.Put
			keep = true
		}
		if openAPIKeyOp(item.Patch) {
			filtered.Patch = item.Patch
			keep = true
		}
		if openAPIKeyOp(item.Delete) {
			filtered.Delete = item.Delete
			keep = true
		}
		if keep {
			next.Paths[p] = filtered
		}
	}
	raw.Paths = next
	return raw
}

func openAPIKeyOp(op *spec.Operation) bool {
	if op == nil || op.Extensions == nil {
		return false
	}
	v, ok := op.Extensions[ExtOpenToAPIKey]
	if !ok {
		return false
	}
	b, ok := v.(bool)
	return ok && b
}
