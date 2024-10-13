package datasource

type PROVIDER string

const (
	PROVIDER_MYSQL    PROVIDER = "mysql"
	PROVIDER_POSTGRES PROVIDER = "postgres"
	PROVIDER_SQLITE   PROVIDER = "sqlite"
)
