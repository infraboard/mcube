package negotiator

import "encoding/json"

const (
	MIME_JSON MIME = "application/json"
)

type jsonImpl struct{}

func (i *jsonImpl) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (i *jsonImpl) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (i *jsonImpl) ContentType() MIME {
	return MIME_JSON
}

func init() {
	Registry(&jsonImpl{})
}
