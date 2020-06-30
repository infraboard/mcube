package logger

// Logger 程序日志接口, 用于适配多种第三方日志插件
type Logger interface {
	StandardLogger
	FormatLogger
	WithMetaLogger
	RecoveryLogger
	CompatibleLogger

	// 用于创建子Logger
	Named(name string) Logger
	With(fields ...Field) Logger
}

// CompatibleLogger todo
type CompatibleLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
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
	Debugw(msg string, fields ...Field)
	Infow(msg string, fields ...Field)
	Warnw(msg string, fields ...Field)
	Errorw(msg string, fields ...Field)
	Fatalw(msg string, fields ...Field)
	Panicw(msg string, fields ...Field)
}

// RecoveryLogger 记录Panice的日志
type RecoveryLogger interface {
	Recover(msg string)
}
