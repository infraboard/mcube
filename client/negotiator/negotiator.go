package negotiator

var negotiators = map[MIME]Negotiator{}

func GetNegotiator(m string) Negotiator {
	n, ok := negotiators[MIME(m)]
	if !ok {
		return negotiators[MIME_TEXT_PLAIN]
	}

	return n
}

type MIME string

func Registry(n Negotiator) {
	negotiators[n.ContentType()] = n
}
