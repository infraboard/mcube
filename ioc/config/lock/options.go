package lock

import "time"

func DefaultOptions() *Options {
	return &Options{
		RetryStrategy: LinearBackoff(1 * time.Second),
	}
}

// Options describe the options for the lock
type Options struct {
	// RetryStrategy allows to customise the lock retry strategy.
	// Default: do not retry
	RetryStrategy RetryStrategy

	// Metadata string.
	Metadata string

	// Token is a unique value that is used to identify the lock. By default, a random tokens are generated. Use this
	// option to provide a custom token instead.
	Token string
}

func (o *Options) getMetadata() string {
	if o != nil {
		return o.Metadata
	}
	return ""
}

func (o *Options) getToken() string {
	if o != nil {
		return o.Token
	}
	return ""
}

func (o *Options) getRetryStrategy() RetryStrategy {
	if o != nil && o.RetryStrategy != nil {
		return o.RetryStrategy
	}
	return NoRetry()
}
