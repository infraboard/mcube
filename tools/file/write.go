package file

import (
	"github.com/infraboard/mcenter/common/format"
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
	return format.Prettify(v)
}
