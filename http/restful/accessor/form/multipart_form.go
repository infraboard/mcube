package form

import (
	"errors"
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/emicklei/go-restful/v3"
)

const (
	MIME_MULTIPART_FORM = "multipart/form-data"
)

const (
	DEFAULT_MEMORY = 32 << 20
)

var (
	// ErrMultiFileHeader multipart.FileHeader invalid
	ErrMultiFileHeader = errors.New("unsupported field type for multipart.FileHeader")

	// ErrMultiFileHeaderLenInvalid array for []*multipart.FileHeader len invalid
	ErrMultiFileHeaderLenInvalid = errors.New("unsupported len of array for []*multipart.FileHeader")
)

// NewEntityAccessorJSON returns a new EntityReaderWriter for accessing Form content.
// This package is already initialized with such an accessor using the MIME_POST_FORM contentType.
func NewEntityAccessorMultipartForm() restful.EntityReaderWriter {
	return entityMultipartFormAccess{ContentType: MIME_MULTIPART_FORM}
}

type entityMultipartFormAccess struct {
	// This is used for setting the Content-Type header when writing
	ContentType string
}

// Read unmarshalls the value from Form
func (e entityMultipartFormAccess) Read(req *restful.Request, v interface{}) error {
	if err := req.Request.ParseMultipartForm(DEFAULT_MEMORY); err != nil {
		return err
	}

	if err := mappingByPtr(v, (*multipartRequest)(req.Request), "form"); err != nil {
		return err
	}

	return mapForm(v, req.Request.PostForm)
}

// Write marshalls the value to Form and set the Content-Type Header.
func (e entityMultipartFormAccess) Write(resp *restful.Response, status int, v interface{}) error {
	return writeJSON(resp, status, restful.MIME_JSON, v)
}

type multipartRequest http.Request

// TrySet tries to set a value by the multipart request with the binding a form file
func (r *multipartRequest) TrySet(value reflect.Value, field reflect.StructField, key string, opt setOptions) (bool, error) {
	if files := r.MultipartForm.File[key]; len(files) != 0 {
		return setByMultipartFormFile(value, field, files)
	}

	return setByForm(value, field, r.MultipartForm.Value, key, opt)
}

func setByMultipartFormFile(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSet bool, err error) {
	switch value.Kind() {
	case reflect.Ptr:
		switch value.Interface().(type) {
		case *multipart.FileHeader:
			value.Set(reflect.ValueOf(files[0]))
			return true, nil
		}
	case reflect.Struct:
		switch value.Interface().(type) {
		case multipart.FileHeader:
			value.Set(reflect.ValueOf(*files[0]))
			return true, nil
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(value.Type(), len(files), len(files))
		isSet, err = setArrayOfMultipartFormFiles(slice, field, files)
		if err != nil || !isSet {
			return isSet, err
		}
		value.Set(slice)
		return true, nil
	case reflect.Array:
		return setArrayOfMultipartFormFiles(value, field, files)
	}
	return false, ErrMultiFileHeader
}

func setArrayOfMultipartFormFiles(value reflect.Value, field reflect.StructField, files []*multipart.FileHeader) (isSet bool, err error) {
	if value.Len() != len(files) {
		return false, ErrMultiFileHeaderLenInvalid
	}
	for i := range files {
		set, err := setByMultipartFormFile(value.Index(i), field, files[i:i+1])
		if err != nil || !set {
			return set, err
		}
	}
	return true, nil
}

func init() {
	restful.RegisterEntityAccessor(MIME_MULTIPART_FORM, NewEntityAccessorMultipartForm())
}
