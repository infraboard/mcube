package yaml

import (
	"github.com/emicklei/go-restful/v3"
	"gopkg.in/yaml.v3"
)

const (
	MIME_YAML = "application/yaml"
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
	decoder := yaml.NewDecoder(req.Request.Body)
	return decoder.Decode(v)
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
	return yaml.NewEncoder(resp).Encode(v)
}

func init() {
	restful.RegisterEntityAccessor(MIME_YAML, NewEntityAccessorJSON(MIME_YAML))
}
