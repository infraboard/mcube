package negotiator

type Negotiator interface {
	ContentType() MIME
	Decoder
	Encoder
}

type Decoder interface {
	Decode(data []byte, v any) error
}

type Encoder interface {
	Encode(v any) ([]byte, error)
}
