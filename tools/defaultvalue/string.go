package defaultvalue

// WithDefault todo
func WithDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}

	return value
}
