//go:generate  mcube enum -m
// Package enum_test for test
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
)

const (
	// E1 todo
	E1 Enum = iota
	// E2 (中文测试) todo
	E2
)

// Status AAA
// BBB
type Status uint

// Enum tood
type Enum uint
