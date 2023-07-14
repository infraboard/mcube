package sense

type Desenser interface {
	DeSense(value string) string
}

func DeSense(value string) string {
	return DefaultDesenser.DeSense(value)
}
