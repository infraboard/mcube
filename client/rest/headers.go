package rest

const (
	ACCEPT_HEADER           = "Accept"
	CONTENT_TYPE_HEADER     = "Content-Type"
	CONTENT_ENCODING_HEADER = "Content-Encoding"
	AUTHORIZATION_HEADER    = "Authorization"
)

func HeaderFilterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}
