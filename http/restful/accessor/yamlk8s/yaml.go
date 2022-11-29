package yamlk8s

import (
	"io"

	"github.com/emicklei/go-restful/v3"
	"sigs.k8s.io/yaml"
)

const (
	MIME_YAML = "application/yaml-k8s"
)

// NewEntityAccessorJSON returns a new EntityReaderWriter for accessing YAML content.
// This package is already initialized with such an accessor using the MIME_YAML contentType.
func NewEntityAccessorJSON(contentType string) restful.EntityReaderWriter {
	return entityYAMLAccess{ContentType: contentType}
}

type entityYAMLAccess struct {
	// This is used for setting the Content-Type header when writing
	ContentType string
}

// Read unmarshalls the value from YAML
func (e entityYAMLAccess) Read(req *restful.Request, v interface{}) error {
	payload, err := io.ReadAll(req.Request.Body)
	if err != nil {
		return err
	}
	defer req.Request.Body.Close()

	return yaml.Unmarshal(payload, v)
}

// Write marshalls the value to YAML and set the Content-Type Header.
func (e entityYAMLAccess) Write(resp *restful.Response, status int, v interface{}) error {
	return writeYAML(resp, status, e.ContentType, v)
}

// write marshalls the value to YAML and set the Content-Type Header.
func writeYAML(resp *restful.Response, status int, contentType string, v interface{}) error {
	if v == nil {
		resp.WriteHeader(status)
		// do not write a nil representation
		return nil
	}

	resp.Header().Set(restful.HEADER_ContentType, contentType)
	resp.WriteHeader(status)

	payload, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	_, err = resp.Write(payload)
	return err
}

func init() {
	restful.RegisterEntityAccessor(MIME_YAML, NewEntityAccessorJSON(MIME_YAML))
}
