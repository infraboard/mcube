package binding

import (
	"bytes"
	"io"
	"net/http"

	"sigs.k8s.io/yaml"
)

type yamlK8sBinding struct{}

func (yamlK8sBinding) Name() string {
	return "yaml"
}

func (yamlK8sBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeK8sYAML(req.Body, obj)
}

func (yamlK8sBinding) BindBody(body []byte, obj interface{}) error {
	return decodeK8sYAML(bytes.NewReader(body), obj)
}

func decodeK8sYAML(r io.Reader, obj interface{}) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, obj); err != nil {
		return err
	}

	return validate(obj)
}
