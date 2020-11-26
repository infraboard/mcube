package kafka

import (
	"errors"
)

var (
	// ErrNoBroker todo
	ErrNoBroker = errors.New("connect kafka error, no broker here")
)
