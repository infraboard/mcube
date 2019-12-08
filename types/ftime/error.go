package ftime

import "errors"

var (
	// ErrUnKnownFormatType format type 不支持
	ErrUnKnownFormatType = errors.New("unknown format type error")
	// ErrUnKnownTimestampLength timestamp length 不支持
	ErrUnKnownTimestampLength = errors.New("unknown timestamp length error")
)
