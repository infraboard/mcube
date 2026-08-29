package apidoc

import (
	"testing"

	"github.com/go-openapi/spec"
)

func TestApplyTagDescriptions(t *testing.T) {
	RegisterTagDescription("应用制品", "## 应用制品\n\n登记、证据与质量门禁。")
	swo := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Paths: &spec.Paths{
				Paths: map[string]spec.PathItem{
					"/a": {
						PathItemProps: spec.PathItemProps{
							Get: &spec.Operation{
								OperationProps: spec.OperationProps{Tags: []string{"应用制品", "其它"}},
							},
						},
					},
				},
			},
		},
	}
	ApplyTagDescriptions(swo)
	var hit *spec.Tag
	for i := range swo.Tags {
		if swo.Tags[i].Name == "应用制品" {
			hit = &swo.Tags[i]
			break
		}
	}
	if hit == nil || hit.Description == "" {
		t.Fatalf("tag description missing: %#v", swo.Tags)
	}
	if !containsAll(hit.Description, "登记", "质量门禁") {
		t.Fatalf("md body incomplete: %s", hit.Description)
	}
}
