package negotiator

import "encoding/json"

type jsonImpl struct{}

func (i *jsonImpl) Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (i *jsonImpl) Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (i *jsonImpl) Name() string {
	return "json"
}

func init() {
	Registry(MIMEJSON, &jsonImpl{})
}
