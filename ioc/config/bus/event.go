package bus

type Event struct {
	Subject string
	Header  map[string][]string
	Data    []byte
}
