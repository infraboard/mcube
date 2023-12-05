package file

import (
	"github.com/infraboard/mcube/v2/tools/pretty"
	"sigs.k8s.io/yaml"
)

func MustToYaml(v any) string {
	b, err := yaml.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func MustToJson(v any) string {
	return pretty.MustToYaml(v)
}
