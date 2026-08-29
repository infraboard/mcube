package apidoc

import (
	"testing"

	"github.com/go-openapi/spec"
)

func TestSanitizeDefinitionName(t *testing.T) {
	cases := map[string]string{
		"map[string]string":       "map_string_string",
		"map[string]interface {}": "map_string_interface",
		"map%5Bstring%5Dstring":   "map_string_string",
		"api.NamespaceSet":        "api.NamespaceSet",
		"types.Set[*User]":        "types.Set_PtrUser",
	}
	for in, want := range cases {
		if got := SanitizeDefinitionName(in); got != want {
			t.Fatalf("%q => %q, want %q", in, got, want)
		}
	}
}

func TestSanitizeSwaggerDefinitions_refEncodeMismatch(t *testing.T) {
	swo := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: spec.Definitions{
				"map[string]string": *spec.MapProperty(spec.StringProperty()),
			},
			Paths: &spec.Paths{
				Paths: map[string]spec.PathItem{
					"/x": {
						PathItemProps: spec.PathItemProps{
							Get: &spec.Operation{
								OperationProps: spec.OperationProps{
									Responses: &spec.Responses{
										ResponsesProps: spec.ResponsesProps{
											StatusCodeResponses: map[int]spec.Response{
												200: {
													ResponseProps: spec.ResponseProps{
														Schema: &spec.Schema{
															SchemaProps: spec.SchemaProps{
																Ref: spec.MustCreateRef("#/definitions/map%5Bstring%5Dstring"),
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	SanitizeSwaggerDefinitions(swo)
	if _, ok := swo.Definitions["map_string_string"]; !ok {
		t.Fatalf("expected sanitized definition key, got %#v", swo.Definitions)
	}
	ref := swo.Paths.Paths["/x"].Get.Responses.StatusCodeResponses[200].Schema.Ref.String()
	if ref != "#/definitions/map_string_string" {
		t.Fatalf("ref not rewritten: %s", ref)
	}
}

func TestSanitizeCollision(t *testing.T) {
	swo := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: spec.Definitions{
				"Foo[Bar]": *spec.StringProperty(),
				"Foo_Bar":  *spec.StringProperty(),
			},
		},
	}
	SanitizeSwaggerDefinitions(swo)
	if len(swo.Definitions) != 2 {
		t.Fatalf("expected 2 defs after collision allocate, got %d %#v", len(swo.Definitions), swo.Definitions)
	}
}

func TestSwagger2ToOpenAPI3(t *testing.T) {
	swo := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger: "2.0",
			Info: &spec.Info{
				InfoProps: spec.InfoProps{Title: "demo", Version: "1.0"},
			},
			Definitions: spec.Definitions{
				"User": *spec.StringProperty(),
			},
			Paths: &spec.Paths{
				Paths: map[string]spec.PathItem{
					"/users": {
						PathItemProps: spec.PathItemProps{
							Get: &spec.Operation{
								OperationProps: spec.OperationProps{
									ID: "ListUsers",
									Responses: &spec.Responses{
										ResponsesProps: spec.ResponsesProps{
											StatusCodeResponses: map[int]spec.Response{
												200: {
													ResponseProps: spec.ResponseProps{
														Description: "ok",
														Schema: &spec.Schema{
															SchemaProps: spec.SchemaProps{
																Ref: spec.MustCreateRef("#/definitions/User"),
															},
														},
													},
												},
											},
										},
									},
								},
								VendorExtensible: spec.VendorExtensible{
									Extensions: spec.Extensions{ExtOpenToAPIKey: true, ExtPerm: "user:read"},
								},
							},
						},
					},
				},
			},
		},
	}
	oa := Swagger2ToOpenAPI3(swo)
	if oa["openapi"] != "3.0.3" {
		t.Fatalf("openapi version: %v", oa["openapi"])
	}
	comps := oa["components"].(map[string]any)
	schemas := comps["schemas"].(map[string]any)
	if _, ok := schemas["User"]; !ok {
		t.Fatal("missing components.schemas.User")
	}
	paths := oa["paths"].(map[string]any)
	users := paths["/users"].(map[string]any)
	get := users["get"].(map[string]any)
	if get["operationId"] != "ListUsers" {
		t.Fatalf("operationId=%v", get["operationId"])
	}
	if get[ExtPerm] != "user:read" {
		t.Fatalf("x-perm=%v", get[ExtPerm])
	}
	resp := get["responses"].(map[string]any)["200"].(map[string]any)
	content := resp["content"].(map[string]any)["application/json"].(map[string]any)
	schema := content["schema"].(map[string]any)
	ref, _ := schema["$ref"].(string)
	if ref != "#/components/schemas/User" {
		t.Fatalf("ref=%s", ref)
	}
}

func TestFilterSwaggerByOpenAPIKey(t *testing.T) {
	swo := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Paths: &spec.Paths{
				Paths: map[string]spec.PathItem{
					"/open": {
						PathItemProps: spec.PathItemProps{
							Get: &spec.Operation{
								VendorExtensible: spec.VendorExtensible{
									Extensions: spec.Extensions{ExtOpenToAPIKey: true},
								},
							},
						},
					},
					"/closed": {
						PathItemProps: spec.PathItemProps{
							Get: &spec.Operation{},
						},
					},
				},
			},
		},
	}
	out := FilterSwaggerByOpenAPIKey(swo)
	if _, ok := out.Paths.Paths["/open"]; !ok {
		t.Fatal("open path missing")
	}
	if _, ok := out.Paths.Paths["/closed"]; ok {
		t.Fatal("closed path should be filtered")
	}
}

func TestJoinURLPath(t *testing.T) {
	got := JoinURLPath("/api/app/v1/apidoc", "/swagger.json")
	if got != "/api/app/v1/apidoc/swagger.json" {
		t.Fatalf("got %s", got)
	}
}

func TestRenderScalarHTML(t *testing.T) {
	html := RenderScalarHTML("/try", "/api/openapi.json", "https://cdn.example/scalar.js")
	if !containsAll(html, "/try", "/api/openapi.json", "https://cdn.example/scalar.js", "Scalar.createApiReference", "doc-switch-link", "generateOperationSlug") {
		t.Fatalf("html missing pieces: %s", html)
	}
	if stringContains(html, "position: fixed") {
		t.Fatal("switch bar must not be fixed overlay")
	}
}

func TestEffectiveUIEngine(t *testing.T) {
	d := &ApiDoc{UIEngine: "redoc"}
	if d.EffectiveUIEngine() != UIEngineScalar {
		t.Fatal("redoc should map to scalar")
	}
	d.UIEngine = "both"
	if d.EffectiveUIEngine() != UIEngineBoth {
		t.Fatal("both")
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !stringContains(s, p) {
			return false
		}
	}
	return true
}

func stringContains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
