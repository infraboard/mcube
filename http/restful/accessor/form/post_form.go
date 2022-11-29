package form

import (
	"encoding/json"

	"github.com/emicklei/go-restful/v3"
)

const (
	MIME_POST_FORM = "application/x-www-form-urlencoded"
)

// NewEntityAccessorJSON returns a new EntityReaderWriter for accessing Form content.
// This package is already initialized with such an accessor using the MIME_POST_FORM contentType.
func NewEntityAccessorPostForm() restful.EntityReaderWriter {
	return entityPostFormAccess{ContentType: MIME_POST_FORM}
}

type entityPostFormAccess struct {
	// This is used for setting the Content-Type header when writing
	ContentType string
}

// Read unmarshalls the value from Form
func (e entityPostFormAccess) Read(req *restful.Request, v interface{}) error {
	if err := req.Request.ParseForm(); err != nil {
		return err
	}

	return mapForm(v, req.Request.PostForm)
}

// Write marshalls the value to Form and set the Content-Type Header.
func (e entityPostFormAccess) Write(resp *restful.Response, status int, v interface{}) error {
	return writeJSON(resp, status, restful.MIME_JSON, v)
}

// write marshalls the value to YAML and set the Content-Type Header.
func writeJSON(resp *restful.Response, status int, contentType string, v interface{}) error {
	if v == nil {
		resp.WriteHeader(status)
		// do not write a nil representation
		return nil
	}

	resp.Header().Set(restful.HEADER_ContentType, contentType)
	resp.WriteHeader(status)
	return json.NewEncoder(resp).Encode(v)
}

func init() {
	restful.RegisterEntityAccessor(MIME_POST_FORM, NewEntityAccessorPostForm())
}
