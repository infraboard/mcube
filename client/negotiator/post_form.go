package negotiator

import (
	"fmt"
	"net/url"
)

const (
	MIME_POST_FORM  MIME = "application/x-www-form-urlencoded"
	MIME_TEXT_PLAIN MIME = "text/plain"
	MIME_TEXT_HTML       = "text/html"
)

func newTextImpl(mime MIME) *textImpl {
	return &textImpl{
		contentType: mime,
	}
}

type textImpl struct {
	contentType MIME
}

func (i *textImpl) Encode(v any) ([]byte, error) {
	switch d := v.(type) {
	case string:
		return []byte(d), nil
	case []byte:
		return d, nil
	case url.Values:
		return []byte(d.Encode()), nil
	default:
		return nil, fmt.Errorf("unknow type: %t", v)
	}
}

func (i *textImpl) Decode(data []byte, v any) error {
	switch i.ContentType() {
	case MIME_POST_FORM:
		d, err := url.ParseQuery(string(data))
		if err != nil {
			return err
		}
		uv, ok := v.(*url.Values)
		if ok {
			*uv = d
		}
	default:
		switch d := v.(type) {
		case *string:
			*d = string(data)
		case *[]byte:
			*d = data
		}
	}

	return nil
}

func (i *textImpl) ContentType() MIME {
	return i.contentType
}

func init() {
	Registry(newTextImpl(MIME_POST_FORM))
	Registry(newTextImpl(MIME_TEXT_PLAIN))
	Registry(newTextImpl(MIME_TEXT_HTML))
}
