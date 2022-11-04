package negotiator

import "gopkg.in/yaml.v3"

type yamlImpl struct{}

func (i *yamlImpl) Encode(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (i *yamlImpl) Decode(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

func (i *yamlImpl) Name() string {
	return "yaml"
}

func init() {
	Registry(MIMEYAML, &yamlImpl{})
}
