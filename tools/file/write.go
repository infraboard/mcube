package file

import (
	"os"

	"github.com/BurntSushi/toml"
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

func MustToToml(key string, value any, path string) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	appConf := map[string]any{key: value}
	toml.NewEncoder(f).Encode(appConf)
	return nil
}
