//go:generate  mcube enum -m
package enum_test

const (
	// Running (running) todo
	Running Status = iota
	// Stopping (stopping) tdo
	Stopping
	// Stopped (stopped) todo
	Stopped
	// Canceled (canceled) todo
	Canceled

	Test11
)

const (
	// Running  todo
	E1 Enum = iota
	// Running (中文测试) todo
	E2
)

// Status AAA
// BBB
type Status uint

type Enum uint