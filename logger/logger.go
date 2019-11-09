package logger

// Meta 日志元数据
type Meta = map[string]interface{}

// Logger 程序日志接口, 用于适配多种第三方日志插件
type Logger interface {
	// 标准的日志打印
	Debug(msgs ...interface{})
	Info(msgs ...interface{})
	Warn(msgs ...interface{})
	Error(msgs ...interface{})
	Fatal(msgs ...interface{})
	Panic(msgs ...interface{})

	// 携带format的日志打印
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})

	// 携带额外的日志meta数据
	Debugw(msg string, m Meta)
	Infow(msg string, m Meta)
	Warnw(msg string, m Meta)
	Errorw(msg string, m Meta)
	Fatalw(msg string, m Meta)
	Panicw(msg string, m Meta)

	// 记录Panice的日志
	Recover(msg string)

	// 用于创建子Logger
	Named(name string) Logger
	With(m Meta) Logger
}
