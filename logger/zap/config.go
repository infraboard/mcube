package zap

import (
	"time"

	"go.uber.org/zap"
)

// Config contains the configuration options for the logger. To create a Config
// from a common.Config use logp/config.Build.
type Config struct {
	Name      string      // Name of the Logger (for default file name).
	JSON      bool        // Write logs as JSON.
	Level     Level       // Logging level (error, warning, info, debug).
	Metas     []zap.Field // Root Logger Metas
	Selectors []string    // Selectors for debug level logging.

	ToStderr    bool
	ToSyslog    bool
	ToFiles     bool
	ToEventLog  bool
	toObserver  bool
	toIODiscard bool

	Files FileConfig

	addCaller   bool // Adds package and line number info to messages.
	development bool // Controls how DPanic behaves.
}

// FileConfig contains the configuration options for the file output.
type FileConfig struct {
	Path            string
	Name            string
	MaxSize         uint
	MaxBackups      uint
	Permissions     uint32
	Interval        time.Duration
	RotateOnStartup bool
	RedirectStderr  bool
}

var defaultConfig = Config{
	Level:   InfoLevel,
	ToFiles: true,
	Files: FileConfig{
		MaxSize:         10 * 1024 * 1024,
		MaxBackups:      7,
		Permissions:     0600,
		Interval:        0,
		RotateOnStartup: true,
	},
	addCaller: true,
}

// DefaultConfig returns the default config options.
func DefaultConfig() Config {
	return defaultConfig
}
