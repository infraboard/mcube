package negotiator

import "gopkg.in/yaml.v3"

const (
	MIME_YAML MIME = "application/yaml"
)

type yamlImpl struct{}

func (i *yamlImpl) Encode(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (i *yamlImpl) Decode(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

func (i *yamlImpl) ContentType() MIME {
	return MIME_YAML
}

func init() {
	Registry(&yamlImpl{})
}
