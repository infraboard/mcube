package desense

type Desenser interface {
	DeSense(value string, args ...string) string
}

var desensers = map[string]Desenser{}

func Registry(name string, d Desenser) {
	desensers[name] = d
}

func Get(name string) Desenser {
	return desensers[name]
}

func Default() Desenser {
	return desensers[DefaultDesenser]
}
