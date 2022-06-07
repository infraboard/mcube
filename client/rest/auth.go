package rest

type AuthType string

const (
	BearerToken AuthType = "bearer_token"
	BasicAuth   AuthType = "basic_auth"
)

type User struct {
	Username string
	Password string
}
