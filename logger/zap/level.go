package zap

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

// Level is a logging priority. Higher levels are more important.
type Level int8

// Logging levels.
const (
	DebugLevel Level = iota - 1
	InfoLevel
	WarnLevel
	ErrorLevel
	PanicLevel
	FatalLevel
)

var levelStrings = map[Level]string{
	DebugLevel: "debug",
	InfoLevel:  "info",
	WarnLevel:  "warning",
	ErrorLevel: "error",
	FatalLevel: "fatal",
	PanicLevel: "panic",
}

var zapLevels = map[Level]zapcore.Level{
	DebugLevel: zapcore.DebugLevel,
	InfoLevel:  zapcore.InfoLevel,
	WarnLevel:  zapcore.WarnLevel,
	ErrorLevel: zapcore.ErrorLevel,
	FatalLevel: zapcore.FatalLevel,
	PanicLevel: zapcore.PanicLevel,
}

// NewLevel 更加string生成Level
func NewLevel(str string) (Level, error) {
	var l Level
	str = strings.ToLower(str)
	for level, name := range levelStrings {
		if name == str {
			l = level
			return l, nil
		}
	}

	return l, errors.Errorf("invalid level '%v'", str)
}

// String returns the name of the logging level.
func (l Level) String() string {
	s, found := levelStrings[l]
	if found {
		return s
	}
	return fmt.Sprintf("Level(%d)", l)
}

// Enabled returns true if given level is enabled.
func (l Level) Enabled(level Level) bool {
	return level >= l
}

// Unpack unmarshals a level string to a Level. This implements
// ucfg.StringUnpacker.
func (l *Level) Unpack(str string) error {
	str = strings.ToLower(str)
	for level, name := range levelStrings {
		if name == str {
			*l = level
			return nil
		}
	}

	return errors.Errorf("invalid level '%v'", str)
}

// zapLevel 对应的zaplevel
func (l Level) zapLevel() zapcore.Level {
	z, found := zapLevels[l]
	if found {
		return z
	}
	return zapcore.InfoLevel
}
