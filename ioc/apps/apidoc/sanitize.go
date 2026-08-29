package apidoc

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-openapi/spec"
)

// SanitizeSwaggerDefinitions 修复 go-restful-openapi 对 map/范型名的处理缺陷，
// 并处理 sanitize 后的 key 碰撞。
func SanitizeSwaggerDefinitions(swo *spec.Swagger) {
	if swo == nil || swo.Definitions == nil {
		return
	}

	rename := map[string]string{}
	used := map[string]string{} // sanitized -> original first owner

	allocate := func(old string) string {
		base := SanitizeDefinitionName(old)
		if base == "" {
			base = "Schema"
		}
		candidate := base
		for i := 2; ; i++ {
			if owner, ok := used[candidate]; !ok || owner == old {
				used[candidate] = old
				return candidate
			}
			candidate = fmt.Sprintf("%s_%d", base, i)
		}
	}

	for old := range swo.Definitions {
		key := old
		if dec, err := url.PathUnescape(old); err == nil && dec != old {
			key = dec
		}
		next := allocate(key)
		if next != old {
			rename[old] = next
		}
	}

	nextDefs := make(spec.Definitions, len(swo.Definitions))
	for old, schema := range swo.Definitions {
		key := old
		if n, ok := rename[old]; ok {
			key = n
		}
		nextDefs[key] = schema
	}
	swo.Definitions = nextDefs
	rewriteSwaggerRefs(swo, rename)
}

// SanitizeDefinitionName 将 OpenAPI definition 名中的特殊字符转为安全标识。
func SanitizeDefinitionName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return name
	}
	if dec, err := url.PathUnescape(name); err == nil {
		name = dec
	}
	replacer := strings.NewReplacer(
		"[", "_",
		"]", "_",
		"{", "_",
		"}", "_",
		"<", "_",
		">", "_",
		" ", "",
		"/", "_",
		"*", "Ptr",
		",", "_",
	)
	out := replacer.Replace(name)
	for strings.Contains(out, "__") {
		out = strings.ReplaceAll(out, "__", "_")
	}
	return strings.Trim(out, "_")
}

func rewriteSwaggerRefs(swo *spec.Swagger, rename map[string]string) {
	raw, err := json.Marshal(swo)
	if err != nil {
		return
	}
	s := string(raw)

	for old, next := range rename {
		s = strings.ReplaceAll(s, "#/definitions/"+old, "#/definitions/"+next)
		s = strings.ReplaceAll(s, "#/definitions/"+url.PathEscape(old), "#/definitions/"+next)
		if dec, err := url.PathUnescape(old); err == nil && dec != old {
			s = strings.ReplaceAll(s, "#/definitions/"+dec, "#/definitions/"+next)
			s = strings.ReplaceAll(s, "#/definitions/"+url.PathEscape(dec), "#/definitions/"+next)
		}
	}

	for _, bad := range []string{
		"map[string]string",
		"map[string]interface {}",
		"map[string]interface{}",
	} {
		safe := SanitizeDefinitionName(bad)
		s = strings.ReplaceAll(s, "#/definitions/"+url.PathEscape(bad), "#/definitions/"+safe)
		s = strings.ReplaceAll(s, "#/definitions/"+bad, "#/definitions/"+safe)
	}

	var out spec.Swagger
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return
	}
	*swo = out
}
