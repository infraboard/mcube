package negotiator

var negotiators map[MIME]Negotiator

func GetNegotiator(m string) Negotiator {
	n, ok := negotiators[MIME(m)]
	if !ok {
		return negotiators[MIMEJSON]
	}

	return n
}

type MIME string

const (
	MIMEJSON MIME = "json"
	MIMEXML  MIME = "xml"
	MIMEYAML MIME = "yaml"
)

func Registry(df MIME, n Negotiator) {
	negotiators[df] = n
}
