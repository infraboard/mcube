package negotiator

import "sigs.k8s.io/yaml"

const (
	MIME_YAML_K8S MIME = "application/yaml-k8s"
)

type yamlK8sImpl struct{}

func (i *yamlK8sImpl) Encode(v any) ([]byte, error) {
	return yaml.Marshal(v)
}

func (i *yamlK8sImpl) Decode(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}

func (i *yamlK8sImpl) ContentType() MIME {
	return MIME_YAML_K8S
}

func init() {
	Registry(&yamlK8sImpl{})
}
