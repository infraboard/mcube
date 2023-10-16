package cache

func WithExpiration(expiration int) SetOption {
	return func(o *options) {
		o.expiration = expiration
	}
}

type options struct {
	expiration int
}

type SetOption func(*options)
