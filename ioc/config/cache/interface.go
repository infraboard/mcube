package cache

type Cache interface {
	Set(key string, value any)
}

type SetOption struct {
}
