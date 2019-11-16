package mock

import (
	"bytes"
)

// NewStandardLogger todo
func NewStandardLogger() *StandardLogger {
	return &StandardLogger{
		Buffer: bytes.NewBuffer([]byte{}),
	}
}

// StandardLogger 用于单元测试
type StandardLogger struct {
	Buffer *bytes.Buffer
}

// Debug todo
func (l *StandardLogger) Debug(msgs ...interface{}) {
	for _, m := range msgs {
		l.Buffer.Write([]byte(m.(string)))
	}
}

// Info todo
func (l *StandardLogger) Info(msgs ...interface{}) {

}

// Warn todo
func (l *StandardLogger) Warn(msgs ...interface{}) {

}

// Error todo
func (l *StandardLogger) Error(msgs ...interface{}) {

}

// Fatal todo
func (l *StandardLogger) Fatal(msgs ...interface{}) {

}

// Panic todo
func (l *StandardLogger) Panic(msgs ...interface{}) {

}
