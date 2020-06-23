//go:generate mcube enum -m

package event

const (
	// Trace (trace)
	Trace Level = iota
	// Debug (debug)
	Debug
	// Info (info)
	Info
	// Warn (warn)
	Warn
	// Error (error)
	Error
	// Critical (critical)
	Critical
	// Disaster (disaster)
	Disaster
)

// Level 事件基本
type Level uint
