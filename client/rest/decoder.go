package rest

type Decoder interface {
	Decode(data []byte, v any) error
}
