package zap

import (
	"time"

	"go.uber.org/zap/zapcore"
)

var baseEncodingConfig = zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     cEncodeTime,
	EncodeDuration: zapcore.NanosDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeName:     zapcore.FullNameEncoder,
}

func buildEncoder(cfg Config) zapcore.Encoder {
	if cfg.JSON {
		return zapcore.NewJSONEncoder(jsonEncoderConfig())
	} else if cfg.ToSyslog {
		return zapcore.NewConsoleEncoder(syslogEncoderConfig())
	} else {
		return zapcore.NewConsoleEncoder(consoleEncoderConfig())
	}
}

func jsonEncoderConfig() zapcore.EncoderConfig {
	return baseEncodingConfig
}

const (
	TimeFieldFormat = "2006-01-02 15:04:05"
)

// cEncodeTime 自定义时间格式显示
func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(TimeFieldFormat))
}

func consoleEncoderConfig() zapcore.EncoderConfig {
	c := baseEncodingConfig
	c.EncodeLevel = zapcore.CapitalLevelEncoder
	c.EncodeName = bracketedNameEncoder
	return c
}

func syslogEncoderConfig() zapcore.EncoderConfig {
	c := consoleEncoderConfig()
	// Time is added by syslog.
	c.TimeKey = ""
	return c
}

func bracketedNameEncoder(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + loggerName + "]")
}
