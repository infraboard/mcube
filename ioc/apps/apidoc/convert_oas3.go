package apidoc

import (
	"encoding/json"
	"strings"

	"github.com/go-openapi/spec"
)

// Swagger2ToOpenAPI3 将 Swagger 2.0 转为 OpenAPI 3.0.3（够用的实用转换，非完整语义覆盖）。
func Swagger2ToOpenAPI3(swo *spec.Swagger) map[string]any {
	out := map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":   "API",
			"version": "0.0.0",
		},
		"paths":      map[string]any{},
		"components": map[string]any{},
	}
	if swo == nil {
		return out
	}
	if swo.Info != nil {
		info := map[string]any{
			"title":   swo.Info.Title,
			"version": swo.Info.Version,
		}
		if swo.Info.Description != "" {
			info["description"] = swo.Info.Description
		}
		if swo.Info.License != nil {
			info["license"] = map[string]any{
				"name": swo.Info.License.Name,
				"url":  swo.Info.License.URL,
			}
		}
		out["info"] = info
	}

	servers := []any{}
	scheme := "http"
	if len(swo.Schemes) > 0 {
		scheme = swo.Schemes[0]
	}
	host := strings.TrimSpace(swo.Host)
	base := strings.TrimSpace(swo.BasePath)
	if host != "" {
		url := scheme + "://" + host
		if base != "" && base != "/" {
			url += strings.TrimRight(base, "/")
		}
		servers = append(servers, map[string]any{"url": url})
	} else if base != "" {
		servers = append(servers, map[string]any{"url": base})
	}
	if len(servers) > 0 {
		out["servers"] = servers
	}

	components := map[string]any{}
	if len(swo.Definitions) > 0 {
		schemas := map[string]any{}
		for name, sch := range swo.Definitions {
			schemas[name] = rewriteRefsToComponents(sch)
		}
		components["schemas"] = schemas
	}
	if len(swo.SecurityDefinitions) > 0 {
		schemes := map[string]any{}
		for name, sec := range swo.SecurityDefinitions {
			if sec == nil {
				continue
			}
			entry := map[string]any{"type": sec.Type}
			switch sec.Type {
			case "apiKey":
				entry["name"] = sec.Name
				entry["in"] = sec.In
			case "oauth2":
				entry["flows"] = map[string]any{}
			case "basic":
				entry["type"] = "http"
				entry["scheme"] = "basic"
			}
			if sec.Description != "" {
				entry["description"] = sec.Description
			}
			// Bearer 习惯：jwt header Authorization
			if sec.Type == "apiKey" && strings.EqualFold(sec.Name, "Authorization") {
				entry = map[string]any{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
					"description":  sec.Description,
				}
			}
			schemes[name] = entry
		}
		components["securitySchemes"] = schemes
	}
	out["components"] = components

	paths := map[string]any{}
	if swo.Paths != nil {
		for p, item := range swo.Paths.Paths {
			ops := map[string]any{}
			addOp := func(method string, op *spec.Operation) {
				if op == nil {
					return
				}
				ops[method] = convertOperation(op)
			}
			addOp("get", item.Get)
			addOp("post", item.Post)
			addOp("put", item.Put)
			addOp("patch", item.Patch)
			addOp("delete", item.Delete)
			addOp("head", item.Head)
			addOp("options", item.Options)
			if len(ops) > 0 {
				paths[p] = ops
			}
		}
	}
	out["paths"] = paths

	if len(swo.Security) > 0 {
		out["security"] = swo.Security
	}
	if len(swo.Tags) > 0 {
		tags := make([]any, 0, len(swo.Tags))
		for _, t := range swo.Tags {
			tags = append(tags, map[string]any{"name": t.Name, "description": t.Description})
		}
		out["tags"] = tags
	}
	return out
}

func convertOperation(op *spec.Operation) map[string]any {
	m := map[string]any{}
	if op.ID != "" {
		m["operationId"] = op.ID
	}
	if op.Summary != "" {
		m["summary"] = op.Summary
	}
	if op.Description != "" {
		m["description"] = op.Description
	}
	if len(op.Tags) > 0 {
		m["tags"] = op.Tags
	}
	if len(op.Consumes) > 0 {
		// OAS3 用 requestBody content；这里简化保留扩展
		m["x-consumes"] = op.Consumes
	}
	if len(op.Produces) > 0 {
		m["x-produces"] = op.Produces
	}
	params := make([]any, 0)
	var body *spec.Parameter
	for i := range op.Parameters {
		p := op.Parameters[i]
		if p.In == "body" {
			body = &p
			continue
		}
		params = append(params, rewriteRefsToComponents(p))
	}
	if len(params) > 0 {
		m["parameters"] = params
	}
	if body != nil {
		contentType := "application/json"
		if len(op.Consumes) > 0 {
			contentType = op.Consumes[0]
		}
		rb := map[string]any{
			"required": body.Required,
			"content": map[string]any{
				contentType: map[string]any{
					"schema": rewriteRefsToComponents(body.Schema),
				},
			},
		}
		if body.Description != "" {
			rb["description"] = body.Description
		}
		m["requestBody"] = rb
	}
	if op.Responses != nil {
		responses := map[string]any{}
		if op.Responses.Default != nil {
			responses["default"] = convertResponse(op.Responses.Default, op.Produces)
		}
		for code, resp := range op.Responses.StatusCodeResponses {
			r := resp
			responses[itoa(code)] = convertResponse(&r, op.Produces)
		}
		m["responses"] = responses
	}
	if len(op.Security) > 0 {
		m["security"] = op.Security
	}
	for k, v := range op.Extensions {
		if strings.HasPrefix(k, "x-") {
			m[k] = v
		}
	}
	return m
}

func convertResponse(resp *spec.Response, produces []string) map[string]any {
	out := map[string]any{}
	if resp == nil {
		return out
	}
	if resp.Description != "" {
		out["description"] = resp.Description
	} else {
		out["description"] = "OK"
	}
	if resp.Schema != nil {
		ct := "application/json"
		if len(produces) > 0 {
			ct = produces[0]
		}
		out["content"] = map[string]any{
			ct: map[string]any{
				"schema": rewriteRefsToComponents(resp.Schema),
			},
		}
	}
	for k, v := range resp.Extensions {
		if strings.HasPrefix(k, "x-") {
			out[k] = v
		}
	}
	return out
}

func rewriteRefsToComponents(v any) any {
	b, err := json.Marshal(v)
	if err != nil {
		return v
	}
	s := strings.ReplaceAll(string(b), "#/definitions/", "#/components/schemas/")
	var out any
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		return v
	}
	return out
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [16]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}
