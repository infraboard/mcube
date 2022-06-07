package rest

const (
	CONTENT_TYPE_HEADER  = "Content-Type"
	AUTHORIZATION_HEADER = "Authorization"
)

func FilterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}
