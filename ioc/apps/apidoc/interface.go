package apidoc

const (
	AppName = "apidoc"
)

func DefaultApiDoc() *ApiDoc {
	return &ApiDoc{
		BasePath: "",
		JsonPath: "/swagger.json",
		UIPath:   "/ui.html",
	}
}

type ApiDoc struct {
	// Swagger API Doc BASE API URL路径
	BasePath string `json:"base_path" yaml:"base_path" toml:"base_path" env:"BASE_PATH"`
	// Swagger JSON API path
	JsonPath string `json:"json_path" yaml:"json_path" toml:"json_path" env:"JSON_PATH"`
	// Swagger UI path
	UIPath string `json:"ui_path" yaml:"ui_path" toml:"ui_path" env:"UI_PATH"`
}
