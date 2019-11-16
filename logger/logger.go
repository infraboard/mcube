package logger

// Meta 日志元数据
type Meta = map[string]interface{}

// Logger 程序日志接口, 用于适配多种第三方日志插件
type Logger interface {
	StandardLogger
	FormatLogger
	WithMetaLogger
	RecoveryLogger

	// 用于创建子Logger
	Named(name string) Logger
	With(m Meta) Logger
}

// StandardLogger 标准的日志打印
type StandardLogger interface {
	Debug(msgs ...interface{})
	Info(msgs ...interface{})
	Warn(msgs ...interface{})
	Error(msgs ...interface{})
	Fatal(msgs ...interface{})
	Panic(msgs ...interface{})
}

// FormatLogger 携带format的日志打印
type FormatLogger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
}

// WithMetaLogger 携带额外的日志meta数据
type WithMetaLogger interface {
	Debugw(msg string, m Meta)
	Infow(msg string, m Meta)
	Warnw(msg string, m Meta)
	Errorw(msg string, m Meta)
	Fatalw(msg string, m Meta)
	Panicw(msg string, m Meta)
}

// RecoveryLogger 记录Panice的日志
type RecoveryLogger interface {
	Recover(msg string)
}
